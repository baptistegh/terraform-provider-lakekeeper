package provider

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	lakekeeper "github.com/lakekeeper/go-lakekeeper/pkg/client"
)

// resourceClient returns the Lakekeeper client passed by the provider to a resource.
// It returns nil when the provider is not configured yet, and adds an error
// diagnostic when the provider data has an unexpected type.
func resourceClient(req resource.ConfigureRequest, resp *resource.ConfigureResponse) *lakekeeper.Client {
	if req.ProviderData == nil {
		return nil
	}
	data, ok := req.ProviderData.(*LakekeeperResourceData)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *LakekeeperResourceData, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return nil
	}
	return data.Client
}

// datasourceClient returns the Lakekeeper client passed by the provider to a data source.
// It returns nil when the provider is not configured yet, and adds an error
// diagnostic when the provider data has an unexpected type.
func datasourceClient(req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) *lakekeeper.Client {
	if req.ProviderData == nil {
		return nil
	}
	data, ok := req.ProviderData.(*LakekeeperDatasourceData)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *LakekeeperDatasourceData, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return nil
	}
	return data.Client
}
