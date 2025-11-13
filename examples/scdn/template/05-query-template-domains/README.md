# Query Domains Bound to Rule Template

This example demonstrates how to query domains that are bound to a specific rule template.

## Usage

1. Copy `terraform.tfvars.example` to `terraform.tfvars` and fill in your credentials:

```bash
cp terraform.tfvars.example terraform.tfvars
```

2. Edit `terraform.tfvars` with your actual values:

```hcl
access_key  = "your_access_key"
secret_key  = "your_secret_key"
template_id = 123
app_type    = "network_speed"
page        = 1
page_size   = 10
domain      = ""  # Optional: filter by domain name
```

3. Initialize and apply:

```bash
terraform init
terraform plan
terraform apply
```

## Outputs

- `total`: Total number of domains bound to the template
- `domains`: List of domain information (ID, domain name, binding timestamp)

## Notes

- `template_id` and `app_type` are required
- `domain` filter is optional
- Supports pagination with `page` and `page_size`

