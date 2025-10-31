Use this data source to query certificates bound to specific domains.

Example Usage

Query certificates by domains

```hcl
data "edgenext_scdn_certificates_by_domains" "example" {
  domains = ["example.com", "www.example.com"]
}

output "certificates" {
  value = data.edgenext_scdn_certificates_by_domains.example.certificates
}
```

Query and save to file

```hcl
data "edgenext_scdn_certificates_by_domains" "example" {
  domains            = ["example.com"]
  result_output_file = "certificates_by_domains.json"
}
```

