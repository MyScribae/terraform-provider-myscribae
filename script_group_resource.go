package main

import (
	"context"

	sdk "github.com/Pritch009/myscribae-sdk-go"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

var _ resource.Resource = (*scriptGroupResource)(nil)
var _ resource.ResourceWithConfigure = (*scriptGroupResource)(nil)

type scriptGroupResource struct {
	provider *myScribaeProvider
}

type scriptGroupResourceData struct {
	Id          types.String `tfsdk:"id"`
	Uuid        types.String `tfsdk:"uuid"`
	AltID       types.String `tfsdk:"alt_id"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
	Public      types.Bool   `tfsdk:"public"`
}

func newScriptGroupResource() resource.Resource {
	return &scriptGroupResource{}
}

func (e *scriptGroupResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "script_group"
}

func (e *scriptGroupResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	prov := req.ProviderData.(*myScribaeProvider)
	e.provider = prov
}

func (e *scriptGroupResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"alt_id": schema.StringAttribute{
				Description: "The alt id of the script group",
				Required:    true,
			},
			"name": schema.StringAttribute{
				Description: "The name of the script group",
				Required:    true,
			},
			"description": schema.StringAttribute{
				Description: "The description of the script group",
				Required:    true,
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

	sg := e.provider.Client.ScriptGroup(data.AltID.ValueString())

	resultUuid, err := sg.Create(ctx, sdk.ScriptGroupInput{
		Name:        data.Name.ValueString(),
		Description: data.Description.ValueString(),
		Public:      data.Public.ValueBool(),
	})
	if err != nil {
		resp.Diagnostics.AddError("failed to create script group: %s", err.Error())
		return
	}

	diags = resp.State.Set(ctx, &scriptGroupResourceData{
		Id:          basetypes.NewStringValue(resultUuid.String()),
		Uuid:        basetypes.NewStringValue(resultUuid.String()),
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

	sg := e.provider.Client.ScriptGroup(data.AltID.ValueString())
	// Set the data in the response
	profile, err := sg.Read(ctx)
	if err != nil {
		resp.Diagnostics.AddError("failed to get script group: %s", err.Error())
		return
	}

	diags = resp.State.Set(ctx, &scriptGroupResourceData{
		Id:          basetypes.NewStringValue(profile.Uuid.String()),
		Uuid:        basetypes.NewStringValue(profile.Uuid.String()),
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
	data := scriptGroupResourceData{}
	diags := req.Plan.Get(ctx, &data)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}

	sg := e.provider.Client.ScriptGroup(data.AltID.ValueString())
	resultUuid, err := sg.Update(ctx, sdk.ScriptGroupInput{
		Name:        data.Name.ValueString(),
		Description: data.Description.ValueString(),
		Public:      data.Public.ValueBool(),
	})
	if err != nil {
		resp.Diagnostics.AddError("failed to update script group: %s", err.Error())
		return
	}

	data.Uuid = basetypes.NewStringValue(resultUuid.String())
	data.Id = basetypes.NewStringValue(resultUuid.String())

	diags = resp.State.Set(ctx, &data)
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

	sg := e.provider.Client.ScriptGroup(data.AltID.ValueString())
	_, err := sg.Update(ctx, sdk.ScriptGroupInput{
		Public: false,
	})
	if err != nil {
		resp.Diagnostics.AddError("failed to delete script group: %s", err.Error())
		return
	}
}
