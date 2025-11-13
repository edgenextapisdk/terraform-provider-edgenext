Use this data source to query base settings of a specific SCDN domain.

Example Usage

Query domain base settings

```hcl
data "edgenext_scdn_domain_base_settings" "example" {
  domain_id = 12345
}

output "proxy_host" {
  value = data.edgenext_scdn_domain_base_settings.example.proxy_host
}

output "domain_redirect" {
  value = data.edgenext_scdn_domain_base_settings.example.domain_redirect
}
```

Query and save to file

```hcl
data "edgenext_scdn_domain_base_settings" "example" {
  domain_id         = 12345
  result_output_file = "base_settings.json"
}
```

