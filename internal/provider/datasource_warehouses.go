package provider

import (
	"cmp"
	"context"
	"fmt"
	"slices"

	lakekeeper "github.com/baptistegh/go-lakekeeper/pkg/client"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource              = &lakekeeperWarehousesDataSource{}
	_ datasource.DataSourceWithConfigure = &lakekeeperWarehousesDataSource{}
)

func init() { registerDataSource(NewLakekeeperWarehousesDataSource) }

// NewLakekeeperWarehousesDataSource lists a project's warehouses the caller can see.
func NewLakekeeperWarehousesDataSource() datasource.DataSource {
	return &lakekeeperWarehousesDataSource{}
}

type lakekeeperWarehousesDataSource struct{ client *lakekeeper.Client }
type lakekeeperWarehousesDataSourceModel struct {
	ProjectID  types.String                    `tfsdk:"project_id"`
	Warehouses []lakekeeperWarehousesItemModel `tfsdk:"warehouses"`
}
type lakekeeperWarehousesItemModel struct {
	ID        types.String `tfsdk:"id"`
	Name      types.String `tfsdk:"name"`
	Active    types.Bool   `tfsdk:"active"`
	Protected types.Bool   `tfsdk:"protected"`
}

func (d *lakekeeperWarehousesDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_warehouses"
}
func (d *lakekeeperWarehousesDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: `The ` + "`lakekeeper_warehouses`" + ` data source lists a project's warehouses that the provider's identity can see, sorted by name. With an authorizer such as OpenFGA, Lakekeeper filters the list by the caller's permissions. The request sets no status filter, so the list contains whatever Lakekeeper returns by default; use ` + "`active`" + ` to filter.

**Upstream API**: [Lakekeeper REST API docs](https://docs.lakekeeper.io/docs/nightly/api/management/#tag/warehouse/operation/list_warehouses)`,
		Attributes: map[string]schema.Attribute{
			"project_id": schema.StringAttribute{Required: true, MarkdownDescription: "Project ID."},
			"warehouses": schema.ListNestedAttribute{
				Computed:            true,
				MarkdownDescription: "Visible warehouses in the project, sorted by name.",
				NestedObject: schema.NestedAttributeObject{Attributes: map[string]schema.Attribute{
					"id":        schema.StringAttribute{Computed: true, MarkdownDescription: "Warehouse ID."},
					"name":      schema.StringAttribute{Computed: true, MarkdownDescription: "Warehouse name."},
					"active":    schema.BoolAttribute{Computed: true, MarkdownDescription: "Whether the warehouse status is `active`."},
					"protected": schema.BoolAttribute{Computed: true, MarkdownDescription: "Whether the warehouse is protected from deletion."},
				}},
			},
		},
	}
}
func (d *lakekeeperWarehousesDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData != nil {
		d.client = req.ProviderData.(*LakekeeperDatasourceData).Client
	}
}
func (d *lakekeeperWarehousesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state lakekeeperWarehousesDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	project := state.ProjectID.ValueString()
	list, _, err := d.client.WarehouseV1(project).List(ctx, nil)
	if err != nil {
		resp.Diagnostics.AddError("Lakekeeper API error occurred", fmt.Sprintf("Unable to list warehouses in project %s, %v", project, err))
		return
	}
	state.Warehouses = []lakekeeperWarehousesItemModel{}
	for _, w := range list.Warehouses {
		if w.ProjectID != project {
			resp.Diagnostics.AddError("Unexpected warehouse", fmt.Sprintf("Lakekeeper listed warehouse %s from project %s for project %s", w.ID, w.ProjectID, project))
			return
		}
		state.Warehouses = append(state.Warehouses, lakekeeperWarehousesItemModel{
			ID:        types.StringValue(w.ID),
			Name:      types.StringValue(w.Name),
			Active:    types.BoolValue(w.IsActive()),
			Protected: types.BoolValue(w.Protected),
		})
	}
	slices.SortFunc(state.Warehouses, func(a, b lakekeeperWarehousesItemModel) int {
		return cmp.Or(cmp.Compare(a.Name.ValueString(), b.Name.ValueString()), cmp.Compare(a.ID.ValueString(), b.ID.ValueString()))
	})
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
