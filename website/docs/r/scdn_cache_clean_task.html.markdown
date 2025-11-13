---
subcategory: "Security CDN (SCDN)"
layout: "edgenext"
page_title: "EdgeNext: edgenext_scdn_cache_clean_task"
sidebar_current: "docs-edgenext-resource-scdn_cache_clean_task"
description: |-
  Provides a resource to create SCDN cache clean tasks.
---

# edgenext_scdn_cache_clean_task

Provides a resource to create SCDN cache clean tasks.

## Example Usage

### Clean whole site cache

```hcl
resource "edgenext_scdn_cache_clean_task" "example" {
  wholesite = ["example.com", "www.example.com"]
}
```

### Clean special URLs

```hcl
resource "edgenext_scdn_cache_clean_task" "example" {
  specialurl = [
    "https://example.com/page1",
    "https://example.com/page2"
  ]
}
```

### Clean special directories

```hcl
resource "edgenext_scdn_cache_clean_task" "example" {
  specialdir = [
    "/static/",
    "/images/"
  ]
}
```

### Clean by group

```hcl
resource "edgenext_scdn_cache_clean_task" "example" {
  group_id  = 1
  protocol  = "https"
  port      = "443"
  wholesite = ["example.com"]
}
```

## Argument Reference

The following arguments are supported:

* `group_id` - (Optional, Int) Group ID, can refresh cache by group
* `port` - (Optional, String) Website port, only needed for special ports; only valid when refreshing by group
* `protocol` - (Optional, String) Protocol: http/https; only valid when refreshing by group
* `specialdir` - (Optional, List: [`String`]) Special directories to clean
* `specialurl` - (Optional, List: [`String`]) Special URLs to clean
* `wholesite` - (Optional, List: [`String`]) Whole site domains to clean

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `id` - The ID of the cache clean task (generated timestamp)


