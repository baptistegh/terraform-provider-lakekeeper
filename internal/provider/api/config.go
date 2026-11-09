package api

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	retryablehttp "github.com/hashicorp/go-retryablehttp"
	managementv1 "github.com/lakekeeper/go-lakekeeper/pkg/apis/management/v1"
	lakekeeper "github.com/lakekeeper/go-lakekeeper/pkg/client"
	"github.com/lakekeeper/go-lakekeeper/pkg/core"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
)

// managementAPIPath is the path prefix of the Lakekeeper management API.
const managementAPIPath = "/management/v1"

type Config struct {
	BaseURL          string
	Insecure         bool
	CACertFile       string
	ClientTimeout    int
	UserAgent        string
	InitialBootstrap bool

	OIDCClientConfig
}

type OIDCClientConfig struct {
	AuthURL      string
	ClientID     string
	ClientSecret string
	Scopes       []string
}

func (c *Config) NewLakekeeperClient(ctx context.Context) (*lakekeeper.Client, error) {
	if c.AuthURL == "" {
		return nil, errors.New("no OIDC Server URI configured, either use the `oidc_server_uri` provider argument or set it as `LAKEKEEPER_AUTH_URL` environment variable")
	}

	transport, err := c.newTransport()
	if err != nil {
		return nil, err
	}

	oauthConfig := &clientcredentials.Config{
		ClientID:     c.ClientID,
		ClientSecret: c.ClientSecret,
		TokenURL:     c.AuthURL,
		Scopes:       c.Scopes,
	}

	// configure context for OAuth client
	oauthCtx := context.Background()
	httpClient := &http.Client{Timeout: 2 * time.Second}
	oauthCtx = context.WithValue(oauthCtx, oauth2.HTTPClient, httpClient)

	authSource := &core.OAuthClientCredentialsAuthSource{
		TokenSource: oauthConfig.TokenSource(oauthCtx),
	}

	var opts []lakekeeper.Option
	if c.UserAgent != "" {
		opts = append(opts, lakekeeper.WithUserAgent(c.UserAgent))
	}

	// The client adds the management API path itself. Accept an endpoint that
	// already includes it, as earlier provider versions did.
	baseURL := strings.TrimSuffix(strings.TrimRight(c.BaseURL, "/"), managementAPIPath)

	// Bootstrap is not requested through lakekeeper.WithInitialBootstrap: the
	// client runs it during construction, before the TLS transport below is
	// installed, so it would ignore the configured CA file and insecure mode.
	client, err := lakekeeper.NewWithAuthSource(ctx, baseURL, authSource, opts...)
	if err != nil {
		return nil, err
	}

	// The client has no option to supply a transport, so replace the HTTP
	// client of the shared configuration. Every API service reads it per request.
	client.GetConfig().HTTPClient = newHTTPClient(transport, authSource)

	if c.InitialBootstrap {
		if err := bootstrap(ctx, client); err != nil {
			return nil, fmt.Errorf("bootstrap: %w", err)
		}
	}

	return client, nil
}

// newTransport returns an HTTP transport honoring the CA file and insecure settings.
func (c *Config) newTransport() (*http.Transport, error) {
	tlsConfig := &tls.Config{}

	// If a CACertFile has been specified, use that for cert validation
	if c.CACertFile != "" {
		caCert, err := os.ReadFile(c.CACertFile)
		if err != nil {
			return nil, err
		}

		caCertPool := x509.NewCertPool()
		caCertPool.AppendCertsFromPEM(caCert)
		tlsConfig.RootCAs = caCertPool
	}

	// If configured as insecure, turn off SSL verification
	if c.Insecure {
		tlsConfig.InsecureSkipVerify = true
	}

	t := http.DefaultTransport.(*http.Transport).Clone()
	t.TLSClientConfig = tlsConfig
	t.MaxIdleConnsPerHost = 100

	return t, nil
}

// newHTTPClient authenticates every request through authSource and retries
// failed requests with the retryablehttp defaults.
func newHTTPClient(transport http.RoundTripper, authSource core.AuthSource) *http.Client {
	retry := retryablehttp.NewClient()
	retry.HTTPClient = &http.Client{
		Transport: &core.AuthRoundTripper{
			Base:       transport,
			AuthSource: authSource,
		},
	}
	retry.Logger = nil

	return retry.StandardClient()
}

// bootstrap bootstraps the server as an operator application, unless it is
// already bootstrapped.
func bootstrap(ctx context.Context, client *lakekeeper.Client) error {
	info, _, err := client.ServerAPI.GetServerInfo(ctx).Execute()
	if err != nil {
		return fmt.Errorf("server info: %w", err)
	}
	if info.Bootstrapped {
		return nil
	}

	req := managementv1.NewBootstrapRequest(true)
	req.SetIsOperator(true)
	req.SetUserType(managementv1.USERTYPE_APPLICATION)

	_, err = client.ServerAPI.Bootstrap(ctx).BootstrapRequest(*req).Execute()
	return err
}
