package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

type storageCredentialValidator struct{}

func (v storageCredentialValidator) Description(ctx context.Context) string {
	return "Validates storage_profile fields depending on the type"
}

func (v storageCredentialValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

func (v storageCredentialValidator) ValidateObject(ctx context.Context, req validator.ObjectRequest, resp *validator.ObjectResponse) {
	val := req.ConfigValue

	if val.IsNull() || val.IsUnknown() {
		return
	}

	var profile = storageCredentialModel{}

	diags := val.As(ctx, &profile, basetypes.ObjectAsOptions{})
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	switch profile.Type.ValueString() {
	case "s3_access_key":
		if profile.AccessKeyID.IsNull() || profile.AccessKeyID.IsUnknown() || profile.AccessKeyID.ValueString() == "" {
			resp.Diagnostics.AddAttributeError(
				path.Root("access_key_id"),
				"'access_key_id' required for type 's3_access_key'",
				"When 'type' is 's3_access_key', you must set the 'access_key_id' attribute.",
			)
		}
		if profile.SecretAccessKey.IsNull() || profile.SecretAccessKey.IsUnknown() || profile.SecretAccessKey.ValueString() == "" {
			resp.Diagnostics.AddAttributeError(
				path.Root("secret_access_key"),
				"'secret_access_key' required for type 's3_access_key'",
				"When 'type' is 's3_access_key', you must set the 'secret_access_key' attribute.",
			)
		}
	case "s3_aws_system_identity":
		break
	case "s3_cloudflare_r2":
		if profile.AccessKeyID.IsNull() || profile.AccessKeyID.IsUnknown() || profile.AccessKeyID.ValueString() == "" {
			resp.Diagnostics.AddAttributeError(
				path.Root("access_key_id"),
				"'access_key_id' required for type 's3_cloudflare_r2'",
				"When 'type' is 's3_cloudflare_r2', you must set the 'access_key_id' attribute.",
			)
		}
		if profile.SecretAccessKey.IsNull() || profile.SecretAccessKey.IsUnknown() || profile.SecretAccessKey.ValueString() == "" {
			resp.Diagnostics.AddAttributeError(
				path.Root("secret_access_key"),
				"'secret_access_key' required for type 's3_cloudflare_r2'",
				"When 'type' is 's3_cloudflare_r2', you must set the 'secret_access_key' attribute.",
			)
		}
		if profile.AccountID.IsNull() || profile.AccountID.IsUnknown() || profile.AccountID.ValueString() == "" {
			resp.Diagnostics.AddAttributeError(
				path.Root("account_id"),
				"'account_id' required for type 's3_cloudflare_r2'",
				"When 'type' is 's3_cloudflare_r2', you must set the 'account_id' attribute.",
			)
		}
		if profile.Token.IsNull() || profile.Token.IsUnknown() || profile.Token.ValueString() == "" {
			resp.Diagnostics.AddAttributeError(
				path.Root("token"),
				"'token' required for type 's3_cloudflare_r2'",
				"When 'type' is 's3_cloudflare_r2', you must set the 'token' attribute.",
			)
		}
	case "az_client_credentials":
		if profile.ClientID.IsNull() || profile.ClientID.IsUnknown() || profile.ClientID.ValueString() == "" {
			resp.Diagnostics.AddAttributeError(
				path.Root("client_id"),
				"'client_id' required for type 'az_client_credentials'",
				"When 'type' is 'az_client_credentials', you must set the 'client_id' attribute.",
			)
		}
		if profile.ClientSecret.IsNull() || profile.ClientSecret.IsUnknown() || profile.ClientSecret.ValueString() == "" {
			resp.Diagnostics.AddAttributeError(
				path.Root("client_secret"),
				"'client_secret' required for type 'az_client_credentials'",
				"When 'type' is 'az_client_credentials', you must set the 'client_secret' attribute.",
			)
		}
		if profile.TenantID.IsNull() || profile.TenantID.IsUnknown() || profile.TenantID.ValueString() == "" {
			resp.Diagnostics.AddAttributeError(
				path.Root("tenant_id"),
				"'tenant_id' required for type 'az_client_credentials'",
				"When 'type' is 'az_client_credentials', you must set the 'tenant_id' attribute.",
			)
		}
	case "az_shared_access_key":
		if profile.AZKey.IsNull() || profile.AZKey.IsUnknown() || profile.AZKey.ValueString() == "" {
			resp.Diagnostics.AddAttributeError(
				path.Root("az_key"),
				"'az_key' required for type 'az_shared_access_key'",
				"When 'type' is 'az_shared_access_key', you must set the 'az_key' attribute.",
			)
		}
	case "az_azure_system_identity":
		break
	case "gcs_service_account_key":
		if profile.Key == nil {
			resp.Diagnostics.AddAttributeError(
				path.Root("key"),
				"'key' required for type 'gcs_service_account_key'",
				"When 'type' is 'gcs_service_account_key', you must set the 'key' attribute.",
			)
		}
	case "gcs_gcp_system_identity":
		break
	default:
		resp.Diagnostics.AddAttributeError(
			path.Root("type"),
			"Unsupported storage credential type",
			fmt.Sprintf("The given type '%s' is not supported. Valid %v", profile.Type.ValueString(), validStorageCredentialTypes),
		)
	}
}
