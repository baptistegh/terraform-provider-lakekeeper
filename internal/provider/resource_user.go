package provider

import (
	"context"
	"fmt"
	"regexp"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	managementv1 "github.com/lakekeeper/go-lakekeeper/pkg/apis/management/v1"
	lakekeeper "github.com/lakekeeper/go-lakekeeper/pkg/client"
)

// Ensure provider defined types fully satisfy framework interfaces
var (
	_ resource.Resource                = &lakekeeperUserResource{}
	_ resource.ResourceWithConfigure   = &lakekeeperUserResource{}
	_ resource.ResourceWithImportState = &lakekeeperUserResource{}
)

func init() {
	registerResource(NewLakekeeperUserResource)
}

// NewLakekeeperUserResource is a helper function to simplify the provider implementation.
func NewLakekeeperUserResource() resource.Resource {
	return &lakekeeperUserResource{}
}

func (r *lakekeeperUserResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_user"
}

// lakekeeperUserResource defines the resource implementation.
type lakekeeperUserResource struct {
	client *lakekeeper.Client
}

func (r *lakekeeperUserResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: fmt.Sprintf(`The ` + "`lakekeeper_user`" + ` resource allows to manage the lifecycle of a lakekeeper user.

**Upstream API**: [Lakekeeper REST API docs](https://docs.lakekeeper.io/docs/nightly/api/management/#tag/user)`),

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: `The ID of the user. The id must be identical to the subject in JWT tokens, prefixed with` + "`<idp-identifier>~`" + `. For example: ` + "`oidc~1234567890`" + " for OIDC users or `kubernetes~1234567890` for Kubernetes users.",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.RegexMatches(
						regexp.MustCompile("^[^~]+~.+"),
						"The id must be in the format `<idp-identifier>~<id>`, for example `oidc~1234567890`.",
					),
				},
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "The name of the user.",
				Required:            true,
			},
			"email": schema.StringAttribute{
				MarkdownDescription: "The email of the user.",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"user_type": schema.StringAttribute{
				MarkdownDescription: fmt.Sprintf("The type of the user, must be `%s` or `%s`", managementv1.USERTYPE_HUMAN, managementv1.USERTYPE_APPLICATION),
				Required:            true,
				Validators:          []validator.String{stringvalidator.OneOf(string(managementv1.USERTYPE_HUMAN), string(managementv1.USERTYPE_APPLICATION))},
			},
			"created_at": schema.StringAttribute{
				MarkdownDescription: "When the user has been created.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
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

// Configure adds the provider configured client to the resource.
func (r *lakekeeperUserResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.client = resourceClient(req, resp)
}

// Create creates a new upstream resources and adds it into the Terraform state.
func (r *lakekeeperUserResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var state LakekeeperUserDataSourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	user, err := r.client.Users.Create(ctx, provisionUserRequest(state, false))
	if err != nil {
		resp.Diagnostics.AddError("Lakekeeper API error occurred", fmt.Sprintf("Unable to create user %s, %v", state.ID.ValueString(), err))
		return
	}

	refreshUserState(&state, user)

	// Log the creation of the resource
	tflog.Debug(ctx, "created a user", map[string]any{
		"name": state.Name.ValueString(), "id": state.ID.ValueString(),
	})

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Read refreshes the Terraform state with the latest data.
func (r *lakekeeperUserResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state LakekeeperUserDataSourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	id := state.ID.ValueString()

	user, err := r.client.Users.Get(ctx, id)
	if err != nil {
		resp.Diagnostics.AddError("Lakekeeper API error occurred", fmt.Sprintf("Unable to read user %s, %v", id, err))
		return
	}

	refreshUserState(&state, user)

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Updates updates the resource in-place.
func (r *lakekeeperUserResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var state LakekeeperUserDataSourceModel
	var plan LakekeeperUserDataSourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	// Read Terraform state data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)

	if resp.Diagnostics.HasError() {
		return
	}

	if state.ID.IsNull() || state.ID.IsUnknown() || state.ID.ValueString() == "" {
		resp.Diagnostics.AddError("Incorrect resource definition", "Resource was requested to perform an in-place upgrade with a null ID.")
		return
	}

	user, err := r.client.Users.Create(ctx, provisionUserRequest(plan, true))
	if err != nil {
		resp.Diagnostics.AddError("Lakekeeper API error occurred", fmt.Sprintf("Unable to update user: %v", err))
		return
	}

	refreshUserState(&state, user)

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Deletes removes the resource.
func (r *lakekeeperUserResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state LakekeeperUserDataSourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	id := state.ID.ValueString()

	if id == "" {
		resp.Diagnostics.AddError("Incorrect state", "Unable to delete user without ID")
		return
	}

	err := r.client.Users.Delete(ctx, id)
	if err != nil {
		resp.Diagnostics.AddError("Lakekeeper API error occurred", fmt.Sprintf("Unable to delete user %s, %v", id, err))
		return
	}

	resp.State.RemoveResource(ctx)
}

// ImportState imports the resource into the Terraform state.
func (r *lakekeeperUserResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// provisionUserRequest builds the request that creates the user described by m,
// or updates it when it exists and updateIfExists is set.
func provisionUserRequest(m LakekeeperUserDataSourceModel, updateIfExists bool) *managementv1.CreateUserRequest {
	req := managementv1.NewCreateUserRequest()
	req.SetUpdateIfExists(updateIfExists)

	if id := m.ID.ValueStringPointer(); id != nil {
		req.SetId(*id)
	}
	if name := m.Name.ValueStringPointer(); name != nil {
		req.SetName(*name)
	}
	if email := m.Email.ValueStringPointer(); email != nil {
		req.SetEmail(*email)
	}
	if userType := m.UserType.ValueString(); userType != "" {
		req.SetUserType(managementv1.UserType(userType))
	}

	return req
}

// refreshUserState copies the attributes of user into state.
func refreshUserState(state *LakekeeperUserDataSourceModel, user *managementv1.User) {
	state.ID = types.StringValue(user.Id)
	state.Name = types.StringValue(user.Name)
	state.Email = types.StringPointerValue(user.Email.Get())
	state.UserType = types.StringValue(string(user.UserType))
	state.CreatedAt = timeValue(user.CreatedAt)
	state.UpdatedAt = timePointerValue(user.UpdatedAt.Get())
	state.LastUpdatedWith = types.StringValue(string(user.LastUpdatedWith))
}
