package provider

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	managementv1 "github.com/baptistegh/go-lakekeeper/pkg/apis/management/v1"
	permissionv1 "github.com/baptistegh/go-lakekeeper/pkg/apis/management/v1/permission"
	"github.com/baptistegh/go-lakekeeper/pkg/apis/management/v1/storage/profile"
	lakekeeper "github.com/baptistegh/go-lakekeeper/pkg/client"
	"github.com/baptistegh/go-lakekeeper/pkg/core"

	"github.com/baptistegh/terraform-provider-lakekeeper/internal/provider/sdk"

	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
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

type storageProfileWrapper struct {
	S3StorageProfile   *sdk.S3StorageProfileModel   `tfsdk:"s3"`
	ADLSStorageProfile *sdk.ADLSStorageProfileModel `tfsdk:"adls"`
	GCSStorageProfile  *sdk.GCSStorageProfileModel  `tfsdk:"gcs"`
}

// lakekeeperWarehouseResourceModel describes the resource data model.
type lakekeeperWarehouseResourceModel struct {
	ID             types.String            `tfsdk:"id"` // form: project_id:warehouse_id (internal ID)
	WarehouseID    types.String            `tfsdk:"warehouse_id"`
	Name           types.String            `tfsdk:"name"`
	ProjectID      types.String            `tfsdk:"project_id"`
	Protected      types.Bool              `tfsdk:"protected"`
	Active         types.Bool              `tfsdk:"active"`
	ManagedAccess  types.Bool              `tfsdk:"managed_access"`
	DeleteProfile  *sdk.DeleteProfileModel `tfsdk:"delete_profile"`
	StorageProfile *storageProfileWrapper  `tfsdk:"storage_profile"`
}

func (r *lakekeeperWarehouseResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: fmt.Sprintf(`The ` + "`lakekeeper_warehouse`" + ` resource allows to manage the lifecycle of a lakekeeper warehouse.

**Upstream API**: [Lakekeeper REST API docs](https://docs.lakekeeper.io/docs/nightly/api/management/#tag/warehouse)`),
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "The internal ID of this resource. In the form: `{{project_id}}/{{warehouse_id}}`",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"warehouse_id": schema.StringAttribute{
				MarkdownDescription: "The ID of the warehouse.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				MarkdownDescription: `Name of the warehouse to create. Must be unique within a project and may not contain "/"`,
				Required:            true,
				Validators:          []validator.String{stringvalidator.LengthAtLeast(1)},
				PlanModifiers:       []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"project_id": schema.StringAttribute{
				MarkdownDescription: "The project ID to which the warehouse belongs.",
				Required:            true,
				Validators:          []validator.String{stringvalidator.LengthAtLeast(1)},
				PlanModifiers:       []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"protected": schema.BoolAttribute{
				MarkdownDescription: "Whether the warehouse is protected from being deleted. Default is `false`.",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"active": schema.BoolAttribute{
				MarkdownDescription: "Whether the warehouse is active. Default is `true`.",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(true),
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"managed_access": schema.BoolAttribute{
				MarkdownDescription: "Whether the managed access is configured on this warehouse. Default is `false`.",
				Computed:            true,
				Optional:            true,
				Default:             booldefault.StaticBool(false),
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"delete_profile": sdk.DeleteProfileResourceSchema(),
			"storage_profile": schema.SingleNestedAttribute{
				Required:            true,
				MarkdownDescription: "Configure the storage profile. Only one Of `s3`, `adls` or `gcs` must be provided.",
				Attributes: map[string]schema.Attribute{
					"s3": schema.SingleNestedAttribute{
						Optional:            true,
						MarkdownDescription: `S3 storage profile. Suitable for AWS or any service compatible with the S3 API.`,
						Attributes: map[string]schema.Attribute{
							"region": schema.StringAttribute{
								Required:            true,
								MarkdownDescription: "Region to use for S3 requests.",
							},
							"bucket": schema.StringAttribute{
								Required:            true,
								MarkdownDescription: "The bucket name for the storage profile.",
								Validators: []validator.String{
									stringvalidator.LengthBetween(3, 64),
								},
							},
							"sts_enabled": schema.BoolAttribute{
								Required:            true,
								MarkdownDescription: "Whether to enable STS for S3 storage profile. Required if the storage type is `s3`. If enabled, the `sts_role_arn` or `assume_role_arn` must be provided.",
							},
							"key_prefix": schema.StringAttribute{
								Optional:            true,
								MarkdownDescription: "Subpath in the filesystem to use.",
							},
							"allow_alternative_protocols": schema.BoolAttribute{
								Optional:            true,
								Computed:            true,
								MarkdownDescription: "Allow `s3a://`, `s3n://`, `wasbs://` in locations. This is disabled by default. We do not recommend to use this setting except for migration.",
								PlanModifiers: []planmodifier.Bool{
									boolplanmodifier.UseStateForUnknown(),
								},
							},
							"assume_role_arn": schema.StringAttribute{
								Optional:            true,
								MarkdownDescription: "Optional ARN to assume when accessing the bucket from Lakekeeper for S3 storage profile",
							},
							"aws_kms_key_arn": schema.StringAttribute{
								Optional:            true,
								MarkdownDescription: "ARN of the KMS key used to encrypt the S3 bucket, if any.",
							},
							"endpoint": schema.StringAttribute{
								Optional:            true,
								Computed:            true,
								MarkdownDescription: "Optional endpoint to use for S3 requests, if not provided the region will be used to determine the endpoint. If both region and endpoint are provided, the endpoint will be used. Example: `http://s3-de.my-domain.com:9000`",
								Validators: []validator.String{
									stringvalidator.RegexMatches(regexp.MustCompile("/$"), "Endpoint must ends with '/' character"),
								},
								PlanModifiers: []planmodifier.String{
									stringplanmodifier.UseStateForUnknown(),
								},
							},
							"flavor": schema.StringAttribute{
								Optional: true,
								Computed: true,
								Validators: []validator.String{
									stringvalidator.OneOf(string(profile.S3CompatFlavor), string(profile.AWSFlavor)),
								},
								MarkdownDescription: fmt.Sprintf("S3 flavor to use. Defaults to `%s`. One of `%s` `%s`", profile.AWSFlavor, profile.S3CompatFlavor, profile.AWSFlavor),
								PlanModifiers: []planmodifier.String{
									stringplanmodifier.UseStateForUnknown(),
								},
							},
							"path_style_access": schema.BoolAttribute{
								Optional:            true,
								Computed:            true,
								MarkdownDescription: "Path style access for S3 requests. If the underlying S3 supports both, we recommend to not set path_style_access.",
								PlanModifiers: []planmodifier.Bool{
									boolplanmodifier.UseStateForUnknown(),
								},
							},
							"push_s3_delete_disabled": schema.BoolAttribute{
								Optional:            true,
								Computed:            true,
								MarkdownDescription: "Controls whether the `s3.delete-enabled=false` flag is sent to clients.",
								PlanModifiers: []planmodifier.Bool{
									boolplanmodifier.UseStateForUnknown(),
								},
							},
							"remote_signing_url_style": schema.StringAttribute{
								Optional: true,
								Computed: true,
								Validators: []validator.String{
									stringvalidator.OneOf(
										string(profile.AutoSigningURLStyle),
										string(profile.PathSigningURLStyle),
										string(profile.VirtualHostSigningURLStyle),
									),
								},
								MarkdownDescription: fmt.Sprintf("S3 URL style detection mode for remote signing. One of `%s`, `%s`, `%s`. Default: `%s`.", profile.AutoSigningURLStyle, profile.PathSigningURLStyle, profile.VirtualHostSigningURLStyle, profile.AutoSigningURLStyle),
								PlanModifiers: []planmodifier.String{
									stringplanmodifier.UseStateForUnknown(),
								},
							}, "sts_role_arn": schema.StringAttribute{
								Optional: true,
							},
							"sts_token_validity_seconds": schema.Int64Attribute{
								Optional: true,
								Computed: true,
								Validators: []validator.Int64{
									int64validator.AtLeast(0),
								},
								MarkdownDescription: "The validity of the STS tokens in seconds. Default is `3600`.",
								PlanModifiers: []planmodifier.Int64{
									int64planmodifier.UseStateForUnknown(),
								},
							},
							"credential": schema.SingleNestedAttribute{
								Required:            true,
								MarkdownDescription: "Configure the credentials to access the S3 storage. Only one of `access_key`, `cloudflare_r2` or `aws_system_identity` must be provided. This is not available for imported resources.",
								Attributes: map[string]schema.Attribute{
									"access_key": schema.SingleNestedAttribute{
										Optional:            true,
										MarkdownDescription: "Authenticate to the S3 bucket with an Access Key",
										Attributes: map[string]schema.Attribute{
											"access_key_id": schema.StringAttribute{
												Required:            true,
												Sensitive:           true,
												MarkdownDescription: "The access key ID. Required for `aws-access-key` credentials.",
											},
											"secret_access_key": schema.StringAttribute{
												Required:            true,
												Sensitive:           true,
												MarkdownDescription: "The secret access key. Required for `aws-access-key` credentials.",
											},
											"external_id": schema.StringAttribute{
												Optional:            true,
												Sensitive:           true,
												MarkdownDescription: "The external ID.",
											},
										},
									},
									"cloudflare_r2": schema.SingleNestedAttribute{
										Optional:            true,
										MarkdownDescription: "Authenticate to a Cloudflare R2 Bucket",
										Attributes: map[string]schema.Attribute{
											"access_key_id": schema.StringAttribute{
												Required:            true,
												Sensitive:           true,
												MarkdownDescription: "Access key ID used for IO operations of Lakekeeper",
											},
											"secret_access_key": schema.StringAttribute{
												Required:            true,
												Sensitive:           true,
												MarkdownDescription: "Secret key associated with the access key ID",
											},
											"account_id": schema.StringAttribute{
												Required:            true,
												Sensitive:           true,
												MarkdownDescription: "Cloudflare account ID, used to determine the temporary credentials",
											},
											"token": schema.StringAttribute{
												Required:            true,
												Sensitive:           true,
												MarkdownDescription: "Token associated with the access key ID",
											},
										},
									},
									"aws_system_identity": schema.SingleNestedAttribute{
										Optional:            true,
										MarkdownDescription: "Authenticate to the S3 bucket with AWS System Identity",
										Attributes: map[string]schema.Attribute{
											"external_id": schema.StringAttribute{
												Required:            true,
												Sensitive:           true,
												MarkdownDescription: "Required for `aws-system-identity` credentials",
											},
										},
									},
								},
							},
						},
					},
					"adls": schema.SingleNestedAttribute{
						Optional:            true,
						MarkdownDescription: "ADLS storage profile. Suitable for Azure Data Lake Storage Gen2.",
						Attributes: map[string]schema.Attribute{
							"account_name": schema.StringAttribute{
								Required:            true,
								MarkdownDescription: "Name of the azure storage account.",
							},
							"filesystem": schema.StringAttribute{
								Required:            true,
								MarkdownDescription: "Name of the adls filesystem, in blobstorage also known as container.",
								Validators: []validator.String{
									stringvalidator.LengthBetween(3, 64),
								},
							},
							"allow_alternative_protocols": schema.BoolAttribute{
								Optional:            true,
								Computed:            true,
								MarkdownDescription: "Allow alternative protocols such as wasbs:// in locations. This is disabled by default. We do not recommend to use this setting except for migration.",
								PlanModifiers: []planmodifier.Bool{
									boolplanmodifier.UseStateForUnknown(),
								},
							},
							"authority_host": schema.StringAttribute{
								Optional:            true,
								Computed:            true,
								MarkdownDescription: "The authority host to use for authentication. Defaults to `https://login.microsoftonline.com`.",
								PlanModifiers: []planmodifier.String{
									stringplanmodifier.UseStateForUnknown(),
								},
							},
							"host": schema.StringAttribute{
								Optional:            true,
								Computed:            true,
								MarkdownDescription: "The host to use for the storage account. Defaults to `dfs.core.windows.net`.",
								PlanModifiers: []planmodifier.String{
									stringplanmodifier.UseStateForUnknown(),
								},
							},
							"key_prefix": schema.StringAttribute{
								Optional:            true,
								MarkdownDescription: "Subpath in the filesystem to use.",
							},
							"sas_token_validity_seconds": schema.Int64Attribute{
								Optional:            true,
								Computed:            true,
								MarkdownDescription: "The validity of the sas token in seconds. Default is `3600`.",
								Validators: []validator.Int64{
									int64validator.AtLeast(0),
								},
								PlanModifiers: []planmodifier.Int64{
									int64planmodifier.UseStateForUnknown(),
								},
							},
							"credential": schema.SingleNestedAttribute{
								Required:            true,
								MarkdownDescription: "Configure the credentials to access the ADLS storage. One of `shared_access_key`, `client_credentials` or `azure_system_identity` must be provided. This is not available for imported resources.",
								Attributes: map[string]schema.Attribute{
									"shared_access_key": schema.SingleNestedAttribute{
										Optional:            true,
										MarkdownDescription: "Authenticate to ADLS with Shared Access Key",
										Attributes: map[string]schema.Attribute{
											"key": schema.StringAttribute{
												Required:            true,
												Sensitive:           true,
												MarkdownDescription: "The shared access key. Required for `azure-shared-access-key` credentials.",
											},
										},
									},
									"client_credentials": schema.SingleNestedAttribute{
										Optional:            true,
										MarkdownDescription: "Authenticate to ADLS with Client Credentials",
										Attributes: map[string]schema.Attribute{
											"client_id": schema.StringAttribute{
												Required:  true,
												Sensitive: true,
											},
											"client_secret": schema.StringAttribute{
												Required:  true,
												Sensitive: true,
											},
											"tenant_id": schema.StringAttribute{
												Required:  true,
												Sensitive: true,
											},
										},
									},
									"azure_system_identity": schema.SingleNestedAttribute{
										Optional:            true,
										MarkdownDescription: "Authenticate to ADLS with Azure System Identity",
										Attributes: map[string]schema.Attribute{
											"enabled": schema.BoolAttribute{
												Computed:            true,
												MarkdownDescription: "This is just an helper to check if the Azure System Identity is activated for this storage profile.",
												Default:             booldefault.StaticBool(true),
											},
										},
									},
								},
							},
						},
					},
					"gcs": schema.SingleNestedAttribute{
						Optional:            true,
						MarkdownDescription: "GCS storage profile. Designed for use with Google Cloud Storage",
						Attributes: map[string]schema.Attribute{
							"bucket": schema.StringAttribute{
								Required:            true,
								MarkdownDescription: "The bucket name.",
								Validators: []validator.String{
									stringvalidator.LengthBetween(3, 64),
								},
							},
							"key_prefix": schema.StringAttribute{
								Optional:            true,
								MarkdownDescription: "Subpath in the filesystem to use.",
							},
							"credential": schema.SingleNestedAttribute{
								Required:            true,
								MarkdownDescription: "Configure the credentials to access the GCS storage. One of `service_account_key` or `gcp_system_identity` must be provided. This is not available for imported resources.",
								Attributes: map[string]schema.Attribute{
									"service_account_key": schema.SingleNestedAttribute{
										Optional:            true,
										MarkdownDescription: "Authenticate to the GCS bucket with a Service Account Key",
										Attributes: map[string]schema.Attribute{
											"key": schema.StringAttribute{
												Required:            true,
												Sensitive:           true,
												MarkdownDescription: "Required for `service-account-key` credentials.",
											},
										},
									},
									"gcp_system_identity": schema.SingleNestedAttribute{
										Optional:            true,
										MarkdownDescription: "Authenticate to the GCS bucket with GCP System Identity",
										Attributes: map[string]schema.Attribute{
											"enabled": schema.BoolAttribute{
												Computed:            true,
												MarkdownDescription: "This is just an helper to check if the GCP System Identity is activated for this storage profile.",
												Default:             booldefault.StaticBool(true),
											},
										},
									},
								},
							},
						},
					},
				},
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
	var state lakekeeperWarehouseResourceModel
	var plan types.Object

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(plan.As(ctx, &state, basetypes.ObjectAsOptions{
		UnhandledNullAsEmpty:    false,
		UnhandledUnknownAsEmpty: false,
	})...)
	if resp.Diagnostics.HasError() {
		return
	}

	opts, err := state.toWarehouseCreateRequest()
	if err != nil {
		resp.Diagnostics.AddError("Error decoding state to model", fmt.Sprintf("Incorrect Warehouse creation request, %v", err))
		return
	}

	w, _, err := r.client.WarehouseV1(state.ProjectID.ValueString()).Create(ctx, opts)
	if err != nil {
		resp.Diagnostics.AddError("Lakekeeper API error occurred",
			fmt.Sprintf("Unable to create warehouse, %v", err))
		return
	}

	// set protection
	if state.Protected.ValueBool() {
		_, _, err := r.client.WarehouseV1(state.ProjectID.ValueString()).SetWarehouseProtection(ctx, w.ID, &managementv1.SetProtectionOptions{Protected: state.Protected.ValueBool()})
		if err != nil {
			resp.Diagnostics.AddError("Lakekeeper API error occurred.",
				fmt.Sprintf("Unable to set protection to %t for warehouse %s, %v", state.Protected.ValueBool(), w.ID, err),
			)
		}
	}

	// set inactive
	if !state.Active.ValueBool() {
		_, err := r.client.WarehouseV1(state.ProjectID.ValueString()).Deactivate(ctx, w.ID)
		if err != nil {
			resp.Diagnostics.AddError("Lakekeeper API error occurred.",
				fmt.Sprintf("Unable to deactivate warehouse %s, %v", w.ID, err),
			)
		}
	}

	warehouse, _, err := r.client.WarehouseV1(state.ProjectID.ValueString()).Get(ctx, w.ID)
	if err != nil {
		resp.Diagnostics.AddError("Lakekeeper API error occurred",
			fmt.Sprintf("Unable to read warehouse, %v", err))
		return
	}

	if !state.ManagedAccess.IsNull() {
		_, err := r.client.PermissionV1().WarehousePermission().SetManagedAccess(ctx, warehouse.ID, &permissionv1.SetWarehouseManagedAccessOptions{
			ManagedAccess: state.ManagedAccess.ValueBool(),
		})
		if err != nil {
			resp.Diagnostics.AddError("Lakekeeper API error occurred",
				fmt.Sprintf("Unable to set managed access, %v", err))
			return
		}
	}

	diags := state.RefreshFromSettings(warehouse, nil)
	if diags.HasError() {
		resp.Diagnostics = append(resp.Diagnostics, diags...)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Read refreshes the Terraform state with the latest data.
func (r *lakekeeperWarehouseResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state lakekeeperWarehouseResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	projectID, warehouseID := splitInternalID(state.ID)
	warehouse, _, err := r.client.WarehouseV1(projectID).Get(ctx, warehouseID)
	if err != nil {
		resp.Diagnostics.AddError("Lakekeeper API error occurred", fmt.Sprintf("Unable to read warehouse %s in project %s, %s", warehouseID, projectID, err))
		return
	}

	diags := state.RefreshFromSettings(warehouse, nil)
	if diags.HasError() {
		resp.Diagnostics = append(resp.Diagnostics, diags...)
		return
	}

	// get managed access property
	m, _, err := r.client.PermissionV1().WarehousePermission().GetAuthzProperties(ctx, warehouse.ID)
	if err != nil {
		resp.Diagnostics.AddError("Lakekeeper API error occurred", fmt.Sprintf("Unable to read warehouse %s authorization properties in project %s, %s", warehouseID, projectID, err))
	}
	state.ManagedAccess = types.BoolValue(m.ManagedAccess)

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Updates updates the resource in-place.
func (r *lakekeeperWarehouseResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state lakekeeperWarehouseResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	projectID, warehouseID := splitInternalID(state.ID)

	if plan.Active.ValueBool() {
		if _, err := r.client.WarehouseV1(projectID).Activate(ctx, warehouseID); err != nil {
			resp.Diagnostics.AddError("Lakekeeper API error occurred", fmt.Sprintf("Unable to activate warehouse %s in project %s, %v", warehouseID, projectID, err))
			return
		}
	} else {
		if _, err := r.client.WarehouseV1(projectID).Deactivate(ctx, warehouseID); err != nil {
			resp.Diagnostics.AddError("Lakekeeper API error occurred", fmt.Sprintf("Unable to deactivate warehouse %s in project %s, %v", warehouseID, projectID, err))
			return
		}
	}

	// Rename the warehouse if the name field is different
	if plan.Name.ValueString() != state.Name.ValueString() {
		if _, err := r.client.WarehouseV1(projectID).Rename(ctx, warehouseID, &managementv1.RenameWarehouseOptions{
			NewName: plan.Name.ValueString(),
		}); err != nil {
			resp.Diagnostics.AddError("Lakekeeper API error occurred", fmt.Sprintf("Unable to rename warehouse %s in project %s, %v", warehouseID, projectID, err))
			return
		}
	}

	// Set warehouse protection if the protected field is different
	if plan.Protected.ValueBool() != state.Protected.ValueBool() {
		if _, _, err := r.client.WarehouseV1(projectID).SetWarehouseProtection(ctx, warehouseID, &managementv1.SetProtectionOptions{Protected: plan.Protected.ValueBool()}); err != nil {
			resp.Diagnostics.AddError("Lakekeeper API error occurred", fmt.Sprintf("Unable to set protection for warehouse %s in project %s, %v", warehouseID, projectID, err))
			return
		}
		state.Protected = plan.Protected
	}

	opts, err := plan.toWarehouseCreateRequest()
	if err != nil {
		resp.Diagnostics.AddError("Error decoding plan to model", fmt.Sprintf("Incorrect Warehouse update request, %v", err))
		return
	}

	// Update the delete profile
	if !plan.DeleteProfile.Type.Equal(state.DeleteProfile.Type) || !plan.DeleteProfile.ExpirationSeconds.Equal(state.DeleteProfile.ExpirationSeconds) {
		if _, err := r.client.WarehouseV1(projectID).UpdateDeleteProfile(ctx, warehouseID, &managementv1.UpdateDeleteProfileOptions{
			DeleteProfile: *opts.DeleteProfile,
		}); err != nil {
			resp.Diagnostics.AddError("Lakekeeper API error occurred", fmt.Sprintf("Unable to update delete profile for warehouse %s in project %s, %v", warehouseID, projectID, err))
			return
		}
	}

	// Update the storage profile and its storage credential
	if _, err := r.client.WarehouseV1(projectID).UpdateStorageProfile(ctx, warehouseID, &managementv1.UpdateStorageProfileOptions{
		StorageProfile:    opts.StorageProfile,
		StorageCredential: &opts.StorageCredential,
	}); err != nil {
		resp.Diagnostics.AddError("Lakekeeper API error occurred", fmt.Sprintf("Unable to update storage profile for warehouse %s in project %s, %v", warehouseID, projectID, err))
		return
	}

	// Update the authorization property
	if _, err := r.client.PermissionV1().WarehousePermission().SetManagedAccess(ctx, warehouseID, &permissionv1.SetWarehouseManagedAccessOptions{
		ManagedAccess: plan.ManagedAccess.ValueBool(),
	}); err != nil {
		resp.Diagnostics.AddError("Lakekeeper API error occurred", fmt.Sprintf("Unable to set authorization properties for warehouse %s in project %s, %v", warehouseID, projectID, err))
		return
	}
	state.ManagedAccess = plan.ManagedAccess

	// Refresh the state with the updated warehouse settings
	warehouse, _, err := r.client.WarehouseV1(projectID).Get(ctx, warehouseID)
	if err != nil {
		resp.Diagnostics.AddError("Lakekeeper API error occurred", fmt.Sprintf("Unable to read warehouse %s in project %s, %v", warehouseID, projectID, err))
		return
	}

	diags := state.RefreshFromSettings(warehouse, &plan)
	if diags.HasError() {
		resp.Diagnostics = append(resp.Diagnostics, diags...)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Deletes removes the resource.
func (r *lakekeeperWarehouseResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state lakekeeperWarehouseResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	projectID, warehouseID := splitInternalID(state.ID)

	opts := managementv1.DeleteWarehouseOptions{
		Force: core.Ptr(true),
	}

	if _, err := r.client.WarehouseV1(projectID).Delete(ctx, warehouseID, &opts); err != nil {
		resp.Diagnostics.AddError("Lakekeeper API error occurred", fmt.Sprintf("Unable to delete warehouse %s in project %s, %v", warehouseID, projectID, err))
		return
	}

	resp.State.RemoveResource(ctx)
}

// ImportState imports the resource into the Terraform state.
func (r *lakekeeperWarehouseResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Expected format: "project_id/warehouse_id"
	parts := strings.Split(req.ID, "/")
	if len(parts) != 2 {
		resp.Diagnostics.AddError(
			"Invalid import ID format",
			"Expected format: project_id/warehouse_id",
		)
		return
	}

	resp.State.SetAttribute(ctx, path.Root("project_id"), parts[0])
	resp.State.SetAttribute(ctx, path.Root("warehouse_id"), parts[1])

	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
