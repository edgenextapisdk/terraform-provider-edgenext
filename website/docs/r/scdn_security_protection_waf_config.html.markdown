---
subcategory: "Security CDN (SCDN)"
layout: "edgenext"
page_title: "EdgeNext: edgenext_scdn_security_protection_waf_config"
sidebar_current: "docs-edgenext-resource-scdn_security_protection_waf_config"
description: |-
  Provides a resource to manage SCDN security protection WAF configuration.
---

# edgenext_scdn_security_protection_waf_config

Provides a resource to manage SCDN security protection WAF configuration.

## Example Usage

### Configure WAF protection

```hcl
resource "edgenext_scdn_security_protection_waf_config" "example" {
  business_id = 12345

  waf_rule_config {
    status    = "on"
    ai_status = "on"
    waf_level = "strict"
    waf_mode  = "block"
  }
}
```

## Argument Reference

The following arguments are supported:

* `business_id` - (Required, Int, ForceNew) Business ID
* `waf_intercept_page` - (Optional, List) WAF intercept page configuration
* `waf_rule_config` - (Optional, List) WAF rule configuration

The `waf_intercept_page` object supports the following:

* `content` - (Optional, String) Custom content
* `status` - (Optional, String) Status: on, off
* `type` - (Optional, String) Page type: custom, default, keep

The `waf_rule_config` object supports the following:

* `ai_status` - (Optional, String) AI status: on, off
* `status` - (Optional, String) Status: on, off, keep
* `waf_level` - (Optional, String) Protection level: general, strict, keep
* `waf_mode` - (Optional, String) Protection mode: off, active, block, ban, keep

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

SCDN security protection WAF configuration can be imported using the business ID:

```shell
terraform import edgenext_scdn_security_protection_waf_config.example 12345
```

