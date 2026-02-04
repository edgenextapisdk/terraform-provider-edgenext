---
subcategory: "Security CDN (SCDN)"
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

* `business_id` - (Required, Int, ForceNew) Business ID (template ID) to bind domains to.
* `bind_business_ids` - (Optional, List: [`Int`]) List of business IDs to bind.
* `domain_ids` - (Optional, List: [`Int`]) List of domain IDs to bind to the template.
* `group_ids` - (Optional, List: [`Int`]) Group ID list. If both group_ids and domain_ids are provided, the intersection of domains from the groups and the domain IDs will be used for binding.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `fail_domains` - Failed domains


## Import

Import is not supported for this resource because:
1. There is no unique identifier for a specific bind relationship
2. The API does not provide enough information to reconstruct all resource attributes
3. Bind relationships are many-to-many between templates and domains
4. The resource ID format does not represent a specific binding

When attempting to import, you will receive an error with the above explanation.

