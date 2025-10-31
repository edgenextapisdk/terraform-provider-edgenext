# Manage Cache Rule Status

This example demonstrates how to enable or disable cache rules.

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
rule_ids     = [3009, 3028]  # List of rule IDs to update
status       = 1  # 1 for enabled, 2 for disabled
```

3. Initialize and apply:

```bash
terraform init
terraform plan
terraform apply
```

## Notes

- This resource can update the status of multiple rules at once.
- Status values: `1` for enabled, `2` for disabled.
- Global cache rule IDs cannot be updated.

## Outputs

- `updated_ids`: List of rule IDs that were successfully updated

