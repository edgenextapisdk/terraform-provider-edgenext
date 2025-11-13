Use this data source to query domains bound to a specific SCDN rule template.

Example Usage

Query template domains

```hcl
data "edgenext_scdn_rule_template_domains" "example" {
  id       = 12345
  app_type = "network_speed"
  page     = 1
  page_size = 100
}

output "domain_count" {
  value = data.edgenext_scdn_rule_template_domains.example.total
}

output "domains" {
  value = data.edgenext_scdn_rule_template_domains.example.list
}
```

Query with domain filter

```hcl
data "edgenext_scdn_rule_template_domains" "example" {
  id       = 12345
  app_type = "network_speed"
  domain   = "example.com"
}

output "filtered_domains" {
  value = data.edgenext_scdn_rule_template_domains.example.list
}
```

Query and save to file

```hcl
data "edgenext_scdn_rule_template_domains" "example" {
  id                = 12345
  app_type          = "network_speed"
  result_output_file = "template_domains.json"
}
```

