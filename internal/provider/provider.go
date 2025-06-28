package provider

import (
	"context"
	"fmt"
	"os"

	"github.com/baptistegh/terraform-provider-lakekeeper/lakekeeper"
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
	CACertFile       types.String `tfsdk:"cacert_file"`
	Insecure         types.Bool   `tfsdk:"insecure"`
	InitialLogin     types.Bool   `tfsdk:"initial_login"`
	InitialBootstrap types.Bool   `tfsdk:"initial_bootstrap"`
}

type (
	LakekeeperClientOptionApplyFunc = func(lakekeeper.Config) lakekeeper.Config
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
				MarkdownDescription: "Lakekeeper base url",
				Required:            true,
			},
			"auth_url": schema.StringAttribute{
				MarkdownDescription: "OIDC Token endpoint",
				Required:            true,
			},
			"client_id": schema.StringAttribute{
				MarkdownDescription: "OIDC Client ID",
				Required:            true,
			},
			"client_secret": schema.StringAttribute{
				MarkdownDescription: "OIDC Client Secret",
				Required:            true,
				Sensitive:           true,
			},
			"cacert_file": schema.StringAttribute{
				MarkdownDescription: "This is a file containing the ca cert to verify the lakekeeper instance. This is available for use when working with a locally-issued or self-signed certificate chain.",
				Optional:            true,
			},
			"insecure": schema.BoolAttribute{
				MarkdownDescription: "When set to true this disables SSL verification of the connection to the Lakekeeper instance.",
				Optional:            true,
			},
			"initial_login": schema.BoolAttribute{
				MarkdownDescription: "When set to true ...",
				Optional:            true,
			},
			"initial_bootstrap": schema.BoolAttribute{
				MarkdownDescription: "When set to true ...",
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
			"The provider cannot create the Lakekeeper API client as there is an unknow configuration value for the Lakekeeper Base URL. "+
				"Either apply the source of the value first, set the token attribute value statically in the configuration, or use the LAKEKEEPER_ENDPOINT environment variable.",
		)
	}

	if config.Endpoint.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("auth_url"),
			"Unknown OIDC authenticate URL",
			"The provider cannot create the Lakekeeper API client as there is an unknow configuration value for the OIDC authenticate endpoint. "+
				"Either apply the source of the value first, set the auth_url attribute value statically in the configuration, or use the LAKEKEEPER_AUTH_URL environment variable.",
		)
	}

	if config.Endpoint.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("client_id"),
			"Unknown OIDC authenticate endpoint",
			"The provider cannot create the Lakekeeper API client as there is an unknow configuration value for the OIDC authenticate endpoint. "+
				"Either apply the source of the value first, set the client_id attribute value statically in the configuration, or use the LAKEKEEPER_CLIENT_ID environment variable.",
		)
	}

	if config.Endpoint.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("client_secret"),
			"Unknown OIDC authenticate endpoint",
			"The provider cannot create the Lakekeeper API client as there is an unknow configuration value for the OIDC authenticate endpoint. "+
				"Either apply the source of the value first, set the client_secret attribute value statically in the configuration, or use the LAKEKEEPER_CLIENT_SECRET environment variable.",
		)
	}

	if config.CACertFile.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("cacert_file"),
			"Unknown Lakekeeper CA Certificate File",
			"The provider cannot create the Lakekeeper API client as there is an unknown configuration value for the Lakekeeper CA Certificate File. "+
				"Either apply the source of the value first, set the token attribute value statically in the configuration.",
		)
	}

	if config.Insecure.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("insecure"),
			"Unknown Lakekeeper Insecure Flag Value",
			"The provider cannot create the Lakekeeper API client as there is an unknown configuration value for the Lakekeeper Insecure flag. "+
				"Either apply the source of the value first, set the token attribute value statically in the configuration.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	// Provider Configuration containing the values after evaluation of defaults etc.
	// Initialized with the defaults which get overridden later if config is set.
	evaluatedConfig := lakekeeper.Config{
		BaseURL: os.Getenv("LAKEKEEPER_ENDPOINT"),
		ClientCredentials: &lakekeeper.ClientCredentials{
			AuthURL:      os.Getenv("LAKEKEEPER_AUTH_URL"),
			ClientID:     os.Getenv("LAKEKEEPER_CLIENT_ID"),
			ClientSecret: os.Getenv("LAKEKEEPER_CLIENT_SECRET"),
		},
		CACertFile: "",
		Insecure:   false,
	}

	if !config.Endpoint.IsNull() {
		evaluatedConfig.BaseURL = config.Endpoint.ValueString()
	}
	if !config.AuthURL.IsNull() {
		evaluatedConfig.ClientCredentials.AuthURL = config.AuthURL.ValueString()
	}
	if !config.ClientID.IsNull() {
		evaluatedConfig.ClientCredentials.ClientID = config.ClientID.ValueString()
	}

	if !config.ClientSecret.IsNull() {
		evaluatedConfig.ClientCredentials.ClientSecret = config.ClientSecret.ValueString()
	}

	if !config.CACertFile.IsNull() {
		evaluatedConfig.CACertFile = config.CACertFile.ValueString()
	}

	if !config.Insecure.IsNull() {
		evaluatedConfig.Insecure = config.Insecure.ValueBool()
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

func newLakekeeperClient(config lakekeeper.Config, tfVersion, providerVersion string) LakekeeperClientFactory {
	return func(ctx context.Context, configFuncs ...LakekeeperClientOptionApplyFunc) (*lakekeeper.Client, error) {
		for _, f := range configFuncs {
			config = f(config)
		}

		// NOTE: there is no helper function for this available yet in the terraform-plugin-framework,
		//       see https://github.com/hashicorp/terraform-plugin-framework/issues/280
		config.UserAgent = fmt.Sprintf("Terraform/%s (+https://www.terraform.io) Terraform-Plugin-Framework terraform-provider-lakekeeper/%s", tfVersion, providerVersion)

		client, err := lakekeeper.NewClient(ctx, &config)
		if err != nil {
			return nil, fmt.Errorf("the provider failed to create a new Lakekeeper Client from the given configuration: %w", err)
		}

		return client, nil
	}
}
