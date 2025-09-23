package main

import (
	"github.com/edgenextapisdk/terraform-provider-edgenext/edgenext"
	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
)

// Version information, set by GoReleaser during build
var (
	version = "dev"
	commit  = "none"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: edgenext.Provider,
	})
}
