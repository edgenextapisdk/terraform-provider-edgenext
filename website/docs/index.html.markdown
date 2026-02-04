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

## Resources and Data Sources

The EdgeNext provider supports the following resource types:

### Content Delivery Network (CDN)

#### Resources

* [`edgenext_cdn_domain`](resources/cdn_domain) - Manage CDN domain configuration
* [`edgenext_cdn_purge`](resources/cdn_purge) - Manage CDN cache purge tasks
* [`edgenext_cdn_prefetch`](resources/cdn_prefetch) - Manage CDN cache prefetch tasks

#### Data Sources

* [`edgenext_cdn_domain`](data-sources/cdn_domain) - Query CDN domain configuration
* [`edgenext_cdn_domains`](data-sources/cdn_domains) - Query CDN domains
* [`edgenext_cdn_purge`](data-sources/cdn_purge) - Query CDN purge task details
* [`edgenext_cdn_purges`](data-sources/cdn_purges) - Query CDN purge tasks
* [`edgenext_cdn_prefetch`](data-sources/cdn_prefetch) - Query CDN prefetch task details
* [`edgenext_cdn_prefetches`](data-sources/cdn_prefetches) - Query CDN prefetch tasks

### SSL Certificate Management (SSL)

#### Resources

* [`edgenext_ssl_certificate`](resources/ssl_certificate) - Manage SSL certificates

#### Data Sources

* [`edgenext_ssl_certificate`](data-sources/ssl_certificate) - Query SSL certificate details
* [`edgenext_ssl_certificates`](data-sources/ssl_certificates) - Query SSL certificates

### Object Storage Service (OSS)

#### Resources

* [`edgenext_oss_bucket`](resources/oss_bucket) - Manage OSS buckets
* [`edgenext_oss_object`](resources/oss_object) - Manage OSS objects
* [`edgenext_oss_object_copy`](resources/oss_object_copy) - Manage OSS object copy

#### Data Sources

* [`edgenext_oss_buckets`](data-sources/oss_buckets) - Query OSS buckets
* [`edgenext_oss_object`](data-sources/oss_object) - Query OSS object details
* [`edgenext_oss_objects`](data-sources/oss_objects) - Query OSS objects

### Security DNS (SDNS)

#### Resources

* [`edgenext_sdns_domain`](resources/sdns_domain) - Manage sdns domain
* [`edgenext_sdns_domain_group`](resources/sdns_domain_group) - Manage sdns domain group
* [`edgenext_sdns_record`](resources/sdns_record) - Manage sdns record

#### Data Sources

* [`edgenext_sdns_domains`](data-sources/sdns_domains) - Query sdns domains
* [`edgenext_sdns_domain_groups`](data-sources/sdns_domain_groups) - Query sdns domain groups
* [`edgenext_sdns_records`](data-sources/sdns_records) - Query sdns records

### Security CDN (SCDN)

#### Resources

* [`edgenext_scdn_domain`](resources/scdn_domain) - Manage SCDN domain configuration
* [`edgenext_scdn_origin`](resources/scdn_origin) - Manage SCDN origin servers
* [`edgenext_scdn_cert_binding`](resources/scdn_cert_binding) - Manage SCDN certificate bindings
* [`edgenext_scdn_domain_base_settings`](resources/scdn_domain_base_settings) - Manage SCDN domain base settings
* [`edgenext_scdn_domain_status`](resources/scdn_domain_status) - Manage SCDN domain status management
* [`edgenext_scdn_domain_node_switch`](resources/scdn_domain_node_switch) - Manage SCDN domain node switching
* [`edgenext_scdn_domain_access_mode`](resources/scdn_domain_access_mode) - Manage SCDN domain access mode
* [`edgenext_scdn_certificate`](resources/scdn_certificate) - Manage SCDN certificates
* [`edgenext_scdn_certificate_apply`](resources/scdn_certificate_apply) - Manage SCDN certificate application
* [`edgenext_scdn_rule_template`](resources/scdn_rule_template) - Manage SCDN rule templates
* [`edgenext_scdn_rule_template_domain_bind`](resources/scdn_rule_template_domain_bind) - Manage SCDN rule template domain bindings
* [`edgenext_scdn_rule_template_domain_unbind`](resources/scdn_rule_template_domain_unbind) - Manage SCDN rule template domain unbindings
* [`edgenext_scdn_rule_template_switch`](resources/scdn_rule_template_switch) - Manage scdn rule template switch
* [`edgenext_scdn_network_speed_config`](resources/scdn_network_speed_config) - Manage SCDN network speed configuration
* [`edgenext_scdn_network_speed_rule`](resources/scdn_network_speed_rule) - Manage SCDN network speed rules
* [`edgenext_scdn_network_speed_rules_sort`](resources/scdn_network_speed_rules_sort) - Manage SCDN network speed rules sorting
* [`edgenext_scdn_cache_rule`](resources/scdn_cache_rule) - Manage SCDN cache rules
* [`edgenext_scdn_cache_rule_status`](resources/scdn_cache_rule_status) - Manage SCDN cache rule status
* [`edgenext_scdn_cache_rules_sort`](resources/scdn_cache_rules_sort) - Manage SCDN cache rules sorting
* [`edgenext_scdn_security_protection_ddos_config`](resources/scdn_security_protection_ddos_config) - Manage SCDN DDoS protection configuration
* [`edgenext_scdn_security_protection_waf_config`](resources/scdn_security_protection_waf_config) - Manage SCDN WAF protection configuration
* [`edgenext_scdn_security_protection_template`](resources/scdn_security_protection_template) - Manage SCDN security protection templates
* [`edgenext_scdn_security_protection_template_domain_bind`](resources/scdn_security_protection_template_domain_bind) - Manage SCDN security protection template domain bindings
* [`edgenext_scdn_security_protection_template_batch_config`](resources/scdn_security_protection_template_batch_config) - Manage SCDN security protection template batch configuration
* [`edgenext_scdn_origin_group`](resources/scdn_origin_group) - Manage SCDN origin groups
* [`edgenext_scdn_origin_group_domain_bind`](resources/scdn_origin_group_domain_bind) - Manage SCDN origin group domain bindings
* [`edgenext_scdn_origin_group_domain_copy`](resources/scdn_origin_group_domain_copy) - Manage SCDN origin group domain copying
* [`edgenext_scdn_cache_clean_task`](resources/scdn_cache_clean_task) - Manage SCDN cache clean tasks
* [`edgenext_scdn_cache_preheat_task`](resources/scdn_cache_preheat_task) - Manage SCDN cache preheat tasks
* [`edgenext_scdn_log_download_task`](resources/scdn_log_download_task) - Manage SCDN log download tasks
* [`edgenext_scdn_log_download_template`](resources/scdn_log_download_template) - Manage SCDN log download templates
* [`edgenext_scdn_log_download_template_status`](resources/scdn_log_download_template_status) - Manage SCDN log download template status
* [`edgenext_scdn_domain_group`](resources/scdn_domain_group) - Manage scdn domain group
* [`edgenext_scdn_user_ip`](resources/scdn_user_ip) - Manage scdn user ip
* [`edgenext_scdn_user_ip_item`](resources/scdn_user_ip_item) - Manage scdn user ip item

#### Data Sources

* [`edgenext_scdn_domain`](data-sources/scdn_domain) - Query SCDN domain details
* [`edgenext_scdn_domains`](data-sources/scdn_domains) - Query SCDN domains
* [`edgenext_scdn_origin`](data-sources/scdn_origin) - Query SCDN origin details
* [`edgenext_scdn_origins`](data-sources/scdn_origins) - Query SCDN origins
* [`edgenext_scdn_domain_base_settings`](data-sources/scdn_domain_base_settings) - Query SCDN domain base settings
* [`edgenext_scdn_access_progress`](data-sources/scdn_access_progress) - Query SCDN access progress options
* [`edgenext_scdn_domain_templates`](data-sources/scdn_domain_templates) - Query SCDN domain templates
* [`edgenext_scdn_brief_domains`](data-sources/scdn_brief_domains) - Query SCDN brief domain information
* [`edgenext_scdn_certificate`](data-sources/scdn_certificate) - Query SCDN certificate details
* [`edgenext_scdn_certificates`](data-sources/scdn_certificates) - Query SCDN certificates
* [`edgenext_scdn_certificates_by_domains`](data-sources/scdn_certificates_by_domains) - Query SCDN certificates by domains
* [`edgenext_scdn_certificate_export`](data-sources/scdn_certificate_export) - Query SCDN certificate export
* [`edgenext_scdn_rule_template`](data-sources/scdn_rule_template) - Query SCDN rule template details
* [`edgenext_scdn_rule_templates`](data-sources/scdn_rule_templates) - Query SCDN rule templates
* [`edgenext_scdn_rule_template_domains`](data-sources/scdn_rule_template_domains) - Query SCDN rule template domains
* [`edgenext_scdn_network_speed_config`](data-sources/scdn_network_speed_config) - Query SCDN network speed configuration
* [`edgenext_scdn_network_speed_rules`](data-sources/scdn_network_speed_rules) - Query SCDN network speed rules
* [`edgenext_scdn_cache_rules`](data-sources/scdn_cache_rules) - Query SCDN cache rules
* [`edgenext_scdn_cache_global_config`](data-sources/scdn_cache_global_config) - Query SCDN cache global configuration
* [`edgenext_scdn_security_protection_ddos_config`](data-sources/scdn_security_protection_ddos_config) - Query SCDN DDoS protection configuration
* [`edgenext_scdn_security_protection_waf_config`](data-sources/scdn_security_protection_waf_config) - Query SCDN WAF protection configuration
* [`edgenext_scdn_security_protection_template`](data-sources/scdn_security_protection_template) - Query SCDN security protection template details
* [`edgenext_scdn_security_protection_templates`](data-sources/scdn_security_protection_templates) - Query SCDN security protection templates
* [`edgenext_scdn_security_protection_template_domains`](data-sources/scdn_security_protection_template_domains) - Query SCDN security protection template domains
* [`edgenext_scdn_security_protection_template_unbound_domains`](data-sources/scdn_security_protection_template_unbound_domains) - Query SCDN security protection template unbound domains
* [`edgenext_scdn_security_protection_member_global_template`](data-sources/scdn_security_protection_member_global_template) - Query SCDN security protection member global template
* [`edgenext_scdn_security_protection_iota`](data-sources/scdn_security_protection_iota) - Query SCDN security protection IOTA
* [`edgenext_scdn_origin_group`](data-sources/scdn_origin_group) - Query SCDN origin group details
* [`edgenext_scdn_origin_groups`](data-sources/scdn_origin_groups) - Query SCDN origin groups
* [`edgenext_scdn_origin_groups_all`](data-sources/scdn_origin_groups_all) - Query SCDN all origin groups
* [`edgenext_scdn_cache_clean_config`](data-sources/scdn_cache_clean_config) - Query SCDN cache clean configuration
* [`edgenext_scdn_cache_clean_tasks`](data-sources/scdn_cache_clean_tasks) - Query SCDN cache clean tasks
* [`edgenext_scdn_cache_clean_task_detail`](data-sources/scdn_cache_clean_task_detail) - Query SCDN cache clean task details
* [`edgenext_scdn_cache_preheat_tasks`](data-sources/scdn_cache_preheat_tasks) - Query SCDN cache preheat tasks
* [`edgenext_scdn_log_download_tasks`](data-sources/scdn_log_download_tasks) - Query SCDN log download tasks
* [`edgenext_scdn_log_download_templates`](data-sources/scdn_log_download_templates) - Query SCDN log download templates
* [`edgenext_scdn_log_download_fields`](data-sources/scdn_log_download_fields) - Query SCDN log download fields
* [`edgenext_scdn_domain_groups`](data-sources/scdn_domain_groups) - Query scdn domain groups
* [`edgenext_scdn_domain_group_domains`](data-sources/scdn_domain_group_domains) - Query scdn domain group domains
* [`edgenext_scdn_user_ips`](data-sources/scdn_user_ips) - Query scdn user ips
* [`edgenext_scdn_user_ip_items`](data-sources/scdn_user_ip_items) - Query scdn user ip items

