# Template Domain Bind Example

This example demonstrates how to bind domains to a security protection template.

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
business_id = 123
domain_ids = [102008]
group_ids = []
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

6. To unbind (remove from state):

```bash
terraform destroy
```

## Configuration Options

- `business_id`: The template ID to bind domains to
- `domain_ids`: List of domain IDs to bind
- `group_ids`: List of group IDs to bind (optional)
- `bind_business_ids`: List of business IDs to bind (optional)

## Output

The example will output:
- `business_id`: The template ID
- `fail_domains`: Map of domains that failed to bind (if any)

