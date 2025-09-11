package main

import (
	"github.com/edgenextapisdk/terraform-provider-edgenext/edgenext"
	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
)

// 版本信息，由 GoReleaser 在构建时设置
var (
	version = "dev"
	commit  = "none"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: edgenext.Provider,
	})
}
