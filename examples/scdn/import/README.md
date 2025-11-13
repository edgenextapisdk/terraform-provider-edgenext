# Import Existing Resources

This example demonstrates how to import existing EdgeNext SCDN resources into Terraform state using the `terraform import` command.

## Prerequisites

1. Valid EdgeNext API credentials (access_key and secret_key)
2. Existing resources in EdgeNext SCDN that you want to import
3. Resource IDs from the EdgeNext console or API

## Supported Import Formats

### 1. Domain Resource

```bash
terraform import edgenext_scdn_domain.example <domain_id>
```

**Example:**
```bash
terraform import edgenext_scdn_domain.example 102008
```

### 2. Origin Group Resource

```bash
terraform import edgenext_scdn_origin_group.example <origin_group_id>
```

**Example:**
```bash
terraform import edgenext_scdn_origin_group.example 85
```

### 3. Certificate Resource

```bash
terraform import edgenext_scdn_certificate.example <certificate_id>
```

**Example:**
```bash
terraform import edgenext_scdn_certificate.example 456
```

### 4. Certificate Binding Resource

```bash
terraform import edgenext_scdn_cert_binding.example "<domain_id>:<certificate_id>"
```

**Example:**
```bash
terraform import edgenext_scdn_cert_binding.example "123:456"
```

### 5. Origin Resource

```bash
terraform import edgenext_scdn_origin.example <origin_id>
```

**Note:** After importing, you may need to manually set the `domain_id` field.

## Files

- `main.tf` - Terraform configuration with resource definitions
- `terraform.tfvars.example` - Example configuration file with resource IDs

## Usage

### Step 1: Prepare Configuration

1. Copy `terraform.tfvars.example` to `terraform.tfvars`:
   ```bash
   cp terraform.tfvars.example terraform.tfvars
   ```

2. Update `terraform.tfvars` with:
   - Your API credentials (`access_key`, `secret_key`)
   - Resource IDs you want to import (`domain_id`, `origin_group_id`, etc.)

**Note:** The `main.tf` file contains placeholder values for required fields (like `origins` blocks, `ca_name`, etc.). These are necessary for Terraform validation but will be replaced with actual values after import.

### Step 2: Initialize Terraform

```bash
terraform init
```

### Step 3: Import Resources

Import each resource using the appropriate command:

```bash
# Import domain
terraform import edgenext_scdn_domain.example $(terraform output -raw domain_id_to_import)

# Or directly with ID
terraform import edgenext_scdn_domain.example 102008

# Import origin group
terraform import edgenext_scdn_origin_group.example 85

# Import certificate
terraform import edgenext_scdn_certificate.example 456
```

### Step 4: Verify Import

After importing, verify the state:

```bash
# View imported resources
terraform show

# Check specific resource
terraform state show edgenext_scdn_domain.example
```

### Step 5: Verify Configuration

Run `terraform plan` to check if the configuration matches the imported state:

```bash
terraform plan
```

**Expected Result:** 
- Should show "No changes" or only computed field differences
- **Important:** If you see "destroy and then create replacement" warnings, this means the placeholder values in `main.tf` don't match the imported state. You MUST update `main.tf` before running `terraform apply` to avoid destroying and recreating resources.

### Step 6: Update Configuration (Required)

After importing, you **must** update your `main.tf` to match the imported state:

1. Run `terraform show` to see the actual state:
   ```bash
   terraform show
   ```

2. Update `main.tf` with the correct values from the output:
   - Replace placeholder `origins` blocks with actual values
   - Update `ca_name` with the actual certificate name
   - Update `domain`, `name`, and other fields with actual values

3. Run `terraform plan` to verify:
   ```bash
   terraform plan
   ```
   Should show "No changes" or only computed field differences.

**Important:** The placeholder values in `main.tf` are only for validation. After import, you **MUST** update them with actual values from `terraform show`. 

**Warning:** If you run `terraform apply` without updating the placeholder values, Terraform may try to destroy and recreate resources to match the configuration, which could cause data loss or service interruption.

## Import Workflow Example

### Example: Importing a Domain

1. **Get the domain ID** from EdgeNext console or API

2. **Create the resource definition** in `main.tf`:
   ```hcl
   resource "edgenext_scdn_domain" "example" {
     domain = "example.com"
     # Other fields will be populated after import
   }
   ```

3. **Import the resource**:
   ```bash
   terraform import edgenext_scdn_domain.example 102008
   ```

4. **Verify the import**:
   ```bash
   terraform show
   ```

5. **Check configuration consistency**:
   ```bash
   terraform plan
   ```

6. **Update configuration** if needed based on `terraform show` output

## Common Issues

### Issue 1: Import fails with "resource not found"

**Solution:**
- Verify the resource ID is correct
- Check if the resource exists in EdgeNext console
- Ensure you have proper API credentials and permissions

### Issue 2: `terraform plan` shows many changes after import

**Possible Causes:**
- Configuration doesn't match actual resource state
- Some fields are computed and may differ
- Default values in config differ from actual values

**Solution:**
1. Run `terraform show` to see actual state
2. Update `main.tf` to match the actual state
3. For computed fields, you can ignore differences

### Issue 3: Import succeeds but Read fails

**Solution:**
- Check API credentials
- Verify network connectivity
- Check Terraform logs for detailed error messages

## Tips

1. **Import one resource at a time** to avoid confusion
2. **Always run `terraform plan`** after import to verify consistency
3. **Keep a backup** of your state file before importing
4. **Document resource IDs** for future reference
5. **Use `terraform state list`** to see all imported resources

## Next Steps

After successfully importing resources:

1. Commit your Terraform configuration to version control
2. Use `terraform plan` and `terraform apply` for future changes
3. Consider using Terraform workspaces for different environments
4. Set up CI/CD pipelines for automated infrastructure management

## Related Examples

- `01-create-domain/` - Create a new domain
- `02-add-origin/` - Add origin to domain
- `03-query-domain/` - Query domain information

