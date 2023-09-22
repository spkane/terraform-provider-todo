package todo

import (
	"context"
	"strconv"

	// Todo API Libraries
	"github.com/spkane/todo-for-terraform/client"
	"github.com/spkane/todo-for-terraform/client/todos"
	"github.com/spkane/todo-for-terraform/models"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &todoResource{}
	_ resource.ResourceWithConfigure   = &todoResource{}
	_ resource.ResourceWithImportState = &todoResource{}
)

// NewTodoResource is a helper function to simplify the provider implementation.
func NewTodoResource() resource.Resource {
	return &todoResource{}
}

// todoResource is the resource implementation.
type todoResource struct {
	client *client.TodoList
}

// todoResourceModel maps the resource schema data.
type todoResourceModel struct {
	ID          types.Int64  `tfsdk:"id"`
	Description types.String `tfsdk:"description"`
	Completed   types.Bool   `tfsdk:"completed"`
}

// Configure adds the provider configured client to the resource.
func (r *todoResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	r.client = req.ProviderData.(*client.TodoList)
}

// Metadata returns the resource type name.
func (r *todoResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_todo"
}

// Schema defines the schema for the resource.
func (r *todoResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manage a todo.",
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Description: "The unique identifier for the todo.",
				Computed:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"description": schema.StringAttribute{
				Description: "The description for the todo.",
				Required:    true,
			},
			"completed": schema.BoolAttribute{
				Description: "The completed status for the todo.",
				Required:    true,
			},
		},
	}
}

func (r *todoResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	// If our ID was a string then we could do this
	//resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)

	id, err := strconv.ParseInt(req.ID, 10, 64)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error importing todo",
			"Could not import todo, unexpected error (ID should be an integer): "+err.Error(),
		)
		return
	}
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), id)...)
}

// Create a new resource
func (r *todoResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	tflog.Debug(ctx, "Preparing to create todo resource")
	// Retrieve values from plan
	var plan todoResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	description := plan.Description.ValueString()
	completed := plan.Completed.ValueBool()

	todo := models.Item{
		ID:          int64(0),
		Description: &description,
		Completed:   &completed,
	}

	params := todos.NewAddOneParams()
	params.SetBody(&todo)

	// Create new todo
	result, err := r.client.Todos.AddOne(params)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating todo",
			"Could not create todo, unexpected error: "+err.Error(),
		)
		return
	}

	// Map response body to schema and populate Computed attribute values
	resultDescription := result.GetPayload().Description
	resultCompleted := result.GetPayload().Completed
	plan.ID = types.Int64Value(result.GetPayload().ID)
	plan.Description = types.StringValue(*resultDescription)
	plan.Completed = types.BoolValue(*resultCompleted)

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Debug(ctx, "Created todo resource", map[string]any{"success": true})
}

// Read resource information
func (r *todoResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	tflog.Debug(ctx, "Preparing to read todo resource")
	// Get current state
	var state todoResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get refreshed todo value from Todo
	params := todos.NewFindTodoParams()
	params.SetID(state.ID.ValueInt64())
	result, err := r.client.Todos.FindTodo(params)
	if err != nil {
		tflog.Debug(ctx, "Error Reading Todo", map[string]interface{}{
			"ID":    state.ID.String(),
			"Error": err.Error()})
		resp.State.RemoveResource(ctx)
		return
	}

	todo := result.GetPayload()

	// Overwrite items with refreshed state
	state = todoResourceModel{
		ID:          types.Int64Value(todo[0].ID),
		Description: types.StringValue(*todo[0].Description),
		Completed:   types.BoolValue(*todo[0].Completed),
	}

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Debug(ctx, "Finished reading todo resource", map[string]any{"success": true})
}

func (r *todoResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	tflog.Debug(ctx, "Preparing to update todo resource")
	// Retrieve values from plan
	var plan todoResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	description := plan.Description.ValueString()
	completed := plan.Completed.ValueBool()

	todo := models.Item{
		Description: &description,
		Completed:   &completed,
	}

	params := todos.NewUpdateOneParams()
	params.SetBody(&todo)
	params.SetID(plan.ID.ValueInt64())

	// Update existing todo
	_, err := r.client.Todos.UpdateOne(params)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating Todo",
			"Could not update todo, unexpected error: "+err.Error(),
		)
		return
	}

	// Fetch updated items from GetTodo as UpdateTodo items are not
	// populated.
	readParams := todos.NewFindTodoParams()
	readParams.SetID(plan.ID.ValueInt64())
	result, err := r.client.Todos.FindTodo(readParams)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Todo",
			"Could not read todo ID "+plan.ID.String()+": "+err.Error(),
		)
		return
	}

	readTodo := result.GetPayload()

	// Overwrite items with refreshed state
	plan = todoResourceModel{
		ID:          types.Int64Value(readTodo[0].ID),
		Description: types.StringValue(*readTodo[0].Description),
		Completed:   types.BoolValue(*readTodo[0].Completed),
	}

	// Set refreshed state
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Debug(ctx, "Updated todo resource", map[string]any{"success": true})
}

func (r *todoResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	tflog.Debug(ctx, "Preparing to delete todo resource")
	// Retrieve values from state
	var state todoResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Delete existing todo
	params := todos.NewDestroyOneParams()
	params.SetID(state.ID.ValueInt64())
	_, err := r.client.Todos.DestroyOne(params)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting todo",
			"Could not delete todo, unexpected error: "+err.Error(),
		)
		return
	}
	tflog.Debug(ctx, "Deleted todo resource", map[string]any{"success": true})
}
