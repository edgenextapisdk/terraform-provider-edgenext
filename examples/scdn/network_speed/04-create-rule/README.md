# Create Network Speed Rule

This example demonstrates how to create a network speed rule.

## Usage

1. Copy `terraform.tfvars.example` to `terraform.tfvars` and fill in your credentials:

```bash
cp terraform.tfvars.example terraform.tfvars
```

2. Edit `terraform.tfvars` with your actual values:

```hcl
access_key   = "your_access_key"
secret_key   = "your_secret_key"
business_id  = 777
business_type = "tpl"
config_group  = "customized_req_headers_rule"
```

3. Initialize and apply:

```bash
terraform init
terraform plan
terraform apply
```

## Notes

- This example creates a customized request headers rule
- You can create different rule types by changing `config_group` and the corresponding rule block
- Supported rule types:
  - `custom_page`: Custom error page
  - `upstream_uri_change_rule`: Upstream URI rewrite
  - `resp_headers_rule`: Custom response headers
  - `customized_req_headers_rule`: Custom request headers

## Outputs

- `rule_id`: The ID of the created rule


