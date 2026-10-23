package provider

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"strings"

	lakekeeper "github.com/baptistegh/go-lakekeeper/pkg/client"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/apache/iceberg-go/catalog"
)

// Ensure provider defined types fully satisfy framework interfaces
var (
	_ resource.Resource                = &lakekeeperNamespaceResource{}
	_ resource.ResourceWithConfigure   = &lakekeeperNamespaceResource{}
	_ resource.ResourceWithImportState = &lakekeeperNamespaceResource{}
)

func init() {
	registerResource(NewLakekeeperNamespaceResource)
}

// NewLakekeeperNamespaceResource is a helper function to simplify the provider implementation.
func NewLakekeeperNamespaceResource() resource.Resource {
	return &lakekeeperNamespaceResource{}
}

func (r *lakekeeperNamespaceResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_namespace"
}

// lakekeeperNamespaceResource defines the resource implementation.
type lakekeeperNamespaceResource struct {
	client *lakekeeper.Client
}

// lakekeeperNamespaceResourceModel describes the resource data model.
type lakekeeperNamespaceResourceModel struct {
	ID            types.String `tfsdk:"id"`
	NamespaceID   types.String `tfsdk:"namespace_id"`
	ManagedAccess types.Bool   `tfsdk:"managed_access"`
	ProjectID     types.String `tfsdk:"project_id"`
	WarehouseName types.String `tfsdk:"warehouse_name"`
	Name          types.String `tfsdk:"name"`
}

func (r *lakekeeperNamespaceResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: fmt.Sprintf(`The ` + "`lakekeeper_namespace`" + ` resource allows to manage the lifecycle of an Iceberg Namespace inside Lakekeeper.

**Upstream API**: [Lakekeeper REST API docs](https://docs.lakekeeper.io/docs/nightly/api/catalog/#tag/Catalog-API/operation/createNamespace)`),

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "The ID the namespace. In the form `{{project_id}}/{{warehouse_name}}/{{name}}`",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"namespace_id":   schema.StringAttribute{Computed: true, MarkdownDescription: "Namespace UUID used by the permissions API."},
			"managed_access": schema.BoolAttribute{Optional: true, MarkdownDescription: "Centralize grants on descendants. Requires OpenFGA when set. Leaving this unset avoids the managed-access endpoint. Removing the setting preserves its remote value."},
			"project_id": schema.StringAttribute{
				MarkdownDescription: "The internal ID of the project where the namespace is located.",
				Required:            true,
			},
			"warehouse_name": schema.StringAttribute{
				MarkdownDescription: "The name of the warehouse where the namespace is located.",
				Required:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "The name of the namespace.",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
		},
	}
}

// Configure adds the provider configured client to the resource.
func (r *lakekeeperNamespaceResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	resourceData := req.ProviderData.(*LakekeeperResourceData)
	r.client = resourceData.Client
}

// Create creates a new upstream resources and adds it into the Terraform state.
func (r *lakekeeperNamespaceResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var state lakekeeperNamespaceResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	project_id := state.ProjectID.ValueString()
	warehouse_name := state.WarehouseName.ValueString()

	cat, err := r.client.CatalogV1(ctx, project_id, warehouse_name)
	if err != nil {
		resp.Diagnostics.AddError("Could not create the Iceberg Catalog client", fmt.Sprintf("Unable to initialize the client, %v", err.Error()))
		return
	}

	name := state.Name.ValueString()
	namespace := catalog.ToIdentifier(name)

	if err := cat.CreateNamespace(ctx, namespace, nil); err != nil {
		resp.Diagnostics.AddError("REST Catalog API error occurred", fmt.Sprintf("Unable to create namespace %s, %v", name, err.Error()))
		return
	}

	state.ID = types.StringValue(fmt.Sprintf("%s/%s/%s", project_id, warehouse_name, name))

	// Save catalog identity before any optional authorization call.
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	if err := r.refreshNamespace(ctx, &state); err != nil {
		resp.Diagnostics.AddError("Cannot read created namespace", err.Error())
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	if err := r.setManagedAccess(ctx, &state); err != nil {
		resp.Diagnostics.AddError("Cannot set managed access", err.Error())
		return
	}

	// Log the creation of the resource
	tflog.Debug(ctx, "created a namespace", map[string]any{
		"id": state.ID.ValueString(),
	})

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Read refreshes the Terraform state with the latest data.
func (r *lakekeeperNamespaceResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state lakekeeperNamespaceResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.refreshNamespace(ctx, &state); err != nil {
		if errors.Is(err, catalog.ErrNoSuchNamespace) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Cannot read namespace", err.Error())
		return
	}
	if !state.ManagedAccess.IsNull() {
		var props struct {
			Managed *bool `json:"managed-access"`
		}
		if err := namespaceRequest(ctx, r.client, state.ProjectID.ValueString(), state.NamespaceID.ValueString(), "", "GET", nil, &props); err != nil {
			resp.Diagnostics.AddError("Cannot read managed access", err.Error())
			return
		}
		if props.Managed == nil {
			resp.Diagnostics.AddError("Cannot read managed access", "Response has no managed-access field")
			return
		}
		state.ManagedAccess = types.BoolValue(*props.Managed)
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Updates updates the resource in-place.
func (r *lakekeeperNamespaceResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state lakekeeperNamespaceResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	plan.ID = state.ID
	plan.NamespaceID = state.NamespaceID
	if err := r.setManagedAccess(ctx, &plan); err != nil {
		resp.Diagnostics.AddError("Cannot set managed access", err.Error())
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

// Deletes removes the resource.
func (r *lakekeeperNamespaceResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state lakekeeperNamespaceResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	project_id := state.ProjectID.ValueString()
	warehouse_name := state.WarehouseName.ValueString()

	cat, err := r.client.CatalogV1(ctx, project_id, warehouse_name)
	if err != nil {
		resp.Diagnostics.AddError("Could not create the Iceberg Catalog client", fmt.Sprintf("Unable to initialize the client, %v", err.Error()))
		return
	}

	name := state.Name.ValueString()

	if err := cat.DropNamespace(ctx, catalog.ToIdentifier(name)); err != nil {
		resp.Diagnostics.AddError("REST Catalog API error occurred", fmt.Sprintf("Unable to delete namespace %s, %v", name, err.Error()))
		return
	}

	resp.State.RemoveResource(ctx)
}

// ImportState imports the resource into the Terraform state.
func (r *lakekeeperNamespaceResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Expected format: "project_id/warehouse_name/name"
	parts := strings.Split(req.ID, "/")
	if len(parts) != 3 {
		resp.Diagnostics.AddError(
			"Invalid import ID format",
			"Expected format: project_id/warehouse_name/name",
		)
		return
	}

	resp.State.SetAttribute(ctx, path.Root("project_id"), parts[0])
	resp.State.SetAttribute(ctx, path.Root("warehouse_name"), parts[1])
	resp.State.SetAttribute(ctx, path.Root("name"), parts[2])

	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r *lakekeeperNamespaceResource) refreshNamespace(ctx context.Context, state *lakekeeperNamespaceResourceModel) error {
	cat, err := r.client.CatalogV1(ctx, state.ProjectID.ValueString(), state.WarehouseName.ValueString())
	if err != nil {
		return err
	}
	props, err := cat.LoadNamespaceProperties(ctx, catalog.ToIdentifier(state.Name.ValueString()))
	if err != nil {
		return err
	}
	id := props["namespace_id"]
	if _, err := uuid.Parse(id); err != nil {
		return fmt.Errorf("catalog returned an invalid namespace_id")
	}
	state.NamespaceID = types.StringValue(id)
	return nil
}
func (r *lakekeeperNamespaceResource) setManagedAccess(ctx context.Context, state *lakekeeperNamespaceResourceModel) error {
	if state.ManagedAccess.IsNull() {
		return nil
	}
	return namespaceRequest(ctx, r.client, state.ProjectID.ValueString(), state.NamespaceID.ValueString(), "/managed-access", "POST", map[string]bool{"managed-access": state.ManagedAccess.ValueBool()}, nil)
}
