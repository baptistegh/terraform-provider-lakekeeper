package provider

import (
	"context"
	"fmt"

	lakekeeper "github.com/baptistegh/go-lakekeeper/pkg/client"
	"github.com/baptistegh/go-lakekeeper/pkg/core"
	tftypes "github.com/baptistegh/terraform-provider-lakekeeper/internal/provider/types"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &LakekeeperWarehouseDataSource{}
	_ datasource.DataSourceWithConfigure = &LakekeeperWarehouseDataSource{}
)

func init() {
	registerDataSource(NewLakekeeperWarehouseDataSource)
}

// NewLakekeeperWarehouseDataSource is a helper function to simplify the provider implementation.
func NewLakekeeperWarehouseDataSource() datasource.DataSource {
	return &LakekeeperWarehouseDataSource{}
}

// LakekeeperWarehouseDataSource is the data source implementation.
type LakekeeperWarehouseDataSource struct {
	client *lakekeeper.Client
}

// LakekeeperWarehouseDataSourceModel describes the data source data model.
type lakekeeperWarehouseDataSourceModel struct {
	ID             types.String                 `tfsdk:"id"` // form: project_id:warehouse_id (internal ID)
	WarehouseID    types.String                 `tfsdk:"warehouse_id"`
	Name           types.String                 `tfsdk:"name"`
	ProjectID      types.String                 `tfsdk:"project_id"` // Optional, if not provided, the default project will be used.
	Protected      types.Bool                   `tfsdk:"protected"`
	Active         types.Bool                   `tfsdk:"active"`
	StorageProfile *tftypes.StorageProfileModel `tfsdk:"storage_profile"`
	DeleteProfile  *tftypes.DeleteProfileModel  `tfsdk:"delete_profile"`
}

// Metadata returns the data source type name.
func (d *LakekeeperWarehouseDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_warehouse"
}

// Schema defines the schema for the data source.
func (d *LakekeeperWarehouseDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: fmt.Sprintf(`The ` + "`lakekeeper_warehouse`" + ` data source retrieves information a lakekeeper warehouse.

**Currently the datasource can only read from the default project**

**Upstream API**: [Lakekeeper REST API docs](https://docs.lakekeeper.io/docs/nightly/api/management/#tag/warehouse/operation/get_warehouse)`),

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "The internal ID the warehouse. In the form: <project_id>:<warehouse_id>",
				Computed:            true,
			},
			"warehouse_id": schema.StringAttribute{
				MarkdownDescription: "The ID the warehouse.",
				Required:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Name of the warehouse.",
				Computed:            true,
			},
			"project_id": schema.StringAttribute{
				MarkdownDescription: "The project ID to which the warehouse belongs. If not provided, the default project will be used.",
				Required:            true,
			},
			"protected": schema.BoolAttribute{
				MarkdownDescription: "Whether the warehouse is protected from being deleted.",
				Computed:            true,
			},
			"active": schema.BoolAttribute{
				MarkdownDescription: "Whether the warehouse is active.",
				Computed:            true,
			},
			"storage_profile": tftypes.StorageProfileDatasourceSchema(),
			"delete_profile":  tftypes.DeleteProfileDatasourceSchema(),
		},
	}
}

// Configure adds the provider configured client to the data source.
func (d *LakekeeperWarehouseDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	datasource := req.ProviderData.(*LakekeeperDatasourceData)
	d.client = datasource.Client
}

// Read refreshes the Terraform state with the latest data.
func (d *LakekeeperWarehouseDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state lakekeeperWarehouseDataSourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	id := state.WarehouseID.ValueString()
	projectID := state.ProjectID.ValueString()

	warehouse, _, err := d.client.WarehouseV1(projectID).Get(id, core.WithContext(ctx))
	if err != nil {
		resp.Diagnostics.AddError("Lakekeeper API error occurred", fmt.Sprintf("Unable to read warehouse %s, %v", state.Name.ValueString(), err))
		return
	}

	diags := state.RefreshFromSettings(warehouse)
	if diags.HasError() {
		resp.Diagnostics = append(resp.Diagnostics, diags...)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
