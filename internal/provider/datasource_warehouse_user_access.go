package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	managementv1 "github.com/lakekeeper/go-lakekeeper/pkg/apis/management/v1"
	lakekeeper "github.com/lakekeeper/go-lakekeeper/pkg/client"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &LakekeeperWarehouseUserAccessDataSource{}
	_ datasource.DataSourceWithConfigure = &LakekeeperWarehouseUserAccessDataSource{}
)

func init() {
	registerDataSource(NewLakekeeperWarehouseUserAccessDataSource)
}

// NewLakekeeperWarehouseUserAccessDataSource is a helper function to simplify the provider implementation.
func NewLakekeeperWarehouseUserAccessDataSource() datasource.DataSource {
	return &LakekeeperWarehouseUserAccessDataSource{}
}

// LakekeeperWarehouseUserAccessDataSource is the data source implementation.
type LakekeeperWarehouseUserAccessDataSource struct {
	client *lakekeeper.Client
}

// LakekeeperWarehouseUserAccessDataSourceModel describes the data source data model.
type lakekeeperWarehouseUserAccessDataSourceModel struct {
	ID             types.String `tfsdk:"id"`
	WarehouseID    types.String `tfsdk:"warehouse_id"`
	UserID         types.String `tfsdk:"user_id"`
	AllowedActions types.Set    `tfsdk:"allowed_actions"`
}

// Metadata returns the data source type name.
func (d *LakekeeperWarehouseUserAccessDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_warehouse_user_access"
}

// Schema defines the schema for the data source.
func (d *LakekeeperWarehouseUserAccessDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: fmt.Sprintf(`The ` + "`lakekeeper_warehouse_user_access`" + ` data source retrieves the accesses a user can have on a warehouse.

**Upstream API**: [Lakekeeper REST API docs](https://docs.lakekeeper.io/docs/nightly/api/management/#tag/permissions/operation/get_warehouse_access_by_id)`),

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "The internal ID of this data source, in the form <warehouse_id>:<user_id>.",
				Computed:            true,
			},
			"warehouse_id": schema.StringAttribute{
				MarkdownDescription: "ID of the warehouse.",
				Required:            true,
			},

			"user_id": schema.StringAttribute{
				MarkdownDescription: "ID of the user.",
				Required:            true,
			},
			"allowed_actions": schema.SetAttribute{
				MarkdownDescription: `List of the user's allowed actions on the warehouse. The possible values are ` +
					enumValuesList(managementv1.AllowedWarehouseActionEnumValues),
				Computed:    true,
				ElementType: types.StringType,
			},
		},
	}
}

// Configure adds the provider configured client to the data source.
func (d *LakekeeperWarehouseUserAccessDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	d.client = datasourceClient(req, resp)
}

// Read refreshes the Terraform state with the latest data.
func (d *LakekeeperWarehouseUserAccessDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state lakekeeperWarehouseUserAccessDataSourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	state.ID = types.StringValue(fmt.Sprintf("%s/%s", state.WarehouseID.ValueString(), state.UserID.ValueString()))

	// The access endpoint is deprecated in favor of the authorizer-actions
	// endpoint, which only returns OpenFGA relations.
	access, err := d.client.Permissions.Warehouses.Access(ctx, state.WarehouseID.ValueString(), &lakekeeper.GetWarehouseAccessByIDOptions{PrincipalUser: state.UserID.ValueString()})
	if err != nil {
		resp.Diagnostics.AddError("Lakekeeper API error occurred", fmt.Sprintf("Unable to read access for user %s on warehouse %s, %v", state.UserID.ValueString(), state.WarehouseID.ValueString(), err))
		return
	}

	actions, diags := types.SetValueFrom(ctx, types.StringType, enumStrings(access.AllowedActions))
	if diags.HasError() {
		resp.Diagnostics = append(resp.Diagnostics, diags...)
		return
	}

	state.AllowedActions = actions

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
