---
subcategory: "Elastic Load Balancer (ELB)"
layout: "edgenext"
page_title: "EdgeNext: edgenext_elb_certificates"
sidebar_current: "docs-edgenext-datasource-elb_certificates"
description: |-
  Use this data source to query EdgeNext ELB certificates.
---

# edgenext_elb_certificates

Use this data source to query EdgeNext ELB certificates.

## Example Usage

```hcl
data "edgenext_elb_certificates" "example" {
  limit = 10
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Optional, String) Filter by certificate name. Use empty string to omit the filter.
* `page_num` - (Optional, Int) Page number for the list request.
* `page_size` - (Optional, Int) Page size for the list request.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `certificates` - Certificates returned by the API.
  * `container_id` - Barbican container ID.
  * `container_ref` - Barbican container reference URL.
  * `created` - Creation time as Unix timestamp (seconds).
  * `expiration` - Expiration time as Unix timestamp (seconds).
  * `name` - Certificate name.
  * `status` - Certificate status.
  * `type` - Container type (for example certificate).
  * `updated` - Last update time as Unix timestamp (seconds).
* `total` - Total number of certificates matching the query.


