Use this data source to query a list of SCDN domains with optional filters.

Example Usage

Query all domains

```hcl
data "edgenext_scdn_domains" "all" {
  page      = 1
  page_size = 100
}

output "domain_count" {
  value = data.edgenext_scdn_domains.all.total
}

output "domain_names" {
  value = [for domain in data.edgenext_scdn_domains.all.domains : domain.domain]
}
```

Query domains with filters

```hcl
data "edgenext_scdn_domains" "filtered" {
  domain          = "example"
  access_progress = "online"
  protect_status  = "scdn"
  page            = 1
  page_size       = 50
}

output "filtered_domains" {
  value = data.edgenext_scdn_domains.filtered.domains
}
```

Query domains and save to file

```hcl
data "edgenext_scdn_domains" "all" {
  page                = 1
  page_size           = 100
  result_output_file = "domains.json"
}
```

