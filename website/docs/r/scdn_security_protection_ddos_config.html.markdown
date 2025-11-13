---
subcategory: "Security CDN (SCDN)"
layout: "edgenext"
page_title: "EdgeNext: edgenext_scdn_security_protection_ddos_config"
sidebar_current: "docs-edgenext-resource-scdn_security_protection_ddos_config"
description: |-
  Provides a resource to manage SCDN security protection DDoS configuration.
---

# edgenext_scdn_security_protection_ddos_config

Provides a resource to manage SCDN security protection DDoS configuration.

## Example Usage

### Configure DDoS protection

```hcl
resource "edgenext_scdn_security_protection_ddos_config" "example" {
  business_id = 12345

  application_ddos_protection {
    status                = "on"
    ai_cc_status          = "on"
    type                  = "normal"
    need_attack_detection = 1
    ai_status             = "on"
  }
}
```

## Argument Reference

The following arguments are supported:

* `business_id` - (Required, Int, ForceNew) Business ID
* `application_ddos_protection` - (Optional, List) Application layer DDoS protection configuration
* `visitor_authentication` - (Optional, List) Visitor authentication configuration

The `application_ddos_protection` object supports the following:

* `ai_cc_status` - (Optional, String) AI protection status: on, off
* `ai_status` - (Optional, String) AI status: on, off
* `need_attack_detection` - (Optional, Int) Attack detection switch: 0 or 1
* `status` - (Optional, String) Status: on, off, keep
* `type` - (Optional, String) Protection type: default, normal, strict, captcha, keep

The `visitor_authentication` object supports the following:

* `auth_token` - (Optional, String) Authentication token
* `pass_still_check` - (Optional, Int) Pass still check: 0 or 1
* `status` - (Optional, String) Status: on, off

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

SCDN security protection DDoS configuration can be imported using the business ID:

```shell
terraform import edgenext_scdn_security_protection_ddos_config.example 12345
```

