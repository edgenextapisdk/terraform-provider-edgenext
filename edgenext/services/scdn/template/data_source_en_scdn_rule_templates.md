Use this data source to query a list of SCDN rule templates with optional filters.

Example Usage

Query all rule templates

```hcl
data "edgenext_scdn_rule_templates" "all" {
  page      = 1
  page_size = 100
}

output "template_count" {
  value = data.edgenext_scdn_rule_templates.all.total
}

output "templates" {
  value = data.edgenext_scdn_rule_templates.all.list
}
```

Query templates with filters

```hcl
data "edgenext_scdn_rule_templates" "filtered" {
  page      = 1
  page_size = 100
  name      = "my-template"
  app_type  = "network_speed"
}

output "filtered_templates" {
  value = data.edgenext_scdn_rule_templates.filtered.list
}
```

Query and save to file

```hcl
data "edgenext_scdn_rule_templates" "all" {
  page                = 1
  page_size           = 100
  result_output_file  = "templates.json"
}
```

