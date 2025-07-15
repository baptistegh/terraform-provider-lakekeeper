package types

import (
	"regexp"

	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	dschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	rschema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type StorageProfileModel struct {
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

func StorageProfileResourceSchema() rschema.SingleNestedAttribute {
	return rschema.SingleNestedAttribute{
		MarkdownDescription: "Whether the warehouse is active.",
		Required:            true,
		Attributes: map[string]rschema.Attribute{
			"type": rschema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					stringvalidator.OneOf("gcs", "adls", "s3"),
				},
				MarkdownDescription: "The type of the storage profile. Supported values are `gcs`, `adls`, and `s3`.",
			},
			"account_name": rschema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "The account name for ADLS storage profile. Required if type is `adls`.",
			},
			"authority_host": rschema.StringAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "The authority host for ADLS storage profile. Defaults to `https://login.microsoftonline.com`.",
			},
			"filesystem": rschema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Name of the adls filesystem, in blobstorage also known as container. Required if type is `adls`.",
			},
			"host": rschema.StringAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "The host for ADLS storage profile. Defaults to `dfs.core.windows.net`.",
			},
			"sas_token_validity_seconds": rschema.Int64Attribute{
				Optional: true,
				Computed: true,
				Validators: []validator.Int64{
					int64validator.AtLeast(0),
				},
				MarkdownDescription: "The validity of the sts tokens in seconds. Default is `3600`.",
			},
			"assume_role_arn": rschema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Optional ARN to assume when accessing the bucket from Lakekeeper for S3 storage profile",
			},
			"aws_kms_key_arn": rschema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "ARN of the KMS key used to encrypt the S3 bucket, if any.",
			},
			"region": rschema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Region to use for S3 requests. Required if type is `s3`.",
			},
			"endpoint": rschema.StringAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Optional endpoint to use for S3 requests, if not provided the region will be used to determine the endpoint. If both region and endpoint are provided, the endpoint will be used. Example: `http://s3-de.my-domain.com:9000`",
				Validators: []validator.String{
					stringvalidator.RegexMatches(regexp.MustCompile("/$"), "Endpoint must ends with '/' character"),
				},
			},
			"flavor": rschema.StringAttribute{
				Optional: true,
				Computed: true,
				Validators: []validator.String{
					stringvalidator.OneOf("aws", "s3-compat"),
				},
				MarkdownDescription: "S3 flavor to use. Defaults to `aws`.",
			},
			"path_style_access": rschema.BoolAttribute{
				Optional:            true,
				MarkdownDescription: "Path style access for S3 requests. If the underlying S3 supports both, we recommend to not set path_style_access.",
			},
			"push_s3_delete_disabled": rschema.BoolAttribute{
				Computed:            true,
				Optional:            true,
				MarkdownDescription: "Controls whether the `s3.delete-enabled=false` flag is sent to clients.",
			},
			"remote_signing_url_style": rschema.StringAttribute{
				Optional: true,
				Computed: true,
				Validators: []validator.String{
					stringvalidator.OneOf("path-style", "virtual-host", "auto"),
				},
				MarkdownDescription: "S3 URL style detection mode for remote signing. One of `auto`, `path-style`, `virtual-host`. Default: `auto`.",
			},
			"sts_enabled": rschema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Whether to enable STS for S3 storage profile. Required if the storage type is `s3`. If enabled, the `sts_role_arn` or `assume_role_arn` must be provided.",
			},
			"sts_role_arn": rschema.StringAttribute{
				Optional: true,
			},
			"sts_token_validity_seconds": rschema.Int64Attribute{
				Optional: true,
				Computed: true,
				Validators: []validator.Int64{
					int64validator.AtLeast(0),
				},
				MarkdownDescription: "The validity of the STS tokens in seconds. Default is `3600`.",
			},
			"allow_alternative_protocols": rschema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Allow `s3a://`, `s3n://`, `wasbs://` in locations. This is disabled by default. We do not recommend to use this setting except for migration.",
			},
			"bucket": rschema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "The bucket name for the storage profile. Required if type is `gcs` or `s3`.",
				Validators: []validator.String{
					stringvalidator.LengthBetween(3, 64),
				},
			},
			"key_prefix": rschema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Subpath in the filesystem to use.",
			},
		},
		Validators: []validator.Object{storageProfileValidator{}},
	}

}

func StorageProfileDatasourceSchema() dschema.SingleNestedAttribute {
	return dschema.SingleNestedAttribute{
		MarkdownDescription: "Whether the warehouse is active.",
		Computed:            true,
		Attributes: map[string]dschema.Attribute{
			"type": dschema.StringAttribute{
				Computed: true,
				Validators: []validator.String{
					stringvalidator.OneOf("gcs", "adls", "s3"),
				},
				MarkdownDescription: "The type of the storage profile. Supported values are `gcs`, `adls`, and `s3`.",
			},
			"account_name": dschema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The account name for ADLS storage profile. Required if type is `adls`.",
			},
			"authority_host": dschema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The authority host for ADLS storage profile. Defaults to `https://login.microsoftonline.com`.",
			},
			"filesystem": dschema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Name of the adls filesystem, in blobstorage also known as container. Required if type is `adls`.",
			},
			"host": dschema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The host for ADLS storage profile. Defaults to `dfs.core.windows.net`.",
			},
			"sas_token_validity_seconds": dschema.Int64Attribute{
				Computed:            true,
				MarkdownDescription: "The validity of the sts tokens in seconds. Default is `3600`.",
			},
			"assume_role_arn": dschema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Optional ARN to assume when accessing the bucket from Lakekeeper for S3 storage profile",
			},
			"aws_kms_key_arn": dschema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "ARN of the KMS key used to encrypt the S3 bucket, if any.",
			},
			"region": dschema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Region to use for S3 requests. Required if type is `s3`.",
			},
			"endpoint": dschema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Optional endpoint to use for S3 requests, if not provided the region will be used to determine the endpoint. If both region and endpoint are provided, the endpoint will be used. Example: `http://s3-de.my-domain.com:9000`",
				Validators: []validator.String{
					stringvalidator.RegexMatches(regexp.MustCompile("/$"), "Endpoint must ends with '/' character"),
				},
			},
			"flavor": dschema.StringAttribute{
				Computed: true,
				Validators: []validator.String{
					stringvalidator.OneOf("aws", "s3-compat"),
				},
				MarkdownDescription: "S3 flavor to use. Defaults to `aws`.",
			},
			"path_style_access": dschema.BoolAttribute{
				Computed:            true,
				MarkdownDescription: "Path style access for S3 requests. If the underlying S3 supports both, we recommend to not set path_style_access.",
			},
			"push_s3_delete_disabled": dschema.BoolAttribute{
				Computed:            true,
				MarkdownDescription: "Controls whether the `s3.delete-enabled=false` flag is sent to clients.",
			},
			"remote_signing_url_style": dschema.StringAttribute{
				Computed: true,
				Validators: []validator.String{
					stringvalidator.OneOf("path-style", "virtual-host", "auto"),
				},
				MarkdownDescription: "S3 URL style detection mode for remote signing. One of `auto`, `path-style`, `virtual-host`. Default: `auto`.",
			},
			"sts_enabled": dschema.BoolAttribute{
				Computed:            true,
				MarkdownDescription: "Whether to enable STS for S3 storage profile. Required if the storage type is `s3`. If enabled, the `sts_role_arn` or `assume_role_arn` must be provided.",
			},
			"sts_role_arn": dschema.StringAttribute{
				Computed: true,
			},
			"sts_token_validity_seconds": dschema.Int64Attribute{
				Computed: true,
				Validators: []validator.Int64{
					int64validator.AtLeast(0),
				},
				MarkdownDescription: "The validity of the STS tokens in seconds. Default is `3600`.",
			},
			"allow_alternative_protocols": dschema.BoolAttribute{
				Computed:            true,
				MarkdownDescription: "Allow `s3a://`, `s3n://`, `wasbs://` in locations. This is disabled by default. We do not recommend to use this setting except for migration.",
			},
			"bucket": dschema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The bucket name for the storage profile. Required if type is `gcs` or `s3`.",
				Validators: []validator.String{
					stringvalidator.LengthBetween(3, 64),
				},
			},
			"key_prefix": dschema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Subpath in the filesystem to use.",
			},
		},
		Validators: []validator.Object{storageProfileValidator{}},
	}

}
