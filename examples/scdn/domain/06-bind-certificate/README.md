# Example 6: Bind SSL Certificate to SCDN Domain

This example demonstrates how to bind an SSL certificate to an existing SCDN domain.

## Prerequisites

1. Valid EdgeNext API credentials (access_key and secret_key)
2. An existing SCDN domain ID (you can get this from Example 1)
3. An existing SSL certificate ID (from EdgeNext SSL certificate management)

## Files

- `main.tf` - Main Terraform configuration
- `terraform.tfvars.example` - Example configuration file

## Usage

1. Copy `terraform.tfvars.example` to `terraform.tfvars`:
   ```bash
   cp terraform.tfvars.example terraform.tfvars
   ```

2. Update `terraform.tfvars` with:
   - Your API credentials
   - `domain_id` - The ID of an existing SCDN domain (you can get this from Example 1 output)
   - `ca_id` - The ID of an SSL certificate (from EdgeNext SSL certificate management)
   - `enabled` - Set to `true` to bind the certificate, `false` to unbind it (default: `true`)

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

## Outputs

After successful binding, you'll see:
- `binding_id` - The unique identifier for the certificate binding (format: `domain_id-ca_id`), or `null` if unbound
- `domain_id` - The domain ID
- `ca_id` - The certificate ID
- `binding_status` - Current binding status ("bound" or "unbound")

## Binding and Unbinding

### Bind Certificate

Set `enabled = true` in `terraform.tfvars` and apply:

```bash
terraform apply
```

### Unbind Certificate

There are two ways to unbind:

**Method 1: Set enabled to false**
```bash
# Edit terraform.tfvars and set enabled = false
terraform apply
```

**Method 2: Destroy the resource**
```bash
terraform destroy
```

Both methods will unbind the certificate from the domain.

## Complete Workflow Example

You can also combine multiple examples to create a complete workflow:

```hcl
# First create a domain (from Example 1)
resource "edgenext_scdn_domain" "example" {
  domain         = "example.com"
  group_id       = 1
  protect_status = "scdn"
  app_type       = "web"
  
  origins {
    protocol        = 0
    listen_ports    = [80, 443]
    origin_protocol = 0
    load_balance    = 1
    origin_type     = 0
    
    records {
      view     = "default"
      value    = "1.2.3.4"
      port     = 80
      priority = 10
    }
  }
}

# Then bind a certificate to it
resource "edgenext_scdn_cert_binding" "example" {
  domain_id = edgenext_scdn_domain.example.id
  ca_id     = 67890  # Your SSL certificate ID
}
```

## Notes

- Certificate binding requires both the domain and certificate to exist
- The binding ID is automatically generated in the format: `{domain_id}-{ca_id}`
- Use `enabled = false` to unbind without destroying the resource completely
- One domain can only have one certificate bound at a time (binding a new certificate will replace the existing one)
- When `enabled = false`, the resource is not created, effectively unbinding the certificate

## Example: Toggle Certificate Binding

You can toggle certificate binding by changing the `enabled` variable:

```bash
# Bind the certificate
terraform apply -var="enabled=true"

# Unbind the certificate
terraform apply -var="enabled=false"
```
