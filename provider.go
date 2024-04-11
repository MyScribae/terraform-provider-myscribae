package main

import (
	"context"

	myscribae_sdk "github.com/Pritch009/myscribae-sdk-go"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

type myScribaeProvider struct {
	Client myscribae_sdk.Provider
}

var _ provider.Provider = (*myScribaeProvider)(nil)

func New() func() provider.Provider {
	return func() provider.Provider {
		return &myScribaeProvider{}
	}
}

func (p *myScribaeProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
}

func (p *myScribaeProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "myscribae"
}

func (p *myScribaeProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		// NewDataSource
	}
}

func (p *myScribaeProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		newProviderResource,
		newScriptGroupResource,
		newScriptResource,
	}
}

func (p *myScribaeProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"api_key": schema.StringAttribute{
				Description: "The API key for the provider",
				Required:    true,
			},
			"secret_key": schema.StringAttribute{
				Description: "The secret key for the provider",
				Required:    true,
			},
		},
	}
}

func (mp *myScribaeProvider) New() *myScribaeProvider {
	return &myScribaeProvider{}
}
