Use this data source to query a list of SSL certificates.

Example Usage

Query all SSL certificates

```hcl
data "edgenext_ssl_certificates" "example" {
  page_number = 1
  page_size   = 100
  output_file = "ssl_certs.json"
}
```

Query SSL certificates with pagination

```hcl
data "edgenext_ssl_certificates" "example" {
  page_number = 1
  page_size   = 50
}
```
