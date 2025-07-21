package provider

import (
	"context"
	"fmt"

	lakekeeper "github.com/baptistegh/go-lakekeeper/pkg/client"
	tftypes "github.com/baptistegh/terraform-provider-lakekeeper/internal/provider/types"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &LakekeeperProjectAssignmentsDataSource{}
	_ datasource.DataSourceWithConfigure = &LakekeeperProjectAssignmentsDataSource{}
)

func init() {
	registerDataSource(NewLakekeeperProjectAssignmentsDataSource)
}

// NewLakekeeperProjectAssignmentsDataSource is a helper function to simplify the provider implementation.
func NewLakekeeperProjectAssignmentsDataSource() datasource.DataSource {
	return &LakekeeperProjectAssignmentsDataSource{}
}

// LakekeeperProjectAssignmentsDataSource is the data source implementation.
type LakekeeperProjectAssignmentsDataSource struct {
	client *lakekeeper.Client
}

// LakekeeperProjectAssignmentsDataSourceModel describes the data source data model.
type lakekeeperProjectAssignmentsDataSourceModel struct {
	ID          types.String         `tfsdk:"id"`
	Assignments []tftypes.Assignment `tfsdk:"assignments"`
}

// Metadata returns the data source type name.
func (d *LakekeeperProjectAssignmentsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_project_assignments"
}

// Schema defines the schema for the data source.
func (d *LakekeeperProjectAssignmentsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: fmt.Sprintf(`The ` + "`lakekeeper_project_assignments`" + ` data source retrieves information about the assignments of a project.

**Upstream API**: [Lakekeeper REST API docs](https://docs.lakekeeper.io/docs/nightly/api/management/#tag/permissions/operation/get_project_assignments_by_id)`),

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "ID of the project.",
				Required:            true,
			},
			"assignments": schema.ListNestedAttribute{
				MarkdownDescription: "List of assignments.",
				Computed:            true,
				NestedObject:        tftypes.AssignmentDataSourceType(),
			},
		},
	}
}

// Configure adds the provider configured client to the data source.
func (d *LakekeeperProjectAssignmentsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	datasource := req.ProviderData.(*LakekeeperDatasourceData)
	d.client = datasource.Client
}

// Read refreshes the Terraform state with the latest data.
func (d *LakekeeperProjectAssignmentsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state lakekeeperProjectAssignmentsDataSourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	id := state.ID.ValueString()
	state.ID = types.StringValue(id)

	assignments, _, err := d.client.PermissionV1().ProjectPermission().GetAssignments(ctx, id, nil)
	if err != nil {
		resp.Diagnostics.AddError("Lakekeeper API error occurred", fmt.Sprintf("Unable to read assignments for project %s, %v", state.ID.ValueString(), err))
		return
	}

	state.Assignments = make([]tftypes.Assignment, len(assignments.Assignments))
	for i, a := range assignments.Assignments {
		state.Assignments[i] = tftypes.Assignment{
			AssigneeID:   types.StringValue(a.Assignee.Value),
			AssigneeType: types.StringValue(string(a.Assignee.Type)),
			Assignment:   types.StringValue(string(a.Assignment)),
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
