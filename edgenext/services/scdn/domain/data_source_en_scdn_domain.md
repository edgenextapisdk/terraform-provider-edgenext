Use this data source to query details of a specific SCDN domain.

Example Usage

Query domain by name

```hcl
data "edgenext_scdn_domain" "example" {
  domain = "example.com"
}

output "domain_id" {
  value = data.edgenext_scdn_domain.example.id
}

output "domain_status" {
  value = data.edgenext_scdn_domain.example.access_progress
}
```

Query domain and save to file

```hcl
data "edgenext_scdn_domain" "example" {
  domain             = "example.com"
  result_output_file = "domain.json"
}
```

