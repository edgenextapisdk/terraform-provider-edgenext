# List SCDN Rule Templates

This example demonstrates how to list SCDN rule templates with various filter options.

## Usage

1. Copy `terraform.tfvars.example` to `terraform.tfvars` and fill in your credentials:

```bash
cp terraform.tfvars.example terraform.tfvars
```

2. Edit `terraform.tfvars` with your actual values:

```hcl
access_key = "your_access_key"
secret_key = "your_secret_key"
page       = 1
page_size  = 10
name       = "my-template"  # Optional: filter by name
domain     = ""             # Optional: filter by domain
app_type   = "network_speed" # Optional: filter by app type
```

3. Initialize and apply:

```bash
terraform init
terraform plan
terraform apply
```

## Outputs

- `total`: Total number of rule templates
- `templates`: List of rule templates with their details

## Notes

- All filter parameters are optional
- `page_size` maximum is 1000
- If no filters are provided, all templates will be returned (subject to pagination)

