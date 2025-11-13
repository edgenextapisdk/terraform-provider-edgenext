Use this data source to query SCDN security protection DDoS configuration.

Example Usage

Query DDoS config

```hcl
data "edgenext_scdn_security_protection_ddos_config" "example" {
  business_id = 12345
}

output "ddos_config" {
  value = data.edgenext_scdn_security_protection_ddos_config.example.application_ddos_protection
}
```

Query specific config keys

```hcl
data "edgenext_scdn_security_protection_ddos_config" "example" {
  business_id = 12345
  keys        = ["application_ddos_protection", "network_ddos_protection"]
}
```

Query and save to file

```hcl
data "edgenext_scdn_security_protection_ddos_config" "example" {
  business_id        = 12345
  result_output_file = "ddos_config.json"
}
```

