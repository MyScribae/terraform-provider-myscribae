package main

import (
	"context"

	sdk "github.com/MyScribae/myscribae-sdk-go/provider"
	"github.com/MyScribae/myscribae-terraform-provider/validators"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

var _ resource.Resource = (*myscribaeProviderResource)(nil)
var _ resource.ResourceWithConfigure = (*myscribaeProviderResource)(nil)

type myscribaeProviderResource struct {
	provider *myScribaeProvider
}

type myscribaeProviderResourceData struct {
	Id             types.String `tfsdk:"id"`
	Uuid           types.String `tfsdk:"uuid"`
	Name           types.String `tfsdk:"name"`
	AltID          types.String `tfsdk:"alt_id"`
	Description    types.String `tfsdk:"description"`
	LogoUrl        types.String `tfsdk:"logo_url"`
	BannerUrl      types.String `tfsdk:"banner_url"`
	Url            types.String `tfsdk:"url"`
	Color          types.String `tfsdk:"color"`
	Public         types.Bool   `tfsdk:"public"`
	AccountService types.Bool   `tfsdk:"account_service"`
}

func newProviderResource() resource.Resource {
	return &myscribaeProviderResource{}
}

func (e *myscribaeProviderResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "provider"
}

func (e *myscribaeProviderResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	prov := req.ProviderData.(*myScribaeProvider)
	e.provider = prov
}

func (e *myscribaeProviderResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Description: "The name of the provider",
				Required:    true,
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 100),
				},
			},
			"alt_id": schema.StringAttribute{
				Description: "The alt id of the provider",
				Required:    false,
				Validators: []validator.String{
					validators.NewAltIdValidator(),
				},
			},
			"description": schema.StringAttribute{
				Description: "The description of the provider",
				Required:    true,
				Validators: []validator.String{
					stringvalidator.LengthBetween(3, 500),
				},
			},
			"logo_url": schema.StringAttribute{
				Description: "The logo url of the provider",
				Required:    false,
				Validators: []validator.String{
					validators.NewUrlValidator(),
				},
			},
			"banner_url": schema.StringAttribute{
				Description: "The banner url of the provider",
				Required:    false,
				Validators: []validator.String{
					validators.NewUrlValidator(),
				},
			},
			"url": schema.StringAttribute{
				Description: "The url of the provider",
				Required:    false,
				Validators: []validator.String{
					validators.NewUrlValidator(),
				},
			},
			"color": schema.StringAttribute{
				Description: "The color of the provider",
				Required:    false,
				Validators: []validator.String{
					validators.NewColorValidator(),
				},
			},
			"public": schema.BoolAttribute{
				Description: "The public status of the provider",
				Required:    true,
			},
			"account_service": schema.BoolAttribute{
				Description: "The account service status of the provider",
				Required:    true,
			},
		},
	}
}

func (e *myscribaeProviderResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	data := myscribaeProviderResourceData{}
	if err := req.Plan.Get(ctx, &data); err != nil {
		resp.Diagnostics.Append(err...)
		return
	}

	resultUuid, err := e.provider.Client.Update(ctx, sdk.ProviderProfileInput{
		AltID:          data.AltID.ValueString(),
		Name:           data.Name.ValueString(),
		Description:    data.Description.ValueString(),
		LogoUrl:        data.LogoUrl.ValueStringPointer(),
		BannerUrl:      data.BannerUrl.ValueStringPointer(),
		Url:            data.Url.ValueStringPointer(),
		Color:          data.Color.ValueStringPointer(),
		Public:         data.Public.ValueBool(),
		AccountService: data.AccountService.ValueBool(),
	})

	if err != nil {
		resp.Diagnostics.Append(
			[]diag.Diagnostic{
				diag.NewErrorDiagnostic(
					"failed to create provider",
					err.Error(),
				),
			}...,
		)
		return
	}

	data.Id = basetypes.NewStringValue(resultUuid.String())
	data.Uuid = basetypes.NewStringValue(resultUuid.String())

	diags := resp.State.Set(ctx, data)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}
}

func (e *myscribaeProviderResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	data := myscribaeProviderResourceData{}
	if err := req.State.Get(ctx, &data); err != nil {
		resp.Diagnostics.Append(err...)
		return
	}

	profile, err := e.provider.Client.Read(ctx)
	if err != nil {
		resp.Diagnostics.AddError(
			"failed to get provider profile",
			err.Error(),
		)
		return
	}

	data = myscribaeProviderResourceData{
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

	d := resp.State.Set(ctx, data)
	resp.Diagnostics.Append(d...)
}

func (e *myscribaeProviderResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	data := myscribaeProviderResourceData{}
	if err := req.Plan.Get(ctx, &data); err != nil {
		resp.Diagnostics.Append(err...)
		return
	}

	resultUuid, err := e.provider.Client.Update(ctx, sdk.ProviderProfileInput{
		AltID:          data.AltID.ValueString(),
		Name:           data.Name.ValueString(),
		Description:    data.Description.ValueString(),
		LogoUrl:        data.LogoUrl.ValueStringPointer(),
		BannerUrl:      data.BannerUrl.ValueStringPointer(),
		Url:            data.Url.ValueStringPointer(),
		Color:          data.Color.ValueStringPointer(),
		Public:         data.Public.ValueBool(),
		AccountService: data.AccountService.ValueBool(),
	})

	if err != nil {
		resp.Diagnostics.AddError(
			"failed to update provider",
			err.Error(),
		)
		return
	}

	data.Id = basetypes.NewStringValue(resultUuid.String())
	data.Uuid = basetypes.NewStringValue(resultUuid.String())

	diags := resp.State.Set(ctx, data)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}
}

func (e *myscribaeProviderResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	data := myscribaeProviderResourceData{}
	if err := req.State.Get(ctx, &data); err != nil {
		resp.Diagnostics.Append(err...)
		return
	}

	// update provider to make it private
	err := e.provider.Client.SetPublic(ctx, false)
	if err != nil {
		resp.Diagnostics.Append(
			[]diag.Diagnostic{
				diag.NewErrorDiagnostic(
					"failed to delete provider",
					err.Error(),
				),
			}...,
		)
		return
	}
}
