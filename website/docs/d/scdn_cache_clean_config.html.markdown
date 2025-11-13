---
subcategory: "Security CDN (SCDN)"
layout: "edgenext"
page_title: "EdgeNext: edgenext_scdn_cache_clean_config"
sidebar_current: "docs-edgenext-datasource-scdn_cache_clean_config"
description: |-
  Use this data source to query SCDN cache clean configuration.
---

# edgenext_scdn_cache_clean_config

Use this data source to query SCDN cache clean configuration.

## Example Usage

### Query cache clean config

```hcl
data "edgenext_scdn_cache_clean_config" "example" {
}

output "wholesite" {
  value = data.edgenext_scdn_cache_clean_config.example.wholesite
}

output "specialurl" {
  value = data.edgenext_scdn_cache_clean_config.example.specialurl
}
```

### Query and save to file

```hcl
data "edgenext_scdn_cache_clean_config" "example" {
  result_output_file = "cache_clean_config.json"
}
```

## Argument Reference

The following arguments are supported:

* `result_output_file` - (Optional, String) Used to save results to a file

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - Config ID
* `specialdir` - Special directory config
* `specialurl` - Special URL config
* `wholesite` - Whole site config


