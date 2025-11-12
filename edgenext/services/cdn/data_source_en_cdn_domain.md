Use this data source to query detailed information of CDN domain configuration.

Example Usage

Query CDN domain configuration by domain name

```hcl
data "edgenext_cdn_domain" "example" {
  domain = "example.com"
}
```

Query specific configuration items for a CDN domain

```hcl
data "edgenext_cdn_domain" "example" {
  domain = "example.com"
  config_item = [
    "origin",
    "https", 
    "cache_rule",
    "referer"
  ]
  output_file = "domain_config.json"
}
```
