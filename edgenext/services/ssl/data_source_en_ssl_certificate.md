Use this data source to query detailed information of SSL certificate.

Example Usage

Query SSL certificate by certificate ID

```hcl
data "edgenext_ssl_certificate" "example" {
  cert_id = "ssl-cert-123456"
  output_file = "ssl_cert.json"
}
```

Query SSL certificate with minimal information

```hcl
data "edgenext_ssl_certificate" "example" {
  cert_id = "ssl-cert-123456"
}
```
