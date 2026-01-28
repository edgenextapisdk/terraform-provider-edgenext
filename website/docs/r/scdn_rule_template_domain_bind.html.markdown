---
subcategory: "Security CDN (SCDN)"
layout: "edgenext"
page_title: "EdgeNext: edgenext_scdn_rule_template_domain_bind"
sidebar_current: "docs-edgenext-resource-scdn_rule_template_domain_bind"
description: |-
  Provides a resource to bind domains to an SCDN rule template.
---

# edgenext_scdn_rule_template_domain_bind

Provides a resource to bind domains to an SCDN rule template.

## Example Usage

### Bind domains to template

```hcl
resource "edgenext_scdn_rule_template_domain_bind" "example" {
  template_id = 12345
  domain_ids  = [67890, 11111]
}
```

## Argument Reference

The following arguments are supported:

* `domain_ids` - (Required, List: [`Int`], ForceNew) List of domain IDs to bind to the template
* `template_id` - (Required, Int, ForceNew) The ID of the rule template

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The unique identifier for this bind operation


## Import

SCDN rule template domain bindings can be imported using the template ID and domain IDs:

```shell
terraform import edgenext_scdn_rule_template_domain_bind.example 12345-67890,11111
```

