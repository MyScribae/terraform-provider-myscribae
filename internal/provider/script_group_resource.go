package provider

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/myscribae/myscribae-sdk-go/provider"
	"github.com/myscribae/myscribae-terraform-provider/validators"
)

var _ resource.Resource = (*scriptGroupResource)(nil)
var _ resource.ResourceWithConfigure = (*scriptGroupResource)(nil)

type scriptGroupResource struct {
	terraformProvider *myScribaeProvider
	myscribaeProvider *provider.Provider
	scriptGroup       *provider.ScriptGroup
}

type scriptGroupResourceData struct {
	ProviderId  types.String `tfsdk:"provider_id"`
	Id          types.String `tfsdk:"id"`
	Uuid        types.String `tfsdk:"uuid"`
	AltID       types.String `tfsdk:"alt_id"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
	Public      types.Bool   `tfsdk:"public"`
}

func (e *scriptGroupResource) MakeClient(ctx context.Context, providerId string, altId string) error {
	providerUuid, err := uuid.Parse(providerId)
	if err != nil {
		return err
	}

	e.myscribaeProvider = &provider.Provider{
		Uuid:   providerUuid,
		Client: e.terraformProvider.Client,
	}
	e.scriptGroup, err = e.myscribaeProvider.ScriptGroup(altId)
	if err != nil {
		return err
	}

	return nil
}

func newScriptGroupResource() resource.Resource {
	return &scriptGroupResource{}
}

func (e *scriptGroupResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "myscribae_script_group"
}

func (e *scriptGroupResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (e *scriptGroupResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The id of the script group",
				Computed:    true,
			},
			"provider_id": schema.StringAttribute{
				Description: "The provider id of the script group",
				Required:    true,
			},
			"uuid": schema.StringAttribute{
				Description: "The uuid of the script",
				Computed:    true,
				Required:    false,
				Optional:    false,
			},
			"alt_id": schema.StringAttribute{
				Description: "The alt id of the script group",
				Required:    true,
				Validators: []validator.String{
					validators.NewAltIdValidator(true),
				},
			},
			"name": schema.StringAttribute{
				Description: "The name of the script group",
				Required:    true,
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 100),
				},
			},
			"description": schema.StringAttribute{
				Description: "The description of the script group",
				Required:    true,
				Validators: []validator.String{
					stringvalidator.LengthBetween(3, 500),
				},
			},
			"public": schema.BoolAttribute{
				Description: "Is the script group public",
				Required:    true,
			},
		},
	}
}

func (e *scriptGroupResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	data := &scriptGroupResourceData{}

	diags := req.Plan.Get(ctx, data)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}

	if err := e.MakeClient(ctx, data.ProviderId.ValueString(), data.AltID.ValueString()); err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("failed to create script group client for create: %s", err), err.Error())
		return
	}

	resultUuid, err := e.scriptGroup.Create(ctx, provider.CreateScriptGroupInput{
		Name:        data.Name.ValueString(),
		Description: data.Description.ValueString(),
		Public:      data.Public.ValueBool(),
	})
	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("failed to create script group: %s", err), err.Error())
		return
	}

	diags = resp.State.Set(ctx, &scriptGroupResourceData{
		Id:          basetypes.NewStringValue(resultUuid.String()),
		Uuid:        basetypes.NewStringValue(resultUuid.String()),
		ProviderId:  basetypes.NewStringValue(data.ProviderId.String()),
		AltID:       data.AltID,
		Name:        data.Name,
		Description: data.Description,
		Public:      data.Public,
	})
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}
}

func (e *scriptGroupResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	data := &scriptGroupResourceData{}

	// Read the data from the provider
	diags := req.State.Get(ctx, data)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}
	if err := e.MakeClient(ctx, data.ProviderId.ValueString(), data.AltID.ValueString()); err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("failed to create script group client for read: %s", err), err.Error())
		return
	}

	// Set the data in the response
	profile, err := e.scriptGroup.Read(ctx)
	if err != nil {
		resp.Diagnostics.AddError("failed to get script group", err.Error())
		return
	}

	diags = resp.State.Set(ctx, &scriptGroupResourceData{
		Id:          basetypes.NewStringValue(profile.Uuid.String()),
		Uuid:        basetypes.NewStringValue(profile.Uuid.String()),
		ProviderId:  data.ProviderId,
		AltID:       basetypes.NewStringValue(profile.AltID),
		Name:        basetypes.NewStringValue(profile.Name),
		Description: basetypes.NewStringValue(profile.Description),
		Public:      basetypes.NewBoolValue(profile.Public),
	})
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}
}

func (e *scriptGroupResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	state := scriptGroupResourceData{}
	diags := req.State.Get(ctx, &state)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}

	data := scriptGroupResourceData{}
	diags = req.Plan.Get(ctx, &data)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}

	if err := e.MakeClient(ctx, data.ProviderId.ValueString(), data.AltID.ValueString()); err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("failed to create script group client for update: %s", err), err.Error())
		return
	}

	_, err := e.scriptGroup.Update(ctx, provider.UpdateScriptGroupInput{
		Name:        data.Name.ValueStringPointer(),
		Description: data.Description.ValueStringPointer(),
		Public:      data.Public.ValueBoolPointer(),
	})
	if err != nil {
		resp.Diagnostics.AddError("failed to update script group: %s", err.Error())
		return
	}

	state.Name = data.Name
	state.Description = data.Description
	state.Public = data.Public

	diags = resp.State.Set(ctx, &state)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}
}

func (e *scriptGroupResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	data := scriptGroupResourceData{}
	diags := req.State.Get(ctx, &data)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}

	if err := e.MakeClient(ctx, data.ProviderId.ValueString(), data.AltID.ValueString()); err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("failed to create script group client for delete: %s", err), err.Error())
		return
	}

	var public = false
	_, err := e.scriptGroup.Update(ctx, provider.UpdateScriptGroupInput{
		Public: &public,
	})
	if err != nil {
		resp.Diagnostics.AddError("failed to delete script group: %s", err.Error())
		return
	}
}
