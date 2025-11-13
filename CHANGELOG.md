# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [1.2.1] - 2025-11-13

### Changed
- **Documentation Improvements**: Updated SCDN documentation structure and naming conventions
  - Changed SCDN subdirectory format to "Security CDN (SCDN) - {SubGroup}" in provider.md
  - Replaced "Secure Content Delivery Network (SCDN)" with "Security CDN (SCDN)" across all documentation
  - Updated documentation generator to support subdirectory grouping with display names
  - Updated sidebar and preview page to show SCDN subdirectories correctly
  - Regenerated all 74 documentation files with updated naming

## [1.2.0] - 2025-11-13

### Added

#### SCDN (Security CDN) Module
- **Complete SCDN Domain Management**: Full lifecycle management for SCDN domains
  - Domain creation, configuration, and deletion
  - Origin server configuration with multiple protocols and load balancing
  - Domain status management and access mode switching
  - Certificate binding and management
  - Domain templates and base settings configuration

- **Cache Management**: Comprehensive cache configuration and operations
  - Cache rules with flexible TTL and priority settings
  - Global cache configuration
  - Cache rule status management and sorting
  - Cache cleaning (purge) operations with task tracking
  - Cache preheating (preload) operations
  - Support for URL and directory-based cache operations

- **Certificate Management**: SSL/TLS certificate lifecycle for SCDN
  - Certificate creation, update, and deletion
  - Certificate export functionality
  - Certificate application and binding to domains
  - Query certificates by domains
  - Support for multiple certificate types

- **Origin Group Management**: Advanced origin server grouping
  - Origin group creation and configuration
  - Domain binding to origin groups
  - Origin group domain copy functionality
  - Load balancing and failover configuration

- **Network Speed Optimization**: Intelligent network acceleration
  - Network speed configuration management
  - Network speed rules with priority control
  - Rule sorting and optimization

- **Security Protection**: Enterprise-grade security features
  - DDoS protection configuration
  - WAF (Web Application Firewall) configuration
  - Security protection templates
  - Template-based batch configuration
  - Domain binding to security templates
  - IOTA (Intelligent Origin Traffic Analysis) support
  - Member global template management

- **Template Management**: Rule template system
  - Rule template creation and management
  - Template domain binding and unbinding
  - Template-based configuration inheritance

- **Log Download**: Access log management
  - Log download task creation and management
  - Log download template configuration
  - Template status management
  - Log field querying

- **SCDN Resources** (31 resources):
  - `edgenext_scdn_domain`: SCDN domain management
  - `edgenext_scdn_origin`: Origin server management
  - `edgenext_scdn_certificate`: Certificate management
  - `edgenext_scdn_cert_binding`: Certificate binding
  - `edgenext_scdn_certificate_apply`: Certificate application
  - `edgenext_scdn_domain_status`: Domain status management
  - `edgenext_scdn_domain_access_mode`: Access mode switching
  - `edgenext_scdn_domain_base_settings`: Base settings configuration
  - `edgenext_scdn_domain_node_switch`: Node switching
  - `edgenext_scdn_cache_rule`: Cache rule management
  - `edgenext_scdn_cache_rule_status`: Cache rule status
  - `edgenext_scdn_cache_rules_sort`: Cache rules sorting
  - `edgenext_scdn_cache_clean_task`: Cache cleaning tasks
  - `edgenext_scdn_cache_preheat_task`: Cache preheating tasks
  - `edgenext_scdn_origin_group`: Origin group management
  - `edgenext_scdn_origin_group_domain_bind`: Origin group domain binding
  - `edgenext_scdn_origin_group_domain_copy`: Origin group domain copy
  - `edgenext_scdn_network_speed_config`: Network speed configuration
  - `edgenext_scdn_network_speed_rule`: Network speed rules
  - `edgenext_scdn_network_speed_rules_sort`: Network speed rules sorting
  - `edgenext_scdn_security_protection_ddos_config`: DDoS protection
  - `edgenext_scdn_security_protection_waf_config`: WAF configuration
  - `edgenext_scdn_security_protection_template`: Security protection templates
  - `edgenext_scdn_security_protection_template_batch_config`: Batch security configuration
  - `edgenext_scdn_security_protection_template_domain_bind`: Security template domain binding
  - `edgenext_scdn_rule_template`: Rule template management
  - `edgenext_scdn_rule_template_domain_bind`: Rule template domain binding
  - `edgenext_scdn_rule_template_domain_unbind`: Rule template domain unbinding
  - `edgenext_scdn_log_download_task`: Log download tasks
  - `edgenext_scdn_log_download_template`: Log download templates
  - `edgenext_scdn_log_download_template_status`: Log template status

- **SCDN Data Sources** (35 data sources):
  - `data.edgenext_scdn_domain`: Query domain configuration
  - `data.edgenext_scdn_domains`: Query multiple domains
  - `data.edgenext_scdn_brief_domains`: Query brief domain list
  - `data.edgenext_scdn_origin`: Query origin configuration
  - `data.edgenext_scdn_origins`: Query multiple origins
  - `data.edgenext_scdn_access_progress`: Query domain access progress
  - `data.edgenext_scdn_domain_base_settings`: Query domain base settings
  - `data.edgenext_scdn_domain_templates`: Query domain templates
  - `data.edgenext_scdn_certificate`: Query certificate details
  - `data.edgenext_scdn_certificates`: Query multiple certificates
  - `data.edgenext_scdn_certificates_by_domains`: Query certificates by domains
  - `data.edgenext_scdn_certificate_export`: Export certificate
  - `data.edgenext_scdn_cache_rules`: Query cache rules
  - `data.edgenext_scdn_cache_global_config`: Query global cache configuration
  - `data.edgenext_scdn_cache_clean_config`: Query cache clean configuration
  - `data.edgenext_scdn_cache_clean_tasks`: Query cache clean tasks
  - `data.edgenext_scdn_cache_clean_task_detail`: Query cache clean task details
  - `data.edgenext_scdn_cache_preheat_tasks`: Query cache preheat tasks
  - `data.edgenext_scdn_origin_group`: Query origin group
  - `data.edgenext_scdn_origin_groups`: Query multiple origin groups
  - `data.edgenext_scdn_origin_groups_all`: Query all origin groups
  - `data.edgenext_scdn_network_speed_config`: Query network speed configuration
  - `data.edgenext_scdn_network_speed_rules`: Query network speed rules
  - `data.edgenext_scdn_security_protection_ddos_config`: Query DDoS protection config
  - `data.edgenext_scdn_security_protection_waf_config`: Query WAF configuration
  - `data.edgenext_scdn_security_protection_template`: Query security protection template
  - `data.edgenext_scdn_security_protection_templates`: Query multiple templates
  - `data.edgenext_scdn_security_protection_template_domains`: Query template domains
  - `data.edgenext_scdn_security_protection_template_unbound_domains`: Query unbound domains
  - `data.edgenext_scdn_security_protection_iota`: Query IOTA configuration
  - `data.edgenext_scdn_security_protection_member_global_template`: Query member global template
  - `data.edgenext_scdn_rule_template`: Query rule template
  - `data.edgenext_scdn_rule_templates`: Query multiple rule templates
  - `data.edgenext_scdn_rule_template_domains`: Query rule template domains
  - `data.edgenext_scdn_log_download_tasks`: Query log download tasks
  - `data.edgenext_scdn_log_download_templates`: Query log download templates
  - `data.edgenext_scdn_log_download_fields`: Query log download fields

#### OSS (Object Storage Service) Module
- **Bucket Management**: Complete S3-compatible bucket lifecycle
  - Bucket creation, configuration, and deletion
  - ACL (Access Control List) management
  - Automatic CORS configuration
  - Force destroy option for bucket cleanup

- **Object Management**: Comprehensive object operations
  - Object upload from file or inline content
  - Object download and metadata retrieval
  - Object deletion
  - Custom metadata support
  - HTTP headers configuration (cache-control, content-type, content-encoding, etc.)
  - ETag verification for integrity

- **Object Copy**: Advanced object copying capabilities
  - Cross-bucket object copying
  - Same-bucket object duplication
  - Metadata preservation or replacement
  - ACL control during copy operations

- **OSS Resources** (3 resources):
  - `edgenext_oss_bucket`: Bucket management
  - `edgenext_oss_object`: Object management
  - `edgenext_oss_object_copy`: Object copy operations

- **OSS Data Sources** (3 data sources):
  - `data.edgenext_oss_buckets`: Query bucket list
  - `data.edgenext_oss_object`: Query object details
  - `data.edgenext_oss_objects`: Query object list

### Changed
- **Documentation Improvements**: Enhanced documentation generation and formatting
  - Automatic directory creation in documentation generator
  - Improved documentation structure and organization
  - Complete English documentation for all modules
  - Removed Chinese comments in favor of English documentation

- **Type Handling**: Improved API response type handling
  - Added `FlexibleInt` type to handle inconsistent API response types (string/number)
  - Better type conversion and validation

### Technical Improvements
- **SCDN Client**: New dedicated SCDN API client implementation
- **OSS Client**: S3-compatible client using AWS SDK v2
- **Enhanced Error Handling**: More detailed error messages and validation
- **Improved State Management**: Better state tracking for complex resources

## [1.0.0] - 2025-09-11

### Added
- **CDN Domain Configuration Management**: Complete lifecycle management for CDN domains
  - Create, read, update, delete CDN domains
  - Support for multiple domain types (page, download, video_demand, dynamic)
  - Support for multiple acceleration areas (mainland_china, outside_mainland_china, global)
  
- **Comprehensive CDN Configuration**: 57+ configuration types supported
  - Origin server configuration with failover support
  - Cache rules with flexible TTL settings
  - Referer whitelist/blacklist protection
  - IP whitelist/blacklist access control
  - HTTPS/SSL configuration
  - Custom response headers
  - Error page customization
  - Rate limiting and speed control
  - Geographic access restrictions
  - And many more...

- **Smart Incremental Updates**: Intelligent configuration management
  - Only updates changed configuration items
  - Preserves existing configurations when not explicitly managed
  - Prevents unnecessary API calls and service disruptions

- **Cache Management**: URL and directory cache operations
  - Cache purging with status tracking
  - URL preloading (push) with batch support
  - Support for file extensions and directory patterns

- **SSL Certificate Management**: Complete certificate lifecycle
  - RSA and ECC certificate support
  - Secure private key handling
  - Certificate validation and status tracking

- **Data Sources**: Query existing resources
  - Domain configuration queries
  - Cache operation status queries
  - Certificate information retrieval

- **Comprehensive Documentation**: Full resource documentation
  - Detailed parameter descriptions
  - Usage examples for all resources
  - Best practices and troubleshooting guides

### Technical Features
- **Type-Safe API Conversion**: Robust bidirectional data conversion between Terraform schema and API formats
- **Error Handling**: Comprehensive error handling with detailed messages
- **State Management**: Accurate state tracking with drift detection
- **Field Mapping**: Automatic handling of field name differences between Terraform and API
- **Default Value Management**: Smart handling of optional parameters and API defaults

### Resources
- `edgenext_cdn_domain`: Complete CDN domain and configuration management
- `edgenext_cdn_purge`: Cache purging operations
- `edgenext_cdn_push`: URL preloading operations
- `edgenext_ssl_certificate`: SSL certificate management

### Data Sources
- `data.edgenext_cdn_domain`: Query domain configurations
- `data.edgenext_cdn_domains`: Query multiple domain configurations
- `data.edgenext_cdn_purge`: Query purge operation status
- `data.edgenext_cdn_purges`: Query multiple purge operations
- `data.edgenext_cdn_push`: Query push operation status
- `data.edgenext_cdn_pushes`: Query multiple push operations
- `data.edgenext_ssl_certificate`: Query certificate information
- `data.edgenext_ssl_certificates`: Query multiple certificates

[Unreleased]: https://github.com/edgenextapisdk/terraform-provider-edgenext/compare/v1.2.1...HEAD
[1.2.1]: https://github.com/edgenextapisdk/terraform-provider-edgenext/compare/v1.2.0...v1.2.1
[1.2.0]: https://github.com/edgenextapisdk/terraform-provider-edgenext/compare/v1.0.0...v1.2.0
[1.0.0]: https://github.com/edgenextapisdk/terraform-provider-edgenext/releases/tag/v1.0.0
