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
	_ datasource.DataSource              = &LakekeeperProjectRoleAccessDataSource{}
	_ datasource.DataSourceWithConfigure = &LakekeeperProjectRoleAccessDataSource{}
)

func init() {
	registerDataSource(NewLakekeeperProjectRoleAccessDataSource)
}

// NewLakekeeperProjectRoleAccessDataSource is a helper function to simplify the provider implementation.
func NewLakekeeperProjectRoleAccessDataSource() datasource.DataSource {
	return &LakekeeperProjectRoleAccessDataSource{}
}

// LakekeeperProjectRoleAccessDataSource is the data source implementation.
type LakekeeperProjectRoleAccessDataSource struct {
	client *lakekeeper.Client
}

// LakekeeperProjectRoleAccessDataSourceModel describes the data source data model.
type lakekeeperProjectRoleAccessDataSourceModel struct {
	ID             types.String `tfsdk:"id"`
	ProjectID      types.String `tfsdk:"project_id"`
	RoleID         types.String `tfsdk:"role_id"`
	AllowedActions types.Set    `tfsdk:"allowed_actions"`
}

// Metadata returns the data source type name.
func (d *LakekeeperProjectRoleAccessDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_project_role_access"
}

// Schema defines the schema for the data source.
func (d *LakekeeperProjectRoleAccessDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: fmt.Sprintf(`The ` + "`lakekeeper_project_role_access`" + ` data source retrieves the accesses a role can have on a project.

**Upstream API**: [Lakekeeper REST API docs](https://docs.lakekeeper.io/docs/nightly/api/management/#tag/permissions/operation/get_project_access_by_id)`),

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "The internal ID of this data source, in the form `{{project_id}}:{{role_id}}`.",
				Computed:            true,
			},
			"project_id": schema.StringAttribute{
				MarkdownDescription: "ID of the project.",
				Required:            true,
			},

			"role_id": schema.StringAttribute{
				MarkdownDescription: "ID of the role.",
				Required:            true,
			},
			"allowed_actions": schema.SetAttribute{
				MarkdownDescription: `List of the role's allowed actions on the project. The possible values are ` +
					"`assume` `can_grant_assignee` `can_change_ownership` `delete` `update` `read` `create_namespace`",
				Computed:    true,
				ElementType: types.StringType,
			},
		},
	}
}

// Configure adds the provider configured client to the data source.
func (d *LakekeeperProjectRoleAccessDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	datasource := req.ProviderData.(*LakekeeperDatasourceData)
	d.client = datasource.Client
}

// Read refreshes the Terraform state with the latest data.
func (d *LakekeeperProjectRoleAccessDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state lakekeeperProjectRoleAccessDataSourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	state.ID = types.StringValue(fmt.Sprintf("%s:%s", state.ProjectID.ValueString(), state.RoleID.ValueString()))

	access, _, err := d.client.PermissionV1().ProjectPermission().GetAccess(ctx, state.ProjectID.ValueString(), &permissionv1.GetProjectAccessOptions{
		PrincipalRole: core.Ptr(state.RoleID.ValueString()),
	})
	if err != nil {
		resp.Diagnostics.AddError("Lakekeeper API error occurred", fmt.Sprintf("Unable to read acess for role %s on project %s, %v", state.RoleID.ValueString(), state.ID.ValueString(), err))
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
