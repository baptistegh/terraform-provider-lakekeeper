package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	lakekeeper "github.com/lakekeeper/go-lakekeeper/pkg/client"
)

func TestResourceClient(t *testing.T) {
	c := &lakekeeper.Client{}
	for _, tc := range []struct {
		name      string
		data      any
		want      *lakekeeper.Client
		wantError bool
	}{
		{name: "unconfigured", data: nil},
		{name: "configured", data: &LakekeeperResourceData{Client: c}, want: c},
		{name: "unexpected type", data: &LakekeeperDatasourceData{Client: c}, wantError: true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			resp := &resource.ConfigureResponse{}
			got := resourceClient(resource.ConfigureRequest{ProviderData: tc.data}, resp)
			if got != tc.want {
				t.Errorf("got client %p, want %p", got, tc.want)
			}
			if resp.Diagnostics.HasError() != tc.wantError {
				t.Errorf("HasError() = %v, want %v: %v", resp.Diagnostics.HasError(), tc.wantError, resp.Diagnostics)
			}
		})
	}
}

func TestDatasourceClient(t *testing.T) {
	c := &lakekeeper.Client{}
	for _, tc := range []struct {
		name      string
		data      any
		want      *lakekeeper.Client
		wantError bool
	}{
		{name: "unconfigured", data: nil},
		{name: "configured", data: &LakekeeperDatasourceData{Client: c}, want: c},
		{name: "unexpected type", data: &LakekeeperResourceData{Client: c}, wantError: true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			resp := &datasource.ConfigureResponse{}
			got := datasourceClient(datasource.ConfigureRequest{ProviderData: tc.data}, resp)
			if got != tc.want {
				t.Errorf("got client %p, want %p", got, tc.want)
			}
			if resp.Diagnostics.HasError() != tc.wantError {
				t.Errorf("HasError() = %v, want %v: %v", resp.Diagnostics.HasError(), tc.wantError, resp.Diagnostics)
			}
		})
	}
}
