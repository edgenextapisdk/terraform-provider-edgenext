---
subcategory: "Security CDN (SCDN)"
layout: "edgenext"
page_title: "EdgeNext: edgenext_scdn_domain"
sidebar_current: "docs-edgenext-datasource-scdn_domain"
description: |-
  Use this data source to query details of a specific SCDN domain.
---

# edgenext_scdn_domain

Use this data source to query details of a specific SCDN domain.

## Example Usage

### Query domain by name

```hcl
data "edgenext_scdn_domain" "example" {
  domain = "example.com"
}

output "domain_id" {
  value = data.edgenext_scdn_domain.example.id
}

output "domain_status" {
  value = data.edgenext_scdn_domain.example.access_progress
}
```

### Query domain and save to file

```hcl
data "edgenext_scdn_domain" "example" {
  domain             = "example.com"
  result_output_file = "domain.json"
}
```

## Argument Reference

The following arguments are supported:

* `domain_id` - (Optional, String) The ID of the domain to query (deprecated, use id instead)
* `domain` - (Optional, String) The domain name to query (either domain, id, or domain_id must be provided)
* `id` - (Optional, String) The ID of the domain to query (either domain, id, or domain_id must be provided). Also returned as the computed ID.
* `result_output_file` - (Optional, String) Used to save results to a file

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `access_mode` - The access mode (ns or cname)
* `access_progress_desc` - The description of the access progress status
* `access_progress` - The access progress status
* `app_type` - The application type
* `ca_id` - The certificate ID
* `ca_status` - The certificate binding status
* `cname` - The CNAME information
  * `master` - The master CNAME record
  * `slaves` - The slave CNAME records
* `created_at` - The creation timestamp
* `ei_forward_status` - The explicit/implicit forwarding status
* `exclusive_resource_id` - The ID of the exclusive resource package
* `group_id` - The ID of the domain group
* `has_origin` - Whether the domain has origin configuration
* `origins` - The origin server configuration
  * `id` - The ID of the origin
  * `listen_port` - The listening port of the origin server
  * `load_balance` - The load balancing method
  * `origin_protocol` - The origin protocol
  * `origin_type` - The origin type
  * `protocol` - The origin protocol
  * `records` - The origin records
    * `port` - The port of the record
    * `priority` - The priority of the record
    * `value` - The value of the record
    * `view` - The view of the record
* `pri_domain` - The primary domain
* `protect_status` - The edge node type
* `remark` - The remark for the domain
* `tpl_id` - The template ID applied to the domain
* `tpl_recommend` - The template recommendation status
* `updated_at` - The last update timestamp
* `use_my_cname` - The CNAME resolution status
* `use_my_dns` - The DNS hosting status


