Use this data source to query templates bound to a specific SCDN domain.

Example Usage

Query domain templates

```hcl
data "edgenext_scdn_domain_templates" "example" {
  domain_id = 12345
}

output "binded_templates" {
  value = data.edgenext_scdn_domain_templates.example.binded_templates
}
```

Query and save to file

```hcl
data "edgenext_scdn_domain_templates" "example" {
  domain_id         = 12345
  result_output_file = "templates.json"
}
```

