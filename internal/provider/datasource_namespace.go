package provider

import (
	"context"
	"fmt"
	"net/http"
	"slices"
	"strings"

	lakekeeper "github.com/baptistegh/go-lakekeeper/pkg/client"
	"github.com/baptistegh/go-lakekeeper/pkg/core"
	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource              = &lakekeeperNamespaceDataSource{}
	_ datasource.DataSourceWithConfigure = &lakekeeperNamespaceDataSource{}
)

func init() { registerDataSource(NewLakekeeperNamespaceDataSource) }

// NewLakekeeperNamespaceDataSource reads a namespace without managing its lifecycle.
func NewLakekeeperNamespaceDataSource() datasource.DataSource {
	return &lakekeeperNamespaceDataSource{}
}

type lakekeeperNamespaceDataSource struct{ client *lakekeeper.Client }
type lakekeeperNamespaceDataSourceModel struct {
	ID            types.String `tfsdk:"id"`
	ProjectID     types.String `tfsdk:"project_id"`
	WarehouseID   types.String `tfsdk:"warehouse_id"`
	Name          types.String `tfsdk:"name"`
	NamespaceID   types.String `tfsdk:"namespace_id"`
	ManagedAccess types.Bool   `tfsdk:"managed_access"`
}

func (d *lakekeeperNamespaceDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_namespace"
}
func (d *lakekeeperNamespaceDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Read an existing namespace by its project, warehouse UUID and dotted path. Verifies the warehouse belongs to the project. Requires OpenFGA to read managed access. No resources or permissions are changed.",
		Attributes: map[string]schema.Attribute{
			"id":             schema.StringAttribute{Computed: true, MarkdownDescription: "Composite project/warehouse/namespace UUID identity."},
			"project_id":     schema.StringAttribute{Required: true, MarkdownDescription: "Expected project UUID."},
			"warehouse_id":   schema.StringAttribute{Required: true, MarkdownDescription: "Exact warehouse UUID; names are not used to select warehouses."},
			"name":           schema.StringAttribute{Required: true, MarkdownDescription: "Dotted namespace path, including all parent levels."},
			"namespace_id":   schema.StringAttribute{Computed: true, MarkdownDescription: "Namespace UUID for permission assignments."},
			"managed_access": schema.BoolAttribute{Computed: true, MarkdownDescription: "Whether managed access is enabled directly on this namespace."},
		},
	}
}
func (d *lakekeeperNamespaceDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData != nil {
		d.client = req.ProviderData.(*LakekeeperDatasourceData).Client
	}
}
func (d *lakekeeperNamespaceDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state lakekeeperNamespaceDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	id, managed, err := readExistingNamespace(ctx, d.client, state.ProjectID.ValueString(), state.WarehouseID.ValueString(), state.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Cannot read namespace", err.Error())
		return
	}
	state.NamespaceID = types.StringValue(id)
	state.ManagedAccess = types.BoolValue(managed)
	state.ID = types.StringValue(state.ProjectID.ValueString() + "/" + state.WarehouseID.ValueString() + "/" + id)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func readExistingNamespace(ctx context.Context, client *lakekeeper.Client, project, warehouse, name string) (string, bool, error) {
	for _, id := range []string{project, warehouse} {
		if _, err := uuid.Parse(id); err != nil {
			return "", false, fmt.Errorf("project and warehouse must be UUIDs")
		}
	}
	parts := strings.Split(name, ".")
	for _, part := range parts {
		if part == "" || strings.ContainsAny(part, "/\x1f") {
			return "", false, fmt.Errorf("invalid namespace path")
		}
	}
	w, _, err := client.WarehouseV1(project).Get(ctx, warehouse)
	if err != nil {
		return "", false, err
	}
	if w.ID != warehouse || w.ProjectID != project || !w.IsActive() {
		return "", false, fmt.Errorf("warehouse is inactive or does not belong to the expected project")
	}
	// Use the verified UUID directly, including for nested namespaces. A warehouse
	// rename or a replacement with the same name cannot redirect this lookup.
	req, err := client.NewRequest(ctx, http.MethodGet, "/", nil, []core.RequestOptionFunc{core.WithHeader("x-project-id", project)})
	if err != nil {
		return "", false, err
	}
	req.URL.Path = strings.TrimSuffix(client.BaseURL().Path, "/management/v1") + "/catalog/v1/" + warehouse + "/namespaces/" + strings.Join(parts, "\x1f")
	req.URL.RawPath = ""
	var result struct {
		Namespace  []string          `json:"namespace"`
		Properties map[string]string `json:"properties"`
	}
	if _, err := client.Do(req, &result); err != nil {
		return "", false, err
	}
	if !slices.Equal(result.Namespace, parts) {
		return "", false, fmt.Errorf("catalog returned a different namespace path")
	}
	id := result.Properties["namespace_id"]
	if _, err := uuid.Parse(id); err != nil {
		return "", false, fmt.Errorf("catalog returned an invalid namespace_id")
	}
	var props struct {
		Managed *bool `json:"managed-access"`
	}
	if err := namespaceRequest(ctx, client, project, id, "", http.MethodGet, nil, &props); err != nil {
		return "", false, err
	}
	if props.Managed == nil {
		return "", false, fmt.Errorf("response has no managed-access field")
	}
	return id, *props.Managed, nil
}
