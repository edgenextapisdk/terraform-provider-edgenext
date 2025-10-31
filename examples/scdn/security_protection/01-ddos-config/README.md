# DDoS Protection Configuration Example

This example demonstrates how to configure DDoS protection for a business.

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

6. To destroy the configuration:

```bash
terraform destroy
```

## Configuration Options

- `business_id`: The business ID to configure DDoS protection for
- `application_ddos_protection`: Application layer DDoS protection settings
  - `status`: Protection status (on, off, keep)
  - `ai_cc_status`: AI protection status (on, off)
  - `type`: Protection type (default, normal, strict, captcha, keep)
  - `need_attack_detection`: Attack detection switch (0 or 1)
  - `ai_status`: AI status (on, off)
- `visitor_authentication`: Visitor authentication settings
  - `status`: Authentication status (on, off)
  - `auth_token`: Authentication token
  - `pass_still_check`: Pass still check (0 or 1)

