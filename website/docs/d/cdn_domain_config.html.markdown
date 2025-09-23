---
subcategory: "Content Delivery Network(CDN)"
layout: "edgenext"
page_title: "EdgeNext: edgenext_cdn_domain"
sidebar_current: "docs-edgenext-datasource-cdn_domain"
description: |-
  Use this data source to query detailed information of CDN domain configuration.
---

# edgenext_cdn_domain

Use this data source to query detailed information of CDN domain configuration.

## Example Usage

### Query CDN domain configuration by domain name

```hcl
data "edgenext_cdn_domain" "example" {
  domain = "example.com"
}
```

### Query specific configuration items for a CDN domain

```hcl
data "edgenext_cdn_domain" "example" {
  domain = "example.com"
  config_item = [
    "origin",
    "https",
    "cache_rule",
    "referer"
  ]
  result_output_file = "domain_config.json"
}
```

## Argument Reference

The following arguments are supported:

* `domain` - (Required, String) Domain to query configuration for.
* `config_item` - (Optional, List: [`String`]) Configuration items.
* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `area` - Acceleration area: mainland_china (mainland China), outside_mainland_china (outside mainland China), global (global).
* `cname` - CNAME.
* `config` - Configuration.
* `create_time` - Creation time.
* `https` - HTTPS.
* `icp_num` - ICP filing number.
* `icp_status` - ICP filing status.
* `status` - Domain status.
* `type` - Domain type: page (web page), download (download), video_demand (video on demand), dynamic (dynamic).
* `update_time` - Update time.


