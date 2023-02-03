# Configure the connection details for the Todo API server
provider "todo" {
  host    = "127.0.0.1"
  port    = 8080
  schema  = "http"
  apipath = "/"
}

# Create a new Todo item
resource "todo_todo" "test" {
  description = "Go Shopping"
  completed   = false
}

# Read in a existing Todo item
data "todo_todo" "example" {
  id = 1
}
