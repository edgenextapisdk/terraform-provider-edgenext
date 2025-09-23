# EdgeNext SSL Certificate Services

This package provides Terraform resources and data sources for managing EdgeNext SSL certificates.

## Resources

### SSL Certificate
- **Resource**: `edgenext_ssl_certificate` (`ResourceEdgenextSslCertificate`)
- **File**: `resource_en_ssl_certificate.go`
- **Description**: Manage SSL certificates including upload, configuration, and lifecycle management
- **Key Features**:
  - Upload SSL certificates and private keys
  - Certificate validation and verification
  - Domain association management
  - Certificate lifecycle tracking

## Data Sources

### SSL Certificate
- **Data Source**: `edgenext_ssl_certificate` (`DataSourceEdgenextSslCertificate`)
- **File**: `data_source_en_ssl_certificate.go`
- **Description**: Query SSL certificate details by certificate ID
- **Key Features**:
  - Retrieve certificate information
  - View associated domains
  - Export certificate details to file

### SSL Certificates List
- **Data Source**: `edgenext_ssl_certificates` (`DataSourceEdgenextSslCertificates`)
- **File**: `data_source_en_ssl_certificate.go`
- **Description**: Query a list of SSL certificates
- **Key Features**:
  - List all SSL certificates
  - Pagination support
  - Filter and search capabilities

## File Structure

```
edgenext/services/ssl/
├── README.md                              # This documentation
├── service_en_ssl_certificate.go          # SSL service client implementation
├── service_en_ssl_certificate_test.go     # SSL service tests
├── resource_en_ssl_certificate.go         # SSL certificate resource implementation
├── resource_en_ssl_certificate.md         # SSL certificate resource documentation
├── data_source_en_ssl_certificate.go      # SSL certificate data sources implementation
├── data_source_en_ssl_certificate.md      # SSL certificate data source documentation
└── data_source_en_ssl_certificates.md     # SSL certificates list data source documentation
```

## Usage Examples

### Basic SSL Certificate Upload

```hcl
resource "edgenext_ssl_certificate" "example" {
  name        = "example-com-cert"
  certificate = file("path/to/certificate.crt")
  key         = file("path/to/private.key")
}
```

### SSL Certificate with Certificate Content

```hcl
resource "edgenext_ssl_certificate" "example" {
  name = "example-com-cert"
  
  certificate = <<-EOT
-----BEGIN CERTIFICATE-----
MIIDXTCCAkWgAwIBAgIJAJC1HiIAZAiIMA0GCSqGSIb3DQEBBQUAMEUxCzAJBgNV
BAYTAkFVMRMwEQYDVQQIDApTb21lLVN0YXRlMSEwHwYDVQQKDBhJbnRlcm5ldCBX
...
-----END CERTIFICATE-----
EOT

  key = <<-EOT
-----BEGIN PRIVATE KEY-----
MIIEvQIBADANBgkqhkiG9w0BAQEFAASCBKcwggSjAgEAAoIBAQC7S2IgnLVkjpQR
RIiTq86f1o3d6nOF4eU4h95tX3YfL8s6eSgfDNF2nDG8VQZ4m8Sv1JHbYDrDJ8Ac
...
-----END PRIVATE KEY-----
EOT
}
```

### Query SSL Certificate by ID

```hcl
data "edgenext_ssl_certificate" "example" {
  cert_id = "ssl-cert-123456"
  result_output_file = "ssl_cert.json"
}
```

### Query SSL Certificates List

```hcl
data "edgenext_ssl_certificates" "example" {
  page_number = 1
  page_size   = 100
  result_output_file = "ssl_certs.json"
}
```

## Certificate Management

### Certificate Format Requirements

SSL certificates must be in PEM format:
- Certificate: Standard X.509 certificate in PEM format
- Private Key: RSA or ECDSA private key in PEM format
- Certificate Chain: Include intermediate certificates if required

### Certificate Validation

The service performs automatic validation:
- Certificate format validation
- Private key matching
- Certificate chain verification
- Domain validation (if applicable)

### Security Features

- Private keys are stored securely and marked as sensitive
- Certificate content diff suppression for format variations
- Automatic detection of certificate changes

## API Integration

This package integrates with the EdgeNext SSL API to provide:

- **Certificate Upload**: Upload SSL certificates and private keys
- **Certificate Management**: Update, delete, and manage certificates
- **Certificate Query**: Retrieve certificate details and status
- **Domain Association**: Manage certificate-domain relationships

## Error Handling

All resources and data sources include comprehensive error handling for:
- Certificate format validation errors
- API connectivity issues
- Certificate not found scenarios
- Private key mismatch errors

## Documentation Generation

This package includes markdown documentation for each resource and data source that is automatically processed by the documentation generation tool to create Terraform Registry compatible documentation.

## Testing

Run tests with:
```bash
go test ./edgenext/services/ssl/...
```

## Related Services

- [CDN Services](../cdn/README.md) - CDN domain and cache management
