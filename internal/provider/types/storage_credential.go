package types

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var ValidStorageCredentialTypes = []string{"s3_access_key", "s3_aws_system_identity", "s3_cloudflare_r2",
	"az_client_credentials", "az_shared_access_key", "az_azure_system_identity",
	"gcs_service_account_key", "gcs_gcp_system_identity"}

type StorageCredentialModel struct {
	Type            types.String `tfsdk:"type"`
	ExternalID      types.String `tfsdk:"external_id"`
	AccessKeyID     types.String `tfsdk:"access_key_id"`
	SecretAccessKey types.String `tfsdk:"secret_access_key"`
	AccountID       types.String `tfsdk:"account_id"`
	Token           types.String `tfsdk:"token"`
	AZKey           types.String `tfsdk:"az_key"`
	ClientID        types.String `tfsdk:"client_id"`
	ClientSecret    types.String `tfsdk:"client_secret"`
	TenantID        types.String `tfsdk:"tenant_id"`
	Key             *GCSKeyModel `tfsdk:"key"`
}

type GCSKeyModel struct {
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

func StorageCredentialSchema() schema.SingleNestedAttribute {
	return schema.SingleNestedAttribute{
		MarkdownDescription: "The credentials used to access the storage. This is required for the warehouse to be able to access the storage profile.",
		Optional:            true,
		Attributes: map[string]schema.Attribute{
			"type": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					stringvalidator.OneOf("s3_access_key", "s3_aws_system_identity", "s3_cloudflare_r2",
						"az_client_credentials", "az_shared_access_key", "az_azure_system_identity",
						"gcs_service_account_key", "gcs_gcp_system_identity")},
				MarkdownDescription: "This is the type of credential to use. Available values are `s3_access_key` `s3_aws_system_identity` `s3_cloudflare_r2` `az_client_credentials` `az_shared_access_key` `az_azure_system_identity` `gcs_service_account_key` `gcs_gcp_system_identity`.",
			},
			"external_id": schema.StringAttribute{
				Optional: true,
			},
			"access_key_id": schema.StringAttribute{
				Optional:            true,
				Sensitive:           true,
				MarkdownDescription: "Required if type is `s3_access_key` or `s3_cloudflare_r2`",
			},
			"secret_access_key": schema.StringAttribute{
				Optional:            true,
				Sensitive:           true,
				MarkdownDescription: "Required if type is `s3_access_key` or `s3_cloudflare_r2`",
			},
			"account_id": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Required if type is `s3_cloudflare_r2`",
			},
			"token": schema.StringAttribute{
				Optional:            true,
				Sensitive:           true,
				MarkdownDescription: "Required if type is `s3_cloudflare_r2`",
			},
			"az_key": schema.StringAttribute{
				Optional:            true,
				Sensitive:           true,
				MarkdownDescription: "Required if type is `az_shared_access_key`",
			},
			"client_id": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Required if type is `az_client_credentials`",
			},
			"client_secret": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Required if type is `az_client_credentials`",
			},
			"tenant_id": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Required if type is `az_client_credentials`",
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
				MarkdownDescription: "Required if type is `gcs_service_account_key`",
			},
		},
		Validators: []validator.Object{storageCredentialValidator{}},
	}
}
