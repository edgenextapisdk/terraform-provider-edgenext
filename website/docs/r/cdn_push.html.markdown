---
subcategory: "Content Delivery Network(CDN)"
layout: "edgenext"
page_title: "EdgeNext: edgenext_cdn_push"
sidebar_current: "docs-edgenext-resource-cdn_push"
description: |-
  Provides a resource to create and manage CDN cache push tasks.
---

# edgenext_cdn_push

Provides a resource to create and manage CDN cache push tasks.

## Example Usage

### Basic CDN cache push (URLs)

```hcl
resource "edgenext_cdn_push" "example" {
  type = "url"
  urls = [
    "https://example.com/static/image1.jpg",
    "https://example.com/static/image2.jpg"
  ]
}
```

### CDN cache push (directories)

```hcl
resource "edgenext_cdn_push" "example" {
  type = "dir"
  urls = [
    "https://example.com/static/css/",
    "https://example.com/static/js/"
  ]
}
```

## Argument Reference

The following arguments are supported:

* `type` - (Required, String, ForceNew) URL type for push: dir(directory), url(URL)
* `urls` - (Required, List: [`String`], ForceNew) List of URLs/directories to refresh, maximum 500 URLs per request

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `list` - List of successfully submitted URLs
  * `complete_time` - Completion time
  * `create_time` - Creation time
  * `id` - URL ID
  * `status` - Status
  * `type` - URL type
  * `url` - URL/Directory
* `task_id` - Task ID for this submission
* `total` - Number of successfully submitted URLs/directories


## Import

CDN push tasks can be imported using the task ID:

```shell
terraform import edgenext_cdn_push.example push-task-123456
```

