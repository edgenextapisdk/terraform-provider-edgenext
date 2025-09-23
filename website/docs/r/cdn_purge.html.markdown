---
subcategory: "Content Delivery Network(CDN)"
layout: "edgenext"
page_title: "EdgeNext: edgenext_cdn_purge"
sidebar_current: "docs-edgenext-resource-cdn_purge"
description: |-
  Provides a resource to create and manage CDN cache purge tasks.
---

# edgenext_cdn_purge

Provides a resource to create and manage CDN cache purge tasks.

## Example Usage

### Basic CDN cache purge

```hcl
resource "edgenext_cdn_purge" "example" {
  urls = [
    "https://example.com/static/old-image.jpg",
    "https://example.com/static/old-style.css"
  ]
}
```

### CDN cache purge with multiple URLs

```hcl
resource "edgenext_cdn_purge" "example" {
  urls = [
    "https://example.com/api/data.json",
    "https://example.com/static/images/photo1.jpg",
    "https://example.com/static/css/style.css"
  ]
}
```

## Argument Reference

The following arguments are supported:

* `urls` - (Required, List: [`String`], ForceNew) List of URLs to purge, maximum 500 URLs per request

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

CDN purge tasks can be imported using the task ID:

```shell
terraform import edgenext_cdn_purge.example purge-task-123456
```

