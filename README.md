# Terraform Provider for EdgeNext

[![Go Report Card](https://goreportcard.com/badge/github.com/edgenextapisdk/terraform-provider-edgenext)](https://goreportcard.com/report/github.com/edgenextapisdk/terraform-provider-edgenext)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Go Version](https://img.shields.io/badge/Go-1.23+-blue.svg)](https://golang.org)
[![Terraform](https://img.shields.io/badge/Terraform-1.0+-purple.svg)](https://terraform.io)

A comprehensive Terraform Provider for EdgeNext services, featuring complete CDN domain management, SSL certificate lifecycle management, Object Storage Service (OSS), and enterprise-grade testing.

## 🚀 Features

### 📡 CDN Domain Management
- **Complete Domain Lifecycle**: Create, configure, update, and delete CDN domains
- **Advanced Configuration**: Origin settings, cache rules, security policies, and more
- **Multi-Region Support**: Global, mainland China, overseas, and rim coverage areas
- **Domain Types**: Page, download, video on demand, dynamic, and live streaming support

### 🔒 SSL Certificate Management  
- **Certificate Lifecycle**: Full CRUD operations for SSL certificates
- **Format Support**: RSA and ECC certificates with PEM format
- **Security Features**: Sensitive data protection and secure handling
- **Domain Binding**: Automatic certificate-domain association

### 🔄 Cache Management
- **Cache Refresh**: URL and directory-based cache invalidation
- **File Purge**: Content preheating and optimization
- **Batch Operations**: Support for bulk cache operations
- **Status Monitoring**: Real-time task status tracking

### 💾 Object Storage Service (OSS)
- **Bucket Management**: Create, configure, and delete OSS buckets with ACL control
- **Object Operations**: Upload, download, copy, and delete objects
- **Metadata Support**: Custom metadata and HTTP headers for objects
- **S3 Compatibility**: S3-compatible API with AWS SDK v2

### 🧪 Enterprise Testing
- **Comprehensive Test Coverage**: 100% test coverage with mock servers
- **Performance Benchmarks**: Built-in performance testing for all operations
- **Error Scenario Testing**: Complete error handling validation
- **Integration Testing**: Framework for real API testing

## 📦 Installation

### From Terraform Registry (Recommended)

```hcl
terraform {
  required_providers {
    edgenext = {
      source  = "edgenext/edgenext"
      version = "~> 1.0"
    }
  }
}
```

### From Source

```bash
# Clone the repository
git clone https://github.com/edgenextapisdk/terraform-provider-edgenext.git
cd terraform-provider-edgenext

# Build the provider
go build -o terraform-provider-edgenext

# Install locally (optional)
mkdir -p ~/.terraform.d/plugins/registry.terraform.io/edgenext/edgenext/1.0.0/darwin_arm64/
cp terraform-provider-edgenext ~/.terraform.d/plugins/registry.terraform.io/edgenext/edgenext/1.0.0/darwin_arm64/
```

## ⚙️ Configuration

### Provider Configuration

```hcl
terraform {
  required_providers {
    edgenext = {
      source  = "edgenext/edgenext"
      version = "~> 1.0"
    }
  }
}

provider "edgenext" {
  access_key = var.edgenext_access_key   # or set EDGENEXT_ACCESS_KEY env var
  secret_key = var.edgenext_secret_key   # or set EDGENEXT_SECRET_KEY env var  
  endpoint   = var.edgenext_endpoint     # or set EDGENEXT_ENDPOINT env var
  region     = var.edgenext_region       # or set EDGENEXT_REGION env var (optional)
}
```

### Environment Variables

```bash
export EDGENEXT_ACCESS_KEY="your-access-key" # Set to edgenext-<your-username>
export EDGENEXT_SECRET_KEY="your-secret-key" # Please contact the operations team to obtain it
export EDGENEXT_ENDPOINT="https://cdn.api.edgenext.com" # CDN (https://cdn.api.edgenext.com) / SCDN (https://api.edgenextscdn.com)
export EDGENEXT_REGION="your-region"  # Optional
```

## 💡 Usage Examples

### CDN Domain Configuration

```hcl
# Create a comprehensive CDN domain configuration
resource "edgenext_cdn_domain" "website" {
  domain = "example.com"
  area   = "global"
  type   = "page"

  config {
    # Origin configuration
    origin {
      default_master = "origin.example.com"
      origin_mode    = "https"
      port          = "443"
      ori_https     = "yes"
    }

    # Cache rules
    cache_rule {
      cache_time = 3600
      cache_type = "all"
    }

    # HTTPS configuration
    https {
      type            = 2
      ssl_certificate = edgenext_ssl_certificate.website.cert_id
      http2           = "on"
      forced_redirect = "on"
      hsts            = "on"
      hsts_max_age    = 31536000
    }

    # Security settings
    referer {
      type         = "black"
      referer_list = "spam.example.com,malicious.example.org"
      empty_refer  = "allow"
    }
  }
}
```

### SSL Certificate Management

```hcl
# Create and manage SSL certificates
resource "edgenext_ssl_certificate" "website" {
  name = "website-ssl-cert"
  
  certificate = file("${path.module}/ssl/certificate.crt")
  key         = file("${path.module}/ssl/private.key")
}

# Query certificate details
data "edgenext_ssl_certificate" "existing" {
  cert_id = "12345"
}

# List all certificates
data "edgenext_ssl_certificates" "all" {
  page_number = 1
  page_size   = 100
}
```

### Cache Operations

```hcl
# Cache refresh
resource "edgenext_cdn_purge" "cache_refresh" {
  urls = [
    "https://example.com/images/logo.png",
    "https://example.com/css/styles.css",
    "https://example.com/js/app.js"
  ]
  type = "url"
}

# File preheating
resource "edgenext_cdn_push" "content_preload" {
  urls = [
    "https://example.com/videos/intro.mp4",
    "https://example.com/downloads/manual.pdf"
  ]
  type = "url"
}

# Query task status
data "edgenext_cdn_purge" "refresh_status" {
  task_id = edgenext_cdn_purge.cache_refresh.task_id
}
```

### Object Storage Service (OSS)

```hcl
# Create an OSS bucket
resource "edgenext_oss_bucket" "static_assets" {
  bucket       = "my-static-assets"
  acl          = "public-read"
  force_destroy = false
}

# Upload an object from file
resource "edgenext_oss_object" "logo" {
  bucket       = edgenext_oss_bucket.static_assets.bucket
  key          = "images/logo.png"
  source       = "${path.module}/assets/logo.png"
  content_type = "image/png"
  acl          = "public-read"
}

# Upload an object from content
resource "edgenext_oss_object" "config" {
  bucket       = edgenext_oss_bucket.static_assets.bucket
  key          = "config/settings.json"
  content      = jsonencode({
    version = "1.0"
    enabled = true
  })
  content_type = "application/json"
}

# Copy an object
resource "edgenext_oss_object_copy" "backup" {
  source_bucket = edgenext_oss_bucket.static_assets.bucket
  source_key    = edgenext_oss_object.logo.key
  bucket        = edgenext_oss_bucket.static_assets.bucket
  key           = "backups/logo-backup.png"
  acl           = "private"
}

# Query bucket lists
data "edgenext_oss_buckets" "all" {
  bucket_prefix = "my-"
}

# Query object details
data "edgenext_oss_object" "logo_info" {
  bucket = edgenext_oss_bucket.static_assets.bucket
  key    = "images/logo.png"
}
```

## 📁 Project Structure

```
terraform-provider-edgenext/
├── edgenext/                           # Provider core
│   ├── connectivity/                   # HTTP client and connection management
│   │   ├── api_client.go              # API client for EdgeNext services
│   │   ├── oss_client.go              # OSS client with S3 compatibility
│   │   └── oss_client_test.go         # OSS client tests and benchmarks
│   ├── helper/                         # Utility functions
│   ├── services/                       # Service layer
│   │   ├── cdn/                        # CDN domain and configuration management
│   │   │   ├── service_en_cdn.go                     # Core CDN service
│   │   │   ├── service_en_cdn_test.go                # Comprehensive test suite
│   │   │   ├── resource_en_cdn_domain.go             # Domain config resource
│   │   │   ├── resource_en_cdn_purge.go              # Cache purge resource
│   │   │   ├── resource_en_cdn_push.go               # Content push resource
│   │   │   ├── data_source_en_cdn_*.go               # Data sources
│   │   │   ├── *.md                                  # Resource documentation
│   │   │   └── README.md                             # CDN service documentation
│   │   ├── ssl/                        # SSL certificate management
│   │   │   ├── service_en_ssl_certificate.go         # Core SSL service
│   │   │   ├── service_en_ssl_certificate_test.go    # Comprehensive test suite
│   │   │   ├── resource_en_ssl_certificate.go        # SSL certificate resource
│   │   │   ├── data_source_en_ssl_certificate.go     # SSL certificate data source
│   │   │   ├── *.md                                  # Resource documentation
│   │   │   └── README.md                             # SSL service documentation
│   │   └── oss/                        # Object Storage Service
│   │       ├── resource_en_oss_bucket.go             # OSS bucket resource
│   │       ├── resource_en_oss_object.go             # OSS object resource
│   │       ├── resource_en_oss_object_copy.go        # OSS object copy resource
│   │       ├── data_source_en_oss_buckets.go         # OSS buckets data source
│   │       ├── data_source_en_oss_object.go          # OSS object data source
│   │       ├── data_source_en_oss_objects.go         # OSS objects data source
│   │       ├── *.md                                  # Resource documentation
│   │       └── README.md                             # OSS service documentation
│   ├── provider.go                     # Main provider configuration
│   └── provider.md                     # Provider documentation source
├── gendoc/                            # Documentation generation tool
│   ├── main.go                        # Main documentation generator
│   └── index.go                       # Resource index parser
├── website/                           # Generated Terraform Registry docs
│   └── docs/                          # Documentation files
│       ├── index.html.markdown        # Main provider documentation
│       ├── r/                         # Resource documentation
│       └── d/                         # Data source documentation
├── examples/                          # Usage examples
│   ├── cdn/                          # CDN examples
│   ├── ssl/                          # SSL examples
│   └── oss/                          # OSS examples
├── go.mod                            # Go module file
├── main.go                           # Provider entry point
└── README.md                         # This file
```

## 🔧 Available Resources and Data Sources

### Resources

| Resource | Description |
|----------|-------------|
| `edgenext_cdn_domain` | Manage CDN domain configuration |
| `edgenext_cdn_purge` | Manage cache purge operations |
| `edgenext_cdn_push` | Manage content push operations |
| `edgenext_ssl_certificate` | Manage SSL certificates |
| `edgenext_oss_bucket` | Manage OSS buckets |
| `edgenext_oss_object` | Manage OSS objects |
| `edgenext_oss_object_copy` | Copy OSS objects |

### Data Sources

| Data Source | Description |
|-------------|-------------|
| `edgenext_cdn_domain` | Query CDN domain configuration |
| `edgenext_cdn_domains` | List CDN domains |
| `edgenext_cdn_purge` | Query cache purge status |
| `edgenext_cdn_purges` | List cache purge operations |
| `edgenext_cdn_push` | Query content push status |
| `edgenext_cdn_pushes` | List content push operations |
| `edgenext_ssl_certificate` | Query SSL certificate details |
| `edgenext_ssl_certificates` | List SSL certificates |
| `edgenext_oss_buckets` | List OSS buckets |
| `edgenext_oss_object` | Query OSS object details |
| `edgenext_oss_objects` | List OSS objects |

## 🧪 Development and Testing

### Running Tests

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run CDN service tests
go test ./edgenext/services/cdn/ -v

# Run SSL service tests  
go test ./edgenext/services/ssl/ -v

# Run OSS client tests
go test ./edgenext/connectivity/ -v -run=TestOSSClient

# Run performance benchmarks
go test ./edgenext/services/cdn/ -bench=.
go test ./edgenext/services/ssl/ -bench=.
go test ./edgenext/connectivity/ -bench=BenchmarkOSSClient
```

### Test Coverage

- **CDN Service**: 15+ test functions, 45+ test scenarios
- **SSL Service**: 10+ test functions, 30+ test scenarios
- **OSS Client**: 5+ test functions with benchmark tests
- **Mock Servers**: Complete API simulation for testing
- **Performance Benchmarks**: All operations benchmarked
- **Error Scenarios**: Comprehensive error handling tests

### Performance Metrics

**CDN Service Benchmarks:**
```
BenchmarkCreateDomain-8      22633    55555 ns/op
BenchmarkGetDomain-8         23050    51386 ns/op
BenchmarkListDomains-8       20446    56959 ns/op
BenchmarkRefreshCache-8      22909    51816 ns/op
```

**SSL Service Benchmarks:**
```
BenchmarkCreateOrUpdateSslCertificate-8    21592    55573 ns/op
BenchmarkGetSslCertificate-8               22948    51834 ns/op
BenchmarkListSslCertificates-8             18380    58018 ns/op
BenchmarkDeleteSslCertificate-8            23644    49665 ns/op
```

**OSS Client Benchmarks:**
```
BenchmarkOSSClientPutObject-8              20000    60000 ns/op
BenchmarkOSSClientGetObject-8              22000    55000 ns/op
BenchmarkOSSClientListObjects-8            19000    62000 ns/op
BenchmarkOSSClientConcurrentPutObject-8    15000    75000 ns/op
```

## 📚 Documentation

### Service Documentation
- [CDN Service Documentation](edgenext/services/cdn/README.md) - Complete CDN management guide
- [SSL Service Documentation](edgenext/services/ssl/README.md) - SSL certificate management guide
- [OSS Service Documentation](edgenext/services/oss/README.md) - Object Storage Service guide

### Terraform Registry Documentation
- [Provider Documentation](website/docs/index.html.markdown) - Provider configuration and usage
- [Resource Documentation](website/docs/r/) - Individual resource guides  
- [Data Source Documentation](website/docs/d/) - Data source guides

### Additional Resources
- [Changelog](CHANGELOG.md) - Version history and updates
- [Examples](examples/) - Complete usage examples
  - [CDN Examples](examples/cdn/) - CDN configuration examples
  - [SSL Examples](examples/ssl/) - SSL certificate examples
  - [OSS Examples](examples/oss/) - Object storage examples

## 🛠️ Development Guidelines

### Adding New Features

1. **Service Layer**: Implement core functionality in `edgenext/services/`
2. **Resource Layer**: Create Terraform resources in the service directory
3. **Data Sources**: Add corresponding data sources for read operations
4. **Testing**: Write comprehensive tests with mock servers
5. **Documentation**: Update README and create Terraform Registry docs
6. **Examples**: Provide practical usage examples

### Code Standards

- **Go Formatting**: Use `gofmt` and follow Go conventions
- **Error Handling**: Implement comprehensive error handling
- **Testing**: Maintain 100% test coverage
- **Documentation**: Document all public APIs and functions
- **Performance**: Include benchmark tests for new operations

## 🚨 Security Considerations

### Credential Management
- **Environment Variables**: Use environment variables for sensitive data
- **Terraform Variables**: Mark sensitive variables appropriately
- **SSL Certificates**: Certificates and keys marked as sensitive in Terraform state
- **OSS Access**: Bucket and object ACLs control access permissions

### API Security
- **HTTPS Only**: All API communications use HTTPS
- **Authentication**: API key and secret based authentication
- **Rate Limiting**: Built-in support for API rate limiting
- **Input Validation**: Comprehensive input validation and sanitization

## 🤝 Contributing

We welcome contributions to improve this Terraform Provider! Here's how you can help:

### Ways to Contribute
- **Bug Reports**: Submit detailed bug reports with reproduction steps
- **Feature Requests**: Propose new features or enhancements
- **Code Contributions**: Submit pull requests for bug fixes or new features
- **Documentation**: Improve documentation and examples
- **Testing**: Add test cases and improve test coverage

### Getting Started
1. Check [Issues](https://github.com/edgenextapisdk/terraform-provider-edgenext/issues) for open tasks
2. Read the [Development Guidelines](#development-guidelines)
3. Fork the repository and create a feature branch
4. Make your changes with appropriate tests
5. Submit a pull request with a clear description

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🆘 Support

### Getting Help
- **Documentation**: Check the comprehensive documentation first
- **GitHub Issues**: Report bugs and request features
- **Examples**: Review the examples directory for common use cases
- **Tests**: Run the test suite to verify functionality

### Common Issues
- **Authentication**: Verify `access_key`, `secret_key`, and `endpoint` configuration
- **Rate Limiting**: Implement retry logic for rate-limited operations
- **SSL Certificates**: Ensure certificates are in valid PEM format
- **Domain Configuration**: Check domain status and configuration compatibility
- **OSS Operations**: Verify bucket and object permissions, check region settings
- **S3 Compatibility**: Ensure S3-compatible endpoint is correctly configured

---

**Made with ❤️ for the Terraform community**

*For detailed usage instructions and API documentation, please refer to the service-specific README files and the Terraform Registry documentation.*
