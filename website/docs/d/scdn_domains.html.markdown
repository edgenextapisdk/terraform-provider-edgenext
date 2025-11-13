---
subcategory: "Secure Content Delivery Network (SCDN)"
layout: "edgenext"
page_title: "EdgeNext: edgenext_scdn_domains"
sidebar_current: "docs-edgenext-datasource-scdn_domains"
description: |-
  Use this data source to query a list of SCDN domains with optional filters.
---

# edgenext_scdn_domains

Use this data source to query a list of SCDN domains with optional filters.

## Example Usage

### Query all domains

```hcl
data "edgenext_scdn_domains" "all" {
  page      = 1
  page_size = 100
}

output "domain_count" {
  value = data.edgenext_scdn_domains.all.total
}

output "domain_names" {
  value = [for domain in data.edgenext_scdn_domains.all.domains : domain.domain]
}
```

### Query domains with filters

```hcl
data "edgenext_scdn_domains" "filtered" {
  domain          = "example"
  access_progress = "online"
  protect_status  = "scdn"
  page            = 1
  page_size       = 50
}

output "filtered_domains" {
  value = data.edgenext_scdn_domains.filtered.domains
}
```

### Query domains and save to file

```hcl
data "edgenext_scdn_domains" "all" {
  page               = 1
  page_size          = 100
  result_output_file = "domains.json"
}
```

## Argument Reference

The following arguments are supported:

* `access_mode` - (Optional, String) Filter by access mode
* `access_progress` - (Optional, String) Filter by access progress status
* `ca_status` - (Optional, String) Filter by certificate binding status
* `domain` - (Optional, String) Filter by domain name (fuzzy search)
* `exclusive_resource_id` - (Optional, Int) Filter by exclusive resource package ID
* `group_id` - (Optional, Int) Filter by domain group ID
* `origin_ip` - (Optional, String) Filter by origin IP
* `page_size` - (Optional, Int) The page size for pagination
* `page` - (Optional, Int) The page number for pagination
* `protect_status` - (Optional, String) Filter by edge node type
* `remark` - (Optional, String) Filter by remark (fuzzy search)
* `result_output_file` - (Optional, String) Used to save results to a file

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `domains` - The list of domains
  * `access_mode` - The access mode
  * `access_progress_desc` - The description of the access progress status
  * `access_progress` - The access progress status
  * `ca_id` - The certificate ID
  * `ca_status` - The certificate binding status
  * `cname` - The CNAME information
    * `master` - The master CNAME record
    * `slaves` - The slave CNAME records
  * `created_at` - The creation timestamp
  * `domain` - The domain name
  * `ei_forward_status` - The explicit/implicit forwarding status
  * `exclusive_resource_id` - The exclusive resource package ID
  * `has_origin` - Whether the domain has origin configuration
  * `id` - The ID of the domain
  * `pri_domain` - The primary domain
  * `protect_status` - The edge node type
  * `remark` - The remark for the domain
  * `updated_at` - The last update timestamp
  * `use_my_cname` - The CNAME resolution status
  * `use_my_dns` - The DNS hosting status
* `total` - The total number of domains


