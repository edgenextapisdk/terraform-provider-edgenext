Use this data source to query domains not bound to any SCDN security protection template.

Example Usage

Query unbound domains

```hcl
data "edgenext_scdn_security_protection_template_unbound_domains" "example" {
  page      = 1
  page_size = 20
}

output "unbound_domain_count" {
  value = data.edgenext_scdn_security_protection_template_unbound_domains.example.total
}

output "unbound_domains" {
  value = data.edgenext_scdn_security_protection_template_unbound_domains.example.list
}
```

Query with domain filter

```hcl
data "edgenext_scdn_security_protection_template_unbound_domains" "example" {
  domain = "example.com"
}
```

Query and save to file

```hcl
data "edgenext_scdn_security_protection_template_unbound_domains" "example" {
  result_output_file = "unbound_domains.json"
}
```

