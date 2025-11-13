---
subcategory: "Secure Content Delivery Network (SCDN)"
layout: "edgenext"
page_title: "EdgeNext: edgenext_scdn_certificate"
sidebar_current: "docs-edgenext-datasource-scdn_certificate"
description: |-
  Use this data source to query details of a specific SCDN certificate.
---

# edgenext_scdn_certificate

Use this data source to query details of a specific SCDN certificate.

## Example Usage

### Query certificate by ID

```hcl
data "edgenext_scdn_certificate" "example" {
  id = "12345"
}

output "certificate_name" {
  value = data.edgenext_scdn_certificate.example.ca_name
}

output "expiry_time" {
  value = data.edgenext_scdn_certificate.example.issuer_expiry_time
}
```

### Query and save to file

```hcl
data "edgenext_scdn_certificate" "example" {
  id                 = "12345"
  result_output_file = "certificate.json"
}
```

## Argument Reference

The following arguments are supported:

* `id` - (Required, String) The certificate ID
* `result_output_file` - (Optional, String) Used to save results to a file

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `apply_status` - The application status: 1-applying, 2-issued, 3-review failed, 4-uploaded
* `authentication_usable_domain` - The authentication usable domain
* `binded` - Whether the certificate is bound
* `ca_domain` - The domains in the certificate
* `ca_name` - The certificate name
* `ca_sn` - The certificate serial number
* `ca_type_domain` - The certificate domain type: 1-single domain, 2-multiple domains, 3-wildcard domain
* `ca_type` - The certificate type: 1-upload, 2-lets apply
* `city` - The city
* `code` - The application error code
* `country` - The country
* `created_at` - The creation timestamp
* `issuer_expiry_time_auto_renew_status` - The certificate auto-renew status
* `issuer_expiry_time_desc` - The certificate expiry time description
* `issuer_expiry_time` - The certificate expiry time
* `issuer_object` - The issuer object
* `issuer_organization_element` - The issuer organization element
* `issuer_organization` - The issuer organization
* `issuer_start_time` - The certificate start time
* `issuer` - The certificate issuer
* `member_id` - The member ID
* `msg` - The application error message
* `province` - The province
* `renew_status` - The renewal status: 1-default, 2-renewing, 3-renewal failed, 4-renewal successful
* `serial_number` - The certificate serial number
* `updated_at` - The last update timestamp
* `use_organization_element` - The use organization element
* `use_organization` - The use organization


