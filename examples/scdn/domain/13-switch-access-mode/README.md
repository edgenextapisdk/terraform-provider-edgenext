# Example 13: Switch Domain Access Mode

This example demonstrates how to switch the access mode of an SCDN domain between NS and CNAME.

## Prerequisites

1. Valid EdgeNext API credentials (access_key and secret_key)
2. An existing SCDN domain ID

## Files

- `main.tf` - Main Terraform configuration
- `terraform.tfvars.example` - Example configuration file

## Usage

1. Copy `terraform.tfvars.example` to `terraform.tfvars`:
   ```bash
   cp terraform.tfvars.example terraform.tfvars
   ```

2. Update `terraform.tfvars` with:
   - Your credentials
   - `domain_id` - The ID of an existing domain
   - `access_mode` - The access mode you want to switch to: `ns` or `cname`

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

After successful switch, you'll see:
- `domain_id` - The domain ID
- `access_mode` - Current access mode

## Access Mode Options

### CNAME Mode (default)
```hcl
access_mode = "cname"
```
Use CNAME records for domain access. Requires CNAME configuration in DNS.

### NS Mode
```hcl
access_mode = "ns"
```
Use NS (Name Server) records for domain access. Requires NS delegation.

## Notes

- Access mode switch operation may take some time to propagate
- The change affects how the domain is accessed via DNS
- Destroying this resource does not revert the access mode switch
- Ensure DNS is properly configured for the selected access mode

