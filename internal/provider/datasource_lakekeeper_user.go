package provider

import (
	"context"
	"fmt"

	lakekeeper "github.com/baptistegh/go-lakekeeper/pkg/client"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &LakekeeperUserDataSource{}
	_ datasource.DataSourceWithConfigure = &LakekeeperUserDataSource{}
)

func init() {
	registerDataSource(NewLakekeeperUserDataSource)
}

// NewLakekeeperUserDataSource is a helper function to simplify the provider implementation.
func NewLakekeeperUserDataSource() datasource.DataSource {
	return &LakekeeperUserDataSource{}
}

// LakekeeperUserDataSource is the data source implementation.
type LakekeeperUserDataSource struct {
	client *lakekeeper.Client
}

// LakekeeperUserDataSourceModel describes the data source data model.
type LakekeeperUserDataSourceModel struct {
	ID              types.String `tfsdk:"id"`
	Name            types.String `tfsdk:"name"`
	Email           types.String `tfsdk:"email"`
	UserType        types.String `tfsdk:"user_type"`
	CreatedAt       types.String `tfsdk:"created_at"`
	UpdatedAt       types.String `tfsdk:"updated_at"`
	LastUpdatedWith types.String `tfsdk:"last_updated_with"`
}

// Metadata returns the data source type name.
func (d *LakekeeperUserDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_user"
}

// Schema defines the schema for the data source.
func (d *LakekeeperUserDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: `The ` + "`lakekeeper_user`" + ` data source retrieves information about a user.

**Upstream API**: [Lakekeeper REST API docs](https://docs.lakekeeper.io/docs/nightly/api/management/#tag/user/operation/get_user)`,

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: `The ID of the user. The id must be identical to the subject in JWT tokens, prefixed with` + "`<idp-identifier>~`" + `. For example: ` + "`oidc~1234567890`" + ` for OIDC users or kubernetes~1234567890 for Kubernetes users.`,
				Required:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "The name of the user.",
				Computed:            true,
			},
			"email": schema.StringAttribute{
				MarkdownDescription: "The email of the user.",
				Computed:            true,
			},
			"user_type": schema.StringAttribute{
				MarkdownDescription: "The type of the user (human/..)",
				Computed:            true,
			},
			"created_at": schema.StringAttribute{
				MarkdownDescription: "When the user has been created.",
				Computed:            true,
			},
			"updated_at": schema.StringAttribute{
				MarkdownDescription: "When the user has last been modified.",
				Computed:            true,
			},
			"last_updated_with": schema.StringAttribute{
				MarkdownDescription: "The endpoint who last modified the user.",
				Computed:            true,
			},
		},
	}
}

// Configure adds the provider configured client to the data source.
func (d *LakekeeperUserDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	datasource := req.ProviderData.(*LakekeeperDatasourceData)
	d.client = datasource.Client
}

// Read refreshes the Terraform state with the latest data.
func (d *LakekeeperUserDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state LakekeeperUserDataSourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	user, _, err := d.client.UserV1().Get(ctx, state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Lakekeeper API error occurred", fmt.Sprintf("Unable to read user %s, %v", state.ID.ValueString(), err))
		return
	}

	state.ID = types.StringValue(user.ID)
	state.Name = types.StringValue(user.Name)
	state.Email = types.StringPointerValue(user.Email)
	state.UserType = types.StringValue(string(user.UserType))
	state.CreatedAt = types.StringValue(user.CreatedAt)
	state.UpdatedAt = types.StringPointerValue(user.UpdatedAt)
	state.LastUpdatedWith = types.StringValue(user.LastUpdatedWith)

	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}
