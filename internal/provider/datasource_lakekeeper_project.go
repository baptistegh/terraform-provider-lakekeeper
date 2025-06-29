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
	_ datasource.DataSource              = &LakekeeperProjectDataSource{}
	_ datasource.DataSourceWithConfigure = &LakekeeperProjectDataSource{}
)

func init() {
	registerDataSource(NewLakekeeperProjectDataSource)
}

// NewlakekeeperApplicationDataSource is a helper function to simplify the provider implementation.
func NewLakekeeperProjectDataSource() datasource.DataSource {
	return &LakekeeperProjectDataSource{}
}

// lakekeeperMetadataDataSource is the data source implementation.
type LakekeeperProjectDataSource struct {
	client *lakekeeper.Client
}

// lakekeeperMetadataDataSourceModel describes the data source data model.
type LakekeeperProjectDataSourceModel struct {
	ID   types.String `tfsdk:"id"`
	Name types.String `tfsdk:"name"`
}

// Metadata returns the data source type name.
func (d *LakekeeperProjectDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_project"
}

// Schema defines the schema for the data source.
func (d *LakekeeperProjectDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: `The ` + "`lakekeeper_project`" + ` data source retrieves information about a lakekeeper project.

**Upstream API**: [Lakekeeper REST API docs](https://docs.lakekeeper.io/docs/nightly/api/management/#tag/project/operation/get_default_project)`,

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "The ID of the project, it can be computed if the project name if used to find",
				Computed:            true,
				Optional:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "The project name, if the Name and the ID is given, the project will be found by ID",
				Computed:            true,
				Optional:            true,
			},
		},
	}
}

// Configure adds the provider configured client to the data source.
func (d *LakekeeperProjectDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	datasource := req.ProviderData.(*LakekeeperDatasourceData)
	d.client = datasource.Client
}

// Read refreshes the Terraform state with the latest data.
func (d *LakekeeperProjectDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state LakekeeperProjectDataSourceModel
	var project *lakekeeper.Project

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}
	// Make API call to read applications
	if !state.Name.IsUnknown() {
		var err error
		project, err = d.client.GetProjectByName(ctx, state.Name.ValueString())
		if err != nil {
			resp.Diagnostics.AddError("Lakekeeper API error occurred", fmt.Sprintf("Unable to read server project: %s", err.Error()))
			return
		}
	}

	if !state.ID.IsUnknown() {
		var err error
		project, err = d.client.GetProjectByID(ctx, state.ID.ValueString())
		if err != nil {
			resp.Diagnostics.AddError("Lakekeeper API error occurred", fmt.Sprintf("Unable to read server project: %s", err.Error()))
			return
		}
	}

	state.ID = types.StringValue(project.ID)
	state.Name = types.StringValue(project.Name)

	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}
