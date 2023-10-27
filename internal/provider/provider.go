// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"
	"os"

	"github.com/cerebrotech/terraform-provider-ravel/internal/client"
	"github.com/cerebrotech/terraform-provider-ravel/internal/data_sources"
	"github.com/cerebrotech/terraform-provider-ravel/internal/http"
	"github.com/cerebrotech/terraform-provider-ravel/internal/resources"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

const (
	EnvRavelToken = "RAVEL_TOKEN"
	EnvRavelURL   = "RAVEL_URL"
)

// Ensure RavelProvider satisfies various provider interfaces.
var _ provider.Provider = &RavelProvider{}

// RavelProvider defines the provider implementation.
type RavelProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

// RavelProviderModel describes the provider data model.
type RavelProviderModel struct {
	URL   types.String `tfsdk:"url"`
	Token types.String `tfsdk:"token"`
}

func (p RavelProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "ravel"
	resp.Version = p.version
}

func (p RavelProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"url": schema.StringAttribute{
				Description: fmt.Sprintf("Host URL for Ravel. Can be defined from env var %s", EnvRavelURL),
				Optional:    true,
			},
			"token": schema.StringAttribute{
				Description: fmt.Sprintf("The access token for API operations. Can be defined from env var %s", EnvRavelToken),
				Optional:    true,
				Sensitive:   true,
			},
		},
	}
}

func (p RavelProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var config RavelProviderModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)

	if resp.Diagnostics.HasError() {
		return
	}

	token := os.Getenv(EnvRavelToken)
	url := os.Getenv(EnvRavelURL)

	if !config.Token.IsNull() {
		token = config.Token.ValueString()
	}

	if !config.URL.IsNull() {
		url = config.URL.ValueString()
	}

	if token == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("token"),
			"Missing Ravel API Token",
			"The provider cannot create the Ravel API client as there is a missing or empty value for the Ravel API token. "+
				"Set the token value in the configuration or use the RAVEL_TOKEN environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if url == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("url"),
			"Missing Ravel API URL",
			"The provider cannot create the Ravel API client as there is a missing or empty value for the Ravel API URL. "+
				"Set the token value in the configuration or use the RAVEL_URL environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	httpClient := http.New(url, p.version, token)
	ravelClient := client.New(httpClient)

	resp.DataSourceData = ravelClient
	resp.ResourceData = ravelClient
}

func (p RavelProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		resources.NewConfigurationResource,
	}
}

func (p RavelProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		data_sources.NewConfigurationDataSource,
	}
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &RavelProvider{
			version: version,
		}
	}
}
