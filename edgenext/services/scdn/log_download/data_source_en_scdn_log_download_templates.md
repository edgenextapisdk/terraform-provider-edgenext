Use this data source to query SCDN log download templates.

Example Usage

Query all log download templates

```hcl
data "edgenext_scdn_log_download_templates" "all" {
  page     = 1
  per_page = 20
}

output "template_count" {
  value = data.edgenext_scdn_log_download_templates.all.total
}

output "templates" {
  value = data.edgenext_scdn_log_download_templates.all.templates
}
```

Query with filters

```hcl
data "edgenext_scdn_log_download_templates" "filtered" {
  page     = 1
  per_page = 20
  search_terms = [
    {
      key   = "template_name"
      value = "my-template"
    }
  ]
}
```

Query and save to file

```hcl
data "edgenext_scdn_log_download_templates" "all" {
  page                = 1
  per_page            = 20
  result_output_file  = "log_download_templates.json"
}
```

