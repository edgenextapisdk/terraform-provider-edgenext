package main

import (
	"github.com/edgenextapisdk/terraform-provider-edgenext/edgenext"
	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: edgenext.Provider,
	})
}
