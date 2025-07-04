package provider

import (
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	allDataSources []func() datasource.DataSource
	allResources   []func() resource.Resource
)

// registerDataSource may be called during package initialization to register a new data source with the provider.
func registerDataSource(fn func() datasource.DataSource) {
	allDataSources = append(allDataSources, fn)
}

// registerResource may be called during package initialization to register a new resource with the provider.
func registerResource(fn func() resource.Resource) {
	allResources = append(allResources, fn)
}

func splitInternalID(s types.String) (string, string) {
	splitted := strings.Split(s.ValueString(), ":")
	return splitted[0], splitted[1]
}
