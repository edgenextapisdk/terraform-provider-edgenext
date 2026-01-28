# edgenext_dns_record

Provides a DNS record resource.

## Example Usage

```hcl
resource "edgenext_dns_record" "www" {
  domain_id = 123
  name      = "www"
  type      = "A"
  view      = "any"
  value     = "1.2.3.4"
  ttl       = 600
}
```

## Argument Reference

The following arguments are supported:

* `domain_id` - (Required, ForceNew) The ID of the domain.
* `name` - (Required) The name of the record.
* `type` - (Required) The type of the record (A, CNAME, etc.).
* `view` - (Required) The view/line for the record.
* `value` - (Required) The value of the record.
* `mx` - (Optional) MX priority. Default is 0.
* `ttl` - (Optional) TTL in seconds. Default is 600.
* `remark` - (Optional) Remark for the record.
* `status` - (Optional) Status of the record (1 for enabled, 2 for paused). Default is 1.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of the record.
