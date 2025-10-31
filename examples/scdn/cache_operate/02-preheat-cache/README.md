# SCDN Cache Preheat Task Example

This example demonstrates how to create a cache preheat task in EdgeNext SCDN.

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

preheat_url = [
  "http://www.example.com/a.jpg",
  "http://www.example.com/b.jpg",
  "http://www.example.com/c.jpg"
]
```

3. Initialize and apply:

```bash
terraform init
terraform plan
terraform apply
```

## Resource Configuration

- `preheat_url`: (Required) List of URLs to preheat
- `group_id`: (Optional) Group ID, can refresh cache by group
- `protocol`: (Optional) Protocol: http/https; only valid when refreshing by group
- `port`: (Optional) Website port, only needed for special ports; only valid when refreshing by group

## Outputs

- `task_id`: The created cache preheat task ID
- `error_url`: List of URLs with preheat errors (if any)

## Notes

- Cache preheat tasks are one-time operations and cannot be deleted
- The task ID will be available in the task list after creation
- Check the `error_url` output to see if any URLs failed to preheat

