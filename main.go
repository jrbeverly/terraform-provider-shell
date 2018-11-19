package main

import (
	"github.com/Brightspace/terraform-provider-shell/shell"
	"github.com/hashicorp/terraform/plugin"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: shell.Provider,
	})
}
