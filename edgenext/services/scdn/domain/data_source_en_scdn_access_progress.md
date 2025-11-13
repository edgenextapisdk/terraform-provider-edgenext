Use this data source to query available SCDN access progress status options.

Example Usage

Query access progress options

```hcl
data "edgenext_scdn_access_progress" "example" {
}

output "progress_options" {
  value = data.edgenext_scdn_access_progress.example.progress
}
```

Query and save to file

```hcl
data "edgenext_scdn_access_progress" "example" {
  result_output_file = "access_progress.json"
}
```

