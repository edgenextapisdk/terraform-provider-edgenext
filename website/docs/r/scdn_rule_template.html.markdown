---
subcategory: "Security CDN (SCDN)"
layout: "edgenext"
page_title: "EdgeNext: edgenext_scdn_rule_template"
sidebar_current: "docs-edgenext-resource-scdn_rule_template"
description: |-
  Provides a resource to create and manage SCDN rule templates.
---

# edgenext_scdn_rule_template

Provides a resource to create and manage SCDN rule templates.

## Example Usage

### Create rule template

```hcl
resource "edgenext_scdn_rule_template" "example" {
  name        = "my-template"
  description = "My rule template"
  app_type    = "network_speed"
}
```

### Create template with domain binding

```hcl
resource "edgenext_scdn_rule_template" "example" {
  name     = "my-template"
  app_type = "network_speed"

  bind_domain {
    domain_ids = [12345, 67890]
  }
}
```

### Create template from existing template

```hcl
resource "edgenext_scdn_rule_template" "example" {
  name        = "my-new-template"
  app_type    = "network_speed"
  from_tpl_id = 11111
}
```

## Argument Reference

The following arguments are supported:

* `app_type` - (Required, String) The application type (e.g., 'network_speed')
* `name` - (Required, String) The rule template name
* `bind_domain` - (Optional, List) Domain binding information
* `description` - (Optional, String) The rule template description
* `from_tpl_id` - (Optional, Int) Existing template ID to copy from
* `template_id` - (Optional, String) The template ID for updating an existing template. If provided, this will update the template instead of creating a new one.
* `tpl_type` - (Optional, String) The template type

The `bind_domain` object supports the following:

* `all_domain` - (Optional, Bool) If true, bind to all domains
* `domain_group_ids` - (Optional, List) List of domain group IDs to bind
* `domain_ids` - (Optional, List) List of domain IDs to bind
* `domains` - (Optional, List) List of domain names to bind
* `is_bind` - (Optional, Bool) Whether to bind domains

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `created_at` - The template creation timestamp
* `id` - The ID of the rule template


## Import

SCDN rule templates can be imported using the template ID:

```shell
terraform import edgenext_scdn_rule_template.example 12345
```

