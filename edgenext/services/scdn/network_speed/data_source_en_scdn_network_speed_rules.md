Use this data source to query SCDN network speed rules.

Example Usage

Query custom page rules

```hcl
data "edgenext_scdn_network_speed_rules" "example" {
  business_id   = 12345
  business_type = "tpl"
  config_group  = "custom_page"
}

output "rule_count" {
  value = data.edgenext_scdn_network_speed_rules.example.total
}

output "rules" {
  value = data.edgenext_scdn_network_speed_rules.example.list
}
```

Query upstream URI change rules

```hcl
data "edgenext_scdn_network_speed_rules" "example" {
  business_id   = 12345
  business_type = "tpl"
  config_group  = "upstream_uri_change_rule"
}
```

Query and save to file

```hcl
data "edgenext_scdn_network_speed_rules" "example" {
  business_id        = 12345
  business_type      = "tpl"
  config_group       = "custom_page"
  result_output_file = "network_speed_rules.json"
}
```

