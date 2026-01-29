# Bind Domains to Existing Group Example

This example demonstrates how to bind domains to an **existing** SCDN Domain Group.

## ⚠️ Important: Import Required

This example requires you to **import an existing group** before applying changes. Without importing, Terraform will create a new group!

## Prerequisites

1. An existing domain group (create via console or example `01-manage-groups`)
2. The group ID and name (get from console or example `02-query-groups`)

## Usage

### Step 1: Configure Variables

Copy and edit the tfvars file:

```bash
cp terraform.tfvars.example terraform.tfvars
```

Edit `terraform.tfvars` with:
- Your credentials
- The existing group ID and name
- Domains to bind

### Step 2: Initialize Terraform

```bash
terraform init
```

### Step 3: Import the Existing Group

**This step is required!**

```bash
terraform import edgenext_scdn_domain_group.existing <group_id>
```

Example:
```bash
terraform import edgenext_scdn_domain_group.existing 111
```

### Step 4: Apply Changes

```bash
terraform plan
terraform apply
```

## What This Does

1. Imports the existing group into Terraform state
2. Binds the specified domains to the group
3. Outputs the bound domains

## Troubleshooting

**Problem: "A new group was created instead of updating the existing one"**

This happens when you run `terraform apply` without first running `terraform import`. Delete the new group and start from Step 3.

**Problem: "Group name mismatch"**

Make sure the `group_name` variable matches the actual name of the group you imported. You can get this from:
- The EdgeNext console
- Running example `02-query-groups`
