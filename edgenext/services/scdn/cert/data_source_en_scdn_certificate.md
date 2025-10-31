Use this data source to query details of a specific SCDN certificate.

Example Usage

Query certificate by ID

```hcl
data "edgenext_scdn_certificate" "example" {
  id = "12345"
}

output "certificate_name" {
  value = data.edgenext_scdn_certificate.example.ca_name
}

output "expiry_time" {
  value = data.edgenext_scdn_certificate.example.issuer_expiry_time
}
```

Query and save to file

```hcl
data "edgenext_scdn_certificate" "example" {
  id               = "12345"
  result_output_file = "certificate.json"
}
```

