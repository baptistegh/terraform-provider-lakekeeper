package provider

import (
	"context"
	"fmt"

	permissionv1 "github.com/baptistegh/go-lakekeeper/pkg/apis/management/v1/permission"
	lakekeeper "github.com/baptistegh/go-lakekeeper/pkg/client"
	"github.com/baptistegh/go-lakekeeper/pkg/core"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
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
					"`create_namespace` `delete` `modify_storage` `modify_storage_credential` `get_config` " +
					"`get_metadata` `list_namespaces` `include_in_list` `deactivate` `activate` `rename` `list_deleted_tabulars` " +
					"`read_assignments` `grant_create` `grant_describe` `grant_modify` `grant_select` `grant_pass_grants` " +
					"`grant_manage_grants` `change_ownership`",
				Computed:    true,
				ElementType: types.StringType,
			},
		},
	}
}

// Configure adds the provider configured client to the data source.
func (d *LakekeeperWarehouseUserAccessDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	datasource := req.ProviderData.(*LakekeeperDatasourceData)
	d.client = datasource.Client
}

// Read refreshes the Terraform state with the latest data.
func (d *LakekeeperWarehouseUserAccessDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state lakekeeperWarehouseUserAccessDataSourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	state.ID = types.StringValue(fmt.Sprintf("%s:%s", state.WarehouseID.ValueString(), state.UserID.ValueString()))

	access, _, err := d.client.PermissionV1().WarehousePermission().GetAccess(ctx, state.WarehouseID.ValueString(), &permissionv1.GetWarehouseAccessOptions{
		PrincipalUser: core.Ptr(state.UserID.ValueString()),
	})
	if err != nil {
		resp.Diagnostics.AddError("Lakekeeper API error occurred", fmt.Sprintf("Unable to read acess for user %s on warehouse %s, %v", state.UserID.ValueString(), state.ID.ValueString(), err))
		return
	}

	actions, diags := types.SetValueFrom(ctx, types.StringType, access.AllowedActions)
	if diags.HasError() {
		resp.Diagnostics = append(resp.Diagnostics, diags...)
		return
	}

	state.AllowedActions = actions

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
