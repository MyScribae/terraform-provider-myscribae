package main

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

type scriptGroupDataSource struct {
	provider *myScribaeProvider
}

var _ datasource.DataSource = (*scriptGroupDataSource)(nil)
var _ datasource.DataSourceWithConfigure = (*scriptGroupDataSource)(nil)

func newScriptGroupDataSource() datasource.DataSource {
	return &scriptGroupDataSource{}
}

func (e *scriptGroupDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = "data_myscribae_script_group"
}

func (e *scriptGroupDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	prov := req.ProviderData.(*myScribaeProvider)
	e.provider = prov
}

type scriptGroupResourceConfig struct {
	AltID string `tfsdk:"alt_id"`
}

func (e *scriptGroupDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"alt_id": schema.StringAttribute{
				Description: "The alt id of the script group",
				Required:    true,
			},
		},
	}
}

func (e *scriptGroupDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	data := &scriptGroupResourceConfig{}
	if diags := req.Config.Get(ctx, data); diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}

	sg := e.provider.Client.ScriptGroup(		data.AltID)

	profile, err := sg.Read(ctx)
	if err != nil {
		resp.Diagnostics.AddError("error reading script group", err.Error())
		return
	}

	diags := resp.State.Set(ctx, &scriptGroupResourceData{
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
