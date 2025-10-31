Use this data source to query details of a specific SCDN security protection template.

Example Usage

Query security protection template

```hcl
data "edgenext_scdn_security_protection_template" "example" {
  business_id = 12345
}

output "template_name" {
  value = data.edgenext_scdn_security_protection_template.example.name
}

output "bind_domain_count" {
  value = data.edgenext_scdn_security_protection_template.example.bind_domain_count
}
```

Query and save to file

```hcl
data "edgenext_scdn_security_protection_template" "example" {
  business_id        = 12345
  result_output_file = "template.json"
}
```

