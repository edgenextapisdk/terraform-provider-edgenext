Provides a resource to create and manage SCDN rule templates.

Example Usage

Create rule template

```hcl
resource "edgenext_scdn_rule_template" "example" {
  name        = "my-template"
  description = "My rule template"
  app_type    = "network_speed"
}
```

Create template with domain binding

```hcl
resource "edgenext_scdn_rule_template" "example" {
  name        = "my-template"
  app_type    = "network_speed"

  bind_domain {
    domain_ids = [12345, 67890]
  }
}
```

Create template from existing template

```hcl
resource "edgenext_scdn_rule_template" "example" {
  name         = "my-new-template"
  app_type     = "network_speed"
  from_tpl_id  = 11111
}
```

Import

SCDN rule templates can be imported using the template ID:

```shell
terraform import edgenext_scdn_rule_template.example 12345
```

