---
subcategory: "Content Delivery Network (CDN)"
layout: "edgenext"
page_title: "EdgeNext: edgenext_cdn_prefetch"
sidebar_current: "docs-edgenext-resource-cdn_prefetch"
description: |-
  Provides a resource to create and manage CDN cache prefetch tasks.
---

# edgenext_cdn_prefetch

Provides a resource to create and manage CDN cache prefetch tasks.

## Example Usage

### Basic CDN cache prefetch

```hcl
resource "edgenext_cdn_prefetch" "example" {
  urls = [
    "https://example.com/static/old-image.jpg",
    "https://example.com/static/old-style.css"
  ]
}
```

### CDN cache prefetch with multiple URLs

```hcl
resource "edgenext_cdn_prefetch" "example" {
  urls = [
    "https://example.com/api/data.json",
    "https://example.com/static/images/photo1.jpg",
    "https://example.com/static/css/style.css"
  ]
}
```

## Argument Reference

The following arguments are supported:

* `urls` - (Required, List: [`String`], ForceNew) List of URLs to prefetch, maximum 500 URLs per request

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `list` - List of successfully submitted URLs
  * `complete_time` - Completion time
  * `create_time` - Creation time
  * `id` - URL ID
  * `status` - Status
  * `url` - URL
* `task_id` - Task ID for this submission
* `total` - Number of successfully submitted URLs


## Import

CDN prefetch tasks can be imported using the task ID:

```shell
terraform import edgenext_cdn_prefetch.example prefetch-task-123456
```

