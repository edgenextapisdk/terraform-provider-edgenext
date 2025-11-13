Use this data source to query SCDN cache clean configuration.

Example Usage

Query cache clean config

```hcl
data "edgenext_scdn_cache_clean_config" "example" {
}

output "wholesite" {
  value = data.edgenext_scdn_cache_clean_config.example.wholesite
}

output "specialurl" {
  value = data.edgenext_scdn_cache_clean_config.example.specialurl
}
```

Query and save to file

```hcl
data "edgenext_scdn_cache_clean_config" "example" {
  result_output_file = "cache_clean_config.json"
}
```

