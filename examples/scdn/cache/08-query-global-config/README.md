# Query Global Cache Config

This example demonstrates how to query the default global cache configuration.

## Usage

1. Copy `terraform.tfvars.example` to `terraform.tfvars` and fill in your credentials:

```bash
cp terraform.tfvars.example terraform.tfvars
```

2. Edit `terraform.tfvars` with your actual values:

```hcl
access_key = "your_access_key"
secret_key = "your_secret_key"
endpoint   = "https://api.edgenextscdn.com"
```

3. Initialize and apply:

```bash
terraform init
terraform plan
terraform apply
```

## Notes

- This data source retrieves the default global cache configuration.
- The global cache configuration is read-only and cannot be modified through Terraform.
- The configuration structure is the same as cache rule configurations.

## Outputs

- `global_config`: Complete global cache configuration including ID, name, and all configuration details

