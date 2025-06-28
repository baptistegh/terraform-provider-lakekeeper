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
	_ datasource.DataSource              = &lakekeeperServerInfoDataSource{}
	_ datasource.DataSourceWithConfigure = &lakekeeperServerInfoDataSource{}
)

func init() {
	registerDataSource(NewLakekeeperServerInfoDataSource)
}

// NewlakekeeperApplicationDataSource is a helper function to simplify the provider implementation.
func NewLakekeeperServerInfoDataSource() datasource.DataSource {
	return &lakekeeperServerInfoDataSource{}
}

// lakekeeperMetadataDataSource is the data source implementation.
type lakekeeperServerInfoDataSource struct {
	client *lakekeeper.Client
}

// lakekeeperMetadataDataSourceModel describes the data source data model.
type lakekeeperServerInfoDataSourceModel struct {
	AuthzBackend                 types.String   `tfsdk:"authz_backend"`
	Bootstrapped                 types.Bool     `tfsdk:"bootstrapped"`
	DefaultProjectID             types.String   `tfsdk:"default_project_id"`
	AWSSystemIdentitiesEnabled   types.Bool     `tfsdk:"aws_system_identities_enabled"`
	AzureSystemIdentitiesEnabled types.Bool     `tfsdk:"azure_system_identities_enabled"`
	GCPSystemIdentitiesEnabled   types.Bool     `tfsdk:"gcp_system_identities_enabled"`
	ServerID                     types.String   `tfsdk:"server_id"`
	Version                      types.String   `tfsdk:"version"`
	Queues                       []types.String `tfsdk:"queues"`
}

// Metadata returns the data source type name.
func (d *lakekeeperServerInfoDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_server_info"
}

// Schema defines the schema for the data source.
func (d *lakekeeperServerInfoDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: `The ` + "`lakekeeper_server_info`" + ` data source retrieves information about a lakekeeper instance.

**Upstream API**: [Lakekeeper REST API docs](https://docs.lakekeeper.io/docs/nightly/api/management/#tag/server/operation/get_server_info)`,

		Attributes: map[string]schema.Attribute{
			"authz_backend": schema.StringAttribute{
				MarkdownDescription: "Authorization backend configured",
				Computed:            true,
			},
			"bootstrapped": schema.BoolAttribute{
				MarkdownDescription: "True if the server has been bootstrapped",
				Computed:            true,
			},
			"default_project_id": schema.StringAttribute{
				MarkdownDescription: "The default project ID",
				Computed:            true,
			},
			"aws_system_identities_enabled": schema.BoolAttribute{
				MarkdownDescription: "",
				Computed:            true,
			},
			"azure_system_identities_enabled": schema.BoolAttribute{
				MarkdownDescription: "",
				Computed:            true,
			},
			"gcp_system_identities_enabled": schema.BoolAttribute{
				MarkdownDescription: "",
				Computed:            true,
			},
			"server_id": schema.StringAttribute{
				MarkdownDescription: "The ID of the server",
				Computed:            true,
			},
			"version": schema.StringAttribute{
				MarkdownDescription: "The current version of the running server",
				Computed:            true,
			},
			"queues": schema.ListAttribute{
				ElementType:         types.StringType,
				MarkdownDescription: "",
				Computed:            true,
			},
		},
	}
}

// Configure adds the provider configured client to the data source.
func (d *lakekeeperServerInfoDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	datasource := req.ProviderData.(*LakekeeperDatasourceData)
	d.client = datasource.Client
}

// Read refreshes the Terraform state with the latest data.
func (d *lakekeeperServerInfoDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state lakekeeperServerInfoDataSourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}
	// Make API call to read applications
	serverInfo, err := d.client.GetServerInfo(ctx)
	if err != nil {
		resp.Diagnostics.AddError("Lakekeeper API error occurred", fmt.Sprintf("Unable to read server infos: %s", err.Error()))
		return
	}

	state.AuthzBackend = types.StringValue(serverInfo.AuthzBackend)
	state.Bootstrapped = types.BoolValue(serverInfo.Bootstrapped)
	state.DefaultProjectID = types.StringValue(serverInfo.DefaultProjectID)
	state.AWSSystemIdentitiesEnabled = types.BoolValue(serverInfo.AWSSystemIdentitiesEnabled)
	state.AzureSystemIdentitiesEnabled = types.BoolValue(serverInfo.AzureSystemIdentitiesEnabled)
	state.GCPSystemIdentitiesEnabled = types.BoolValue(serverInfo.GCPSystemIdentitiesEnabled)
	state.ServerID = types.StringValue(serverInfo.ServerID)
	state.Version = types.StringValue(serverInfo.Version)
	state.Queues = flattenQueues(serverInfo.Queues)

	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

func flattenQueues(queues []string) []types.String {
	var typedQueues = []types.String{}
	for _, q := range queues {
		typedQueues = append(typedQueues, types.StringValue(q))
	}
	return typedQueues
}
