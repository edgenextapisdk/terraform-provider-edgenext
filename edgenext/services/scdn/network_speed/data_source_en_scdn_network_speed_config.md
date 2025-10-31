Use this data source to query SCDN network speed configuration.

Example Usage

Query network speed config for template

```hcl
data "edgenext_scdn_network_speed_config" "example" {
  business_id   = 12345
  business_type = "tpl"
}

output "domain_proxy_conf" {
  value = data.edgenext_scdn_network_speed_config.example.domain_proxy_conf
}
```

Query specific config groups

```hcl
data "edgenext_scdn_network_speed_config" "example" {
  business_id   = 12345
  business_type = "tpl"
  config_groups = ["domain_proxy_conf", "upstream_redirect"]
}
```

Query and save to file

```hcl
data "edgenext_scdn_network_speed_config" "example" {
  business_id        = 12345
  business_type      = "tpl"
  result_output_file = "network_speed_config.json"
}
```

