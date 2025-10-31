# Example 11: Manage Domain Status

This example demonstrates how to enable or disable an SCDN domain.

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
   - `enabled` - Set to `true` to enable, `false` to disable

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

After successful operation, you'll see:
- `domain_id` - The domain ID
- `enabled` - Current enabled status

## Enable/Disable Domain

### Enable Domain
```hcl
enabled = true
```

### Disable Domain
```hcl
enabled = false
```

## Notes

- Enabling a domain activates it for traffic
- Disabling a domain stops traffic but preserves configuration
- You can toggle the status by changing the `enabled` parameter and reapplying
- Destroying this resource will disable the domain before removal

