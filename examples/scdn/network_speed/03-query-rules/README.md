# Query Network Speed Rules

This example demonstrates how to query network speed rules for a template.

## Usage

1. Copy `terraform.tfvars.example` to `terraform.tfvars` and fill in your credentials:

```bash
cp terraform.tfvars.example terraform.tfvars
```

2. Edit `terraform.tfvars` with your actual values:

```hcl
access_key   = "your_access_key"
secret_key   = "your_secret_key"
business_id  = 1246
business_type = "tpl"
config_group  = "upstream_uri_change_rule"
```

3. Initialize and apply:

```bash
terraform init
terraform plan
terraform apply
```

## Notes

- This example queries network speed rules for a template
- `business_id` is the template ID for `business_type = "tpl"`
- `business_id` is the user ID for `business_type = "global"`
- `config_group` must be one of:
  - `custom_page`
  - `upstream_uri_change_rule`
  - `resp_headers_rule`
  - `customized_req_headers_rule`

## Outputs

- `total`: Total number of rules
- `rules`: List of network speed rules with all their details

