Provides a resource to unbind domains from an SCDN rule template.

Example Usage

Unbind domains from template

```hcl
resource "edgenext_scdn_rule_template_domain_unbind" "example" {
  template_id = 12345
  domain_ids  = [67890, 11111]
}
```

Import

SCDN rule template domain unbindings can be imported using the template ID and domain IDs:

```shell
terraform import edgenext_scdn_rule_template_domain_unbind.example 12345-67890,11111
```

