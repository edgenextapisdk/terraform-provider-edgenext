Use this data source to query SCDN cache global configuration.

Example Usage

Query cache global config

```hcl
data "edgenext_scdn_cache_global_config" "example" {
}

output "cache_config" {
  value = data.edgenext_scdn_cache_global_config.example.conf
}
```

Query and save to file

```hcl
data "edgenext_scdn_cache_global_config" "example" {
  result_output_file = "cache_global_config.json"
}
```

