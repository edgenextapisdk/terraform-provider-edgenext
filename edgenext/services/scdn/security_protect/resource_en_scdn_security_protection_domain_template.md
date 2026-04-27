Provides a resource to create a domain-level security protection template.

This resource creates a domain-level security protection template for a single domain. Unlike the multi-domain template (`edgenext_scdn_security_protection_template`), this resource:
- Only supports creation for a single domain, not update
- Deletion rebinds the domain to the global template
- Requires `template_source_id` (obtained from global template data source)

Example Usage

Create single domain template

```hcl
# First, get the global template ID
data "edgenext_scdn_security_protection_member_global_template" "global" {}

# Create a domain template for a single domain
resource "edgenext_scdn_security_protection_domain_template" "example" {
  domain_id          = 12345
  template_source_id = data.edgenext_scdn_security_protection_member_global_template.global.template[0].id
}
```

Complete example with domain data source

```hcl
# Get domain information
data "edgenext_scdn_domains" "example" {
  domain = "example.com"
}

# Get the global template ID
data "edgenext_scdn_security_protection_member_global_template" "global" {}

# Create a domain template
resource "edgenext_scdn_security_protection_domain_template" "example" {
  domain_id          = data.edgenext_scdn_domains.example.domains[0].id
  template_source_id = data.edgenext_scdn_security_protection_member_global_template.global.template[0].id
}
```

Argument Reference

The following arguments are supported:

* `domain_id` - (Required) Domain ID to create template for.
* `template_source_id` - (Required) Source template ID. Use data source `edgenext_scdn_security_protection_member_global_template` to get the global template ID.

Attributes Reference

The following attributes are exported:

* `business_id` - Created business ID (template ID) for the domain.

Import

Domain templates can be imported using the resource ID:

```hcl
terraform import edgenext_scdn_security_protection_domain_template.example domain-template-12345