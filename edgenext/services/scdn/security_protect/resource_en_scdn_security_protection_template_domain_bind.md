Provides a resource to bind domains to an SCDN security protection template.

Example Usage

Bind domains to template

```hcl
resource "edgenext_scdn_security_protection_template_domain_bind" "example" {
  business_id = 12345
  domain_ids  = [67890, 11111]
}
```

Bind by group IDs

```hcl
resource "edgenext_scdn_security_protection_template_domain_bind" "example" {
  business_id = 12345
  group_ids   = [1, 2]
}
```

Import

Import is not supported for this resource because:
1. There is no unique identifier for a specific bind relationship
2. The API does not provide enough information to reconstruct all resource attributes
3. Bind relationships are many-to-many between templates and domains
4. The resource ID format does not represent a specific binding

When attempting to import, you will receive an error with the above explanation.

