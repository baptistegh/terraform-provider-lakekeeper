package provider

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	v1 "github.com/lakekeeper/go-lakekeeper/pkg/apis/management/v1"
	lakekeeper "github.com/lakekeeper/go-lakekeeper/pkg/client"
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
				MarkdownDescription: "The ID of the role. in the form `{{project_id}}/{{role_id}}`",
				Computed:            true,
			},
			"role_id": schema.StringAttribute{
				MarkdownDescription: `The internal ID of the role. Configure exactly one of role_id or name.`,
				Optional:            true,
				Computed:            true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
					stringvalidator.ExactlyOneOf(path.MatchRoot("name")),
				},
			},
			"project_id": schema.StringAttribute{
				MarkdownDescription: `The ID of the project the role belongs to.`,
				Required:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Role name within project_id. Matching ignores case, as Lakekeeper does for role names. Name lookup requires one visible match and permission to list the project's roles. Configure exactly one of role_id or name.",
				Optional:            true,
				Computed:            true,
				Validators:          []validator.String{stringvalidator.LengthAtLeast(1)},
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "The description of the role.",
				Computed:            true,
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
func (d *LakekeeperRoleDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	d.client = datasourceClient(req, resp)
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

	var role *v1.Role
	var err error
	if state.RoleID.IsNull() {
		role, err = findRoleByName(ctx, d.client, projectID, state.Name.ValueString())
	} else {
		role, err = d.client.Roles.Get(ctx, projectID, id)
	}
	if err != nil {
		resp.Diagnostics.AddError("Lakekeeper API error occurred", fmt.Sprintf("Unable to read role %s in project %s, %v", id, projectID, err))
		return
	}
	if role.ProjectId != projectID {
		resp.Diagnostics.AddError("Role project mismatch", "The returned role does not belong to the requested project.")
		return
	}

	state.ID = types.StringValue(fmt.Sprintf("%s/%s", role.ProjectId, role.Id))
	state.RoleID = types.StringValue(role.Id)
	state.ProjectID = types.StringValue(role.ProjectId)
	if state.Name.IsNull() {
		state.Name = types.StringValue(role.Name)
	}
	state.CreatedAt = timeValue(role.CreatedAt)

	state.Description = types.StringPointerValue(role.Description.Get())
	state.UpdatedAt = timePointerValue(role.UpdatedAt.Get())

	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

func findRoleByName(ctx context.Context, client *lakekeeper.Client, projectID, name string) (*v1.Role, error) {
	// The role list has no name filter, so page through every role of the project.
	pageToken := ""
	seenTokens := map[string]bool{}
	var match *v1.Role
	for {
		list := client.RoleAPI.ListRoles(ctx).XProjectId(projectID)
		if pageToken != "" {
			list = list.PageToken(pageToken)
		}
		page, _, err := list.Execute()
		if err != nil {
			return nil, fmt.Errorf("list roles in project %s: %w", projectID, apiError(err))
		}
		for _, role := range page.Roles {
			if role.ProjectId != projectID {
				return nil, fmt.Errorf("role list returned a role from another project")
			}
			if !strings.EqualFold(role.Name, name) {
				continue
			}
			if match != nil && match.Id != role.Id {
				return nil, fmt.Errorf("multiple roles named %q in project %s; use role_id", name, projectID)
			}
			match = &role
		}
		token := page.GetNextPageToken()
		if token == "" {
			break
		}
		if seenTokens[token] {
			return nil, fmt.Errorf("role list repeated a page token")
		}
		seenTokens[token] = true
		pageToken = token
	}
	if match == nil {
		return nil, fmt.Errorf("no role named %q in project %s", name, projectID)
	}
	return match, nil
}
