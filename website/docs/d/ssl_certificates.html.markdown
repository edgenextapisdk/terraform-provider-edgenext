---
subcategory: "SSL Certificate Management(SSL)"
layout: "edgenext"
page_title: "EdgeNext: edgenext_ssl_certificates"
sidebar_current: "docs-edgenext-datasource-ssl_certificates"
description: |-
  Use this data source to query a list of SSL certificates.
---

# edgenext_ssl_certificates

Use this data source to query a list of SSL certificates.

## Example Usage

### Query all SSL certificates

```hcl
data "edgenext_ssl_certificates" "example" {
  page_number        = 1
  page_size          = 100
  result_output_file = "ssl_certs.json"
}
```

### Query SSL certificates with pagination

```hcl
data "edgenext_ssl_certificates" "example" {
  page_number = 1
  page_size   = 50
}
```

## Argument Reference

The following arguments are supported:

* `page_number` - (Optional, Int) Page number, must be greater than 0 if specified
* `page_size` - (Optional, Int) Number of items per page, range 1-500 if specified
* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `list` - SSL certificate list
  * `associated_domains` - Associated domain list
  * `cert_expire_time` - Certificate end time
  * `cert_id` - Certificate ID
  * `cert_start_time` - Certificate start time
  * `include_domains` - Included domain list
  * `name` - Certificate name


