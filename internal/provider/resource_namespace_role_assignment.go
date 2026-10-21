package provider

import (
	"context"
	"fmt"
	"regexp"
	"sort"
	"strings"

	lakekeeper "github.com/baptistegh/go-lakekeeper/pkg/client"
	"github.com/baptistegh/go-lakekeeper/pkg/core"
	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
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
		MarkdownDescription: "Manages all direct namespace permissions for one role. Other roles and users are preserved. Requires Lakekeeper's OpenFGA authorization backend. Import existing grants before managing them.",
		Attributes: map[string]schema.Attribute{
			"id":           schema.StringAttribute{Computed: true, MarkdownDescription: "project_id/namespace_id/role_id", PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()}},
			"project_id":   identity("Project containing both the namespace and role."),
			"namespace_id": identity("Namespace UUID, available as lakekeeper_namespace.example.namespace_id."),
			"role_id":      identity("Role UUID."),
			"assignments": schema.SetAttribute{Required: true, ElementType: types.StringType,
				MarkdownDescription: "Direct permissions for this role. An empty set revokes them. Removing this resource revokes its permissions.",
				Validators:          []validator.Set{setvalidator.ValueStringsAre(stringvalidator.OneOf("ownership", "pass_grants", "manage_grants", "describe", "select", "create", "modify"))},
			},
		},
	}
}
func (r *namespaceRoleAssignmentResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData != nil {
		r.client = req.ProviderData.(*LakekeeperResourceData).Client
	}
}

type namespaceAssignment struct {
	Type string `json:"type"`
	Role string `json:"role,omitempty"`
	User string `json:"user,omitempty"`
}

func namespaceRequest(ctx context.Context, client *lakekeeper.Client, project, namespace, suffix, method string, body, result any) error {
	req, err := client.NewRequest(ctx, method, "/permissions/namespace/"+namespace+suffix, body, []core.RequestOptionFunc{core.WithHeader("x-project-id", project)})
	if err != nil {
		return err
	}
	_, apiErr := client.Do(req, result)
	if apiErr != nil {
		return apiErr
	}
	return nil
}
func namespaceRoleAssignments(ctx context.Context, client *lakekeeper.Client, project, namespace, role string) ([]string, error) {
	var result struct {
		Assignments *[]namespaceAssignment `json:"assignments"`
	}
	if err := namespaceRequest(ctx, client, project, namespace, "/assignments", "GET", nil, &result); err != nil {
		return nil, err
	}
	if result.Assignments == nil {
		return nil, fmt.Errorf("namespace assignments response has no assignments list")
	}
	found := map[string]bool{}
	for _, a := range *result.Assignments {
		if a.Role == role && a.User == "" {
			found[a.Type] = true
		}
	}
	values := make([]string, 0, len(found))
	for a := range found {
		values = append(values, a)
	}
	sort.Strings(values)
	return values, nil
}

// Reconcile only this role's direct grants. Delete before adding; refresh recovers a partial apply.
func setNamespaceRoleAssignments(ctx context.Context, client *lakekeeper.Client, project, namespace, role string, desired []string) error {
	current, err := namespaceRoleAssignments(ctx, client, project, namespace, role)
	if err != nil {
		return err
	}
	before, after := map[string]bool{}, map[string]bool{}
	for _, a := range current {
		before[a] = true
	}
	for _, a := range desired {
		after[a] = true
	}
	for _, op := range []struct {
		key      string
		from, to map[string]bool
	}{{"deletes", before, after}, {"writes", after, before}} {
		changes := []namespaceAssignment{}
		for a := range op.from {
			if !op.to[a] {
				changes = append(changes, namespaceAssignment{Type: a, Role: role})
			}
		}
		sort.Slice(changes, func(i, j int) bool { return changes[i].Type < changes[j].Type })
		if len(changes) > 0 {
			if err := namespaceRequest(ctx, client, project, namespace, "/assignments", "POST", map[string]any{op.key: changes}, nil); err != nil {
				return err
			}
		}
	}
	return nil
}
func (r *namespaceRoleAssignmentResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan namespaceRoleAssignmentModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	// Create adopts existing direct grants only when they already match; import is required otherwise.
	current, err := namespaceRoleAssignments(ctx, r.client, plan.ProjectID.ValueString(), plan.NamespaceID.ValueString(), plan.RoleID.ValueString())
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
	if err := setNamespaceRoleAssignments(ctx, r.client, plan.ProjectID.ValueString(), plan.NamespaceID.ValueString(), plan.RoleID.ValueString(), desired); err != nil {
		resp.Diagnostics.AddError("Cannot create namespace grants", err.Error())
	}
}
func (r *namespaceRoleAssignmentResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state namespaceRoleAssignmentModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	grants, err := namespaceRoleAssignments(ctx, r.client, state.ProjectID.ValueString(), state.NamespaceID.ValueString(), state.RoleID.ValueString())
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
	if err := setNamespaceRoleAssignments(ctx, r.client, plan.ProjectID.ValueString(), plan.NamespaceID.ValueString(), plan.RoleID.ValueString(), desired); err != nil {
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
	if err := setNamespaceRoleAssignments(ctx, r.client, state.ProjectID.ValueString(), state.NamespaceID.ValueString(), state.RoleID.ValueString(), nil); err != nil {
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
