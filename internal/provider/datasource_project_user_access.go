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
					enumValuesList(managementv1.AllowedProjectActionEnumValues),
				Computed:    true,
				ElementType: types.StringType,
			},
		},
	}
}

// Configure adds the provider configured client to the data source.
func (d *LakekeeperProjectUserAccessDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	d.client = datasourceClient(req, resp)
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

	// The access endpoint is deprecated in favor of the authorizer-actions
	// endpoint, which only returns OpenFGA relations.
	access, _, err := d.client.PermissionsOpenfgaAPI.GetProjectAccessById(ctx, state.ProjectID.ValueString()).
		PrincipalUser(state.UserID.ValueString()).
		Execute()
	if err != nil {
		resp.Diagnostics.AddError("Lakekeeper API error occurred", fmt.Sprintf("Unable to read acess for user %s on project %s, %v", state.UserID.ValueString(), state.ID.ValueString(), apiError(err)))
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
