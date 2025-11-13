---
subcategory: "Security CDN (SCDN)"
layout: "edgenext"
page_title: "EdgeNext: edgenext_scdn_security_protection_ddos_config"
sidebar_current: "docs-edgenext-datasource-scdn_security_protection_ddos_config"
description: |-
  Use this data source to query SCDN security protection DDoS configuration.
---

# edgenext_scdn_security_protection_ddos_config

Use this data source to query SCDN security protection DDoS configuration.

## Example Usage

### Query DDoS config

```hcl
data "edgenext_scdn_security_protection_ddos_config" "example" {
  business_id = 12345
}

output "ddos_config" {
  value = data.edgenext_scdn_security_protection_ddos_config.example.application_ddos_protection
}
```

### Query specific config keys

```hcl
data "edgenext_scdn_security_protection_ddos_config" "example" {
  business_id = 12345
  keys        = ["application_ddos_protection", "network_ddos_protection"]
}
```

### Query and save to file

```hcl
data "edgenext_scdn_security_protection_ddos_config" "example" {
  business_id        = 12345
  result_output_file = "ddos_config.json"
}
```

## Argument Reference

The following arguments are supported:

* `business_id` - (Required, Int) Business ID
* `keys` - (Optional, List: [`String`]) Specify config keys
* `result_output_file` - (Optional, String) Used to save results to a file

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `application_ddos_protection` - Application layer DDoS protection configuration
  * `ai_cc_status` - AI protection status: on, off
  * `ai_status` - AI status: on, off
  * `id` - ID
  * `need_attack_detection` - Attack detection switch: 0 or 1
  * `status` - Status: on, off, keep
  * `type` - Protection type: default, normal, strict, captcha, keep
* `visitor_authentication` - Visitor authentication configuration
  * `auth_token` - Authentication token
  * `id` - ID
  * `pass_still_check` - Pass still check: 0 or 1
  * `status` - Status: on, off


