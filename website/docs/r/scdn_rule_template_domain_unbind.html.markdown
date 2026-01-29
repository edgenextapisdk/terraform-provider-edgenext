---
subcategory: "Security CDN (SCDN)"
layout: "edgenext"
page_title: "EdgeNext: edgenext_scdn_rule_template_domain_unbind"
sidebar_current: "docs-edgenext-resource-scdn_rule_template_domain_unbind"
description: |-
  Provides a resource to unbind domains from an SCDN rule template.
---

# edgenext_scdn_rule_template_domain_unbind

Provides a resource to unbind domains from an SCDN rule template.

## Example Usage

### Unbind domains from template

```hcl
resource "edgenext_scdn_rule_template_domain_unbind" "example" {
  template_id = 12345
  domain_ids  = [67890, 11111]
}
```

## Argument Reference

The following arguments are supported:

* `domain_ids` - (Required, List: [`Int`], ForceNew) List of domain IDs to unbind from the template
* `template_id` - (Required, Int, ForceNew) The ID of the rule template

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The unique identifier for this unbind operation


## Import

SCDN rule template domain unbindings can be imported using the template ID and domain IDs:

```shell
terraform import edgenext_scdn_rule_template_domain_unbind.example 12345-67890,11111
```

