# Delete Cache Rule

This example demonstrates how to delete a cache rule.

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
rule_id      = 3009  # The cache rule ID to delete
```

3. Import the existing rule first:

```bash
terraform init
terraform import edgenext_scdn_cache_rule.example <business_id>-<business_type>-<rule_id>
# Example: terraform import edgenext_scdn_cache_rule.example 1246-tpl-3009
```

4. Destroy the resource:

```bash
terraform plan
terraform destroy
```

## Notes

- You must import the existing rule before you can delete it.
- Global cache rule IDs cannot be deleted.
- This operation is irreversible. Make sure you want to delete the rule before proceeding.

