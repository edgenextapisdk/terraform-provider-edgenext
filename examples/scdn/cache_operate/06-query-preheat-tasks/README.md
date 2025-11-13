# Query SCDN Cache Preheat Tasks Example

This example demonstrates how to query cache preheat task list in EdgeNext SCDN.

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

page     = 1
per_page = 20

# Optional filters
start_time = "2022-01-01 00:00:00"
end_time   = "2022-12-31 23:59:59"
status     = "2"
```

3. Initialize and apply:

```bash
terraform init
terraform plan
terraform apply
```

## Configuration

- `page`: Page number (default: 1)
- `per_page`: Items per page (default: 20)
- `start_time`: (Optional) Start time, format: YYYY-MM-DD HH:II:SS
- `end_time`: (Optional) End time, format: YYYY-MM-DD HH:II:SS
- `status`: (Optional) Status: 1-executing, 2-completed

## Outputs

- `total`: Total number of tasks
- `tasks`: Task list with details

