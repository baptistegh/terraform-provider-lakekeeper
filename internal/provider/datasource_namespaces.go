package provider

import (
	"cmp"
	"context"
	"fmt"
	"slices"
	"strings"

	lakekeeper "github.com/baptistegh/go-lakekeeper/pkg/client"
	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource              = &lakekeeperNamespacesDataSource{}
	_ datasource.DataSourceWithConfigure = &lakekeeperNamespacesDataSource{}
)

func init() { registerDataSource(NewLakekeeperNamespacesDataSource) }

// NewLakekeeperNamespacesDataSource lists the direct children of a warehouse or namespace.
func NewLakekeeperNamespacesDataSource() datasource.DataSource {
	return &lakekeeperNamespacesDataSource{}
}

type lakekeeperNamespacesDataSource struct{ client *lakekeeper.Client }
type lakekeeperNamespacesDataSourceModel struct {
	WarehouseID types.String                    `tfsdk:"warehouse_id"`
	Parent      []string                        `tfsdk:"parent"`
	Namespaces  []lakekeeperNamespacesItemModel `tfsdk:"namespaces"`
}
type lakekeeperNamespacesItemModel struct {
	ID    types.String `tfsdk:"id"`
	Parts []string     `tfsdk:"parts"`
}

func (d *lakekeeperNamespacesDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_namespaces"
}
func (d *lakekeeperNamespacesDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: `The ` + "`lakekeeper_namespaces`" + ` data source lists the direct child namespaces of a warehouse, or of a namespace in it, that the provider's identity can see, sorted by path. It follows every page of the Iceberg REST catalog list. With an authorizer such as OpenFGA, Lakekeeper filters the list by the caller's permissions.

**Upstream API**: [Lakekeeper REST API docs](https://docs.lakekeeper.io/docs/nightly/api/catalog/#tag/Catalog-API/operation/listNamespaces)`,
		Attributes: map[string]schema.Attribute{
			"warehouse_id": schema.StringAttribute{Required: true, MarkdownDescription: "Warehouse UUID."},
			"parent":       schema.ListAttribute{Optional: true, ElementType: types.StringType, MarkdownDescription: "Namespace levels of the parent whose children to list. Lists the warehouse's top-level namespaces when unset. Each level must be non-empty."},
			"namespaces": schema.ListNestedAttribute{
				Computed:            true,
				MarkdownDescription: "Visible direct child namespaces, sorted by path.",
				NestedObject: schema.NestedAttributeObject{Attributes: map[string]schema.Attribute{
					"id":    schema.StringAttribute{Computed: true, MarkdownDescription: "Namespace UUID."},
					"parts": schema.ListAttribute{Computed: true, ElementType: types.StringType, MarkdownDescription: "Namespace path levels."},
				}},
			},
		},
	}
}
func (d *lakekeeperNamespacesDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	d.client = datasourceClient(req, resp)
}
func (d *lakekeeperNamespacesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state lakekeeperNamespacesDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	namespaces, err := listNamespaces(ctx, d.client, state.WarehouseID.ValueString(), state.Parent)
	if err != nil {
		resp.Diagnostics.AddError("Cannot list namespaces", err.Error())
		return
	}
	state.Namespaces = namespaces
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

type listNamespacesOptions struct {
	Parent      string `url:"parent,omitempty"`
	PageToken   string `url:"pageToken,omitempty"`
	ReturnUUIDs bool   `url:"returnUuids"`
}

func listNamespaces(ctx context.Context, client *lakekeeper.Client, warehouse string, parentParts []string) ([]lakekeeperNamespacesItemModel, error) {
	opts := listNamespacesOptions{ReturnUUIDs: true}
	if len(parentParts) > 0 {
		parsed, err := parseNamespacePath(parentParts)
		if err != nil {
			return nil, err
		}
		opts.Parent = parsed
	}
	namespaces := []lakekeeperNamespacesItemModel{}
	seen := map[string]bool{}
	for {
		req, err := newCatalogRequest(ctx, client, warehouse, "/namespaces", opts)
		if err != nil {
			return nil, err
		}
		var page struct {
			Namespaces     [][]string `json:"namespaces"`
			NamespaceUUIDs []string   `json:"namespace-uuids"`
			NextPageToken  *string    `json:"next-page-token"`
		}
		if _, err := client.Do(req, &page); err != nil {
			return nil, fmt.Errorf("list namespaces: %w", err)
		}
		if len(page.NamespaceUUIDs) != len(page.Namespaces) {
			return nil, fmt.Errorf("catalog returned %d namespace UUIDs for %d namespaces", len(page.NamespaceUUIDs), len(page.Namespaces))
		}
		for i, parts := range page.Namespaces {
			if len(parts) != len(parentParts)+1 || !slices.Equal(parts[:len(parentParts)], parentParts) {
				return nil, fmt.Errorf("catalog returned namespace %q outside the requested parent", strings.Join(parts, "."))
			}
			id := page.NamespaceUUIDs[i]
			if _, err := uuid.Parse(id); err != nil {
				return nil, fmt.Errorf("catalog returned invalid namespace UUID %q", id)
			}
			namespaces = append(namespaces, lakekeeperNamespacesItemModel{
				ID:    types.StringValue(id),
				Parts: parts,
			})
		}
		if page.NextPageToken == nil || *page.NextPageToken == "" {
			break
		}
		if seen[*page.NextPageToken] {
			return nil, fmt.Errorf("namespace list repeated a page token")
		}
		seen[*page.NextPageToken] = true
		opts.PageToken = *page.NextPageToken
	}
	slices.SortFunc(namespaces, func(a, b lakekeeperNamespacesItemModel) int {
		return cmp.Compare(strings.Join(a.Parts, "."), strings.Join(b.Parts, "."))
	})
	return namespaces, nil
}
