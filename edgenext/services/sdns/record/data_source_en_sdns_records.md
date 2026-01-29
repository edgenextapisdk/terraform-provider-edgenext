Use this data source to query a list of SDNS DNS records.

Example Usage

Query SDNS DNS records

```hcl
data "edgenext_sdns_records" "example" {
  domain_id = 12345
}
```

Attributes Reference

The following attributes are exported:

* `records` - List of matched records
  * `id` - The ID of the record
  * `name` - The record name
  * `type` - The record type (A, CNAME, etc.)
  * `view` - The record view/line
  * `value` - The record value
  * `mx` - MX priority
  * `ttl` - TTL in seconds
  * `status` - Status of the record
  * `remark` - Remark for the record
