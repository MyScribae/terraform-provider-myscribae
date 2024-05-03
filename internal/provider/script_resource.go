package provider

import (
	"context"
	"math"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/myscribae/myscribae-sdk-go/provider"
	"github.com/myscribae/myscribae-sdk-go/utilities"
	"github.com/myscribae/myscribae-terraform-provider/validators"
)

var _ resource.Resource = (*scriptResource)(nil)
var _ resource.ResourceWithConfigure = (*scriptResource)(nil)

type scriptResource struct {
	terraformProvider *myScribaeProvider
	myscribaeProvider *provider.Provider
	script            *provider.Script
}

type scriptResourceData struct {
	ProviderID       types.String `tfsdk:"provider_id"`
	ScriptGroupID    types.String `tfsdk:"script_group_id"`
	Id               types.String `tfsdk:"id"`
	Uuid             types.String `tfsdk:"uuid"`
	AltID            types.String `tfsdk:"alt_id"`
	Name             types.String `tfsdk:"name"`
	Description      types.String `tfsdk:"description"`
	Recurrence       types.String `tfsdk:"recurrence"`
	PriceInCents     types.Int64  `tfsdk:"price_in_cents"`
	SlaSec           types.Int64  `tfsdk:"sla_sec"`
	TokenLifetimeSec types.Int64  `tfsdk:"token_lifetime_sec"`
	Public           types.Bool   `tfsdk:"public"`
}

func newScriptResource() resource.Resource {
	return &scriptResource{}
}

func (e *scriptResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "myscribae_script"
}

func (e *scriptResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (e *scriptResource) MakeClient(ctx context.Context, providerId string, scriptGroupId string, altId string) error {
	providerUuid, err := uuid.Parse(providerId)
	if err != nil {
		return err
	}

	scriptGroupAltID, err := utilities.NewAltUuid(scriptGroupId)
	if err != nil {
		return err
	}

	e.myscribaeProvider = &provider.Provider{
		Uuid:   providerUuid,
		Client: e.terraformProvider.Client,
	}
	e.script, err = e.myscribaeProvider.Script(scriptGroupAltID, altId)
	return err
}

func (e *scriptResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The id of the script",
				Computed:    true,
				Required:    false,
				Optional:    false,
			},
			"provider_id": schema.StringAttribute{
				Description: "The provider id of the script",
				Required:    true,
				Validators: []validator.String{
					validators.NewUuidValidator(false),
				},
			},
			"script_group_id": schema.StringAttribute{
				Description: "The script group uuid",
				Required:    true,
				Validators: []validator.String{
					validators.NewUuidValidator(false),
				},
			},
			"alt_id": schema.StringAttribute{
				Description: "The alt id of the script",
				Required:    true,
				Validators: []validator.String{
					validators.NewAltIdValidator(true),
				},
			},
			"uuid": schema.StringAttribute{
				Description: "The uuid of the script",
				Computed:    true,
				Required:    false,
				Optional:    false,
			},
			"name": schema.StringAttribute{
				Description: "The name of the script",
				Required:    true,
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 100),
				},
			},
			"description": schema.StringAttribute{
				Description: "The description of the script",
				Required:    true,
				Validators: []validator.String{
					stringvalidator.LengthBetween(3, 500),
				},
			},
			"recurrence": schema.StringAttribute{
				Description: "The recurrence of the script",
				Required:    true,
				Validators: []validator.String{
					validators.NewRecurrenceValidator(),
				},
			},
			"price_in_cents": schema.Int64Attribute{
				Description: "The price in cents of the script (minimum 1)",
				Required:    true,
				Validators: []validator.Int64{
					int64validator.AtLeast(1),
				},
			},
			"sla_sec": schema.Int64Attribute{
				Description: "The SLA in seconds of the script (minimum 2400)",
				Required:    true,
				Validators: []validator.Int64{
					int64validator.AtLeast(2400),
				},
			},
			"token_lifetime_sec": schema.Int64Attribute{
				Description: "The token lifetime in seconds of the script (minimum 600)",
				Required:    true,
				Validators: []validator.Int64{
					int64validator.AtLeast(600),
				},
			},
			"public": schema.BoolAttribute{
				Description: "Is the script public",
				Required:    true,
			},
		},
	}
}

func (e *scriptResource) Plan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	data := scriptResourceData{}
	diags := req.State.Get(ctx, &data)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}

	recurrence := data.Recurrence.ValueStringPointer()
	// Check if recurrence
	if recurrence != nil && *recurrence != "" {
		diags = resp.Plan.SetAttribute(ctx, path.Root("recurrence"), data.Recurrence)
		if diags.HasError() {
			resp.Diagnostics.Append(diags...)
			return
		}
	}

}

func (e *scriptResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	data := scriptResourceData{}
	diags := req.Plan.Get(ctx, &data)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}

	if err := e.MakeClient(ctx, data.ProviderID.ValueString(), data.ScriptGroupID.ValueString(), data.AltID.ValueString()); err != nil {
		resp.Diagnostics.AddError(
			"failed to create client",
			err.Error(),
		)
	}

	recurrence, err := utilities.NewRecurrence(data.Recurrence.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"failed to parse recurrence",
			err.Error(),
		)
		return
	}

	priceInCents := uint64(data.PriceInCents.ValueInt64())
	if priceInCents > math.MaxUint32 {
		resp.Diagnostics.AddError(
			"price_in_cents is too large",
			"price_in_cents must be less than 4294967296",
		)
		return
	}

	slaSec := data.SlaSec.ValueInt64()
	if slaSec > math.MaxUint32 {
		resp.Diagnostics.AddError(
			"sla_sec is too large",
			"sla_sec must be less than 4294967296",
		)
		return
	}

	tokenLifetimeSec := data.TokenLifetimeSec.ValueInt64()
	if tokenLifetimeSec > math.MaxUint32 {
		resp.Diagnostics.AddError(
			"token_lifetime_sec is too large",
			"token_lifetime_sec must be less than 4294967296",
		)
	}

	resultUuid, err := e.script.Create(ctx, provider.CreateScriptInput{
		Name:             data.Name.ValueString(),
		Description:      data.Description.ValueString(),
		Recurrence:       *recurrence,
		PriceInCents:     utilities.NewCentValue(uint(priceInCents)),
		SlaSec:           utilities.NewUInt(uint(slaSec)),
		TokenLifetimeSec: utilities.NewUInt(uint(tokenLifetimeSec)),
		Public:           data.Public.ValueBool(),
	})
	if err != nil {
		resp.Diagnostics.AddError(
			"failed to create script",
			err.Error(),
		)
		return
	}

	diags = resp.State.Set(ctx, &scriptResourceData{
		Id:               basetypes.NewStringValue(resultUuid.String()),
		Uuid:             basetypes.NewStringValue(resultUuid.String()),
		ProviderID:       data.ProviderID,
		ScriptGroupID:    data.ScriptGroupID,
		AltID:            data.AltID,
		Name:             data.Name,
		Description:      data.Description,
		PriceInCents:     data.PriceInCents,
		SlaSec:           data.SlaSec,
		TokenLifetimeSec: data.TokenLifetimeSec,
		Public:           data.Public,
	})
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}
}

func (e *scriptResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	stateData := scriptResourceData{}
	diags := req.State.Get(ctx, &stateData)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}

	if err := e.MakeClient(ctx, stateData.ProviderID.ValueString(), stateData.ScriptGroupID.ValueString(), stateData.AltID.ValueString()); err != nil {
		resp.Diagnostics.AddError(
			"failed to create client",
			err.Error(),
		)
		return
	}

	profile, err := e.script.Read(ctx)
	if err != nil {
		resp.Diagnostics.AddError(
			"failed to get script profile",
			err.Error(),
		)
		return
	}

	diags = resp.State.Set(ctx, &scriptResourceData{
		Id:               basetypes.NewStringValue(profile.Uuid.String()),
		Uuid:             basetypes.NewStringValue(profile.Uuid.String()),
		ScriptGroupID:    stateData.ScriptGroupID,
		ProviderID:       stateData.ProviderID,
		AltID:            basetypes.NewStringValue(profile.AltID),
		Name:             basetypes.NewStringValue(profile.Name),
		Description:      basetypes.NewStringValue(profile.Description),
		Recurrence:       basetypes.NewStringValue(profile.Recurrence),
		PriceInCents:     basetypes.NewInt64Value(int64(profile.PriceInCents)),
		SlaSec:           basetypes.NewInt64Value(int64(profile.SlaSec)),
		TokenLifetimeSec: basetypes.NewInt64Value(int64(profile.TokenLifetimeSec)),
		Public:           basetypes.NewBoolValue(profile.Public),
	})
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}
}

func (e *scriptResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	stateData := scriptResourceData{}
	diags := req.State.Get(ctx, &stateData)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}

	planData := scriptResourceData{}
	diags = req.Plan.Get(ctx, &planData)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}

	if err := e.MakeClient(ctx, planData.ProviderID.ValueString(), planData.ScriptGroupID.ValueString(), planData.AltID.ValueString()); err != nil {
		resp.Diagnostics.AddError(
			"failed to create client",
			err.Error(),
		)
	}

	priceInCents := uint64(planData.PriceInCents.ValueInt64())
	if priceInCents > math.MaxUint32 {
		resp.Diagnostics.AddError(
			"price_in_cents is too large",
			"price_in_cents must be less than 4294967296",
		)
		return
	}

	slaSec := planData.SlaSec.ValueInt64()
	if slaSec > math.MaxUint32 {
		resp.Diagnostics.AddError(
			"sla_sec is too large",
			"sla_sec must be less than 4294967296",
		)
		return
	}

	tokenLifetimeSec := planData.TokenLifetimeSec.ValueInt64()
	if tokenLifetimeSec > math.MaxUint32 {
		resp.Diagnostics.AddError(
			"token_lifetime_sec is too large",
			"token_lifetime_sec must be less than 4294967296",
		)
	}

	var (
		_priceInCents     = utilities.NewCentValue(uint(priceInCents))
		_slaSec           = utilities.NewUInt(uint(slaSec))
		_tokenLifetimeSec = utilities.NewUInt(uint(tokenLifetimeSec))
	)

	resultUuid, err := e.script.Update(ctx, provider.UpdateScriptInput{
		Name:             planData.Name.ValueStringPointer(),
		Description:      planData.Description.ValueStringPointer(),
		PriceInCents:     &_priceInCents,
		SlaSec:           &_slaSec,
		TokenLifetimeSec: &_tokenLifetimeSec,
		Public:           planData.Public.ValueBoolPointer(),
	})
	if err != nil {
		resp.Diagnostics.AddError(
			"failed to update script",
			err.Error(),
		)
		return
	}

	stateData.Name = planData.Name
	stateData.Description = planData.Description
	stateData.PriceInCents = planData.PriceInCents
	stateData.SlaSec = planData.SlaSec
	stateData.TokenLifetimeSec = planData.TokenLifetimeSec
	stateData.Public = planData.Public
	stateData.Uuid = basetypes.NewStringValue(resultUuid.String())
	stateData.Id = basetypes.NewStringValue(resultUuid.String())

	diags = resp.State.Set(ctx, &stateData)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}
}

func (e *scriptResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	data := scriptResourceData{}
	diags := req.State.Get(ctx, &data)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}

	if err := e.MakeClient(ctx, data.ProviderID.ValueString(), data.ScriptGroupID.ValueString(), data.AltID.ValueString()); err != nil {
		resp.Diagnostics.AddError(
			"failed to create client",
			err.Error(),
		)
	}

	err := e.script.Delete(ctx)
	if err != nil {
		resp.Diagnostics.AddError(
			"failed to delete script",
			err.Error(),
		)
		return
	}
}
