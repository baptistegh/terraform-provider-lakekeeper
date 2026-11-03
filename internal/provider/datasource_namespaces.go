package provider

import (
	"cmp"
	"context"
	"fmt"
	"net/http"
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
	Parent      types.String                    `tfsdk:"parent"`
	Namespaces  []lakekeeperNamespacesItemModel `tfsdk:"namespaces"`
}
type lakekeeperNamespacesItemModel struct {
	ID    types.String `tfsdk:"id"`
	Name  types.String `tfsdk:"name"`
	Parts []string     `tfsdk:"parts"`
}

func (d *lakekeeperNamespacesDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_namespaces"
}
func (d *lakekeeperNamespacesDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: `The ` + "`lakekeeper_namespaces`" + ` data source lists the direct child namespaces of a warehouse, or of a namespace in it, that the provider's identity can see, sorted by name. It follows every page of the Iceberg REST catalog list. With an authorizer such as OpenFGA, Lakekeeper filters the list by the caller's permissions.

**Upstream API**: [Lakekeeper REST API docs](https://docs.lakekeeper.io/docs/nightly/api/catalog/#tag/Catalog-API/operation/listNamespaces)`,
		Attributes: map[string]schema.Attribute{
			"warehouse_id": schema.StringAttribute{Required: true, MarkdownDescription: "Warehouse UUID."},
			"parent":       schema.StringAttribute{Optional: true, MarkdownDescription: "Dotted path of the namespace whose children to list. Lists the warehouse's top-level namespaces when unset."},
			"namespaces": schema.ListNestedAttribute{
				Computed:            true,
				MarkdownDescription: "Visible direct child namespaces, sorted by name.",
				NestedObject: schema.NestedAttributeObject{Attributes: map[string]schema.Attribute{
					"id":    schema.StringAttribute{Computed: true, MarkdownDescription: "Namespace UUID."},
					"name":  schema.StringAttribute{Computed: true, MarkdownDescription: "Dotted namespace path, including all parent levels."},
					"parts": schema.ListAttribute{Computed: true, ElementType: types.StringType, MarkdownDescription: "Namespace path levels."},
				}},
			},
		},
	}
}
func (d *lakekeeperNamespacesDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData != nil {
		d.client = req.ProviderData.(*LakekeeperDatasourceData).Client
	}
}
func (d *lakekeeperNamespacesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state lakekeeperNamespacesDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	namespaces, err := listNamespaces(ctx, d.client, state.WarehouseID.ValueString(), state.Parent.ValueString())
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

func listNamespaces(ctx context.Context, client *lakekeeper.Client, warehouse, parent string) ([]lakekeeperNamespacesItemModel, error) {
	if _, err := uuid.Parse(warehouse); err != nil {
		return nil, fmt.Errorf("warehouse must be a UUID")
	}
	opts := listNamespacesOptions{ReturnUUIDs: true}
	var parentParts []string
	if parent != "" {
		parentParts = strings.Split(parent, ".")
		for _, part := range parentParts {
			if part == "" || strings.ContainsAny(part, "/\x1f") {
				return nil, fmt.Errorf("invalid parent namespace path")
			}
		}
		opts.Parent = strings.Join(parentParts, "\x1f")
	}
	namespaces := []lakekeeperNamespacesItemModel{}
	seen := map[string]bool{}
	for {
		req, err := client.NewRequest(ctx, http.MethodGet, "/", opts, nil)
		if err != nil {
			return nil, err
		}
		req.URL.Path = strings.TrimSuffix(client.BaseURL().Path, "/management/v1") + "/catalog/v1/" + warehouse + "/namespaces"
		req.URL.RawPath = ""
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
			namespaces = append(namespaces, lakekeeperNamespacesItemModel{
				ID:    types.StringValue(page.NamespaceUUIDs[i]),
				Name:  types.StringValue(strings.Join(parts, ".")),
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
		return cmp.Compare(a.Name.ValueString(), b.Name.ValueString())
	})
	return namespaces, nil
}
