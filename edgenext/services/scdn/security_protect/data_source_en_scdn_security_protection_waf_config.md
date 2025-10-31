Use this data source to query SCDN security protection WAF configuration.

Example Usage

Query WAF config

```hcl
data "edgenext_scdn_security_protection_waf_config" "example" {
  business_id = 12345
}

output "waf_config" {
  value = data.edgenext_scdn_security_protection_waf_config.example.waf_rule_config
}
```

Query specific config keys

```hcl
data "edgenext_scdn_security_protection_waf_config" "example" {
  business_id = 12345
  keys        = ["waf_rule_config"]
}
```

Query and save to file

```hcl
data "edgenext_scdn_security_protection_waf_config" "example" {
  business_id        = 12345
  result_output_file = "waf_config.json"
}
```

