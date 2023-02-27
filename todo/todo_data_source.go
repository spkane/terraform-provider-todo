package todo

import (
	"context"

	// Todo API Libraries
	"github.com/spkane/todo-for-terraform/client"
	"github.com/spkane/todo-for-terraform/client/todos"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &todoDataSource{}
	_ datasource.DataSourceWithConfigure = &todoDataSource{}
)

// NewTodoDataSource is a helper function to simplify the provider implementation.
func NewTodoDataSource() datasource.DataSource {
	return &todoDataSource{}
}

// todoDataSource is the data source implementation.
type todoDataSource struct {
	client *client.TodoList
}

// todoDataSourceModel maps the data source schema data.
type todoDataSourceModel struct {
	ID          types.Int64  `tfsdk:"id"`
	Description types.String `tfsdk:"description"`
	Completed   types.Bool   `tfsdk:"completed"`
}

// Configure adds the provider configured client to the data source.
func (d *todoDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.client = req.ProviderData.(*client.TodoList)

}

// Metadata returns the data source type name.
func (d *todoDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_todo"
}

// Schema defines the schema for the data source.
func (d *todoDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Fetch a todo.",
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Description: "The unique identifier for the todo.",
				Required:    true,
			},
			"description": schema.StringAttribute{
				Description: "The description for the todo.",
				Computed:    true,
			},
			"completed": schema.BoolAttribute{
				Description: "The completed status for the todo.",
				Computed:    true,
			},
		},
	}
}

// Read refreshes the Terraform state with the latest data.
func (d *todoDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	tflog.Debug(ctx, "Preparing to read todo data source")
	var state todoDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)

	params := todos.NewFindTodoParams()
	params.SetID(state.ID.ValueInt64())
	result, err := d.client.Todos.FindTodo(params)

	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Todo",
			err.Error(),
		)
		return
	}

	todo := result.GetPayload()

	// Map response body to model
	state = todoDataSourceModel{
		ID:          types.Int64Value(int64(todo[0].ID)),
		Description: types.StringValue(*todo[0].Description),
		Completed:   types.BoolValue(*todo[0].Completed),
	}

	// Set state
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	tflog.Debug(ctx, "Finished reading todo data source", map[string]any{"success": true})
}
