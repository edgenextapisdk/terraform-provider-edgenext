Use this data source to query all SCDN origin groups by protection status.

Example Usage

Query all origin groups for SCDN nodes

```hcl
data "edgenext_scdn_origin_groups_all" "example" {
  protect_status = "scdn"
}

output "origin_group_count" {
  value = data.edgenext_scdn_origin_groups_all.example.total
}

output "origin_groups" {
  value = data.edgenext_scdn_origin_groups_all.example.origin_groups
}
```

Query all origin groups for exclusive nodes

```hcl
data "edgenext_scdn_origin_groups_all" "example" {
  protect_status = "exclusive"
}
```

Query and save to file

```hcl
data "edgenext_scdn_origin_groups_all" "example" {
  protect_status     = "scdn"
  result_output_file = "origin_groups_all.json"
}
```

