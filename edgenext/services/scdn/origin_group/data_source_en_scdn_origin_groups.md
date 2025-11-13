Use this data source to query a list of SCDN origin groups with optional filters.

Example Usage

Query all origin groups

```hcl
data "edgenext_scdn_origin_groups" "all" {
  page      = 1
  page_size = 20
}

output "origin_group_count" {
  value = data.edgenext_scdn_origin_groups.all.total
}

output "origin_groups" {
  value = data.edgenext_scdn_origin_groups.all.origin_groups
}
```

Query with name filter

```hcl
data "edgenext_scdn_origin_groups" "filtered" {
  name      = "my-group"
  page      = 1
  page_size = 20
}
```

Query and save to file

```hcl
data "edgenext_scdn_origin_groups" "all" {
  page                = 1
  page_size           = 20
  result_output_file  = "origin_groups.json"
}
```

