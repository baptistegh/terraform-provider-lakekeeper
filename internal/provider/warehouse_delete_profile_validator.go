package provider

import (
	"context"
	"fmt"

	"github.com/baptistegh/terraform-provider-lakekeeper/lakekeeper"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

type deleteProfileValidator struct{}

func (v deleteProfileValidator) Description(ctx context.Context) string {
	return "Validates delete_profile fields depending on the type"
}

func (v deleteProfileValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

func (v deleteProfileValidator) ValidateObject(ctx context.Context, req validator.ObjectRequest, resp *validator.ObjectResponse) {
	val := req.ConfigValue

	if val.IsNull() || val.IsUnknown() {
		return
	}

	var profile = deleteProfileModel{}

	diags := val.As(ctx, &profile, basetypes.ObjectAsOptions{})
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	switch profile.Type.ValueString() {
	case "soft":
		if profile.ExpirationSeconds.IsNull() || profile.ExpirationSeconds.IsUnknown() {
			resp.Diagnostics.AddAttributeError(
				path.Root("expiration_seconds"),
				"'expiration_seconds' required for type 'soft'",
				"When 'type' is 'soft', you must set the 'expiration_seconds' attribute.",
			)
		}
	case "hard":
		if !profile.ExpirationSeconds.IsNull() && !profile.ExpirationSeconds.IsUnknown() {
			resp.Diagnostics.AddAttributeError(
				path.Root("expiration_seconds"),
				"'expiration_seconds' can't be set for type 'hard'",
				"When 'type' is 'hard', 'expiration_seconds' is not used.",
			)
		}
	default:
		resp.Diagnostics.AddAttributeError(
			path.Root("type"),
			"Unsupported delete profile type",
			fmt.Sprintf("The given type '%s' is not supported. Valid %v", profile.Type.ValueString(), lakekeeper.ValidDeleteProfileTypes),
		)
	}
}
