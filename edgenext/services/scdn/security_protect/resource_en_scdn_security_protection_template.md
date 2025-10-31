Provides a resource to create and manage SCDN security protection templates.

Example Usage

Create security protection template

```hcl
resource "edgenext_scdn_security_protection_template" "example" {
  name   = "my-security-template"
  remark = "My security protection template"
}
```

Create template with domain binding

```hcl
resource "edgenext_scdn_security_protection_template" "example" {
  name      = "my-security-template"
  domain_ids = [12345, 67890]
}
```

Create template from existing template

```hcl
resource "edgenext_scdn_security_protection_template" "example" {
  name             = "my-new-template"
  template_source_id = 11111
}
```

Import

SCDN security protection templates can be imported using the template ID:

```shell
terraform import edgenext_scdn_security_protection_template.example 12345
```

