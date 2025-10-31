# Update SCDN Rule Template

This example demonstrates how to update an existing SCDN rule template.

## Usage

1. Copy `terraform.tfvars.example` to `terraform.tfvars` and fill in your credentials:

```bash
cp terraform.tfvars.example terraform.tfvars
```

2. Edit `terraform.tfvars` with your actual values:

```hcl
access_key    = "your_access_key"
secret_key    = "your_secret_key"
template_id   = "123"
template_name = "updated-template-name"
description   = "Updated template description"
app_type      = "network_speed"
```

3. Apply the changes directly (no import needed):

```bash
terraform init
terraform plan
terraform apply
```

## Notes

- You can update directly by providing `template_id` in the resource - no import needed
- Alternatively, you can import the existing template first: `terraform import edgenext_scdn_rule_template.example <template_id>`
- Only `name` and `description` can be updated
- `app_type` is required but cannot be changed

