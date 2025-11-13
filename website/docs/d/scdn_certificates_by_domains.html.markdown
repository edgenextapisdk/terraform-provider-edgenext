---
subcategory: "Secure Content Delivery Network (SCDN)"
layout: "edgenext"
page_title: "EdgeNext: edgenext_scdn_certificates_by_domains"
sidebar_current: "docs-edgenext-datasource-scdn_certificates_by_domains"
description: |-
  Use this data source to query certificates bound to specific domains.
---

# edgenext_scdn_certificates_by_domains

Use this data source to query certificates bound to specific domains.

## Example Usage

### Query certificates by domains

```hcl
data "edgenext_scdn_certificates_by_domains" "example" {
  domains = ["example.com", "www.example.com"]
}

output "certificates" {
  value = data.edgenext_scdn_certificates_by_domains.example.certificates
}
```

### Query and save to file

```hcl
data "edgenext_scdn_certificates_by_domains" "example" {
  domains            = ["example.com"]
  result_output_file = "certificates_by_domains.json"
}
```

## Argument Reference

The following arguments are supported:

* `domains` - (Required, List: [`String`]) The list of domain names to query certificates for
* `result_output_file` - (Optional, String) Used to save results to a file

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `certificates` - The list of certificates
  * `apply_status` - The application status
  * `binded` - Whether the certificate is bound
  * `ca_domain` - The domains in the certificate
  * `ca_name` - The certificate name
  * `ca_type_domain` - The certificate domain type
  * `ca_type` - The certificate type
  * `code` - The application error code
  * `id` - The certificate ID
  * `issuer_expiry_time_auto_renew_status` - The certificate auto-renew status
  * `issuer_expiry_time_desc` - The certificate expiry time description
  * `issuer_expiry_time` - The certificate expiry time
  * `issuer_start_time` - The certificate start time
  * `issuer` - The certificate issuer
  * `member_id` - The member ID
  * `msg` - The application error message
  * `renew_status` - The renewal status


