package provider

import (
	"context"
	"github.com/hashicorp/nomad/api"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func New() provider.Provider {
	return &NomadCustomDriverProvider{}
}

type NomadCustomDriverProvider struct {
	client *api.Client
	provider.Provider
}

func (p *NomadCustomDriverProvider) Metadata(_ context.Context,
	_ provider.MetadataRequest,
	resp *provider.MetadataResponse) {
	resp.TypeName = "nomad-driver"
}

func (p *NomadCustomDriverProvider) Schema(_ context.Context,
	_ provider.SchemaRequest,
	resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"nomad_address": schema.StringAttribute{
				Required: true,
			},
		},
	}
}

func (p *NomadCustomDriverProvider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		func() resource.Resource {
			return NewDriverResource(p.client)
		},
	}
}

func (p *NomadCustomDriverProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{}
}

func (p *NomadCustomDriverProvider) Configure(
	ctx context.Context,
	req provider.ConfigureRequest,
	resp *provider.ConfigureResponse) {
	var address types.String
	diags := req.Config.GetAttribute(ctx, path.Path{}.AtName("nomad_address"), &address)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	client, err := api.NewClient(&api.Config{Address: address.ValueString()})
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Create Nomad API Client",
			err.Error(),
		)
		return
	}

	p.client = client
}
