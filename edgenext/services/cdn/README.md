# EdgeNext CDN Services

This package provides Terraform resources and data sources for managing EdgeNext CDN (Content Delivery Network) services.

## Resources

### CDN Domain Configuration
- **Resource**: `edgenext_cdn_domain` (`ResourceEdgenextCdnDomainConfig`)
- **File**: `resource_en_cdn_domain.go`
- **Description**: Manage CDN domain configurations including origin settings, cache rules, HTTPS configuration, and various CDN features

### CDN Cache Push
- **Resource**: `edgenext_cdn_push` (`ResourceEdgenextCdnPush`)
- **File**: `resource_en_cdn_push.go`  
- **Description**: Manage CDN cache push tasks to pre-populate cache with content

### CDN Cache Purge
- **Resource**: `edgenext_cdn_purge` (`ResourceEdgenextCdnPurge`)
- **File**: `resource_en_cdn_purge.go`
- **Description**: Manage CDN cache purge tasks to invalidate cached content

## Data Sources

### CDN Domain Configuration
- **Data Source**: `edgenext_cdn_domain` (`DataSourceEdgenextCdnDomainConfig`)
- **File**: `data_source_en_cdn_domain.go`
- **Description**: Query CDN domain configuration details

### CDN Domains List
- **Data Source**: `edgenext_cdn_domains` (`DataSourceEdgenextCdnDomains`)
- **File**: `data_source_en_cdn_domain.go`
- **Description**: Query a list of CDN domains

### CDN Cache Push
- **Data Source**: `edgenext_cdn_push` (`DataSourceEdgenextCdnPush`)
- **File**: `data_source_en_cdn_push.go`
- **Description**: Query CDN cache push task details

### CDN Cache Push Tasks
- **Data Source**: `edgenext_cdn_pushes` (`DataSourceEdgenextCdnPushes`) 
- **File**: `data_source_en_cdn_push.go`
- **Description**: Query a list of CDN cache push tasks

### CDN Cache Purge
- **Data Source**: `edgenext_cdn_purge` (`DataSourceEdgenextCdnPurge`)
- **File**: `data_source_en_cdn_purge.go`
- **Description**: Query CDN cache purge task details

### CDN Cache Purge Tasks
- **Data Source**: `edgenext_cdn_purges` (`DataSourceEdgenextCdnPurges`)
- **File**: `data_source_en_cdn_purge.go`
- **Description**: Query a list of CDN cache purge tasks

## File Structure

```
edgenext/services/cdn/
├── README.md                           # This documentation
├── service_en_cdn.go                   # CDN service client implementation
├── service_en_cdn_test.go              # CDN service tests
├── resource_en_cdn_domain.go           # CDN domain resource implementation
├── resource_en_cdn_domain.md           # CDN domain resource documentation
├── data_source_en_cdn_domain.go        # CDN domain data sources implementation
├── data_source_en_cdn_domain.md        # CDN domain data source documentation
├── data_source_en_cdn_domains.md       # CDN domains list data source documentation
├── resource_en_cdn_push.go             # CDN push resource implementation
├── resource_en_cdn_push.md             # CDN push resource documentation
├── data_source_en_cdn_push.go          # CDN push data sources implementation
├── data_source_en_cdn_push.md          # CDN push data source documentation
├── data_source_en_cdn_pushes.md        # CDN push tasks data source documentation
├── resource_en_cdn_purge.go            # CDN purge resource implementation
├── resource_en_cdn_purge.md            # CDN purge resource documentation
├── data_source_en_cdn_purge.go         # CDN purge data sources implementation
├── data_source_en_cdn_purge.md         # CDN purge data source documentation
└── data_source_en_cdn_purges.md        # CDN purge tasks data source documentation
```

## Usage Examples

### Basic CDN Domain Configuration

```hcl
resource "edgenext_cdn_domain" "example" {
  domain = "example.com"
  area   = "global"
  type   = "page"
  
  config {
    origin {
      default_master = "origin.example.com"
      origin_mode    = "default"
    }
  }
}
```

### CDN Cache Push

```hcl
resource "edgenext_cdn_push" "example" {
  type = "url"
  urls = [
    "https://example.com/static/image1.jpg",
    "https://example.com/static/image2.jpg"
  ]
}
```

### CDN Cache Purge

```hcl
resource "edgenext_cdn_purge" "example" {
  urls = [
    "https://example.com/static/old-image.jpg",
    "https://example.com/static/old-style.css"
  ]
}
```

### Query CDN Domain Configuration

```hcl
data "edgenext_cdn_domain" "example" {
  domain = "example.com"
  result_output_file = "domain_config.json"
}
```

## Related Services

- [SSL Services](../ssl/README.md) - SSL certificate management
