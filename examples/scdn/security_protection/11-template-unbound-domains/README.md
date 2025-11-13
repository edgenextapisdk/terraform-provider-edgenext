# Template Unbound Domains Query Example

This example demonstrates how to query domains that are not bound to any security protection template.

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
domain = ""
page = 1
page_size = 20
```

3. Initialize Terraform:

```bash
terraform init
```

4. Review the execution plan:

```bash
terraform plan
```

5. Apply the configuration:

```bash
terraform apply
```

## Configuration Options

- `domain`: Domain filter (optional)
- `page`: Page number (default: 1)
- `page_size`: Page size (default: 20)
- `member_id`: Member ID (optional)

## Output

The example will output:
- `total`: Total number of unbound domains
- `domains`: List of unbound domains with details (ID, domain, type, created_at, remark)

