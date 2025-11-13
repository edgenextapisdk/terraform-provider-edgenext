Use this data source to query details of a specific SCDN origin server.

Example Usage

Query origin by ID

```hcl
data "edgenext_scdn_origin" "example" {
  origin_id = 12345
  domain_id = 67890
}

output "origin_protocol" {
  value = data.edgenext_scdn_origin.example.protocol
}

output "origin_records" {
  value = data.edgenext_scdn_origin.example.records
}
```

Query origin and save to file

```hcl
data "edgenext_scdn_origin" "example" {
  origin_id         = 12345
  domain_id         = 67890
  result_output_file = "origin.json"
}
```

