package api

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"net/http"
	"os"
	"time"

	managementv1 "github.com/baptistegh/go-lakekeeper/pkg/apis/management/v1"
	lakekeeper "github.com/baptistegh/go-lakekeeper/pkg/client"
	"github.com/baptistegh/go-lakekeeper/pkg/core"
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

	if c.InitialBootstrap {
		opts = append(opts, lakekeeper.WithInitialBootstrapV1Enabled(
			true, true, core.Ptr(managementv1.ApplicationUserType),
		))
	}

	if c.UserAgent != "" {
		opts = append(opts, lakekeeper.WithUserAgent(c.UserAgent))
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

	client, err := lakekeeper.NewAuthSourceClient(ctx, &core.OAuthTokenSource{
		TokenSource: oauthConfig.TokenSource(oauthCtx),
	}, c.BaseURL, opts...)
	if err != nil {
		return nil, err
	}

	return client, nil
}
