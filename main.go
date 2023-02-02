package main

import (
	"context"

	"github.com/spkane/terraform-provider-todo/todo"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
)

// Provider documentation generation.
//go:generate go run github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs generate --provider-name todo

func main() {
	providerserver.Serve(context.Background(), todo.New, providerserver.ServeOpts{
		Address: "registry.terraform.io/hashicorp/todo",
	})
}
