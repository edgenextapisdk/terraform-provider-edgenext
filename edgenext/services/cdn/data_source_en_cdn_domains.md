Use this data source to query a list of CDN domains.

Example Usage

Query all CDN domains

```hcl
data "edgenext_cdn_domains" "example" {
  page_number = 1
  page_size   = 100
}
```

Query CDN domains with specific status

```hcl
data "edgenext_cdn_domains" "example" {
  domain_status = "serving"
  page_number   = 1
  page_size     = 50
}
```
