Provides a resource to create and manage SCDN log download templates.

Example Usage

Create log download template

```hcl
resource "edgenext_scdn_log_download_template" "example" {
  template_name = "my-template"
  group_name    = "my-group"
  group_id      = 1
  data_source   = "ng"
  download_fields = ["time", "domain", "url"]
}
```

Create template with search terms

```hcl
resource "edgenext_scdn_log_download_template" "example" {
  template_name = "my-template"
  group_name    = "my-group"
  group_id      = 1
  data_source   = "ng"
  download_fields = ["time", "domain", "url"]
  search_terms = [
    {
      key   = "domain"
      value = "example.com"
    }
  ]
}
```

Import

SCDN log download templates can be imported using the template ID:

```shell
terraform import edgenext_scdn_log_download_template.example 12345
```

