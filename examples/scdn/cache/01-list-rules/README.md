# List Cache Rules

This example demonstrates how to list cache rules for a template or domain.

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
```

3. Initialize and apply:

```bash
terraform init
terraform plan
terraform apply
```

## Outputs

- `total`: Total number of cache rules
- `rules`: List of cache rules with all their configurations

