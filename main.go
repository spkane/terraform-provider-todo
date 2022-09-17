package main

import (
	// This uses protocol 5 via the Terraform Plugin SDKv2.
	// Protocol 6 is available by using the newer Terraform Plugin Framework.

	// Upstream Terraform Plugin Library
	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
	// Our local Terraform Provider code
	"github.com/spkane/todo-for-terraform/terraform-provider-todo/todo"
)

// Generate the Terraform provider documentation using `tfplugindocs`:
//go:generate go run github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs

// main is the entrypoint to the terraform plugin
func main() {
	// see: https://github.com/hashicorp/terraform/blob/master/plugin/serve.go
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: todo.Provider})
}
