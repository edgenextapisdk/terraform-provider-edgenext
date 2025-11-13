# Update Cache Rule Configuration

This example demonstrates how to update a cache rule's configuration including TTL, cache key settings, etc.

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
rule_id      = 3028  # The cache rule ID to update
rule_name    = "updated-cache-rule"
expr         = "(http.request.postfix in {\"css\" \"js\" \"txt\"})"
remark       = "Updated cache rule configuration"

conf = {
  # ... (see terraform.tfvars.example for full configuration)
}
```

3. Initialize and apply:

```bash
terraform init
terraform plan
terraform apply
```

## Notes

- This example updates the rule's configuration, including cache settings, TTL, cache key, etc.
- To only update the name/remark, use the `04-update-rule` example instead.

## Outputs

- `rule_id`: The updated cache rule ID

