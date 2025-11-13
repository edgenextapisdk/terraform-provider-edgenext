# Create SCDN Rule Template

This example demonstrates how to create a basic SCDN rule template.

## Usage

1. Copy `terraform.tfvars.example` to `terraform.tfvars` and fill in your credentials:

```bash
cp terraform.tfvars.example terraform.tfvars
```

2. Edit `terraform.tfvars` with your actual values:

```hcl
access_key    = "your_access_key"
secret_key    = "your_secret_key"
template_name = "my-template"
description   = "A test template"
app_type      = "network_speed"
domain_ids    = [101753]  # Optional: domain IDs to bind
```

3. Initialize and apply:

```bash
terraform init
terraform plan
terraform apply
```

## Outputs

- `template_id`: The created rule template ID
- `template_name`: The rule template name
- `created_at`: The template creation timestamp

## Notes

- The `app_type` is required and should match your application type (e.g., "network_speed")
- `domain_ids` is optional. If provided, domains will be bound to the template during creation
- You can also bind domains later using the `edgenext_scdn_rule_template_binding` resource

