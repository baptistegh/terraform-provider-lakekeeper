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
	_ datasource.DataSource              = &LakekeeperProjectUserAccessDataSource{}
	_ datasource.DataSourceWithConfigure = &LakekeeperProjectUserAccessDataSource{}
)

func init() {
	registerDataSource(NewLakekeeperProjectUserAccessDataSource)
}

// NewLakekeeperProjectUserAccessDataSource is a helper function to simplify the provider implementation.
func NewLakekeeperProjectUserAccessDataSource() datasource.DataSource {
	return &LakekeeperProjectUserAccessDataSource{}
}

// LakekeeperProjectUserAccessDataSource is the data source implementation.
type LakekeeperProjectUserAccessDataSource struct {
	client *lakekeeper.Client
}

// LakekeeperProjectUserAccessDataSourceModel describes the data source data model.
type lakekeeperProjectUserAccessDataSourceModel struct {
	ID             types.String `tfsdk:"id"`
	ProjectID      types.String `tfsdk:"project_id"`
	UserID         types.String `tfsdk:"user_id"`
	AllowedActions types.Set    `tfsdk:"allowed_actions"`
}

// Metadata returns the data source type name.
func (d *LakekeeperProjectUserAccessDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_project_user_access"
}

// Schema defines the schema for the data source.
func (d *LakekeeperProjectUserAccessDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: fmt.Sprintf(`The ` + "`lakekeeper_project_user_access`" + ` data source retrieves the accesses a user can have on a project.

**Upstream API**: [Lakekeeper REST API docs](https://docs.lakekeeper.io/docs/nightly/api/management/#tag/permissions/operation/get_project_access_by_id)`),

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "The internal ID of this data source, in the form `{{project_id}}/{{user_id}}`.",
				Computed:            true,
			},
			"project_id": schema.StringAttribute{
				MarkdownDescription: "ID of the project.",
				Required:            true,
			},

			"user_id": schema.StringAttribute{
				MarkdownDescription: "ID of the user.",
				Required:            true,
			},
			"allowed_actions": schema.SetAttribute{
				MarkdownDescription: `List of the user's allowed actions on the project. The possible values are ` +
					fmt.Sprintf(
						"`%s` `%s` `%s` `%s` `%s` `%s` `%s` `%s` `%s` `%s` `%s` `%s` `%s` `%s` `%s` `%s`",
						permissionv1.CreateWarehouse,
						permissionv1.DeleteProject,
						permissionv1.RenameProject,
						permissionv1.ListWarehouses,
						permissionv1.CreateRole,
						permissionv1.ListRoles,
						permissionv1.SearchRoles,
						permissionv1.ReadProjectAssignments,
						permissionv1.GrantProjectRoleCreator,
						permissionv1.GrantProjectCreate,
						permissionv1.GrantProjectDescribe,
						permissionv1.GrantProjectModify,
						permissionv1.GrantProjectSelet,
						permissionv1.GrantProjectAdmin,
						permissionv1.GrantSecurityAdmin,
						permissionv1.GrantDataAdmin,
					),
				Computed:    true,
				ElementType: types.StringType,
			},
		},
	}
}

// Configure adds the provider configured client to the data source.
func (d *LakekeeperProjectUserAccessDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	datasource := req.ProviderData.(*LakekeeperDatasourceData)
	d.client = datasource.Client
}

// Read refreshes the Terraform state with the latest data.
func (d *LakekeeperProjectUserAccessDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state lakekeeperProjectUserAccessDataSourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	state.ID = types.StringValue(fmt.Sprintf("%s/%s", state.ProjectID.ValueString(), state.UserID.ValueString()))

	access, _, err := d.client.PermissionV1().ProjectPermission().GetAccess(ctx, state.ProjectID.ValueString(), &permissionv1.GetProjectAccessOptions{
		PrincipalUser: core.Ptr(state.UserID.ValueString()),
	})
	if err != nil {
		resp.Diagnostics.AddError("Lakekeeper API error occurred", fmt.Sprintf("Unable to read acess for user %s on project %s, %v", state.UserID.ValueString(), state.ID.ValueString(), err))
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
