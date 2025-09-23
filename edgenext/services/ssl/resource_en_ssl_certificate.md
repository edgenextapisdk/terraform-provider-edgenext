Provides a resource to create and manage SSL certificates.

Example Usage

Basic SSL certificate upload

```hcl
resource "edgenext_ssl_certificate" "example" {
  name        = "example-com-cert"
  certificate = file("path/to/certificate.crt")
  key         = file("path/to/private.key")
}
```

SSL certificate with certificate content

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

Import

SSL certificates can be imported using the certificate ID:

```shell
terraform import edgenext_ssl_certificate.example ssl-cert-123456
```
