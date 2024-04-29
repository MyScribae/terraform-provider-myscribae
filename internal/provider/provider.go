package provider

import (
	"context"
	"os"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hasura/go-graphql-client"
	"github.com/myscribae/myscribae-sdk-go/gql"
)

type myScribaeProvider struct {
	ApiToken string
	ApiUrl   string
	Client   *graphql.Client
	Version  string
}

type myScribaeProviderConfig struct {
	ApiToken types.String
}

var _ provider.Provider = (*myScribaeProvider)(nil)

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &myScribaeProvider{
			Version: version,
		}
	}
}

func (p *myScribaeProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	// apiKey := os.Getenv("MYSCRIBAE_API_KEY")
	apiUrl := os.Getenv("MYSCRIBAE_API_URL")
	apiToken := os.Getenv("MYSCRIBAE_API_TOKEN")
	var cfg myScribaeProviderConfig

	diags := req.Config.Get(ctx, &cfg)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}

	if cfg.ApiToken.ValueString() != "" {
		apiToken = cfg.ApiToken.ValueString()
	}

	if apiUrl == "" {
		apiUrl = "https://api.myscribae.com"
	}

	p.ApiUrl = apiUrl
	p.ApiToken = apiToken
	p.Client = gql.CreateGraphQLClient(
		apiUrl,
		&apiToken,
	)

	resp.DataSourceData = p
	resp.ResourceData = p
}

func (p *myScribaeProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "myscribae"
}

func (p *myScribaeProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		newProviderDataSource,
		newScriptGroupDataSource,
		newScriptDataSource,
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
			"api_token": schema.StringAttribute{
				Description: "You must provide an API token to authenticate with the MyScribae API",
				Required:    true,
			},
		},
	}
}

func (mp *myScribaeProvider) New() *myScribaeProvider {
	return &myScribaeProvider{}
}
