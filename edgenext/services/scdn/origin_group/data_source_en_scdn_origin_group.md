Use this data source to query details of a specific SCDN origin group.

Example Usage

Query origin group by ID

```hcl
data "edgenext_scdn_origin_group" "example" {
  origin_group_id = 12345
}

output "origin_group_name" {
  value = data.edgenext_scdn_origin_group.example.name
}

output "origins" {
  value = data.edgenext_scdn_origin_group.example.origins
}
```

Query and save to file

```hcl
data "edgenext_scdn_origin_group" "example" {
  origin_group_id    = 12345
  result_output_file = "origin_group.json"
}
```

