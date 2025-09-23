# EdgeNext Provider Documentation Generator

This directory contains the documentation generation tool for the EdgeNext Terraform Provider.

## Overview

The `gendoc` tool automatically generates Terraform provider documentation from Go source code and markdown description files. It creates documentation that is compatible with the Terraform Registry documentation standards.

## Usage

### Generate Documentation

```bash
# Method 1: Run from gendoc directory
cd gendoc && go run ./...

# Method 2: Use Makefile from project root
make doc

# Method 3: Use faster generation (with pre-built binary)
make doc-faster

# Method 4: Generate GitHub-compatible links
make doc-github
```

### Link Format Options

The documentation generator supports two link formats:

#### Terraform Registry Format (Default)
```bash
make doc
```
Generates links compatible with Terraform Registry:
- Resources: `resources/resource_name`
- Data Sources: `data-sources/data_source_name`

#### GitHub Format
```bash
make doc-github
```
Generates links compatible with GitHub:
- Resources: `r/resource_name.html.markdown`
- Data Sources: `d/data_source_name.html.markdown`

### Manual Usage

You can also run the tool directly:

```bash
# For Terraform Registry (default)
cd gendoc && go run ./... -link-format=terraform

# For GitHub
cd gendoc && go run ./... -link-format=github
```

### Build Documentation Generator Binary

```bash
# Build the binary
make doc-bin-build

# Use the binary
make doc-with-bin
```

## How It Works

### 1. Source Files

The documentation generator reads from:

- `edgenext/provider.md` - Main provider description and resource list
- `edgenext/services/{service}/data_source_en_{name}.md` - Data source descriptions
- `edgenext/services/{service}/resource_en_{name}.md` - Resource descriptions
- Go source code schemas for parameter documentation

### 2. Generated Files

The tool generates:

- `website/docs/index.html.markdown` - Provider main page
- `website/docs/d/{resource}.html.markdown` - Data source documentation
- `website/docs/r/{resource}.html.markdown` - Resource documentation
- `website/edgenext.erb` - Sidebar navigation template

### 3. Documentation Structure

Each generated documentation file includes:

- **Front matter** - Metadata for documentation rendering
- **Description** - Resource/data source overview
- **Example Usage** - Code examples with proper HCL formatting
- **Argument Reference** - All supported parameters with types and descriptions
- **Attributes Reference** - Computed attributes (for data sources)
- **Import** - Import instructions (for resources)

## File Naming Convention

- Data sources: `data_source_en_{resource_name}.md`
- Resources: `resource_en_{resource_name}.md`
- Generated docs: `{resource_name}.html.markdown`

## Quality Checks

The generator performs validation on:

- Description format and ending punctuation
- Required vs optional parameter definitions
- Data source `result_output_file` parameter requirement
- Schema consistency

## Template System

Uses Go templates with the following data:

- `{{.cloud_mark}}` - Provider name (`edgenext`)
- `{{.cloud_title}}` - Provider title (`EdgeNext`)
- `{{.product}}` - Product category
- `{{.name}}` - Full resource name
- `{{.resource}}` - Resource name without prefix
- `{{.description}}` - Full description
- `{{.example}}` - Formatted examples
- `{{.arguments}}` - Parameter documentation
- `{{.attributes}}` - Computed attributes
- `{{.import}}` - Import instructions

## Integration with CI/CD

The documentation generator integrates with:

- **Git hooks** - Pre-commit validation of documentation sync
- **Makefile** - Standardized build commands
- **Formatting checks** - HCL code formatting validation

## Terraform Registry Compatibility

Generated documentation is fully compatible with:

- Terraform Registry documentation standards
- Provider documentation structure requirements
- Markdown formatting guidelines
- Navigation sidebar structure
