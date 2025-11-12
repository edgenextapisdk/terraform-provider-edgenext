---
subcategory: "Content Delivery Network (CDN)"
layout: "edgenext"
page_title: "EdgeNext: edgenext_cdn_domains"
sidebar_current: "docs-edgenext-datasource-cdn_domains"
description: |-
  Use this data source to query a list of CDN domains.
---

# edgenext_cdn_domains

Use this data source to query a list of CDN domains.

## Example Usage

### Query all CDN domains

```hcl
data "edgenext_cdn_domains" "example" {
  page_number = 1
  page_size   = 100
}
```

### Query CDN domains with specific status

```hcl
data "edgenext_cdn_domains" "example" {
  domain_status = "serving"
  page_number   = 1
  page_size     = 50
}
```

## Argument Reference

The following arguments are supported:

* `domain_status` - (Optional, String) Specify the service status of the domain, support specifying multiple service status queries: 
serving：Serving. When querying with serving, the domain whose "status" is "deploying" is in the configuration deployment state.
suspend：Suspended
deleted：Deleted
Default value is all status domain names when not specified.
* `output_file` - (Optional, String) Used to save results.
* `page_number` - (Optional, Int) Get the page number. 
Default value is 1 when not specified.
* `page_size` - (Optional, Int) Page size, value range: 1-500. 
Default value is 100 when not specified.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `list` - Domain list
  * `area` - Acceleration area: mainland_china, outside_mainland_china, global
  * `cname` - CNAME
  * `create_time` - Creation time
  * `https` - HTTPS
  * `icp_num` - ICP filing number
  * `icp_status` - ICP filing status
  * `id` - Domain ID
  * `status` - Domain status
  * `type` - Domain type: page(web page), download(download), video_demand(video demand), dynamic(dynamic)
  * `update_time` - Update time


