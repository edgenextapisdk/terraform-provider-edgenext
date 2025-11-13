Use this data source to query SCDN security protection IOTA (enum key-value pairs).

Example Usage

Query IOTA

```hcl
data "edgenext_scdn_security_protection_iota" "example" {
}

output "iota" {
  value = data.edgenext_scdn_security_protection_iota.example.iota
}
```

Query and save to file

```hcl
data "edgenext_scdn_security_protection_iota" "example" {
  result_output_file = "iota.json"
}
```

