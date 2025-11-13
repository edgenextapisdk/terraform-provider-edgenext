Use this data source to query domains bound to a specific SCDN security protection template.

Example Usage

Query template domains

```hcl
data "edgenext_scdn_security_protection_template_domains" "example" {
  business_id = 12345
  page        = 1
  page_size   = 20
}

output "domain_count" {
  value = data.edgenext_scdn_security_protection_template_domains.example.total
}

output "domains" {
  value = data.edgenext_scdn_security_protection_template_domains.example.list
}
```

Query with domain filter

```hcl
data "edgenext_scdn_security_protection_template_domains" "example" {
  business_id = 12345
  domain      = "example.com"
}
```

Query and save to file

```hcl
data "edgenext_scdn_security_protection_template_domains" "example" {
  business_id        = 12345
  result_output_file = "template_domains.json"
}
```

