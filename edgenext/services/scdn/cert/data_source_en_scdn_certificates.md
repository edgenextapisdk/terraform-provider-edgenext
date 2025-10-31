Use this data source to query a list of SCDN certificates with optional filters.

Example Usage

Query all certificates

```hcl
data "edgenext_scdn_certificates" "all" {
  page     = 1
  per_page = 20
}

output "certificate_count" {
  value = data.edgenext_scdn_certificates.all.total
}

output "certificates" {
  value = data.edgenext_scdn_certificates.all.list
}
```

Query certificates with filters

```hcl
data "edgenext_scdn_certificates" "filtered" {
  page        = 1
  per_page    = 20
  domain      = "example.com"
  binded      = "true"
  apply_status = "2"
}

output "filtered_certificates" {
  value = data.edgenext_scdn_certificates.filtered.list
}
```

Query and save to file

```hcl
data "edgenext_scdn_certificates" "all" {
  page                = 1
  per_page            = 20
  result_output_file  = "certificates.json"
}
```

