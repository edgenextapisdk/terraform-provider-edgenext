---
subcategory: "Security CDN (SCDN)"
layout: "edgenext"
page_title: "EdgeNext: edgenext_scdn_security_protection_template_batch_config"
sidebar_current: "docs-edgenext-resource-scdn_security_protection_template_batch_config"
description: |-
  Provides a resource to batch configure SCDN security protection templates.
---

# edgenext_scdn_security_protection_template_batch_config

Provides a resource to batch configure SCDN security protection templates.

## Example Usage

### Batch configure templates with precise access control

```hcl
resource "edgenext_scdn_security_protection_template_batch_config" "example" {
  template_ids = [12345, 67890]
  domains      = ["example.com"]

  ddos_config {
    application_ddos_protection {
      status                = "on"
      ai_cc_status          = "on"
      type                  = "strict"
      need_attack_detection = 1
      ai_status             = "on"
    }
  }

  waf_rule_config {
    waf_rule_config {
      status    = "on"
      ai_status = "on"
      waf_level = "strict"
      waf_mode  = "block"
    }
  }

  precise_access_control_config {
    action = "add"

    # Anti-CC protection for specific URL
    policies {
      action = "anticc"
      action_data = {
        level = "default"
      }
      type   = "plus"
      status = 1
      sort   = 1

      rules {
        rule_type = "url"
        logic     = "contains"
        data      = jsonencode(["/aaa"])
      }
    }

    # Deny access from specific referer
    policies {
      action = "deny"
      type   = "plus"
      from   = "aR"
      status = 1

      rules {
        rule_type = "referer_domain"
        logic     = "not_equals"
        data      = jsonencode(["home.example.com"])
      }

      rules {
        rule_type = "referer"
        logic     = "len_greater_than"
        data      = jsonencode({ len = 0 })
      }

      rules {
        rule_type = "postfix"
        logic     = "equals"
        data      = jsonencode(["css", "js", "png", "jpg"])
      }
    }

    # Region-based access control
    policies {
      action = "deny"
      type   = "plus"
      from   = "zL"
      status = 1

      rules {
        rule_type = "region"
        logic     = "not_belongs"
        data      = jsonencode({ province = [], country = ["CN", "JP"] })
      }
    }
  }
}
```

### Batch configure for all domains

```hcl
resource "edgenext_scdn_security_protection_template_batch_config" "example" {
  template_ids = [12345]
  all          = 1

  ddos_config {
    application_ddos_protection {
      status = "on"
    }
  }
}
```

## Argument Reference

The following arguments are supported:

* `template_ids` - (Required, List: [`Int`]) Template ID list
* `all` - (Optional, Int) All flag (0 or 1)
* `bot_management_config` - (Optional, List) Bot management configuration
* `ddos_config` - (Optional, List) DDoS protection configuration
* `domain_ids` - (Optional, List: [`Int`]) Domain ID list
* `domains` - (Optional, List: [`String`]) Domain list
* `precise_access_control_config` - (Optional, List) Precise access control configuration
* `waf_rule_config` - (Optional, List) WAF rule configuration

The `application_ddos_protection` object of `ddos_config` supports the following:

* `ai_cc_status` - (Optional, String) AI protection status: on, off
* `ai_status` - (Optional, String) AI status: on, off
* `need_attack_detection` - (Optional, Int) Attack detection switch: 0 or 1
* `status` - (Optional, String) Status: on, off, keep
* `type` - (Optional, String) Protection type: default, normal, strict, captcha, keep

The `bot_management_config` object supports the following:

* `business_id` - (Optional, Int) Business ID
* `ids` - (Optional, List) ID list

The `ddos_config` object supports the following:

* `application_ddos_protection` - (Optional, List) Application layer DDoS protection configuration
* `visitor_authentication` - (Optional, List) Visitor authentication configuration

The `policies` object of `precise_access_control_config` supports the following:

* `action_data` - (Optional, Map) Action data
* `action` - (Optional, String) Policy action
* `from` - (Optional, String) From source
* `id` - (Optional, Int) Policy ID
* `remark` - (Optional, String) Policy remark
* `rule_type` - (Optional, String) Rule type (was 'type')
* `rules` - (Optional, List) Rules list
* `sort` - (Optional, Int) Sort order
* `status` - (Optional, Int) Status
* `type` - (Optional, String) Policy type: plus

The `precise_access_control_config` object supports the following:

* `action` - (Required, String) Action: add, cover
* `policies` - (Optional, List) Policy list

The `rules` object of `policies` supports the following:

* `data` - (Required, String) Rule data (JSON string for array/object, or plain string)
* `logic` - (Required, String) Logic: contains, equals, not_equals, not_belongs, len_greater_than
* `rule_type` - (Required, String) Rule type: url, referer_domain, referer, postfix, region

The `visitor_authentication` object of `ddos_config` supports the following:

* `auth_token` - (Optional, String) Authentication token
* `pass_still_check` - (Optional, Int) Pass still check: 0 or 1
* `status` - (Optional, String) Status: on, off

The `waf_intercept_page` object of `waf_rule_config` supports the following:

* `content` - (Optional, String) Custom content
* `status` - (Optional, String) Status: on, off
* `type` - (Optional, String) Type: custom, default, keep

The `waf_rule_config` object of `waf_rule_config` supports the following:

* `ai_status` - (Optional, String) AI status: on, off
* `status` - (Optional, String) Status: on, off, keep
* `waf_level` - (Optional, String) Protection level: general, strict, keep
* `waf_mode` - (Optional, String) Protection mode: off, active, block, ban, keep
* `waf_strategy_id` - (Optional, Int) WAF strategy ID

The `waf_rule_config` object supports the following:

* `waf_intercept_page` - (Optional, List) WAF intercept page config
* `waf_rule_config` - (Optional, List) WAF rule config

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `fail_templates` - (**Deprecated**) This attribute is deprecated and will be removed in a future version. Please check the apply output or logs for failure details. Failed templates (Deprecated)


## Import

SCDN security protection template batch configuration can be imported using the template IDs:

```hcl
terraform import edgenext_scdn_security_protection_template_batch_config.example 12345,67890

