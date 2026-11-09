package provider

import (
	"context"
	"regexp"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	managementv1 "github.com/lakekeeper/go-lakekeeper/pkg/apis/management/v1"
	lakekeeper "github.com/lakekeeper/go-lakekeeper/pkg/client"
	"github.com/lakekeeper/go-lakekeeper/pkg/permissions"
)

func init() { registerResource(NewLakekeeperNamespaceRoleAssignment) }
func NewLakekeeperNamespaceRoleAssignment() resource.Resource {
	return &namespaceRoleAssignmentResource{}
}

type namespaceRoleAssignmentResource struct{ client *lakekeeper.Client }
type namespaceRoleAssignmentModel struct {
	ID          types.String `tfsdk:"id"`
	ProjectID   types.String `tfsdk:"project_id"`
	NamespaceID types.String `tfsdk:"namespace_id"`
	RoleID      types.String `tfsdk:"role_id"`
	Assignments types.Set    `tfsdk:"assignments"`
}

func (r *namespaceRoleAssignmentResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_namespace_role_assignment"
}
func (r *namespaceRoleAssignmentResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	uuid := stringvalidator.RegexMatches(regexp.MustCompile(`^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`), "must be a UUID")
	identity := func(description string) schema.StringAttribute {
		return schema.StringAttribute{
			Required: true, MarkdownDescription: description, Validators: []validator.String{uuid},
			PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
		}
	}
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages all direct namespace permissions for one role. Other roles and users are preserved. Use one resource per namespace-role pair and import existing grants before managing them.\n\n" +
			"Requires Lakekeeper's OpenFGA authorization backend. Namespace grants apply to descendants according to Lakekeeper's authorization model.\n\n" +
			"Updates revoke obsolete permissions before adding new ones. A failed write can leave partial progress; refresh and reapply. Failed reads and deletes return errors and preserve state. A missing permissions endpoint is an error, including with allow-all.",
		Attributes: map[string]schema.Attribute{
			"id":           schema.StringAttribute{Computed: true, MarkdownDescription: "project_id/namespace_id/role_id", PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()}},
			"project_id":   identity("Project containing both the namespace and role."),
			"namespace_id": identity("Namespace UUID, available as lakekeeper_namespace.example.namespace_id."),
			"role_id":      identity("Role UUID."),
			"assignments": schema.SetAttribute{Required: true, ElementType: types.StringType,
				MarkdownDescription: "Direct permissions for this role, among " + enumValuesList(managementv1.AllowedNamespaceRelationEnumValues) + ". An empty set revokes them. Removing this resource revokes its permissions.",
				Validators:          []validator.Set{setvalidator.ValueStringsAre(stringvalidator.OneOf(enumStrings(managementv1.AllowedNamespaceRelationEnumValues)...))},
			},
		},
	}
}
func (r *namespaceRoleAssignmentResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.client = resourceClient(req, resp)
}

// namespaceRoleAssignments returns the sorted direct grants of a role on a namespace.
func namespaceRoleAssignments(ctx context.Context, client *lakekeeper.Client, namespace, role string) ([]string, error) {
	return namespaceAssignmentScope(client, namespace).relations(ctx, permissions.PrincipalRole, role)
}

// Reconcile only this role's direct grants. Delete before adding; refresh recovers a partial apply.
func setNamespaceRoleAssignments(ctx context.Context, client *lakekeeper.Client, namespace, role string, desired []string) error {
	current, err := namespaceRoleAssignments(ctx, client, namespace, role)
	if err != nil {
		return err
	}

	added, removed := diffStrings(current, desired)
	scope := namespaceAssignmentScope(client, namespace)
	if err := scope.apply(ctx, permissions.PrincipalRole, role, nil, removed); err != nil {
		return err
	}

	return scope.apply(ctx, permissions.PrincipalRole, role, added, nil)
}

func (r *namespaceRoleAssignmentResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan namespaceRoleAssignmentModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	// Existing direct grants must be imported before this resource can manage them.
	current, err := namespaceRoleAssignments(ctx, r.client, plan.NamespaceID.ValueString(), plan.RoleID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Cannot read namespace grants", err.Error())
		return
	}
	if len(current) > 0 {
		resp.Diagnostics.AddError("Namespace grants already exist", "Import the namespace-role assignment before managing existing permissions.")
		return
	}
	plan.ID = types.StringValue(plan.ProjectID.ValueString() + "/" + plan.NamespaceID.ValueString() + "/" + plan.RoleID.ValueString())
	// Retain identity if a remote write succeeds but its response is lost.
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
	var desired []string
	resp.Diagnostics.Append(plan.Assignments.ElementsAs(ctx, &desired, false)...)
	if resp.Diagnostics.HasError() {
		return
	}
	if err := setNamespaceRoleAssignments(ctx, r.client, plan.NamespaceID.ValueString(), plan.RoleID.ValueString(), desired); err != nil {
		resp.Diagnostics.AddError("Cannot create namespace grants", err.Error())
	}
}
func (r *namespaceRoleAssignmentResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state namespaceRoleAssignmentModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	grants, err := namespaceRoleAssignments(ctx, r.client, state.NamespaceID.ValueString(), state.RoleID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Cannot read namespace grants", err.Error())
		return
	}
	values, diags := types.SetValueFrom(ctx, types.StringType, grants)
	resp.Diagnostics.Append(diags...)
	state.Assignments = values
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
func (r *namespaceRoleAssignmentResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan namespaceRoleAssignmentModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	var desired []string
	resp.Diagnostics.Append(plan.Assignments.ElementsAs(ctx, &desired, false)...)
	if resp.Diagnostics.HasError() {
		return
	}
	if err := setNamespaceRoleAssignments(ctx, r.client, plan.NamespaceID.ValueString(), plan.RoleID.ValueString(), desired); err != nil {
		resp.Diagnostics.AddError("Cannot update namespace grants", err.Error())
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}
func (r *namespaceRoleAssignmentResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state namespaceRoleAssignmentModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	if err := setNamespaceRoleAssignments(ctx, r.client, state.NamespaceID.ValueString(), state.RoleID.ValueString(), nil); err != nil {
		resp.Diagnostics.AddError("Cannot revoke namespace grants", err.Error())
		return
	}
	resp.State.RemoveResource(ctx)
}
func (r *namespaceRoleAssignmentResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	parts := strings.Split(req.ID, "/")
	if len(parts) != 3 {
		resp.Diagnostics.AddError("Invalid import ID", "Use project_id/namespace_id/role_id.")
		return
	}
	for i, name := range []string{"project_id", "namespace_id", "role_id"} {
		resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root(name), parts[i])...)
	}
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
