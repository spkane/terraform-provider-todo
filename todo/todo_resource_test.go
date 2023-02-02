package todo

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTodoResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: providerConfig + `
resource "todo_todo" "test" {
	description = "Go Shopping"
	completed = false
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("todo_todo.test", "description", "Go Shopping"),
					resource.TestCheckResourceAttr("todo_todo.test", "completed", "false"),
					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet("todo_todo.test", "id"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "todo_todo.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update and Read testing
			{
				Config: providerConfig + `
resource "todo_todo" "test" {
	description = "Go shopping for avocados"
	completed = true
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("todo_todo.test", "description", "Go shopping for avocados"),
					resource.TestCheckResourceAttr("todo_todo.test", "completed", "true"),
					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet("todo_todo.test", "id"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}
