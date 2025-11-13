# Example 12: Switch Domain Nodes

This example demonstrates how to switch the node type (protect_status) of an SCDN domain.

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
   - `protect_status` - The node type you want to switch to:
     - `back_source` - Back source node
     - `scdn` - SCDN node (default)
     - `exclusive` - Exclusive node (requires exclusive_resource_id)

3. If switching to `exclusive`, also set `exclusive_resource_id`

4. Initialize Terraform:
   ```bash
   terraform init
   ```

5. Review the execution plan:
   ```bash
   terraform plan
   ```

6. Apply the configuration:
   ```bash
   terraform apply
   ```

## Outputs

After successful switch, you'll see:
- `domain_id` - The domain ID
- `protect_status` - Current protect status

## Protect Status Options

### SCDN Node (default)
```hcl
protect_status = "scdn"
```
Standard SCDN protection.

### Back Source Node
```hcl
protect_status = "back_source"
```
Back source node type.

### Exclusive Node
```hcl
protect_status = "exclusive"
exclusive_resource_id = 12345
```
Exclusive node type, requires an exclusive resource package ID.

## Notes

- Node switch operation cannot be easily reverted
- Switching to exclusive requires a valid exclusive_resource_id
- The change may take some time to propagate
- Destroying this resource does not revert the node switch

