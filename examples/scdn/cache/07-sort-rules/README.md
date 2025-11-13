# Sort Cache Rules

This example demonstrates how to sort existing cache rules within a specific business context.

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
rule_ids_to_sort = [
  3028,  # Replace with your actual rule IDs in the desired order
  3009
]
```

3. Initialize and apply:

```bash
terraform init
terraform plan
terraform apply
```

## Notes

- The `edgenext_scdn_cache_rules_sort` resource is used to define the desired order of rules.
- The `ids` attribute should contain a list of rule IDs in the exact order you want them to appear.
- The `sorted_ids` output will show the actual order of rules after the API call.
- This operation is idempotent; applying it multiple times with the same `ids` list will not cause further changes.
- Global cache rule IDs cannot be sorted.
- Deleting this resource from your Terraform state (`terraform destroy` or `terraform state rm`) will not delete the rules themselves, only the sorting configuration from Terraform's management.

## Outputs

- `sorted_ids`: The actual sorted rule IDs after the API call

