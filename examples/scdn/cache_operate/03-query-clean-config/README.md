# Query SCDN Cache Clean Config Example

This example demonstrates how to query cache clean configuration in EdgeNext SCDN.

## Usage

1. Copy `terraform.tfvars.example` to `terraform.tfvars` and fill in your credentials:

```bash
cp terraform.tfvars.example terraform.tfvars
```

2. Edit `terraform.tfvars` with your actual values:

```hcl
access_key = "your-access-key"
secret_key = "your-secret-key"
endpoint   = "https://api.edgenextscdn.com"
```

3. Initialize and apply:

```bash
terraform init
terraform plan
terraform apply
```

## Outputs

- `config_id`: The config ID
- `wholesite`: Whole site config list
- `specialurl`: Special URL config list
- `specialdir`: Special directory config list

