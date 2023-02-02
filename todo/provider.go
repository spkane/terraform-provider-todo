package todo

import (
	"context"
	"os"

	"github.com/go-openapi/runtime"
	httptransport "github.com/go-openapi/runtime/client"
	strfmt "github.com/go-openapi/strfmt"

	// Todo API Libraries
	"github.com/spkane/todo-for-terraform/client"
	"github.com/spkane/todo-for-terraform/client/todos"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure the implementation satisfies the expected interfaces
var (
	_ provider.Provider = &todoProvider{}
)

// New is a helper function to simplify provider server and testing implementation.
func New() provider.Provider {
	return &todoProvider{}
}

// todoProvider is the provider implementation.
type todoProvider struct{}

// todoProviderModel maps provider schema data to a Go type.
type todoProviderModel struct {
	Host    types.String `tfsdk:"host"`
	Port    types.String `tfsdk:"port"`
	Schema  types.String `tfsdk:"schema"`
	APIPath types.String `tfsdk:"apipath"`
}

// Metadata returns the provider type name.
func (p *todoProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "todo"
}

// Schema defines the provider-level schema for configuration data.
func (p *todoProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"host": schema.StringAttribute{
				Optional:    true,
				Description: "The FQDN or IP address for the Todo server (e.g. 127.0.0.1). May also be provided via TODO_HOST environment variable.",
			},
			"port": schema.StringAttribute{
				Optional:    true,
				Description: "The port for the Todo server (e.g. 8080). May also be provided via TODO_PORT environment variable.",
			},
			"schema": schema.StringAttribute{
				Optional:    true,
				Description: "The URL schema for the Todo server (e.g. http). May also be provided via TODO_SCHEMA environment variable.",
			},
			"apipath": schema.StringAttribute{
				Optional:    true,
				Description: "The URL path for the Todo server API (e.g. /). May also be provided via TODO_APIPATH environment variable.",
			},
		},
		Blocks:      map[string]schema.Block{},
		Description: "Interface with the Todo API server (github.com/spkane/todo-for-terraform)",
	}
}

// Configure prepares a Todo API client for data sources and resources.
//
//gocyclo:ignore
func (p *todoProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	tflog.Info(ctx, "Configuring Todo client")

	// Retrieve provider data from configuration
	var config todoProviderModel
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// If practitioner provided a configuration value for any of the
	// attributes, it must be a known value.

	if config.Host.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("host"),
			"Unknown Todo API Host",
			"The provider cannot create the Todo API client as there is an unknown configuration value for the Todo API host. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the TODO_HOST environment variable.",
		)
	}

	if config.Port.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("port"),
			"Unknown Todo API Port",
			"The provider cannot create the Todo API client as there is an unknown configuration value for the Todo API port. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the TODO_PORT environment variable.",
		)
	}

	if config.Schema.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("schema"),
			"Unknown Todo API Schema",
			"The provider cannot create the Todo API client as there is an unknown configuration value for the Todo API schema. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the TODO_SCHEMA environment variable.",
		)
	}

	if config.APIPath.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("apipath"),
			"Unknown Todo API API Path",
			"The provider cannot create the Todo API client as there is an unknown configuration value for the Todo API path. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the TODO_APIPATH environment variable.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	// Default values to environment variables, but override
	// with Terraform configuration value if set.

	host := os.Getenv("TODO_HOST")
	port := os.Getenv("TODO_PORT")
	schema := os.Getenv("TODO_SCHEMA")
	apipath := os.Getenv("TODO_APIPATH")

	if !config.Host.IsNull() {
		host = config.Host.ValueString()
	}

	if !config.Port.IsNull() {
		port = config.Port.ValueString()
	}

	if !config.Schema.IsNull() {
		schema = config.Schema.ValueString()
	}

	if !config.APIPath.IsNull() {
		apipath = config.APIPath.ValueString()
	}

	// If any of the expected configurations are missing, return
	// errors with provider-specific guidance.

	if host == "" {
		resp.Diagnostics.AddAttributeWarning(
			path.Root("host"),
			"Missing Todo API Host (using default value: 127.0.0.1)",
			"The provider is using a default value as there is a missing or empty value for the Todo API host. "+
				"Set the host value in the configuration or use the TODO_HOST environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
		host = "127.0.0.1"
	}

	if port == "" {
		resp.Diagnostics.AddAttributeWarning(
			path.Root("port"),
			"Missing Todo API port (using default value: 8080)",
			"The provider is using a default value as there is a missing or empty value for the Todo API port. "+
				"Set the port value in the configuration or use the TODO_PORT environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
		port = "8080"
	}

	if schema == "" {
		resp.Diagnostics.AddAttributeWarning(
			path.Root("schema"),
			"Missing Todo API Schema (using default value: http)",
			"The provider is using a default value as there is a missing or empty value for the Todo API schema. "+
				"Set the schema value in the configuration or use the TODO_SCHEMA environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
		schema = "http"
	}

	if apipath == "" {
		resp.Diagnostics.AddAttributeWarning(
			path.Root("apipath"),
			"Missing Todo API Path (using default value: /)",
			"The provider is using a default value as there is a missing or empty value for the Todo API path. "+
				"Set the apipath value in the configuration or use the TODO_APIPATH environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
		apipath = "/"
	}

	if resp.Diagnostics.HasError() {
		return
	}

	ctx = tflog.SetField(ctx, "todo_host", host)
	ctx = tflog.SetField(ctx, "todo_port", port)
	ctx = tflog.SetField(ctx, "todo_schema", schema)
	ctx = tflog.SetField(ctx, "todo_apipath", apipath)
	// If we had a sensitive field we could mask it with something like this:
	// ctx = tflog.MaskFieldValuesWithFieldKeys(ctx, "todo_password")

	tflog.Debug(ctx, "Creating Todo client")

	// Create a new Todo client using the configuration values
	hostport := host + ":" + port
	transport := httptransport.New(hostport, apipath, []string{schema})
	transport.Consumers["application/spkane.todo-list.v1+json"] = runtime.JSONConsumer()
	transport.Producers["application/spkane.todo-list.v1+json"] = runtime.JSONProducer()
	// Instantiate the client that we will use to talk to the Todo server
	client := client.New(transport, strfmt.Default)
	params := todos.NewFindTodosParams()
	var limit int32 = 1
	params.SetLimit(&limit)
	// Let's make sure we can talk to the server now
	_, err := client.Todos.FindTodos(params)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Create Todo API Client",
			"An unexpected error occurred when creating the Todo API client. "+
				"If the error is not clear, please contact the provider developers.\n\n"+
				"Todo Client Error: "+err.Error(),
		)
		return
	}

	// Make the Todo client available during DataSource and Resource
	// type Configure methods.
	resp.DataSourceData = client
	resp.ResourceData = client

	tflog.Info(ctx, "Configured Todo client", map[string]any{"success": true})
}

// DataSources defines the data sources implemented in the provider.
func (p *todoProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewTodoDataSource,
	}
}

// Resources defines the resources implemented in the provider.
func (p *todoProvider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewTodoResource,
	}
}
