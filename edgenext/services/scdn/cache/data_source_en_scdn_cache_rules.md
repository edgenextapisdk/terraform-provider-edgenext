Use this data source to query SCDN cache rules.

Example Usage

Query cache rules for template

```hcl
data "edgenext_scdn_cache_rules" "example" {
  business_id   = 12345
  business_type = "tpl"
}

output "rule_count" {
  value = data.edgenext_scdn_cache_rules.example.total
}

output "rules" {
  value = data.edgenext_scdn_cache_rules.example.list
}
```

Query cache rules for domain

```hcl
data "edgenext_scdn_cache_rules" "example" {
  business_id   = 67890
  business_type = "domain"
  page          = 1
  page_size     = 50
}
```

Query and save to file

```hcl
data "edgenext_scdn_cache_rules" "example" {
  business_id        = 12345
  business_type      = "tpl"
  result_output_file = "cache_rules.json"
}
```

