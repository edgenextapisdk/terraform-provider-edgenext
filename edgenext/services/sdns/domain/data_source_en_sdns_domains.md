Use this data source to query a list of SDNS domains.

Example Usage

Query SDNS domains

```hcl
data "edgenext_sdns_domains" "example" {
  domain = "example.com"
}
```

Attributes Reference

The following attributes are exported:

* `domains` - List of matched domains
  * `id` - The ID of the domain
  * `domain` - The domain name
  * `status` - Status of the domain
