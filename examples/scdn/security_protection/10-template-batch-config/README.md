# Template Batch Config Example

This example demonstrates how to batch configure security protection templates.

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
template_ids = [1194]
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

- `template_ids`: List of template IDs to configure
- `ddos_config`: DDoS protection configuration
  - `application_ddos_protection`: Application layer DDoS protection settings
  - `visitor_authentication`: Visitor authentication settings
- `waf_rule_config`: WAF rule configuration
  - `waf_rule_config`: WAF rule config settings
  - `waf_intercept_page`: WAF intercept page settings
- `bot_management_config`: Bot management configuration (optional)
- `precise_access_control_config`: Precise access control configuration (optional)
- `all`: All flag (0 or 1, optional)
- `domains`: Domain list (optional)
- `domain_ids`: Domain ID list (optional)

## Output

The example will output:
- `template_ids`: The list of template IDs that were configured
- `fail_templates`: Map of templates that failed to configure (if any)

