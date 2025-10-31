package networkspeed

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Resources returns all network speed-related resources
func Resources() map[string]*schema.Resource {
	return map[string]*schema.Resource{
		"edgenext_scdn_network_speed_config":     ResourceEdgenextScdnNetworkSpeedConfig(),
		"edgenext_scdn_network_speed_rule":       ResourceEdgenextScdnNetworkSpeedRule(),
		"edgenext_scdn_network_speed_rules_sort": ResourceEdgenextScdnNetworkSpeedRulesSort(),
	}
}

// DataSources returns all network speed-related data sources
func DataSources() map[string]*schema.Resource {
	return map[string]*schema.Resource{
		"edgenext_scdn_network_speed_config": DataSourceEdgenextScdnNetworkSpeedConfig(),
		"edgenext_scdn_network_speed_rules":  DataSourceEdgenextScdnNetworkSpeedRules(),
	}
}
