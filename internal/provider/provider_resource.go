package provider

import (
	"context"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/myscribae/myscribae-sdk-go/provider"
	"github.com/myscribae/myscribae-terraform-provider/validators"
)

var _ resource.Resource = (*myscribaeProviderResource)(nil)
var _ resource.ResourceWithConfigure = (*myscribaeProviderResource)(nil)

type myscribaeProviderResource struct {
	terraformProvider *myScribaeProvider
	myscribaeProvider *provider.Provider
}

type myscribaeProviderResourceData struct {
	Id             types.String `tfsdk:"id"`
	Name           types.String `tfsdk:"name"`
	AltID          types.String `tfsdk:"alt_id"`
	Uuid           types.String `tfsdk:"uuid"`
	Description    types.String `tfsdk:"description"`
	LogoUrl        types.String `tfsdk:"logo_url"`
	BannerUrl      types.String `tfsdk:"banner_url"`
	Url            types.String `tfsdk:"url"`
	Color          types.String `tfsdk:"color"`
	Public         types.Bool   `tfsdk:"public"`
	AccountService types.Bool   `tfsdk:"account_service"`
	SecretKey      types.String `tfsdk:"secret_key"`
	ApiKey         types.String `tfsdk:"api_key"`
}

func newProviderResource() resource.Resource {
	return &myscribaeProviderResource{}
}

func (e *myscribaeProviderResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "myscribae_provider"
}

func (e *myscribaeProviderResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (e *myscribaeProviderResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The id of the provider",
				Computed:    true,
			},
			"alt_id": schema.StringAttribute{
				Description: "The alt id of the provider",
				Optional:    true,
				Validators: []validator.String{
					validators.NewAltIdValidator(false),
					validators.NewXorValidator([]string{
						"uuid",
					}, true),
				},
				PlanModifiers: []planmodifier.String{},
			},
			"uuid": schema.StringAttribute{
				Description: "The uuid of the provider",
				Optional:    true,
				Validators: []validator.String{
					validators.NewUuidValidator(false),
				},
			},
			"name": schema.StringAttribute{
				Description: "The name of the provider",
				Required:    true,
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 100),
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
				Optional:    true,
				Required:    false,
				Validators: []validator.String{
					validators.NewUrlValidator(false),
				},
			},
			"banner_url": schema.StringAttribute{
				Description: "The banner url of the provider",
				Optional:    true,
				Required:    false,
				Validators: []validator.String{
					validators.NewUrlValidator(false),
				},
			},
			"url": schema.StringAttribute{
				Description: "The url of the provider",
				Optional:    true,
				Required:    false,
				Validators: []validator.String{
					validators.NewUrlValidator(false),
				},
			},
			"color": schema.StringAttribute{
				Description: "The color of the provider",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString("#A0A0A0"),
				Validators: []validator.String{
					validators.NewColorValidator(false),
				},
			},
			"public": schema.BoolAttribute{
				Description: "The public status of the provider",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(true),
			},
			"account_service": schema.BoolAttribute{
				Description: "The account service status of the provider",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(true),
			},
			"secret_key": schema.StringAttribute{
				Description: "The secret key of the provider",
				Computed:    true,
				Sensitive:   true,
			},
			"api_key": schema.StringAttribute{
				Description: "The api key of the provider",
				Computed:    true,
				Sensitive:   true,
			},
		},
	}
}

func (e *myscribaeProviderResource) MakeClient(ctx context.Context, providerId string) error {
	providerUuid, err := uuid.Parse(providerId)
	if err != nil {
		return err
	}

	e.myscribaeProvider = &provider.Provider{
		Uuid:   providerUuid,
		Client: e.terraformProvider.Client,
	}

	return nil
}

func (e *myscribaeProviderResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	planData := myscribaeProviderResourceData{}
	if diags := req.Plan.Get(ctx, &planData); diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}

	// if plan has uuid, then we just take over this provider
	// if plan does not have uuid, then we attempt to create one
	// with this provider

	var err error
	if planData.Uuid.IsNull() {
		// create a new provider
		e.myscribaeProvider, err = provider.CreateNewProvider(
			ctx,
			e.terraformProvider.Client,
			&provider.CreateProviderProfileInput{
				AltID:          planData.AltID.ValueStringPointer(),
				Name:           planData.Name.ValueString(),
				Description:    planData.Description.ValueString(),
				LogoUrl:        planData.LogoUrl.ValueStringPointer(),
				BannerUrl:      planData.BannerUrl.ValueStringPointer(),
				Url:            planData.Url.ValueStringPointer(),
				Color:          planData.Color.ValueStringPointer(),
				Public:         planData.Public.ValueBool(),
				AccountService: planData.AccountService.ValueBool(),
			},
		)

		if err != nil {
			resp.Diagnostics.AddError(
				"failed to create provider",
				err.Error(),
			)
			return
		}

		planData.Uuid = basetypes.NewStringValue(e.myscribaeProvider.Uuid.String())
		planData.Id = basetypes.NewStringValue(e.myscribaeProvider.Uuid.String())
	} else {
		// take over this provider
		err = e.MakeClient(ctx, planData.Uuid.ValueString())
		if err != nil {
			resp.Diagnostics.AddError(
				"failed to make client for create",
				err.Error(),
			)
			return
		}

		provUuid, err := e.myscribaeProvider.Update(ctx, provider.UpdateProviderProfileInput{
			AltID:          planData.AltID.ValueStringPointer(),
			Name:           planData.Name.ValueStringPointer(),
			Description:    planData.Description.ValueStringPointer(),
			LogoUrl:        planData.LogoUrl.ValueStringPointer(),
			BannerUrl:      planData.BannerUrl.ValueStringPointer(),
			Url:            planData.Url.ValueStringPointer(),
			Color:          planData.Color.ValueStringPointer(),
			Public:         planData.Public.ValueBoolPointer(),
			AccountService: planData.AccountService.ValueBoolPointer(),
		})
		if err != nil {
			resp.Diagnostics.AddError(
				"failed to update provider",
				err.Error(),
			)
			return
		}

		planData.Id = basetypes.NewStringValue(provUuid.String())
		planData.Uuid = basetypes.NewStringValue(provUuid.String())

		// if we do not have a secret key, which likely, then we must update the secret key and keep it in state
		// this is a one time operation, unless the secret key needs to be reset
		err = e.myscribaeProvider.ResetProviderKeys(ctx)
		if err != nil {
			resp.Diagnostics.AddError(
				"failed to reset provider keys",
				err.Error(),
			)
			return
		}
	}

	if err != nil {
		resp.Diagnostics.AddError(
			"failed to create provider",
			err.Error(),
		)
		return
	}

	state := myscribaeProviderResourceData{
		Id:             planData.Id,
		Uuid:           planData.Uuid,
		Name:           planData.Name,
		AltID:          planData.AltID,
		Description:    planData.Description,
		LogoUrl:        planData.LogoUrl,
		BannerUrl:      planData.BannerUrl,
		Url:            planData.Url,
		Color:          planData.Color,
		Public:         planData.Public,
		AccountService: planData.AccountService,
		SecretKey:      basetypes.NewStringPointerValue(e.myscribaeProvider.SecretKey),
		ApiKey:         basetypes.NewStringPointerValue(e.myscribaeProvider.ApiKey),
	}

	diags := resp.State.Set(ctx, state)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}
}

func (e *myscribaeProviderResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	currentState := myscribaeProviderResourceData{}
	if err := req.State.Get(ctx, &currentState); err != nil {
		resp.Diagnostics.Append(err...)
		return
	}

	if err := e.MakeClient(ctx, currentState.Id.ValueString()); err != nil {
		resp.Diagnostics.AddError(
			"failed to make client for read",
			err.Error(),
		)
		return
	}

	profile, err := e.myscribaeProvider.Read(ctx)
	if err != nil {
		resp.Diagnostics.AddError(
			"failed to get provider profile",
			err.Error(),
		)
		return
	}

	newState := myscribaeProviderResourceData{
		SecretKey:      currentState.SecretKey,
		ApiKey:         currentState.ApiKey,
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

	if d := resp.State.Set(ctx, &newState); d.HasError() {
		resp.Diagnostics.Append(d...)
	}
}

func (e *myscribaeProviderResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	currentState := myscribaeProviderResourceData{}
	if diags := req.State.Get(ctx, &currentState); diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}

	if err := e.MakeClient(ctx, currentState.Id.ValueString()); err != nil {
		resp.Diagnostics.AddError(
			"failed to make client for update",
			err.Error(),
		)
		return
	}

	planData := myscribaeProviderResourceData{}
	if diags := req.Plan.Get(ctx, &planData); diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}

	resultUuid, err := e.myscribaeProvider.Update(ctx, provider.UpdateProviderProfileInput{
		AltID:          planData.AltID.ValueStringPointer(),
		Name:           planData.Name.ValueStringPointer(),
		Description:    planData.Description.ValueStringPointer(),
		LogoUrl:        planData.LogoUrl.ValueStringPointer(),
		BannerUrl:      planData.BannerUrl.ValueStringPointer(),
		Url:            planData.Url.ValueStringPointer(),
		Color:          planData.Color.ValueStringPointer(),
		Public:         planData.Public.ValueBoolPointer(),
		AccountService: planData.AccountService.ValueBoolPointer(),
	})

	if err != nil {
		resp.Diagnostics.AddError(
			"failed to update provider",
			err.Error(),
		)
		return
	}

	planData.Id = basetypes.NewStringValue(resultUuid.String())
	planData.Uuid = basetypes.NewStringValue(resultUuid.String())
	newState := myscribaeProviderResourceData{
		SecretKey:      currentState.SecretKey,
		ApiKey:         currentState.ApiKey,
		Id:             basetypes.NewStringValue(resultUuid.String()),
		Uuid:           basetypes.NewStringValue(resultUuid.String()),
		Name:           planData.Name,
		AltID:          planData.AltID,
		Description:    planData.Description,
		LogoUrl:        planData.LogoUrl,
		BannerUrl:      planData.BannerUrl,
		Url:            planData.Url,
		Color:          planData.Color,
		Public:         planData.Public,
		AccountService: planData.AccountService,
	}

	diags := resp.State.Set(ctx, newState)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}
}

func (e *myscribaeProviderResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	currentState := myscribaeProviderResourceData{}
	if diags := req.State.Get(ctx, &currentState); diags != nil {
		resp.Diagnostics.Append(diags...)
		return
	}

	if err := e.MakeClient(ctx, currentState.Id.ValueString()); err != nil {
		resp.Diagnostics.AddError(
			"failed to make client for delete",
			err.Error(),
		)
		return
	}

	// update provider to make it private
	err := e.myscribaeProvider.SetPublic(ctx, false)
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
