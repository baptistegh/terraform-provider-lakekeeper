package provider

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

var (
	// testAccProtoV6ProviderFactories is used to instantiate a provider during acceptance testing.
	// The factory function is called for each Terraform CLI command to create a provider
	// server that the CLI can connect to and interact with.
	testAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
		"lakekeeper": providerserver.NewProtocol6WithError(New("acctest")()),
	}
)

func TestProvider_OIDCAuth(t *testing.T) {
	loginCall := false
	bootstrapCall := false
	serverInfoCall := false
	mockLakekeeperServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/management/v1/bootstrap" {
			bootstrapCall = true
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusNotModified)
		}
		if r.URL.Path == "/management/v1/info" {
			serverInfoCall = true
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			// nolint - don't need to err check writing the response in the test
			w.Write([]byte(`{
				"version":"0.9.1",
				"bootstrapped":true,
				"server-id":"00000000-0000-0000-0000-000000000000",
				"default-project-id":"00000000-0000-0000-0000-000000000000",
				"authz-backend":"allow-all",
				"aws-system-identities-enabled":false,
				"azure-system-identities-enabled":false,
				"gcp-system-identities-enabled":false,
				"queues":["tabular_expiration","tabular_purge"]
			}`)) // nolint - don't need to err check writing the response in the test
		}
	}))
	defer mockLakekeeperServer.Close()

	mockOIDCServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Early auth calls are made to /management/v1/whoami
		if r.URL.Path == "/token" {
			loginCall = true
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			// nolint - don't need to err check writing the response in the test
			w.Write([]byte(`{
				"access_token": "SlAV32hkKG",
				"token_type": "Bearer",
				"expires_in": 3600
			}`)) // nolint - don't need to err check writing the response in the test
		}
	}))
	defer mockOIDCServer.Close()

	//lintignore:AT001 // Providers don't need check destroy in their tests
	resource.ParallelTest(t, resource.TestCase{
		Steps: []resource.TestStep{
			{
				ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
				//lintignore:AT004 // Explicitly testing a provider configuration
				Config: fmt.Sprintf(`
					provider "lakekeeper" {
						endpoint = "%s"
						auth_url = "%s/token"
						client_id = "test-id" // doesn't matter, we're intercepting the call
						client_secret = "test-secret" // doesn't matter, we're intercepting the call
					}

					data "lakekeeper_server_info" "test" {}
					`, mockLakekeeperServer.URL, mockOIDCServer.URL),
				Check: func(*terraform.State) error {
					if !loginCall {
						return fmt.Errorf("expected a fetch token request")
					}
					if bootstrapCall {
						return fmt.Errorf("did not expect a bootstrap request")
					}
					if !serverInfoCall {
						return fmt.Errorf("expected a server info request")
					}
					return nil
				},
			},
		},
	})

}
