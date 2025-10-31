# Update Cache Rule Name/Remark

This example demonstrates how to update a cache rule's name and remark without changing its configuration.

## Usage

1. Copy `terraform.tfvars.example` to `terraform.tfvars` and fill in your credentials:

```bash
cp terraform.tfvars.example terraform.tfvars
```

2. Edit `terraform.tfvars` with your actual values:

```hcl
access_key   = "your_access_key"
secret_key   = "your_secret_key"
business_id  = 1246  # Your template ID or domain ID
business_type = "tpl"  # 'tpl' for template or 'domain' for domain
rule_id      = 3009  # The cache rule ID to update
rule_name    = "updated-cache-rule-name"
remark       = "Updated remark"
```

3. Initialize and apply:

```bash
terraform init
terraform plan
terraform apply
```

## Notes

- This example only updates the rule's name and remark, not its configuration.
- The `expr` and `conf` fields are kept as-is (using empty/default values).
- To update the configuration, use the `05-update-rule-config` example instead.

## Outputs

- `rule_id`: The updated cache rule ID
- `rule_name`: The updated cache rule name

