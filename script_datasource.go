package main

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

type scriptDataSource struct {
	provider *myScribaeProvider
}

var _ datasource.DataSource = (*scriptDataSource)(nil)
var _ datasource.DataSourceWithConfigure = (*scriptDataSource)(nil)

func newScriptDataSource() datasource.DataSource {
	return &scriptDataSource{}
}

func (e *scriptDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = "myscribae_script"
}

func (e *scriptDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	prov := req.ProviderData.(*myScribaeProvider)
	e.provider = prov
}

type scriptResourceConfig struct {
	AltID         types.String `tfsdk:"alt_id"`
	ScriptGroupId types.String `tfsdk:"script_group_id"`
}

func (e *scriptDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"script_group_id": schema.StringAttribute{
				Description: "The script group id",
				Required:    true,
			},
			"alt_id": schema.StringAttribute{
				Description: "The alt id of the script",
				Required:    true,
			},
		},
	}
}

func (e *scriptDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	data := &scriptResourceConfig{}
	if diags := req.Config.Get(ctx, data); diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}

	s := e.provider.Client.Script(data.ScriptGroupId.ValueString(), data.AltID.ValueString())

	profile, err := s.Read(ctx)
	if err != nil {
		resp.Diagnostics.AddError("error reading script", err.Error())
		return
	}

	diags := resp.State.Set(ctx, &scriptResourceData{
		Id:               basetypes.NewStringValue(profile.Uuid.String()),
		Uuid:             basetypes.NewStringValue(profile.Uuid.String()),
		AltID:            basetypes.NewStringValue(profile.AltID),
		Name:             basetypes.NewStringValue(profile.Name),
		Description:      basetypes.NewStringValue(profile.Description),
		Recurrence:       basetypes.NewStringValue(profile.Recurrence),
		PriceInCents:     basetypes.NewInt64Value(int64(profile.PriceInCents)),
		SlaSec:           basetypes.NewInt64Value(int64(profile.SlaSec)),
		TokenLifetimeSec: basetypes.NewInt64Value(int64(profile.TokenLifetimeSec)),
		Public:           basetypes.NewBoolValue(profile.Public),
	})
	resp.Diagnostics.Append(diags...)
}
