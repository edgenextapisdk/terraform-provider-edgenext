Use this data source to query the member global SCDN security protection template.

Example Usage

Query member global template

```hcl
data "edgenext_scdn_security_protection_member_global_template" "example" {
}

output "template" {
  value = data.edgenext_scdn_security_protection_member_global_template.example.template
}
```

Query and save to file

```hcl
data "edgenext_scdn_security_protection_member_global_template" "example" {
  result_output_file = "member_global_template.json"
}
```

