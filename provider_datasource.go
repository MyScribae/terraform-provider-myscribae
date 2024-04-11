package main

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

type mysribaeProviderDataSource struct {
	provider *myScribaeProvider
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

	prov := req.ProviderData.(*myScribaeProvider)
	e.provider = prov
}

func (e *mysribaeProviderDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{},
	}
}

func (e *mysribaeProviderDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	profile, err := e.provider.Client.Read(
		ctx,
	)
	if err != nil {
		resp.Diagnostics.AddError("error reading provider", err.Error())
		return
	}

	diags := resp.State.Set(ctx, &myscribaeProviderResourceData{
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
	})
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}
}
