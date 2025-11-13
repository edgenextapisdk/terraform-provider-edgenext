Use this data source to query details of a specific SCDN rule template.

Example Usage

Query rule template by ID

```hcl
data "edgenext_scdn_rule_template" "example" {
  id      = "12345"
  app_type = "network_speed"
}

output "template_name" {
  value = data.edgenext_scdn_rule_template.example.name
}

output "bind_domains" {
  value = data.edgenext_scdn_rule_template.example.bind_domains
}
```

Query and save to file

```hcl
data "edgenext_scdn_rule_template" "example" {
  id               = "12345"
  app_type         = "network_speed"
  result_output_file = "template.json"
}
```

