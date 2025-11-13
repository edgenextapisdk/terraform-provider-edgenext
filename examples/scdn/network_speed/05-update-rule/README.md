# Update Network Speed Rule

This example demonstrates how to update an existing network speed rule.

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
config_group  = "customized_req_headers_rule"
rule_id       = 12345
```

3. Initialize and apply:

```bash
terraform init
terraform plan
terraform apply
```

## Notes

- This example updates an existing network speed rule
- `business_id` is the template ID for `business_type = "tpl"`
- `business_id` is the user ID for `business_type = "global"`
- `config_group` must match the rule's original config group
- `rule_id` is the ID of the rule to update
- You can either provide `rule_id` to update directly, or import the existing rule first using:
  ```bash
  terraform import edgenext_scdn_network_speed_rule.example <business_id>-<business_type>-<config_group>-<rule_id>
  ```

## Outputs

- `rule_id`: Updated rule ID

