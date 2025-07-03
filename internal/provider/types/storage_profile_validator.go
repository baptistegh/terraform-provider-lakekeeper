package types

import (
	"context"
	"fmt"

	"github.com/baptistegh/terraform-provider-lakekeeper/lakekeeper/storage"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

type storageProfileValidator struct{}

func (v storageProfileValidator) Description(ctx context.Context) string {
	return "Validates storage_profile fields depending on the type"
}

func (v storageProfileValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

func (v storageProfileValidator) ValidateObject(ctx context.Context, req validator.ObjectRequest, resp *validator.ObjectResponse) {
	val := req.ConfigValue

	if val.IsNull() || val.IsUnknown() {
		return
	}

	var profile = StorageProfileModel{}

	diags := val.As(ctx, &profile, basetypes.ObjectAsOptions{})
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	switch profile.Type.ValueString() {
	case "gcs":
		if profile.Bucket.IsNull() || profile.Bucket.IsUnknown() || profile.Bucket.ValueString() == "" {
			resp.Diagnostics.AddAttributeError(
				path.Root("bucket"),
				"'bucket' required for type 'gcs'",
				"When 'type' is 'gcs', you must set the 'bucket' attribute.",
			)
		}
	case "adls":
		if profile.AccountName.IsNull() || profile.AccountName.IsUnknown() || profile.AccountName.ValueString() == "" {
			resp.Diagnostics.AddAttributeError(
				path.Root("account_name"),
				"'account_name' required for type 'adls'",
				"When 'type' is 'adls', you must set the 'account_name' attribute.",
			)
		}
		if profile.Filesystem.IsNull() || profile.Filesystem.IsUnknown() || profile.Filesystem.ValueString() == "" {
			resp.Diagnostics.AddAttributeError(
				path.Root("filesystem"),
				"'filesystem' required for type 'adls'",
				"When 'type' is 'adls', you must set the 'filesystem' attribute.",
			)
		}
	case "s3":
		if profile.Bucket.IsNull() || profile.Bucket.IsUnknown() || profile.Bucket.ValueString() == "" {
			resp.Diagnostics.AddAttributeError(
				path.Root("bucket"),
				"'bucket' required for type 's3'",
				"When 'type' is 's3', you must set the 'bucket' attribute.",
			)
		}
		if profile.Region.IsNull() || profile.Region.IsUnknown() || profile.Region.ValueString() == "" {
			resp.Diagnostics.AddAttributeError(
				path.Root("region"),
				"'region' required for type 's3'",
				"When 'type' is 's3', you must set the 'region' attribute.",
			)
		}
		if profile.STSEnabled.IsNull() || profile.STSEnabled.IsUnknown() {
			resp.Diagnostics.AddAttributeError(
				path.Root("sts_enabled"),
				"'sts_enabled' required for type 's3'",
				"When 'type' is 's3', you must set the 'sts_enabled' attribute.",
			)
		}
		if profile.STSEnabled.ValueBool() {
			if (profile.STSRoleARN.IsNull() || profile.STSRoleARN.IsUnknown() || profile.STSRoleARN.ValueString() == "") &&
				(profile.AssumeRoleARN.IsNull() || profile.AssumeRoleARN.IsUnknown() || profile.AssumeRoleARN.ValueString() == "") {
				resp.Diagnostics.AddAttributeError(
					path.Root("sts_role_arn"),
					"'assume_role_arn' or 'sts_role_arn' is required for type 's3' with STS enabled",
					"When 'type' is 's3' and 'sts_enabled' is true, you must set the 'assume_role_arn' or 'sts_role_arn' attribute.",
				)
				resp.Diagnostics.AddAttributeError(
					path.Root("assume_role_arn"),
					"'assume_role_arn' or 'sts_role_arn' is required for type 's3' with STS enabled",
					"When 'type' is 's3' and 'sts_enabled' is true, you must set the 'assume_role_arn' or 'sts_role_arn' attribute.",
				)
			}
		}
	default:
		resp.Diagnostics.AddAttributeError(
			path.Root("type"),
			"Unsupported storage profile type",
			fmt.Sprintf("The given type '%s' is not supported. Valid %v", profile.Type.ValueString(), storage.ValidStorageProfileTypes),
		)
	}
}
