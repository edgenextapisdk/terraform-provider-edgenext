---
subcategory: "Security CDN (SCDN)"
layout: "edgenext"
page_title: "EdgeNext: edgenext_scdn_security_protection_template"
sidebar_current: "docs-edgenext-resource-scdn_security_protection_template"
description: |-
  Provides a resource to create and manage SCDN security protection templates.
---

# edgenext_scdn_security_protection_template

Provides a resource to create and manage SCDN security protection templates.

## Example Usage

### Create security protection template

```hcl
resource "edgenext_scdn_security_protection_template" "example" {
  name   = "my-security-template"
  remark = "My security protection template"
}
```

### Create template with domain binding

```hcl
resource "edgenext_scdn_security_protection_template" "example" {
  name       = "my-security-template"
  domain_ids = [12345, 67890]
}
```

### Create template from existing template

```hcl
resource "edgenext_scdn_security_protection_template" "example" {
  name               = "my-new-template"
  template_source_id = 11111
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required, String) Template name
* `bind_all` - (Optional, Bool) Bind all domains
* `business_id` - (Optional, Int) Business ID (template ID). Required for update/delete, computed for create.
* `domain_ids` - (Optional, List: [`Int`]) Domain ID list
* `domains` - (Optional, List: [`String`]) Domain list
* `group_ids` - (Optional, List: [`Int`]) Group ID list
* `remark` - (Optional, String) Template remark
* `template_source_id` - (Optional, Int) Source template ID

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `bind_domain_count` - Bind domain count
* `created_at` - Creation time
* `fail_domains` - Failed domains
* `type` - Template type: global, only_domain, more_domain


## Import

SCDN security protection templates can be imported using the template ID:

```shell
terraform import edgenext_scdn_security_protection_template.example 12345
```

