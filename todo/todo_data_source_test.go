package todo

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTodoDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: providerConfig + `
resource "todo_todo" "test" {
	description = "Go Shopping"
	completed = false
}

data "todo_todo" "test" {
	id = todo_todo.test.id
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify the todo to ensure all attributes are set
					resource.TestCheckResourceAttr("data.todo_todo.test", "description", "Go Shopping"),
					resource.TestCheckResourceAttr("data.todo_todo.test", "completed", "false"),
					// Verify placeholder id attribute
					resource.TestCheckResourceAttrSet("data.todo_todo.test", "id"),
				),
			},
		},
	})
}
