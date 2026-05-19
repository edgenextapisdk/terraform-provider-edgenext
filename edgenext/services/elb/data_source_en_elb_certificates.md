Use this data source to list EdgeNext ELB TLS certificates.

Example Usage

```hcl
data "edgenext_elb_certificates" "all" {
  page_num  = 1
  page_size = 50
  name      = ""
}

output "certificate_refs" {
  value = [for c in data.edgenext_elb_certificates.all.certificates : c.certificate_ref]
}
```
