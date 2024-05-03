package provider

import (
	"context"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
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
			"id": schema.StringAttribute{
				Description: "The id of the script group",
				Computed:    true,
			},
			"provider_id": schema.StringAttribute{
				Description: "The provider id of the script group",
				Required:    true,
			},
			"uuid": schema.StringAttribute{
				Description: "The uuid of the script group",
				Computed:    true,
			},
			"alt_id": schema.StringAttribute{
				Description: "The alt id of the script group",
				Required:    true,
			},
			"name": schema.StringAttribute{
				Description: "The name of the script group",
				Computed:    true,
			},
			"description": schema.StringAttribute{
				Description: "The description of the script group",
				Computed:    true,
			},
			"public": schema.BoolAttribute{
				Description: "The public status of the script group",
				Computed:    true,
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
	e.scriptGroup, err = e.myscribaeProvider.ScriptGroup(altId)
	if err != nil {
		return err
	}
	return nil
}

func (e *scriptGroupDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	data := &scriptGroupResourceData{}
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
		ProviderId:  data.ProviderId,
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
