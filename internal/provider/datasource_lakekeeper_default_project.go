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
	_ datasource.DataSource              = &LakekeeperDefaultProjectDataSource{}
	_ datasource.DataSourceWithConfigure = &LakekeeperDefaultProjectDataSource{}
)

func init() {
	registerDataSource(NewLakekeeperDefaultProjectDataSource)
}

// NewLakekeeperDefaultProjectDataSource is a helper function to simplify the provider implementation.
func NewLakekeeperDefaultProjectDataSource() datasource.DataSource {
	return &LakekeeperDefaultProjectDataSource{}
}

// LakekeeperDefaultProjectDataSource is the data source implementation.
type LakekeeperDefaultProjectDataSource struct {
	client *lakekeeper.Client
}

// LakekeeperDefaultProjectDataSourceModel describes the data source data model.
type LakekeeperDefaultProjectDataSourceModel struct {
	ID   types.String `tfsdk:"id"`
	Name types.String `tfsdk:"name"`
}

// Metadata returns the data source type name.
func (d *LakekeeperDefaultProjectDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_default_project"
}

// Schema defines the schema for the data source.
func (d *LakekeeperDefaultProjectDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: `The ` + "`lakekeeper_default_project`" + ` data source retrieves information about the user's default project.

**Upstream API**: [Lakekeeper REST API docs](https://docs.lakekeeper.io/docs/nightly/api/management/#tag/project/operation/get_project)`,

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "The ID of the project.",
				Computed:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "The name of the project.",
				Computed:            true,
			},
		},
	}
}

// Configure adds the provider configured client to the data source.
func (d *LakekeeperDefaultProjectDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	datasource := req.ProviderData.(*LakekeeperDatasourceData)
	d.client = datasource.Client
}

// Read refreshes the Terraform state with the latest data.
func (d *LakekeeperDefaultProjectDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state LakekeeperDefaultProjectDataSourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Make API call to read default project
	project, _, err := d.client.ProjectV1().GetDefault(ctx)
	if err != nil {
		resp.Diagnostics.AddError("Lakekeeper API error occurred", fmt.Sprintf("Unable to read default project, %v", err))
		return
	}

	state.ID = types.StringValue(project.ID)
	state.Name = types.StringValue(project.Name)

	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}
