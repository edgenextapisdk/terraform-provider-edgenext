---
subcategory: "Security CDN (SCDN)"
layout: "edgenext"
page_title: "EdgeNext: edgenext_scdn_rule_template_switch"
sidebar_current: "docs-edgenext-resource-scdn_rule_template_switch"
description: |-
  Provides a resource to switch domains to a different SCDN Rule Template.
---

# edgenext_scdn_rule_template_switch

Provides a resource to switch domains to a different SCDN Rule Template.

## Example Usage

### Switch domains to a new template

```hcl
resource "edgenext_scdn_rule_template_switch" "example" {
  app_type     = "network_speed"
  new_tpl_id   = 12345
  new_tpl_type = "more_domain"
  domain_ids   = [1001, 1002]
}
```

## Argument Reference

The following arguments are supported:

* `domain_ids` - (Required, List: [`Int`], ForceNew) List of domain IDs to switch templates
* `new_tpl_id` - (Required, Int, ForceNew) New template ID to switch to, when new_tpl_type=global, pass 0
* `new_tpl_type` - (Required, String, ForceNew) New template type
* `app_type` - (Optional, String, ForceNew) Application type the template applies to

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The unique identifier for this switch operation


## Import

Rule Template Switch instances can be imported using an ID with the format `template_id:domain_id1,domain_id2,...`:

```shell
terraform import edgenext_scdn_rule_template_switch.example 12345:1001,1002
```

