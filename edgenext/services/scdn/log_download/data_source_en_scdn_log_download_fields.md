Use this data source to query available SCDN log download fields.

Example Usage

Query all log download fields

```hcl
data "edgenext_scdn_log_download_fields" "all" {
}

output "field_configs" {
  value = data.edgenext_scdn_log_download_fields.all.configs
}
```

Query fields for specific data source

```hcl
data "edgenext_scdn_log_download_fields" "ng" {
  data_source = "ng"
}

output "ng_fields" {
  value = data.edgenext_scdn_log_download_fields.ng.configs
}
```

Query and save to file

```hcl
data "edgenext_scdn_log_download_fields" "all" {
  result_output_file = "log_download_fields.json"
}
```

