# Query SCDN Rule Template

This example demonstrates how to query a specific SCDN rule template by ID.

## Usage

1. Copy `terraform.tfvars.example` to `terraform.tfvars` and fill in your credentials:

```bash
cp terraform.tfvars.example terraform.tfvars
```

2. Edit `terraform.tfvars` with your actual values:

```hcl
access_key  = "your_access_key"
secret_key  = "your_secret_key"
template_id = "123"
app_type    = "network_speed"
```

3. Initialize and apply:

```bash
terraform init
terraform plan
terraform apply
```

## Outputs

- `template_id`: The rule template ID
- `template_name`: The rule template name
- `description`: The rule template description
- `created_at`: The template creation timestamp
- `bind_domains`: List of domains bound to this template

