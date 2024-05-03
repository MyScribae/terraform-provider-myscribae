package provider

import (
	"context"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"

	sdk "github.com/myscribae/myscribae-sdk-go"
	"github.com/myscribae/myscribae-sdk-go/provider"
)

type mysribaeProviderDataSource struct {
	terraformProvider *myScribaeProvider
	myscribaeProvider *provider.Provider
}

type myscribaeProviderDataSourceInput struct {
	ApiKey    types.String `tfsdk:"api_key"`
	SecretKey types.String `tfsdk:"secret_key"`
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
		input myscribaeProviderDataSourceInput
		err   error
	)

	if diags := req.Config.Get(ctx, &input); diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}

	e.myscribaeProvider, err = sdk.NewProvider(provider.ProviderConfig{
		ApiToken:  &e.terraformProvider.ApiToken,
		ApiKey:    input.ApiKey.ValueStringPointer(),
		SecretKey: input.SecretKey.ValueStringPointer(),
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
		SecretKey:      input.SecretKey,
		ApiKey:         input.ApiKey,
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
		AccountService: basetypes.NewBoolValue(profile.AccountService),
	}

	diags := resp.State.Set(ctx, &state)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}
}
