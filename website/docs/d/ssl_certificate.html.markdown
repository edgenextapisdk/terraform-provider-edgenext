---
subcategory: "SSL Certificate Management(SSL)"
layout: "edgenext"
page_title: "EdgeNext: edgenext_ssl_certificate"
sidebar_current: "docs-edgenext-datasource-ssl_certificate"
description: |-
  Use this data source to query detailed information of SSL certificate.
---

# edgenext_ssl_certificate

Use this data source to query detailed information of SSL certificate.

## Example Usage

### Query SSL certificate by certificate ID

```hcl
data "edgenext_ssl_certificate" "example" {
  cert_id            = "ssl-cert-123456"
  result_output_file = "ssl_cert.json"
}
```

### Query SSL certificate with minimal information

```hcl
data "edgenext_ssl_certificate" "example" {
  cert_id = "ssl-cert-123456"
}
```

## Argument Reference

The following arguments are supported:

* `cert_id` - (Required, String) Certificate ID
* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `bind_domains` - List of bound domains
* `cert_expire_time` - Certificate end time
* `cert_start_time` - Certificate start time
* `certificate` - Certificate content
* `key` - Private key content
* `name` - Certificate name


