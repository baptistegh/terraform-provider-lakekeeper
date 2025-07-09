package provider

import (
	"context"
	"fmt"

	v1 "github.com/baptistegh/go-lakekeeper/pkg/apis/v1"
	lakekeeper "github.com/baptistegh/go-lakekeeper/pkg/client"
	"github.com/baptistegh/go-lakekeeper/pkg/core"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces
var (
	_ resource.Resource                = &lakekeeperProjectResource{}
	_ resource.ResourceWithConfigure   = &lakekeeperProjectResource{}
	_ resource.ResourceWithImportState = &lakekeeperProjectResource{}
)

func init() {
	registerResource(NewLakekeeperProjectResource)
}

// NewLakekeeperProjectResource is a helper function to simplify the provider implementation.
func NewLakekeeperProjectResource() resource.Resource {
	return &lakekeeperProjectResource{}
}

func (r *lakekeeperProjectResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_project"
}

// lakekeeperProjectResource defines the resource implementation.
type lakekeeperProjectResource struct {
	client *lakekeeper.Client
}

// lakekeeperProjectResourceModel describes the resource data model.
type lakekeeperProjectResourceModel struct {
	ID   types.String `tfsdk:"id"`
	Name types.String `tfsdk:"name"`
}

func (r *lakekeeperProjectResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: fmt.Sprintf(`The ` + "`lakekeeper_project`" + ` resource allows to manage the lifecycle of a lakekeeper project.

**Upstream API**: [Lakekeeper REST API docs](https://docs.lakekeeper.io/docs/nightly/api/management/#tag/project)`),

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "The ID the project.",
				Computed:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Name of the project.",
				Required:            true,
				Validators:          []validator.String{stringvalidator.LengthAtLeast(1)},
			},
		},
	}
}

// Configure adds the provider configured client to the resource.
func (r *lakekeeperProjectResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	resourceData := req.ProviderData.(*LakekeeperResourceData)
	r.client = resourceData.Client
}

// Create creates a new upstream resources and adds it into the Terraform state.
func (r *lakekeeperProjectResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var state lakekeeperProjectResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	name := state.Name.ValueString()

	opts := v1.CreateProjectOptions{
		Name: name,
	}

	project, _, err := r.client.ProjectV1().Create(&opts, core.WithContext(ctx))
	if err != nil {
		resp.Diagnostics.AddError("Lakekeeper API error occurred", fmt.Sprintf("Unable to create project %s, %v", name, err.Error()))
		return
	}

	state.ID = types.StringValue(project.ID)

	// Log the creation of the resource
	tflog.Debug(ctx, "created an application", map[string]any{
		"name": state.Name.ValueString(), "id": state.ID.ValueString(),
	})

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Read refreshes the Terraform state with the latest data.
func (r *lakekeeperProjectResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state lakekeeperProjectResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	id := state.ID.ValueString()

	project, _, err := r.client.ProjectV1().Get(id, core.WithContext(ctx))
	if err != nil {
		resp.Diagnostics.AddError("Lakekeeper API error occurred", fmt.Sprintf("Unable to read project %s, %v", id, err.Error()))
		return
	}

	state.ID = types.StringValue(project.ID)
	state.Name = types.StringValue(project.Name)

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Updates updates the resource in-place.
func (r *lakekeeperProjectResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state lakekeeperProjectResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	if plan.Name.ValueString() != state.Name.ValueString() {
		tflog.Debug(ctx, "renaming project", map[string]any{
			"id":      state.ID.ValueString(),
			"oldName": state.Name.ValueString(),
			"newName": plan.Name.ValueString(),
		})

		opts := v1.RenameProjectOptions{
			NewName: plan.Name.ValueString(),
		}

		_, err := r.client.ProjectV1().Rename(state.ID.ValueString(), &opts, core.WithContext(ctx))
		if err != nil {
			resp.Diagnostics.AddError("Lakekeeper API error occurred", fmt.Sprintf("Unable to rename project %s, %v", state.ID.ValueString(), err.Error()))
			return
		}
	}

	// Update the state with the new name
	state.Name = plan.Name
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Deletes removes the resource.
func (r *lakekeeperProjectResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state lakekeeperProjectResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	id := state.ID.ValueString()

	_, err := r.client.ProjectV1().Delete(id, core.WithContext(ctx))
	if err != nil {
		resp.Diagnostics.AddError("Lakekeeper API error occurred", fmt.Sprintf("Unable to delete project %s, %v", id, err))
		return
	}

	resp.State.RemoveResource(ctx)
}

// ImportState imports the resource into the Terraform state.
func (r *lakekeeperProjectResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
