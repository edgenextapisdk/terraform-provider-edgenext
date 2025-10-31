# Query SCDN Cache Clean Task Detail Example

This example demonstrates how to query cache clean task detail in EdgeNext SCDN.

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

task_id  = "acd92fee-e257-4948-ad05-c17018a6f29c"
page     = 1
per_page = 20
```

3. Initialize and apply:

```bash
terraform init
terraform plan
terraform apply
```

## Configuration

- `task_id`: (Required) Task ID to query
- `page`: Page number (default: 1)
- `per_page`: Items per page (default: 20)
- `result`: (Optional) Result filter: 1-success, 2-failed, 3-executing

## Outputs

- `total`: Total number of tasks
- `details`: Task detail list with execution results

