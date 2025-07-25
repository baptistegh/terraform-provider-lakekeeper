package provider

import (
	"context"
	"fmt"
	"os"

	lakekeeper "github.com/baptistegh/go-lakekeeper/pkg/client"

	"github.com/baptistegh/terraform-provider-lakekeeper/internal/provider/api"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure LakekeeperProvider satisfies various provider interfaces.
var _ provider.Provider = &LakekeeperProvider{}

// LakekeeperProvider defines the provider implementation.
type LakekeeperProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

// LakekeeperProviderModel describes the provider data model.
type LakekeeperProviderModel struct {
	Endpoint         types.String `tfsdk:"endpoint"`
	AuthURL          types.String `tfsdk:"auth_url"`
	ClientID         types.String `tfsdk:"client_id"`
	ClientSecret     types.String `tfsdk:"client_secret"`
	Scopes           types.List   `tfsdk:"scopes"`
	CACertFile       types.String `tfsdk:"cacert_file"`
	Insecure         types.Bool   `tfsdk:"insecure"`
	InitialBootstrap types.Bool   `tfsdk:"initial_bootstrap"`
}

type (
	LakekeeperClientOptionApplyFunc = func(api.Config) api.Config
	LakekeeperClientFactory         = func(ctx context.Context, configFuncs ...LakekeeperClientOptionApplyFunc) (*lakekeeper.Client, error)
)

// Attributes passed into Datasources from the Provider
type LakekeeperDatasourceData struct {
	Client *lakekeeper.Client
}

// Attributes passed into Resources from the Provider
type LakekeeperResourceData struct {
	Client              *lakekeeper.Client
	NewLakekeeperClient LakekeeperClientFactory
}

func (p *LakekeeperProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "lakekeeper"
	resp.Version = p.version
}

func (p *LakekeeperProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"endpoint": schema.StringAttribute{
				MarkdownDescription: "Lakekeeper endpoint. This is the base URL of the Lakekeeper instance, e.g. `https://lakekeeper.example.com`. It can also be set using the `LAKEKEEPER_ENDPOINT` environment variable.",
				Optional:            true,
			},
			"auth_url": schema.StringAttribute{
				MarkdownDescription: "OIDC Token endpoint. This is the URL of the OIDC authentication endpoint, e.g. `https://auth.example.com/oauth2/token`. It can also be set using the `LAKEKEEPER_AUTH_URL` environment variable.",
				Optional:            true,
			},
			"client_id": schema.StringAttribute{
				MarkdownDescription: "OIDC Client ID. This is the client ID used to authenticate with the OIDC provider, e.g. `my-client-id`. It can also be set using the `LAKEKEEPER_CLIENT_ID` environment variable.",
				Optional:            true,
			},
			"client_secret": schema.StringAttribute{
				MarkdownDescription: "OIDC Client Secret. This is the client secret used to authenticate with the OIDC provider, e.g. `my-client-secret`. It can also be set using the `LAKEKEEPER_CLIENT_SECRET` environment variable.",
				Optional:            true,
				Sensitive:           true,
			},
			"scopes": schema.ListAttribute{
				MarkdownDescription: "OIDC Scope. This is the scopes used to request the OIDC token, default `" + `["lakekeeper"]` + "`.",
				Optional:            true,
				ElementType:         types.StringType,
			},
			"cacert_file": schema.StringAttribute{
				MarkdownDescription: "This is a file containing the ca cert to verify the lakekeeper instance. This is available for use when working with a locally-issued or self-signed certificate chain.",
				Optional:            true,
			},
			"insecure": schema.BoolAttribute{
				MarkdownDescription: "When set to true this disables SSL verification of the connection to the Lakekeeper instance.",
				Optional:            true,
			},
			"initial_bootstrap": schema.BoolAttribute{
				MarkdownDescription: "When set to true, the provider will try to bootstrap the server first. default: `true`.",
				Optional:            true,
			},
		},
	}
}

func (p *LakekeeperProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var config LakekeeperProviderModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Configuration values are now available.
	if config.Endpoint.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("endpoint"),
			"Unknown Lakekeeper Base URL",
			"The provider cannot create the Lakekeeper API client as there is an unknown configuration value for the Lakekeeper Base URL. "+
				"Either apply the source of the value first, set the token attribute value statically in the configuration, or use the LAKEKEEPER_ENDPOINT environment variable.",
		)
	}

	if config.AuthURL.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("auth_url"),
			"Unknown OIDC authenticate URL",
			"The provider cannot create the Lakekeeper API client as there is an unknown configuration value for the OIDC authenticate endpoint. "+
				"Either apply the source of the value first, set the auth_url attribute value statically in the configuration, or use the LAKEKEEPER_AUTH_URL environment variable.",
		)
	}

	if config.ClientID.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("client_id"),
			"Unknown OIDC authenticate endpoint",
			"The provider cannot create the Lakekeeper API client as there is an unknown configuration value for the OIDC authenticate endpoint. "+
				"Either apply the source of the value first, set the client_id attribute value statically in the configuration, or use the LAKEKEEPER_CLIENT_ID environment variable.",
		)
	}

	if config.ClientSecret.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("client_secret"),
			"Unknown OIDC authenticate endpoint",
			"The provider cannot create the Lakekeeper API client as there is an unknown configuration value for the OIDC authenticate endpoint. "+
				"Either apply the source of the value first, set the client_secret attribute value statically in the configuration, or use the LAKEKEEPER_CLIENT_SECRET environment variable.",
		)
	}

	if config.Scopes.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("scope"),
			"Unknown OIDC authenticate endpoint",
			"The provider cannot create the Lakekeeper API client as there is an unknown configuration value for the OIDC authenticate endpoint. "+
				"Either apply the source of the value first, set the scope attribute value statically in the configuration.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	// Provider Configuration containing the values after evaluation of defaults etc.
	// Initialized with the defaults which get overridden later if config is set.
	evaluatedConfig := api.Config{
		BaseURL: os.Getenv("LAKEKEEPER_ENDPOINT"),
		OIDCClientConfig: api.OIDCClientConfig{
			AuthURL:      os.Getenv("LAKEKEEPER_AUTH_URL"),
			ClientID:     os.Getenv("LAKEKEEPER_CLIENT_ID"),
			ClientSecret: os.Getenv("LAKEKEEPER_CLIENT_SECRET"),
			Scopes:       []string{"lakekeeper"},
		},
		InitialBootstrap: true,
	}

	if !config.Endpoint.IsNull() && !config.Endpoint.IsUnknown() {
		evaluatedConfig.BaseURL = config.Endpoint.ValueString()
	}

	if !config.AuthURL.IsNull() && !config.AuthURL.IsUnknown() {
		evaluatedConfig.AuthURL = config.AuthURL.ValueString()
	}

	if !config.ClientID.IsNull() && !config.ClientID.IsUnknown() {
		evaluatedConfig.ClientID = config.ClientID.ValueString()
	}

	if !config.ClientSecret.IsNull() && !config.ClientSecret.IsUnknown() {
		evaluatedConfig.ClientSecret = config.ClientSecret.ValueString()
	}

	if !config.Scopes.IsNull() && !config.Scopes.IsUnknown() {
		resp.Diagnostics.Append(config.Scopes.ElementsAs(ctx, &config.Scopes, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
	}

	if !config.CACertFile.IsNull() && !config.CACertFile.IsUnknown() {
		evaluatedConfig.CACertFile = config.CACertFile.ValueString()
	}

	if !config.Insecure.IsNull() && !config.Insecure.IsUnknown() {
		evaluatedConfig.Insecure = config.Insecure.ValueBool()
	}

	if !config.InitialBootstrap.IsNull() && !config.InitialBootstrap.IsUnknown() {
		evaluatedConfig.InitialBootstrap = config.InitialBootstrap.ValueBool()
	}

	clientFactory := newLakekeeperClient(evaluatedConfig, req.TerraformVersion, p.version)
	lakekeeperClient, err := clientFactory(ctx)
	if err != nil {
		resp.Diagnostics.AddError("Failed to create Lakekeeper Client from provider configuration", err.Error())
		return
	}

	// Attach the client to the response so that it will be available for the Data Sources and Resources
	resp.DataSourceData = &LakekeeperDatasourceData{
		Client: lakekeeperClient,
	}
	resp.ResourceData = &LakekeeperResourceData{
		Client:              lakekeeperClient,
		NewLakekeeperClient: clientFactory,
	}
}

func (p *LakekeeperProvider) Resources(ctx context.Context) []func() resource.Resource {
	return allResources
}

func (p *LakekeeperProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return allDataSources
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &LakekeeperProvider{
			version: version,
		}
	}
}

func newLakekeeperClient(config api.Config, tfVersion, providerVersion string) LakekeeperClientFactory {
	return func(ctx context.Context, configFuncs ...LakekeeperClientOptionApplyFunc) (*lakekeeper.Client, error) {
		for _, f := range configFuncs {
			config = f(config)
		}

		// NOTE: there is no helper function for this available yet in the terraform-plugin-framework,
		//       see https://github.com/hashicorp/terraform-plugin-framework/issues/280
		config.UserAgent = fmt.Sprintf("Terraform/%s (+https://www.terraform.io) Terraform-Plugin-Framework terraform-provider-lakekeeper/%s", tfVersion, providerVersion)

		client, err := config.NewLakekeeperClient(ctx)
		if err != nil {
			return nil, fmt.Errorf("the provider failed to create a new Lakekeeper Client from the given configuration: %w", err)
		}

		return client, nil
	}
}
