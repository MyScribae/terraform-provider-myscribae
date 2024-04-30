package provider

import (
	"context"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/myscribae/myscribae-sdk-go/provider"
)

var _ datasource.DataSource = (*scriptGroupDataSource)(nil)
var _ datasource.DataSourceWithConfigure = (*scriptGroupDataSource)(nil)

type scriptGroupDataSource struct {
	terraformProvider *myScribaeProvider
	myscribaeProvider *provider.Provider
	scriptGroup       *provider.ScriptGroup
}

type scriptGroupResourceConfig struct {
	ProviderID types.String `tfsdk:"provider_id"`
	AltID      types.String `tfsdk:"alt_id"`
}

func newScriptGroupDataSource() datasource.DataSource {
	return &scriptGroupDataSource{}
}

func (e *scriptGroupDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = "myscribae_script_group"
}

func (e *scriptGroupDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (e *scriptGroupDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"provider_id": schema.StringAttribute{
				Description: "The provider id of the script group",
				Required:    true,
			},
			"alt_id": schema.StringAttribute{
				Description: "The alt id of the script group",
				Required:    true,
			},
		},
	}
}

func (e *scriptGroupDataSource) MakeClient(ctx context.Context, providerId string, altId string) error {
	providerUuid, err := uuid.Parse(providerId)
	if err != nil {
		return err
	}

	e.myscribaeProvider = &provider.Provider{
		Uuid:   providerUuid,
		Client: e.terraformProvider.Client,
	}
	e.scriptGroup = e.myscribaeProvider.ScriptGroup(altId)

	return nil
}

func (e *scriptGroupDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	data := &scriptGroupResourceConfig{}
	if diags := req.Config.Get(ctx, data); diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}

	profile, err := e.scriptGroup.Read(ctx)
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
