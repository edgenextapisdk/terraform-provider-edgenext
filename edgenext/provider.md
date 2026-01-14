---
layout: "edgenext"
page_title: "Provider: EdgeNext"
sidebar_current: "docs-edgenext-index"
description: |-
  The EdgeNext provider is used to interact with EdgeNext services.
---

# EdgeNext Provider

The EdgeNext Provider can be used to configure infrastructure in [EdgeNext](https://www.edgenext.com) using the EdgeNext Resource Manager API's. Documentation regarding the Data Sources and Resources supported by the EdgeNext Provider can be found in the navigation to the left.

-> **Note:** This provider requires EdgeNext API credentials (access key and secret key).

## Example Usage

```hcl
terraform {
  required_providers {
    edgenext = {
      source  = "edgenextapisdk/edgenext"
      version = "~> 1.0"
    }
  }
}

# Configure the EdgeNext Provider
provider "edgenext" {
  access_key = var.access_key
  secret_key = var.secret_key
  endpoint   = var.endpoint
  region     = var.region
}
```

## Authentication

The EdgeNext provider offers a flexible means of providing credentials for authentication. The following methods are supported, in order of precedence:

1. **Static credentials** in the provider configuration block
2. **Environment variables**

### Static Credentials

!> **Warning:** Hard-coding credentials into any Terraform configuration is not recommended, and risks secret leakage should this file ever be committed to a public version control system.

Static credentials can be provided by adding `access_key` and `secret_key` in-line in the EdgeNext provider block:

```hcl
provider "edgenext" {
  access_key = "your-access-key" # Set to edgenext-<your-username>
  secret_key = "your-secret-key" # Please contact the operations team to obtain it
  endpoint   = "https://cdn.api.edgenext.com" # CDN (https://cdn.api.edgenext.com) / SCDN (https://api.edgenextscdn.com)
  region     = "us-east-1" # Optional
}
```

### Environment Variables

You can provide your credentials via the `EDGENEXT_ACCESS_KEY`, `EDGENEXT_SECRET_KEY`, `EDGENEXT_ENDPOINT` and `EDGENEXT_REGION` environment variables:

```hcl
provider "edgenext" {}
```

Usage:

```bash
export EDGENEXT_ACCESS_KEY="your-access-key"
export EDGENEXT_SECRET_KEY="your-secret-key"
export EDGENEXT_ENDPOINT="https://cdn.api.edgenext.com"
export EDGENEXT_REGION="us-east-1"
terraform plan
```

## Argument Reference

The following arguments are supported in the `provider` block:

* `access_key` - (Required) EdgeNext access key for authentication. It can also be sourced from the `EDGENEXT_ACCESS_KEY` environment variable.

* `secret_key` - (Required) EdgeNext secret key for authentication. It can also be sourced from the `EDGENEXT_SECRET_KEY` environment variable.

* `endpoint` - (Required) EdgeNext API endpoint address. It can also be sourced from the `EDGENEXT_ENDPOINT` environment variable.

* `region` - (Optional) EdgeNext region. It can also be sourced from the `EDGENEXT_REGION` environment variable.

Resources List

Content Delivery Network (CDN)
Data Source
edgenext_cdn_domain
edgenext_cdn_domains
edgenext_cdn_push
edgenext_cdn_pushes
edgenext_cdn_purge
edgenext_cdn_purges

Resource
edgenext_cdn_domain
edgenext_cdn_push
edgenext_cdn_purge

SSL Certificate Management (SSL)
Data Source
edgenext_ssl_certificate
edgenext_ssl_certificates

Resource
edgenext_ssl_certificate

Object Storage Service (OSS)
Data Source
edgenext_oss_buckets
edgenext_oss_object
edgenext_oss_objects

Resource
edgenext_oss_bucket
edgenext_oss_object
edgenext_oss_object_copy

Security CDN (SCDN)
Data Source
edgenext_scdn_domain
edgenext_scdn_domains
edgenext_scdn_origin
edgenext_scdn_origins
edgenext_scdn_domain_base_settings
edgenext_scdn_access_progress
edgenext_scdn_domain_templates
edgenext_scdn_brief_domains
edgenext_scdn_certificate
edgenext_scdn_certificates
edgenext_scdn_certificates_by_domains
edgenext_scdn_certificate_export
edgenext_scdn_rule_template
edgenext_scdn_rule_templates
edgenext_scdn_rule_template_domains
edgenext_scdn_network_speed_config
edgenext_scdn_network_speed_rules
edgenext_scdn_cache_rules
edgenext_scdn_cache_global_config
edgenext_scdn_security_protection_ddos_config
edgenext_scdn_security_protection_waf_config
edgenext_scdn_security_protection_template
edgenext_scdn_security_protection_templates
edgenext_scdn_security_protection_template_domains
edgenext_scdn_security_protection_template_unbound_domains
edgenext_scdn_security_protection_member_global_template
edgenext_scdn_security_protection_iota
edgenext_scdn_origin_group
edgenext_scdn_origin_groups
edgenext_scdn_origin_groups_all
edgenext_scdn_cache_clean_config
edgenext_scdn_cache_clean_tasks
edgenext_scdn_cache_clean_task_detail
edgenext_scdn_cache_preheat_tasks
edgenext_scdn_log_download_tasks
edgenext_scdn_log_download_templates
edgenext_scdn_log_download_fields

Resource
edgenext_scdn_domain
edgenext_scdn_origin
edgenext_scdn_cert_binding
edgenext_scdn_domain_base_settings
edgenext_scdn_domain_status
edgenext_scdn_domain_node_switch
edgenext_scdn_domain_access_mode
edgenext_scdn_certificate
edgenext_scdn_certificate_apply
edgenext_scdn_rule_template
edgenext_scdn_rule_template_domain_bind
edgenext_scdn_rule_template_domain_unbind
edgenext_scdn_network_speed_config
edgenext_scdn_network_speed_rule
edgenext_scdn_network_speed_rules_sort
edgenext_scdn_cache_rule
edgenext_scdn_cache_rule_status
edgenext_scdn_cache_rules_sort
edgenext_scdn_security_protection_ddos_config
edgenext_scdn_security_protection_waf_config
edgenext_scdn_security_protection_template
edgenext_scdn_security_protection_template_domain_bind
edgenext_scdn_security_protection_template_batch_config
edgenext_scdn_origin_group
edgenext_scdn_origin_group_domain_bind
edgenext_scdn_origin_group_domain_copy
edgenext_scdn_cache_clean_task
edgenext_scdn_cache_preheat_task
edgenext_scdn_log_download_task
edgenext_scdn_log_download_template
edgenext_scdn_log_download_template_status
