package main

import (
	"github.com/hashicorp/terraform/plugin"
	"github.com/terraform-providers/terraform-provider-template/typescript"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: typescript.Provider})
}
