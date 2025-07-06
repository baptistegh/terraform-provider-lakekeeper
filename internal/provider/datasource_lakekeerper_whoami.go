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
	_ datasource.DataSource              = &LakekeeperWhoamiDataSource{}
	_ datasource.DataSourceWithConfigure = &LakekeeperWhoamiDataSource{}
)

func init() {
	registerDataSource(NewLakekeeperWhoamiDataSource)
}

// NewLakekeeperWhoamiDataSource is a helper function to simplify the provider implementation.
func NewLakekeeperWhoamiDataSource() datasource.DataSource {
	return &LakekeeperWhoamiDataSource{}
}

// LakekeeperWhoamiDataSource is the data source implementation.
type LakekeeperWhoamiDataSource struct {
	client *lakekeeper.Client
}

// Metadata returns the data source type name.
func (d *LakekeeperWhoamiDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_whoami"
}

// Schema defines the schema for the data source.
func (d *LakekeeperWhoamiDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: `The ` + "`lakekeeper_whoami`" + ` data source retrieves information about the current user.

**Upstream API**: [Lakekeeper REST API docs](https://docs.lakekeeper.io/docs/nightly/api/management/#tag/user/operation/whoami)`,

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "The ID of the current user.",
				Computed:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "The name of the current user.",
				Computed:            true,
			},
			"email": schema.StringAttribute{
				MarkdownDescription: "The email of the current user.",
				Computed:            true,
			},
			"user_type": schema.StringAttribute{
				MarkdownDescription: "The type of the current user (human/..)",
				Computed:            true,
			},
			"created_at": schema.StringAttribute{
				MarkdownDescription: "When the current user has been created.",
				Computed:            true,
			},
			"updated_at": schema.StringAttribute{
				MarkdownDescription: "When the current user has last been modified.",
				Computed:            true,
			},
			"last_updated_with": schema.StringAttribute{
				MarkdownDescription: "The endpoint who last modified the current user.",
				Computed:            true,
			},
		},
	}
}

// Configure adds the provider configured client to the data source.
func (d *LakekeeperWhoamiDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	datasource := req.ProviderData.(*LakekeeperDatasourceData)
	d.client = datasource.Client
}

// Read refreshes the Terraform state with the latest data.
func (d *LakekeeperWhoamiDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state LakekeeperUserDataSourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	user, _, err := d.client.User.Whoami(lakekeeper.WithContext(ctx))
	if err != nil {
		resp.Diagnostics.AddError("Lakekeeper API error occurred", fmt.Sprintf("Unable to read current user, %v", err))
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
