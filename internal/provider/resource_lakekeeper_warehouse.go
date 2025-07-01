package provider

import (
	"context"
	"fmt"

	"github.com/baptistegh/terraform-provider-lakekeeper/lakekeeper"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

// Ensure provider defined types fully satisfy framework interfaces
var (
	_ resource.Resource                = &lakekeeperWarehouseResource{}
	_ resource.ResourceWithConfigure   = &lakekeeperWarehouseResource{}
	_ resource.ResourceWithImportState = &lakekeeperWarehouseResource{}
)

func init() {
	registerResource(NewLakekeeperWarehouseResource)
}

// NewLakekeeperWarehouseResource is a helper function to simplify the provider implementation.
func NewLakekeeperWarehouseResource() resource.Resource {
	return &lakekeeperWarehouseResource{}
}

func (r *lakekeeperWarehouseResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_warehouse"
}

// lakekeeperWarehouseResource defines the resource implementation.
type lakekeeperWarehouseResource struct {
	client *lakekeeper.Client
}

// lakekeeperWarehouseResourceModel describes the resource data model.
// type lakekeeperWarehouseResourceModel struct {
// 	  ID        types.String `tfsdk:"id"`
// 	  Name      types.String `tfsdk:"name"`
// 	  ProjectID types.String `tfsdk:"project_id"` // Optional, if not provided, the default project will be used.
// 	  Protected types.Bool   `tfsdk:"protected"`
// 	  Active    types.Bool   `tfsdk:"active"`
// 	  // TODO: Design how we will handle StorageProfile and DeleteProfile
// }

func (r *lakekeeperWarehouseResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: fmt.Sprintf(`The ` + "`lakekeeper_warehouse`" + ` resource allows to manage the lifecycle of a lakekeeper warehouse.

**Upstream API**: [Lakekeeper REST API docs](https://docs.lakekeeper.io/docs/nightly/api/management/#tag/warehouse/operation/get_warehouse)`),

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "The ID the warehouse.",
				Computed:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Name of the warehouse.",
				Required:            true,
				Validators:          []validator.String{stringvalidator.LengthAtLeast(1)},
				PlanModifiers:       []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"project_id": schema.StringAttribute{
				MarkdownDescription: "The project ID to which the warehouse belongs. If not provided, the default project will be used.",
				Optional:            true,
				Computed:            true,
				Validators:          []validator.String{stringvalidator.LengthAtLeast(1)},
			},
			"protected": schema.BoolAttribute{
				MarkdownDescription: "Whether the warehouse is protected from being deleted.",
				Required:            true,
			},
			"active": schema.BoolAttribute{
				MarkdownDescription: "Whether the warehouse is active.",
				Required:            true,
			},
		},
	}
}

// Configure adds the provider configured client to the resource.
func (r *lakekeeperWarehouseResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	resourceData := req.ProviderData.(*LakekeeperResourceData)
	r.client = resourceData.Client
}

// Create creates a new upstream resources and adds it into the Terraform state.
func (r *lakekeeperWarehouseResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// TODO
}

// Read refreshes the Terraform state with the latest data.
func (r *lakekeeperWarehouseResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// TODO
}

// Updates updates the resource in-place.
func (r *lakekeeperWarehouseResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	resp.Diagnostics.AddError(
		"Provider Error, report upstream",
		"Somehow the resource was requested to perform an in-place upgrade which is not possible.",
	)
}

// Deletes removes the resource.
func (r *lakekeeperWarehouseResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// TODO
}

// ImportState imports the resource into the Terraform state.
func (r *lakekeeperWarehouseResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
