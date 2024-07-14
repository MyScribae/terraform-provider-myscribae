package provider

import (
	"context"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"

	sdk "github.com/myscribae/myscribae-sdk-go"
	"github.com/myscribae/myscribae-sdk-go/provider"
)

type mysribaeProviderDataSource struct {
	terraformProvider *myScribaeProvider
	myscribaeProvider *provider.Provider
}

var _ datasource.DataSource = (*mysribaeProviderDataSource)(nil)

func newProviderDataSource() datasource.DataSource {
	return &mysribaeProviderDataSource{}
}

func (e *mysribaeProviderDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = "myscribae_provider"
}

func (e *mysribaeProviderDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	prov, ok := req.ProviderData.(*myScribaeProvider)
	if !ok {
		resp.Diagnostics.AddError("invalid provider data", "expected *myScribaeProvider")
		return
	}

	e.terraformProvider = prov
}

func (e *mysribaeProviderDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The id of the provider",
				Computed:    true,
			},
			"alt_id": schema.StringAttribute{
				Description: "The alt id of the provider",
				Computed:    true,
			},
			"uuid": schema.StringAttribute{
				Description: "The uuid of the provider",
				Computed:    true,
			},
			"name": schema.StringAttribute{
				Description: "The name of the provider",
				Computed:    true,
			},
			"description": schema.StringAttribute{
				Description: "The description of the provider",
				Computed:    true,
			},
			"logo_url": schema.StringAttribute{
				Description: "The logo url of the provider",
				Computed:    true,
			},
			"banner_url": schema.StringAttribute{
				Description: "The banner url of the provider",
				Computed:    true,
			},
			"url": schema.StringAttribute{
				Description: "The url of the provider",
				Computed:    true,
			},
			"color": schema.StringAttribute{
				Description: "The color choice of the provider as hex color code (e.g. #000000)",
				Computed:    true,
			},
			"public": schema.BoolAttribute{
				Description: "Is the provider public",
				Computed:    true,
			},
			"account_service": schema.BoolAttribute{
				Description: "Is the provider an account service",
				Computed:    true,
			},
			"secret_key": schema.StringAttribute{
				Description: "The secret key of the provider",
				Optional:    true,
			},
			"api_key": schema.StringAttribute{
				Description: "The api key of the provider",
				Optional:    true,
			},
		},
	}
}

func (e *mysribaeProviderDataSource) MakeClient(
	ctx context.Context,
	providerId string,
	apiKey *string,
	secretKey *string,
) (err error) {
	providerUuid, err := uuid.Parse(providerId)
	if err != nil {
		return err
	}

	e.myscribaeProvider = &provider.Provider{
		ApiUrl:    e.terraformProvider.ApiUrl,
		Uuid:      providerUuid,
		ApiKey:    apiKey,
		SecretKey: secretKey,
	}

	return nil
}

func (e *mysribaeProviderDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var (
		data myscribaeProviderResourceData
		err  error
	)

	if diags := req.Config.Get(ctx, &data); diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}

	e.myscribaeProvider, err = sdk.NewProvider(provider.ProviderConfig{
		ApiToken:  &e.terraformProvider.ApiToken,
		ApiKey:    data.ApiKey.ValueStringPointer(),
		SecretKey: data.SecretKey.ValueStringPointer(),
	})
	if err != nil {
		resp.Diagnostics.AddError("error creating provider client", err.Error())
		return
	}

	profile, err := e.myscribaeProvider.Read(ctx)
	if err != nil {
		resp.Diagnostics.AddError("error reading provider", err.Error())
		return
	}

	state := myscribaeProviderResourceData{
		SecretKey:      data.SecretKey,
		ApiKey:         data.ApiKey,
		Id:             basetypes.NewStringValue(profile.Uuid.String()),
		Uuid:           basetypes.NewStringValue(profile.Uuid.String()),
		Name:           basetypes.NewStringValue(profile.Name),
		AltID:          basetypes.NewStringPointerValue(profile.AltID),
		Description:    basetypes.NewStringValue(profile.Description),
		LogoUrl:        basetypes.NewStringPointerValue(profile.LogoUrl),
		BannerUrl:      basetypes.NewStringPointerValue(profile.BannerUrl),
		Url:            basetypes.NewStringPointerValue(profile.Url),
		Color:          basetypes.NewStringPointerValue(profile.Color),
		Public:         basetypes.NewBoolValue(profile.Public),
		AccountService: basetypes.NewBoolValue(profile.AccountService.Enabled),
	}

	diags := resp.State.Set(ctx, &state)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}
}
