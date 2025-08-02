package sdk

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type Assignment struct {
	AssigneeID   types.String `tfsdk:"assignee_id"`
	AssigneeType types.String `tfsdk:"assignee_type"`
	Assignment   types.String `tfsdk:"assignment"`
}

func AssignmentDataSourceType() schema.NestedAttributeObject {
	return schema.NestedAttributeObject{
		Attributes: map[string]schema.Attribute{
			"assignee_id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: `The ID of this assignee.`,
			},
			"assignee_type": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: `The type of this assignee. Can be ` + "`user` or  `role`.",
			},
			"assignment": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: `The assignment type. Refers to the resource object documentation to see the possible values.`,
			},
		},
	}
}
