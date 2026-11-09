package provider

import (
	"cmp"
	"context"
	"fmt"
	"slices"

	lakekeeper "github.com/baptistegh/go-lakekeeper/pkg/client"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource              = &lakekeeperProjectsDataSource{}
	_ datasource.DataSourceWithConfigure = &lakekeeperProjectsDataSource{}
)

func init() { registerDataSource(NewLakekeeperProjectsDataSource) }

// NewLakekeeperProjectsDataSource lists the projects the caller can see.
func NewLakekeeperProjectsDataSource() datasource.DataSource {
	return &lakekeeperProjectsDataSource{}
}

type lakekeeperProjectsDataSource struct{ client *lakekeeper.Client }
type lakekeeperProjectsDataSourceModel struct {
	Projects []lakekeeperProjectsItemModel `tfsdk:"projects"`
}
type lakekeeperProjectsItemModel struct {
	ID   types.String `tfsdk:"id"`
	Name types.String `tfsdk:"name"`
}

func (d *lakekeeperProjectsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_projects"
}
func (d *lakekeeperProjectsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: `The ` + "`lakekeeper_projects`" + ` data source lists the projects the provider's identity can see, sorted by name. With an authorizer such as OpenFGA, Lakekeeper filters the list by the caller's permissions.

**Upstream API**: [Lakekeeper REST API docs](https://docs.lakekeeper.io/docs/nightly/api/management/#tag/project/operation/list_projects)`,
		Attributes: map[string]schema.Attribute{
			"projects": schema.ListNestedAttribute{
				Computed:            true,
				MarkdownDescription: "Visible projects, sorted by name.",
				NestedObject: schema.NestedAttributeObject{Attributes: map[string]schema.Attribute{
					"id":   schema.StringAttribute{Computed: true, MarkdownDescription: "Project ID."},
					"name": schema.StringAttribute{Computed: true, MarkdownDescription: "Project name."},
				}},
			},
		},
	}
}
func (d *lakekeeperProjectsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	d.client = datasourceClient(req, resp)
}
func (d *lakekeeperProjectsDataSource) Read(ctx context.Context, _ datasource.ReadRequest, resp *datasource.ReadResponse) {
	list, _, err := d.client.ProjectV1().List(ctx)
	if err != nil {
		resp.Diagnostics.AddError("Lakekeeper API error occurred", fmt.Sprintf("Unable to list projects, %v", err))
		return
	}
	state := lakekeeperProjectsDataSourceModel{Projects: []lakekeeperProjectsItemModel{}}
	for _, p := range list.Projects {
		state.Projects = append(state.Projects, lakekeeperProjectsItemModel{ID: types.StringValue(p.ID), Name: types.StringValue(p.Name)})
	}
	slices.SortFunc(state.Projects, func(a, b lakekeeperProjectsItemModel) int {
		return cmp.Or(cmp.Compare(a.Name.ValueString(), b.Name.ValueString()), cmp.Compare(a.ID.ValueString(), b.ID.ValueString()))
	})
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
