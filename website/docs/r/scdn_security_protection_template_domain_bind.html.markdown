---
subcategory: "Secure Content Delivery Network (SCDN)"
layout: "edgenext"
page_title: "EdgeNext: edgenext_scdn_security_protection_template_domain_bind"
sidebar_current: "docs-edgenext-resource-scdn_security_protection_template_domain_bind"
description: |-
  Provides a resource to bind domains to an SCDN security protection template.
---

# edgenext_scdn_security_protection_template_domain_bind

Provides a resource to bind domains to an SCDN security protection template.

## Example Usage

### Bind domains to template

```hcl
resource "edgenext_scdn_security_protection_template_domain_bind" "example" {
  business_id = 12345
  domain_ids  = [67890, 11111]
}
```

### Bind by group IDs

```hcl
resource "edgenext_scdn_security_protection_template_domain_bind" "example" {
  business_id = 12345
  group_ids   = [1, 2]
}
```

## Argument Reference

The following arguments are supported:

* `business_id` - (Required, Int, ForceNew) Business ID (template ID)
* `bind_business_ids` - (Optional, List: [`Int`], ForceNew) Bind business ID list
* `domain_ids` - (Optional, List: [`Int`], ForceNew) Domain ID list
* `group_ids` - (Optional, List: [`Int`], ForceNew) Group ID list

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `fail_domains` - Failed domains


## Import

SCDN security protection template domain bindings can be imported using the template ID and domain IDs:

```shell
terraform import edgenext_scdn_security_protection_template_domain_bind.example 12345-67890,11111
```

