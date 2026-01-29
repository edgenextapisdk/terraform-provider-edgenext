---
subcategory: "Security CDN (SCDN)"
layout: "edgenext"
page_title: "EdgeNext: edgenext_scdn_cache_preheat_task"
sidebar_current: "docs-edgenext-resource-scdn_cache_preheat_task"
description: |-
  Provides a resource to create SCDN cache preheat tasks.
---

# edgenext_scdn_cache_preheat_task

Provides a resource to create SCDN cache preheat tasks.

## Example Usage

### Preheat cache for URLs

```hcl
resource "edgenext_scdn_cache_preheat_task" "example" {
  preheat_url = [
    "https://example.com/page1",
    "https://example.com/page2",
    "https://example.com/page3"
  ]
}
```

### Preheat by group

```hcl
resource "edgenext_scdn_cache_preheat_task" "example" {
  group_id = 1
  protocol = "https"
  port     = "443"
  preheat_url = [
    "https://example.com/page1"
  ]
}
```

## Argument Reference

The following arguments are supported:

* `preheat_url` - (Required, List: [`String`]) Preheat URLs
* `group_id` - (Optional, Int) Group ID, can refresh cache by group
* `port` - (Optional, String) Website port, only needed for special ports; only valid when refreshing by group
* `protocol` - (Optional, String) Protocol: http/https; only valid when refreshing by group

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `error_url` - List of URLs with preheat errors
* `id` - The ID of the preheat task (generated timestamp)


