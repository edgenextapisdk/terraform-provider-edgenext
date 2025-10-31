# SCDN Cache Clean Task Example

This example demonstrates how to create a cache clean task in EdgeNext SCDN.

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

# Choose one or more of the following options:
wholesite  = ["www.example.com", "test.example.com"]
specialurl = ["http://www.example.com/a", "http://test.example.com/a1"]
specialdir = ["http://www.example.com/a/", "http://test.example.com/a1/a2/"]
```

3. Initialize and apply:

```bash
terraform init
terraform plan
terraform apply
```

## Resource Configuration

- `wholesite`: List of whole site domains to clean
- `specialurl`: List of special URLs to clean
- `specialdir`: List of special directories to clean
- `group_id`: (Optional) Group ID, can refresh cache by group
- `protocol`: (Optional) Protocol: http/https; only valid when refreshing by group
- `port`: (Optional) Website port, only needed for special ports; only valid when refreshing by group

## Notes

- At least one of `wholesite`, `specialurl`, or `specialdir` must be provided
- Cache clean tasks are one-time operations and cannot be deleted
- The task ID will be available in the task list after creation

