Provides a resource to bind domains to an SCDN rule template.

Example Usage

Bind domains to template

```hcl
resource "edgenext_scdn_rule_template_domain_bind" "example" {
  template_id = 12345
  domain_ids  = [67890, 11111]
}
```

Import

SCDN rule template domain bindings can be imported using the template ID and domain IDs:

```shell
terraform import edgenext_scdn_rule_template_domain_bind.example 12345-67890,11111
```

