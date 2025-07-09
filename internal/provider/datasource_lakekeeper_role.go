package provider

import (
	"context"
	"fmt"

	"github.com/baptistegh/terraform-provider-lakekeeper/lakekeeper"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &LakekeeperRoleDataSource{}
	_ datasource.DataSourceWithConfigure = &LakekeeperRoleDataSource{}
)

func init() {
	registerDataSource(NewLakekeeperRoleDataSource)
}

// NewLakekeeperRoleDataSource is a helper function to simplify the provider implementation.
func NewLakekeeperRoleDataSource() datasource.DataSource {
	return &LakekeeperRoleDataSource{}
}

// LakekeeperRoleDataSource is the data source implementation.
type LakekeeperRoleDataSource struct {
	client *lakekeeper.Client
}

// LakekeeperRoleDataSourceModel describes the data source data model.
type LakekeeperRoleDataSourceModel struct {
	ID          types.String `tfsdk:"id"`
	RoleID      types.String `tfsdk:"role_id"`
	ProjectID   types.String `tfsdk:"project_id"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`

	CreatedAt types.String `tfsdk:"created_at"`
	UpdatedAt types.String `tfsdk:"updated_at"`
}

// Metadata returns the data source type name.
func (d *LakekeeperRoleDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_role"
}

// Schema defines the schema for the data source.
func (d *LakekeeperRoleDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: fmt.Sprintf(`The ` + "`lakekeeper_role`" + ` data source retrieves information a lakekeeper role.

**Upstream API**: [Lakekeeper REST API docs](https://docs.lakekeeper.io/docs/nightly/api/management/#tag/role/operation/get_role)`),

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: `The ID of the role. In the form <project_id>:<role_id>`,
				Computed:            true,
			},
			"role_id": schema.StringAttribute{
				MarkdownDescription: `The internal ID of the role.`,
				Required:            true,
			},
			"project_id": schema.StringAttribute{
				MarkdownDescription: `The ID of the project the role belongs to.`,
				Required:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "The name of the role.",
				Computed:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "The description of the role.",
				Optional:            true,
			},
			"created_at": schema.StringAttribute{
				MarkdownDescription: "When the role has been created.",
				Computed:            true,
			},
			"updated_at": schema.StringAttribute{
				MarkdownDescription: "When the role has last been modified.",
				Computed:            true,
			},
		},
	}
}

// Configure adds the provider configured client to the data source.
func (d *LakekeeperRoleDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	datasource := req.ProviderData.(*LakekeeperDatasourceData)
	d.client = datasource.Client
}

// Read refreshes the Terraform state with the latest data.
func (d *LakekeeperRoleDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state LakekeeperRoleDataSourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	id := state.RoleID.ValueString()
	projectID := state.ProjectID.ValueString()

	role, _, err := d.client.Role.GetRole(id, projectID, lakekeeper.WithContext(ctx))
	if err != nil {
		resp.Diagnostics.AddError("Lakekeeper API error occurred", fmt.Sprintf("Unable to read role %s in project %s, %v", id, projectID, err))
		return
	}

	state.ID = types.StringValue(fmt.Sprintf("%s:%s", role.ProjectID, role.ID))
	state.RoleID = types.StringValue(role.ID)
	state.ProjectID = types.StringValue(role.ProjectID)
	state.Name = types.StringValue(role.Name)
	state.CreatedAt = types.StringValue(role.CreatedAt)

	state.Description = types.StringPointerValue(role.Description)
	state.UpdatedAt = types.StringPointerValue(role.UpdatedAt)

	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}
