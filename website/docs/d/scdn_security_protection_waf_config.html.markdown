---
subcategory: "Security CDN (SCDN)"
layout: "edgenext"
page_title: "EdgeNext: edgenext_scdn_security_protection_waf_config"
sidebar_current: "docs-edgenext-datasource-scdn_security_protection_waf_config"
description: |-
  Use this data source to query SCDN security protection WAF configuration.
---

# edgenext_scdn_security_protection_waf_config

Use this data source to query SCDN security protection WAF configuration.

## Example Usage

### Query WAF config

```hcl
data "edgenext_scdn_security_protection_waf_config" "example" {
  business_id = 12345
}

output "waf_config" {
  value = data.edgenext_scdn_security_protection_waf_config.example.waf_rule_config
}
```

### Query specific config keys

```hcl
data "edgenext_scdn_security_protection_waf_config" "example" {
  business_id = 12345
  keys        = ["waf_rule_config"]
}
```

### Query and save to file

```hcl
data "edgenext_scdn_security_protection_waf_config" "example" {
  business_id        = 12345
  result_output_file = "waf_config.json"
}
```

## Argument Reference

The following arguments are supported:

* `business_id` - (Required, Int) Business ID
* `keys` - (Optional, List: [`String`]) Specify config keys
* `result_output_file` - (Optional, String) Used to save results to a file

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `waf_intercept_page` - WAF intercept page configuration
  * `content` - Custom content
  * `id` - ID
  * `status` - Status: on, off
  * `type` - Page type: custom, default, keep
* `waf_rule_config` - WAF rule configuration
  * `ai_status` - AI status: on, off
  * `id` - ID
  * `status` - Status: on, off, keep
  * `waf_level` - Protection level: general, strict, keep
  * `waf_mode` - Protection mode: off, active, block, ban, keep


