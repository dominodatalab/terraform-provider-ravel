// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package resources

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/cerebrotech/terraform-provider-ravel/internal/client"
	"github.com/cerebrotech/terraform-provider-ravel/internal/models"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/mapplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &ConfigurationResource{}
var _ resource.ResourceWithImportState = &ConfigurationResource{}

func NewConfigurationResource() resource.Resource {
	return &ConfigurationResource{}
}

type ConfigurationResource struct {
	client *client.RavelClient
}

type ConfigurationSchemaModel struct {
	Version types.String            `tfsdk:"version"`
	Name    types.String            `tfsdk:"name"`
	Scope   map[string]types.String `tfsdk:"scope"`
}

type ConfigurationResourceModel struct {
	Id         types.String              `tfsdk:"id"`
	Version    types.Int64               `tfsdk:"version"`
	Name       types.String              `tfsdk:"name"`
	Labels     map[string]types.String   `tfsdk:"labels"`
	Scope      map[string]types.String   `tfsdk:"scope"`
	Schema     *ConfigurationSchemaModel `tfsdk:"schema"`
	Definition types.String              `tfsdk:"definition"`
}

func (r *ConfigurationResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_configuration"
}

func (r *ConfigurationResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Example resource",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Configuration identifier",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"version": schema.Int64Attribute{
				Computed:            true,
				MarkdownDescription: "Configuration version",
			},
			"name": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Configuration name",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"labels": schema.MapAttribute{
				MarkdownDescription: "Configuration labels (Map<String, String>)",
				ElementType:         types.StringType,
				Optional:            true,
			},
			"scope": schema.MapAttribute{
				MarkdownDescription: "Configuration scope (Map<String, String>)",
				ElementType:         types.StringType,
				Optional:            true,
				PlanModifiers: []planmodifier.Map{
					mapplanmodifier.UseStateForUnknown(),
				},
			},
			"schema": schema.SingleNestedAttribute{
				Optional: true,
				Attributes: map[string]schema.Attribute{
					"name": schema.StringAttribute{
						Required:            true,
						MarkdownDescription: "Schema name",
					},
					"version": schema.StringAttribute{
						Required:            true,
						MarkdownDescription: "Schema version",
					},
					"scope": schema.MapAttribute{
						MarkdownDescription: "Schema scope (Map<String, String>)",
						ElementType:         types.StringType,
						Optional:            true,
					},
				},
			},
			"definition": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Configuration definition",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

func (r *ConfigurationResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*client.RavelClient)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *client.RavelClient, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.client = client
}

func (r *ConfigurationResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	r.Upsert(ctx, &resp.Diagnostics, &req.Plan, &resp.State)
}

func (r *ConfigurationResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	r.Upsert(ctx, &resp.Diagnostics, &req.Plan, &resp.State)
}

func (r *ConfigurationResource) Upsert(ctx context.Context, diagnostics *diag.Diagnostics, plan *tfsdk.Plan, state *tfsdk.State) {
	var data *ConfigurationResourceModel

	// Read Terraform plan data into the model
	diagnostics.Append(plan.Get(ctx, &data)...)

	if diagnostics.HasError() {
		return
	}

	scope := make(map[string]string, len(data.Scope))
	for key, elem := range data.Scope {
		scope[key] = elem.ValueString()
	}

	labels := make(map[string]string, len(data.Labels))
	for key, elem := range data.Labels {
		labels[key] = elem.ValueString()
	}

	meta := models.RavelConfigMeta{
		RavelResourceMeta: models.RavelResourceMeta{
			Name:   data.Name.ValueString(),
			Scope:  scope,
			Labels: labels,
		},
	}

	var schema *models.RavelSchemaMeta
	if data.Schema != nil {
		schemaScope := make(map[string]string, len(data.Schema.Scope))
		for key, elem := range data.Schema.Scope {
			schemaScope[key] = elem.ValueString()
		}

		schema = &models.RavelSchemaMeta{
			RavelResourceMeta: models.RavelResourceMeta{
				Name:  data.Schema.Name.ValueString(),
				Scope: schemaScope,
			},
			Version: data.Schema.Version.ValueString(),
		}
	}

	var definition map[string]interface{}
	if err := json.Unmarshal([]byte(data.Definition.ValueString()), &definition); err != nil {
		return
	}

	createdConf, err := r.client.CreateConfig(ctx, meta, schema, definition)
	if err != nil {
		diagnostics.AddError(
			"Error Creating Ravel configuration",
			err.Error(),
		)
		return
	}

	var createdDef []byte
	createdDef, err = json.Marshal(createdConf.Spec.Def)
	if err != nil {
		return
	}

	data.Id = types.StringValue(createdConf.Id)
	data.Version = types.Int64Value(createdConf.Meta.Version)
	data.Definition = types.StringValue(string(createdDef))

	tflog.Trace(ctx, fmt.Sprintf("upserted a configuration with id: %s and version: %d", createdConf.Id, createdConf.Meta.Version))

	diagnostics.Append(state.Set(ctx, &data)...)
}

func (r *ConfigurationResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *ConfigurationResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	configuration, err := r.client.GetConfigVersion(ctx, data.Id.ValueString(), int(data.Version.ValueInt64()))
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Ravel configuration",
			fmt.Sprintf("Could not read Ravel configuration ID: %s and version: %d. Error: %s ", data.Id.ValueString(), data.Version.ValueInt64(), err.Error()),
		)
		return
	}

	data.Name = types.StringValue(configuration.Meta.Name)

	var labels map[string]types.String
	if configuration.Meta.Labels != nil {
		labels = make(map[string]types.String, len(configuration.Meta.Labels))

		for key, val := range configuration.Meta.Labels {
			labels[key] = types.StringValue(val)
		}
	}
	data.Labels = labels

	var scope map[string]types.String
	if configuration.Meta.Scope != nil {
		scope = make(map[string]types.String, len(configuration.Meta.Scope))

		for key, val := range configuration.Meta.Scope {
			scope[key] = types.StringValue(val)
		}
	}
	data.Scope = scope

	if configuration.Spec.ConfigurationFormat != nil {
		schema := ConfigurationSchemaModel{}

		schema.Name = types.StringValue(configuration.Spec.ConfigurationFormat.Name)
		schema.Version = types.StringValue(configuration.Spec.ConfigurationFormat.Version)

		data.Schema = &schema
	}

	var definition []byte
	definition, err = json.Marshal(configuration.Spec.Def)
	if err != nil {
		return
	}
	data.Definition = types.StringValue(string(definition))

	tflog.Trace(ctx, fmt.Sprintf("read configuration with id: %s and version: %d", data.Id, data.Version.ValueInt64()))

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ConfigurationResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *ConfigurationResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteConfig(ctx, data.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting Ravel configuration",
			fmt.Sprintf("Could not delete Ravel configuration ID: %s. Error: %s ", data.Id.ValueString(), err.Error()),
		)
		return
	}
}

func (r *ConfigurationResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
