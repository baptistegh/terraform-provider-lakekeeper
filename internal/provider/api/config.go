package api

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/baptistegh/terraform-provider-lakekeeper/lakekeeper"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
)

type Config struct {
	BaseURL          string
	Insecure         bool
	CACertFile       string
	ClientTimeout    int
	UserAgent        string
	InitialBootstrap bool
	Headers          map[string]any
	EarlyAuthFail    bool

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

	// Configure TLS/SSL
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

	opts := []lakekeeper.ClientOptionFunc{
		lakekeeper.WithHTTPClient(
			&http.Client{
				Transport: t,
			},
		),
	}

	if c.Headers != nil {
		stringMap := make(map[string]string, len(c.Headers))
		for k, v := range c.Headers {
			stringMap[k] = fmt.Sprintf("%v", v)
		}
		opts = append(opts, lakekeeper.WithRequestOptions(
			lakekeeper.WithHeaders(stringMap),
		))
	}

	if !c.InitialBootstrap {
		opts = append(opts, lakekeeper.WithBootstrapDisabled())
	}

	oauthConfig := &clientcredentials.Config{
		ClientID:     c.ClientID,
		ClientSecret: c.ClientSecret,
		TokenURL:     c.AuthURL,
		Scopes:       c.Scopes,
	}

	initialToken, err := oauthConfig.Token(ctx)
	if err != nil {
		return nil, err
	}

	client, err := lakekeeper.NewAuthSourceClient(lakekeeper.OAuthTokenSource{
		TokenSource: oauth2.ReuseTokenSource(initialToken, oauthConfig.TokenSource(ctx)),
	}, c.BaseURL, opts...)
	if err != nil {
		return nil, err
	}

	// Test the credentials by checking we can get information about the authenticated user.
	if c.EarlyAuthFail {
		_, _, err = client.Server.Info(lakekeeper.WithContext(ctx))
	}

	return client, err

}
