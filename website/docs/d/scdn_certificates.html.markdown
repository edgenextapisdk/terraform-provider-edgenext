---
subcategory: "Secure Content Delivery Network (SCDN)"
layout: "edgenext"
page_title: "EdgeNext: edgenext_scdn_certificates"
sidebar_current: "docs-edgenext-datasource-scdn_certificates"
description: |-
  Use this data source to query a list of SCDN certificates with optional filters.
---

# edgenext_scdn_certificates

Use this data source to query a list of SCDN certificates with optional filters.

## Example Usage

### Query all certificates

```hcl
data "edgenext_scdn_certificates" "all" {
  page     = 1
  per_page = 20
}

output "certificate_count" {
  value = data.edgenext_scdn_certificates.all.total
}

output "certificates" {
  value = data.edgenext_scdn_certificates.all.list
}
```

### Query certificates with filters

```hcl
data "edgenext_scdn_certificates" "filtered" {
  page         = 1
  per_page     = 20
  domain       = "example.com"
  binded       = "true"
  apply_status = "2"
}

output "filtered_certificates" {
  value = data.edgenext_scdn_certificates.filtered.list
}
```

### Query and save to file

```hcl
data "edgenext_scdn_certificates" "all" {
  page               = 1
  per_page           = 20
  result_output_file = "certificates.json"
}
```

## Argument Reference

The following arguments are supported:

* `apply_status` - (Optional, String) Filter by application status: 1-applying, 2-issued, 3-review failed, 4-uploaded
* `binded` - (Optional, String) Filter by binding status: true-bound, false-unbound
* `ca_name` - (Optional, String) Filter by certificate name
* `domain` - (Optional, String) Filter by domain name
* `expiry_time` - (Optional, String) Filter by expiry status: true-expired, false-not expired, inno-about to expire (within 30 days)
* `is_exact_search` - (Optional, String) Whether to use exact search: on-yes, off-no
* `issuer` - (Optional, String) Filter by issuer
* `page` - (Optional, Int) The page number for pagination
* `per_page` - (Optional, Int) The page size for pagination
* `product_flag` - (Optional, String) Filter by product flag
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
* `issuer_list` - The list of issuers
* `total` - The total number of certificates


