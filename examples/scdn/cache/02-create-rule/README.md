# Create Cache Rule

This example demonstrates how to create a cache rule for a template or domain.

## Usage

1. Copy `terraform.tfvars.example` to `terraform.tfvars` and fill in your credentials:

```bash
cp terraform.tfvars.example terraform.tfvars
```

2. Edit `terraform.tfvars` with your actual values:

```hcl
access_key   = "your_access_key"
secret_key   = "your_secret_key"
business_id  = 1246  # Your template ID or domain ID
business_type = "tpl"  # 'tpl' for template or 'domain' for domain
rule_name    = "my-cache-rule"
expr         = "(http.request.uri.path eq \"/test\")"
remark       = "My cache rule description"

conf = {
  # ... (see terraform.tfvars.example for full configuration)
}
```

3. Initialize and apply:

```bash
terraform init
terraform plan
terraform apply
```

## Configuration Options

The `conf` object supports the following options:

- `nocache`: Cache eligibility (true: bypass cache, false: cache)
- `cache_rule`: Edge TTL cache configuration
- `browser_cache_rule`: Browser cache configuration
- `cache_errstatus`: Status code cache configuration
- `cache_url_rewrite`: Custom cache key configuration
- `cache_share`: Cache sharing configuration (required)

## Outputs

- `rule_id`: The created cache rule ID
- `rule_status`: The cache rule status (1: enabled, 2: disabled)

