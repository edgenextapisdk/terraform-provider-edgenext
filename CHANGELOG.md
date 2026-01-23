# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

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
- `edgenext_cdn_prefetch`: Cache prefetch operations
- `edgenext_cdn_purge`: Cache purging operations
- `edgenext_ssl_certificate`: SSL certificate management

### Data Sources
- `data.edgenext_cdn_domain`: Query domain configurations
- `data.edgenext_cdn_domains`: Query multiple domain configurations
- `data.edgenext_cdn_prefetch`: Query prefetch operation status
- `data.edgenext_cdn_prefetches`: Query multiple prefetch operations
- `data.edgenext_cdn_purge`: Query purge operation status
- `data.edgenext_cdn_purges`: Query multiple purge operations
- `data.edgenext_ssl_certificate`: Query certificate information
- `data.edgenext_ssl_certificates`: Query multiple certificates

[Unreleased]: https://github.com/edgenextapisdk/terraform-provider-edgenext/compare/v1.0.0...HEAD
[1.0.0]: https://github.com/edgenextapisdk/terraform-provider-edgenext/releases/tag/v1.0.0
