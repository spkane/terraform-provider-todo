terraform {
  required_providers {
    todo = {
      source  = "spkane/todo"
    }
  }
}

provider "todo" {
  host = "127.0.0.1"
  port = "8080"
  apipath = "/"
  schema = "http"
}

variable "purpose" {
    type = string
    description = "Tag the purpose of this todo"
    default = "test"
}

resource "todo_todo" "test1" {
  count = 5
  description = "${count.index}-1 ${var.purpose} todo"
  completed = false
}

output "todo_1_ids" {
  value = todo_todo.test1.*.id
}
