package provider

import (
	"context"
	"fmt"

	"github.com/baptistegh/terraform-provider-lakekeeper/lakekeeper"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
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
type lakekeeperWarehouseResourceModel struct {
	ID                types.String            `tfsdk:"id"`
	Name              types.String            `tfsdk:"name"`
	ProjectID         types.String            `tfsdk:"project_id"` // Optional, if not provided, the default project will be used.
	Protected         types.Bool              `tfsdk:"protected"`
	Active            types.Bool              `tfsdk:"active"`
	StorageProfile    *storageProfileModel    `tfsdk:"storage_profile"`
	DeleteProfile     *deleteProfileModel     `tfsdk:"delete_profile"`
	StorageCredential *storageCredentialModel `tfsdk:"storage_credential"`
}

type storageProfileModel struct {
	Type types.String `tfsdk:"type"`

	// ADLS specific fields
	AccountName             types.String `tfsdk:"account_name"`
	AuthorityHost           types.String `tfsdk:"authority_host"`
	Filesystem              types.String `tfsdk:"filesystem"`
	Host                    types.String `tfsdk:"host"`
	SASTokenValiditySeconds types.Int64  `tfsdk:"sas_token_validity_seconds"`

	// S3 specific fields
	AssumeRoleARN           types.String `tfsdk:"assume_role_arn"`
	AWSKMSKeyARN            types.String `tfsdk:"aws_kms_key_arn"`
	Region                  types.String `tfsdk:"region"`
	Endpoint                types.String `tfsdk:"endpoint"`
	Flavor                  types.String `tfsdk:"flavor"`
	PathStyleAccess         types.Bool   `tfsdk:"path_style_access"`
	PushS3DeleteDisabled    types.Bool   `tfsdk:"push_s3_delete_disabled"`
	RemoteSigningURLStyle   types.String `tfsdk:"remote_signing_url_style"`
	STSEnabled              types.Bool   `tfsdk:"sts_enabled"`
	STSRoleARN              types.String `tfsdk:"sts_role_arn"`
	STSTokenValiditySeconds types.Int64  `tfsdk:"sts_token_validity_seconds"`

	// S3 + ADLS common fields
	AllowAlternativeProtocols types.Bool `tfsdk:"allow_alternative_protocols"`

	// S3 + GCS common fields
	Bucket types.String `tfsdk:"bucket"`

	// Common fields
	KeyPrefix types.String `tfsdk:"key_prefix"`
}

type deleteProfileModel struct {
	Type              types.String `tfsdk:"type"`
	ExpirationSeconds types.Int32  `tfsdk:"expiration_seconds"`
}

var validStorageCredentialTypes = []string{"s3_access_key", "s3_aws_system_identity", "s3_cloudflare_r2",
	"az_client_credentials", "az_shared_access_key", "az_azure_system_identity",
	"gcs_service_account_key", "gcs_gcp_system_identity"}

type storageCredentialModel struct {
	Type               types.String `tfsdk:"type"`
	AWSAccessKeyID     types.String `tfsdk:"aws_access_key_id"`
	AWSSecretAccessKey types.String `tfsdk:"aws_secret_access_key"`
	ExternalID         types.String `tfsdk:"external_id"`
	AccessKeyID        types.String `tfsdk:"access_key_id"`
	SecretAccessKey    types.String `tfsdk:"secret_access_key"`
	AccountID          types.String `tfsdk:"account_id"`
	Token              types.String `tfsdk:"token"`
	AZKey              types.String `tfsdk:"az_key"`
	ClientID           types.String `tfsdk:"client_id"`
	ClientSecret       types.String `tfsdk:"client_secret"`
	TenantID           types.String `tfsdk:"tenant_id"`
	Key                *gcsKeyModel `tfsdk:"key"`
}

type gcsKeyModel struct {
	AuthProviderX509CertURL types.String `tfsdk:"auth_provider_x509_cert_url"`
	AuthURI                 types.String `tfsdk:"auth_uri"`
	ClientEmail             types.String `tfsdk:"client_email"`
	ClientID                types.String `tfsdk:"client_id"`
	PrivateKey              types.String `tfsdk:"private_key"`
	PrivateKeyID            types.String `tfsdk:"private_key_id"`
	ProjectID               types.String `tfsdk:"project_id"`
	TokenURI                types.String `tfsdk:"token_uri"`
	UniverseDomain          types.String `tfsdk:"universe_domain"`
	ClientX509CertURL       types.String `tfsdk:"client_x509_cert_url"`
	Type                    types.String `tfsdk:"type"`
}

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
			"storage_profile": schema.SingleNestedAttribute{
				MarkdownDescription: "Whether the warehouse is active.",
				Required:            true,
				Attributes: map[string]schema.Attribute{
					"type": schema.StringAttribute{
						Required: true,
						Validators: []validator.String{
							stringvalidator.OneOf("gcs", "adls", "s3"),
						},
						MarkdownDescription: "The type of the storage profile. Supported values are `gcs`, `adls`, and `s3`.",
					},
					"account_name": schema.StringAttribute{
						Optional:            true,
						MarkdownDescription: "The account name for ADLS storage profile. Required if type is `adls`.",
					},
					"authority_host": schema.StringAttribute{
						Optional:            true,
						MarkdownDescription: "The authority host for ADLS storage profile. Defaults to `https://login.microsoftonline.com`.",
					},
					"filesystem": schema.StringAttribute{
						Optional:            true,
						MarkdownDescription: "Name of the adls filesystem, in blobstorage also known as container. Required if type is `adls`.",
					},
					"host": schema.StringAttribute{
						Optional:            true,
						MarkdownDescription: "The host for ADLS storage profile. Defaults to `dfs.core.windows.net`.",
					},
					"sas_token_validity_seconds": schema.Int64Attribute{
						Optional: true,
						Validators: []validator.Int64{
							int64validator.AtLeast(0),
						},
						MarkdownDescription: "The validity of the sts tokens in seconds. Default is `3600`.",
					},
					"assume_role_arn": schema.StringAttribute{
						Optional:            true,
						MarkdownDescription: "Optional ARN to assume when accessing the bucket from Lakekeeper for S3 storage profile",
					},
					"aws_kms_key_arn": schema.StringAttribute{
						Optional:            true,
						MarkdownDescription: "ARN of the KMS key used to encrypt the S3 bucket, if any.",
					},
					"region": schema.StringAttribute{
						Optional:            true,
						MarkdownDescription: "Region to use for S3 requests. Required if type is `s3`.",
					},
					"endpoint": schema.StringAttribute{
						Optional:            true,
						MarkdownDescription: "Optional endpoint to use for S3 requests, if not provided the region will be used to determine the endpoint. If both region and endpoint are provided, the endpoint will be used. Example: `http://s3-de.my-domain.com:9000`",
					},
					"flavor": schema.StringAttribute{
						Optional: true,
						Validators: []validator.String{
							stringvalidator.OneOf("aws", "s3-compat"),
						},
						MarkdownDescription: "S3 flavor to use. Defaults to `aws`.",
					},
					"path_style_access": schema.BoolAttribute{
						Optional:            true,
						MarkdownDescription: "Path style access for S3 requests. If the underlying S3 supports both, we recommend to not set path_style_access.",
					},
					"push_s3_delete_disabled": schema.BoolAttribute{
						Optional:            true,
						MarkdownDescription: "Controls whether the `s3.delete-enabled=false` flag is sent to clients.",
					},
					"remote_signing_url_style": schema.StringAttribute{
						Optional: true,
						Validators: []validator.String{
							stringvalidator.OneOf("path-style", "virtual-host", "auto"),
						},
						MarkdownDescription: "S3 URL style detection mode for remote signing. One of `auto`, `path-style`, `virtual-host`. Default: `auto`.",
					},
					"sts_enabled": schema.BoolAttribute{
						Optional:            true,
						MarkdownDescription: "Whether to enable STS for S3 storage profile. Required if the storage type is `s3`. If enabled, the `sts_role_arn` or `assume_role_arn` must be provided.",
					},
					"sts_role_arn": schema.StringAttribute{
						Optional: true,
					},
					"sts_token_validity_seconds": schema.Int64Attribute{
						Optional: true,
						Validators: []validator.Int64{
							int64validator.AtLeast(0),
						},
						MarkdownDescription: "The validity of the STS tokens in seconds. Default is `3600`.",
					},
					"allow_alternative_protocols": schema.BoolAttribute{
						Optional:            true,
						MarkdownDescription: "Allow `s3a://`, `s3n://`, `wasbs://` in locations. This is disabled by default. We do not recommend to use this setting except for migration.",
					},
					"bucket": schema.StringAttribute{
						Optional:            true,
						MarkdownDescription: "The bucket name for the storage profile. Required if type is `gcs` or `s3`.",
						Validators: []validator.String{
							stringvalidator.LengthBetween(3, 64),
						},
					},
					"key_prefix": schema.StringAttribute{
						Optional:            true,
						MarkdownDescription: "Subpath in the filesystem to use.",
					},
				},
				Validators: []validator.Object{storageProfileValidator{}},
			},
			"delete_profile": schema.SingleNestedAttribute{
				MarkdownDescription: "The delete profile for the warehouse. It can be either a soft or hard delete profile.",
				Required:            true,
				Attributes: map[string]schema.Attribute{
					"type": schema.StringAttribute{
						Required: true,
						Validators: []validator.String{
							stringvalidator.OneOf("soft", "hard"),
						},
					},
					"expiration_seconds": schema.Int32Attribute{
						Optional: true,
					},
				},
				Validators: []validator.Object{deleteProfileValidator{}},
			},
			"storage_credential": schema.SingleNestedAttribute{
				MarkdownDescription: "The credentials used to access the storage. This is required for the warehouse to be able to access the storage profile.",
				Optional:            true,
				Attributes: map[string]schema.Attribute{
					"type": schema.StringAttribute{
						Required: true,
						Validators: []validator.String{
							stringvalidator.OneOf("s3_access_key", "s3_aws_system_identity", "s3_cloudflare_r2",
								"az_client_credentials", "az_shared_access_key", "az_azure_system_identity",
								"gcs_service_account_key", "gcs_gcp_system_identity")},
					},
					"aws_access_key_id": schema.StringAttribute{
						Optional:  true,
						Sensitive: true,
					},
					"aws_secret_access_key": schema.StringAttribute{
						Optional:  true,
						Sensitive: true,
					},
					"external_id": schema.StringAttribute{
						Optional: true,
					},
					"access_key_id": schema.StringAttribute{
						Optional:  true,
						Sensitive: true,
					},
					"secret_access_key": schema.StringAttribute{
						Optional:  true,
						Sensitive: true,
					},
					"account_id": schema.StringAttribute{
						Optional: true,
					},
					"token": schema.StringAttribute{
						Optional:  true,
						Sensitive: true,
					},
					"az_key": schema.StringAttribute{
						Optional:  true,
						Sensitive: true,
					},
					"client_id": schema.StringAttribute{
						Optional: true,
					},
					"client_secret": schema.StringAttribute{
						Optional: true,
					},
					"tenant_id": schema.StringAttribute{
						Optional: true,
					},
					"key": schema.SingleNestedAttribute{
						Optional: true,
						Attributes: map[string]schema.Attribute{
							"auth_provider_x509_cert_url": schema.StringAttribute{
								Required: true,
							},
							"auth_uri": schema.StringAttribute{
								Required: true,
							},
							"client_email": schema.StringAttribute{
								Required: true,
							},
							"client_id": schema.StringAttribute{
								Required: true,
							},
							"private_key": schema.StringAttribute{
								Required:  true,
								Sensitive: true,
							},
							"private_key_id": schema.StringAttribute{
								Required: true,
							},
							"project_id": schema.StringAttribute{
								Required: true,
							},
							"token_uri": schema.StringAttribute{
								Required: true,
							},
							"universe_domain": schema.StringAttribute{
								Required: true,
							},
							"client_x509_cert_url": schema.StringAttribute{
								Required: true,
							},
							"type": schema.StringAttribute{
								Required: true,
							},
						},
					},
				},
				Validators: []validator.Object{storageCredentialValidator{}},
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

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	status := "active"
	if !state.Active.ValueBool() {
		status = "inactive"
	}

	storageProfile := storageProfileModelToResource(*state.StorageProfile)

	_, err := r.client.NewWarehouse(ctx, state.ProjectID.ValueString(), state.Name.ValueString(), state.Protected.ValueBool(), status, storageProfile, nil)
	if err != nil {
		resp.Diagnostics.AddError("Lakekeeper API error occurred", fmt.Sprintf("Unable to create warehouse: %s", err.Error()))
		return
	}
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

func storageProfileModelToResource(s storageProfileModel) lakekeeper.StorageProfile {
	switch s.Type.ValueString() {
	case "s3":
		return &lakekeeper.StorageProfileS3{
			Type:   "s3",
			Bucket: s.Bucket.ValueString(),
		}
	case "gcs":
		return &lakekeeper.StorageProfileGCS{}
	case "adls":
		return &lakekeeper.StorageProfileADLS{}
	default:
		return nil
	}
}
