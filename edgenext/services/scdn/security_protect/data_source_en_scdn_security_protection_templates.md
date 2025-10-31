Use this data source to query a list of SCDN security protection templates.

Example Usage

Query all templates

```hcl
data "edgenext_scdn_security_protection_templates" "example" {
  tpl_type = "global"
  page     = 1
  page_size = 20
}

output "template_count" {
  value = data.edgenext_scdn_security_protection_templates.example.total
}

output "templates" {
  value = data.edgenext_scdn_security_protection_templates.example.list
}
```

Query with filters

```hcl
data "edgenext_scdn_security_protection_templates" "example" {
  tpl_type   = "only_domain"
  search_key = "my-template"
  page       = 1
  page_size  = 20
}
```

Query and save to file

```hcl
data "edgenext_scdn_security_protection_templates" "example" {
  tpl_type          = "global"
  result_output_file = "templates.json"
}
```

