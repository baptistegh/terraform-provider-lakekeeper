package types

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	dschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	rschema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/terraform-community-providers/terraform-plugin-framework-utils/modifiers"
)

var validDeleteProfileTypes = []string{"soft", "hard"}

type DeleteProfileModel struct {
	Type              types.String `tfsdk:"type"`
	ExpirationSeconds types.Int32  `tfsdk:"expiration_seconds"`
}

func (d DeleteProfileModel) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"type":               types.StringType,
		"expiration_seconds": types.Int32Type,
	}
}

func DeleteProfileResourceSchema() rschema.SingleNestedAttribute {
	return rschema.SingleNestedAttribute{
		MarkdownDescription: "The delete profile for the warehouse. It can be either a soft or hard delete profile. Default: `hard`",
		Optional:            true,
		Computed:            true,
		PlanModifiers: []planmodifier.Object{
			modifiers.UnknownAttributesOnUnknown(),
		},
		Attributes: map[string]rschema.Attribute{
			"type": rschema.StringAttribute{
				Optional: true,
				Computed: true,
				Validators: []validator.String{
					stringvalidator.OneOf("soft", "hard"),
				},
				Default: stringdefault.StaticString("hard"),
			},
			"expiration_seconds": rschema.Int32Attribute{
				Optional: true,
			},
		},
		Validators: []validator.Object{deleteProfileValidator{}},
	}
}

func DeleteProfileDatasourceSchema() dschema.SingleNestedAttribute {
	return dschema.SingleNestedAttribute{
		MarkdownDescription: "The delete profile for the warehouse. It can be either a soft or hard delete profile.",
		Computed:            true,
		Attributes: map[string]dschema.Attribute{
			"type": dschema.StringAttribute{
				Computed: true,
				Validators: []validator.String{
					stringvalidator.OneOf("soft", "hard"),
				},
			},
			"expiration_seconds": dschema.Int32Attribute{
				Computed: true,
			},
		},
		Validators: []validator.Object{deleteProfileValidator{}},
	}
}

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

	var deleteProfile DeleteProfileModel

	diags := val.As(ctx, &deleteProfile, basetypes.ObjectAsOptions{})
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	switch deleteProfile.Type.ValueString() {
	case "soft":
		if deleteProfile.ExpirationSeconds.IsNull() {
			resp.Diagnostics.AddAttributeError(
				path.Root("expiration_seconds"),
				"'expiration_seconds' required for type 'soft'",
				"When 'type' is 'soft', you must set the 'expiration_seconds' attribute.",
			)
		}
	case deleteProfile.Type.ValueString():
		if !deleteProfile.ExpirationSeconds.IsNull() {
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
			fmt.Sprintf("The given type '%s' is not supported. Valid %v", deleteProfile.Type.ValueString(), validDeleteProfileTypes),
		)
	}
}
