package main

import "strings"

// getResourceLink returns the appropriate link format based on the target platform
func getResourceLink(resourceName string) string {
	baseName := strings.Replace(resourceName, "edgenext_", "", 1)
	if *linkFormat == "github" {
		return "r/" + baseName + ".html.markdown"
	}
	return "resources/" + baseName
}

// getDataSourceLink returns the appropriate link format based on the target platform  
func getDataSourceLink(dataSourceName string) string {
	baseName := strings.Replace(dataSourceName, "edgenext_", "", 1)
	if *linkFormat == "github" {
		return "d/" + baseName + ".html.markdown"
	}
	return "data-sources/" + baseName
}

const (
	docTPL = `---
subcategory: "{{.product}}"
layout: "{{.cloud_mark}}"
page_title: "{{.cloud_title}}: {{.name}}"
sidebar_current: "docs-{{.cloud_mark}}-{{.dtype}}-{{.resource}}"
description: |-
  {{.description_short}}
---

# {{.name}}

{{.description}}

## Example Usage

{{.example}}

## Argument Reference

The following arguments are supported:

{{.arguments}}
{{if ne .attributes ""}}
## Attributes Reference

In addition to all arguments above, the following attributes are exported:

{{.attributes}}
{{end}}
{{if ne .import ""}}
## Import

{{.import}}
{{end}}
`
	idxTPL = `
<% wrap_layout :inner do %>
    <% content_for :sidebar do %>
        <div class="docs-sidebar hidden-print affix-top" role="complementary">
            <ul class="nav docs-sidenav">
                <li>
                    <a href="/docs/providers/index.html">All Providers</a>
                </li>
                <li>
                    <a href="/docs/providers/{{.cloud_mark}}/index.html">{{.cloud_title}} Provider</a>
                </li>
                {{range .Products}}
                <li>
                    <a href="#">{{.Name}}</a>
                    <ul class="nav">{{if eq .Name "Provider Data Sources"}}{{range $Resource := .DataSources}}
                        <li>
                            <a href="/docs/providers/{{$.cloud_mark}}/d/{{replace $Resource $.cloudPrefix ""}}.html">{{$Resource}}</a>
                        </li>{{end}}{{else}}
                        {{ if .DataSources }}<li>
                            <a href="#">Data Sources</a>
                            <ul class="nav nav-auto-expand">{{range $Resource := .DataSources}}
                                <li>
                                    <a href="/docs/providers/{{$.cloud_mark}}/d/{{replace $Resource $.cloudPrefix ""}}.html">{{$Resource}}</a>
                                </li>{{end}}
                            </ul>
                        </li>{{ end }}
                        <li>
                            <a href="#">Resources</a>
                            <ul class="nav nav-auto-expand">{{range $Resource := .Resources}}
                                <li>
                                    <a href="/docs/providers/{{$.cloud_mark}}/r/{{replace $Resource $.cloudPrefix ""}}.html">{{$Resource}}</a>
                                </li>{{end}}
                            </ul>
                        </li>{{end}}
                    </ul>
                </li>{{end}}
            </ul>
        </div>
    <% end %>
    <%= yield %>
<% end %>
`
	mainPageTPL = `---
layout: "{{.cloud_mark}}"
page_title: "Provider: {{.cloud_title}}"
sidebar_current: "docs-{{.cloud_mark}}-index"
description: |-
  The {{.cloud_title}} provider is used to interact with {{.cloud_title}} CDN services.
---

# {{.cloud_title}} Provider

The {{.cloud_title}} provider is used to interact with many resources supported
by [{{.cloud_title}} CDN](https://www.edgenext.com).
The provider needs to be configured with the proper API credentials before it can be used.

Use the navigation on the left to read about the available resources.

-> **Note:** This provider supports {{.cloud_title}} CDN v2.0+ API.

## Example Usage

### Configure the {{.cloud_title}} Provider

` + "```" + `hcl
terraform {
  required_providers {
    {{.cloud_mark}} = {
      source = "edgenextapisdk/{{.cloud_mark}}"
    }
  }
}

# Configure the {{.cloud_title}} Provider
provider "{{.cloud_mark}}" {
  api_key  = var.api_key
  secret   = var.secret
  endpoint = var.endpoint
}
` + "```" + `

### Configure with environment variables

` + "```" + `hcl
provider "{{.cloud_mark}}" {
  # Configuration options
}
` + "```" + `

` + "```" + `shell
export {{upper .cloud_mark}}_API_KEY="your-api-key"
export {{upper .cloud_mark}}_SECRET="your-secret"  
export {{upper .cloud_mark}}_ENDPOINT="https://api.edgenext.com"
` + "```" + `

### Configure with timeout and retry settings

` + "```" + `hcl
provider "{{.cloud_mark}}" {
  api_key     = var.api_key
  secret      = var.secret
  endpoint    = var.endpoint
  timeout     = 300
  retry_count = 3
}
` + "```" + `

## Authentication

The {{.cloud_title}} provider requires API credentials to authenticate with the {{.cloud_title}} API.

### API Credentials

You can provide your credentials via the following methods:

1. **Static credentials** (not recommended for production):

` + "```" + `hcl
provider "{{.cloud_mark}}" {
  api_key  = "your-api-key"
  secret   = "your-secret"
  endpoint = "https://api.edgenext.com"
}
` + "```" + `

2. **Environment variables** (recommended):

` + "```" + `shell
export {{upper .cloud_mark}}_API_KEY="your-api-key"
export {{upper .cloud_mark}}_SECRET="your-secret"
export {{upper .cloud_mark}}_ENDPOINT="https://api.edgenext.com"
` + "```" + `

3. **Terraform variables**:

` + "```" + `hcl
variable "api_key" {
  description = "{{.cloud_title}} API Key"
  type        = string
  sensitive   = true
}

variable "secret" {
  description = "{{.cloud_title}} Secret"
  type        = string
  sensitive   = true
}

variable "endpoint" {
  description = "{{.cloud_title}} API Endpoint"
  type        = string
}

provider "{{.cloud_mark}}" {
  api_key  = var.api_key
  secret   = var.secret
  endpoint = var.endpoint
}
` + "```" + `

## Argument Reference

The following arguments are supported in the ` + "`" + `provider` + "`" + ` block:

* ` + "`" + `api_key` + "`" + ` - (Required) {{.cloud_title}} API key for authentication. This can also be specified with the ` + "`" + `{{upper .cloud_mark}}_API_KEY` + "`" + ` environment variable.

* ` + "`" + `secret` + "`" + ` - (Required) {{.cloud_title}} secret for authentication. This can also be specified with the ` + "`" + `{{upper .cloud_mark}}_SECRET` + "`" + ` environment variable.

* ` + "`" + `endpoint` + "`" + ` - (Required) {{.cloud_title}} API endpoint address. This can also be specified with the ` + "`" + `{{upper .cloud_mark}}_ENDPOINT` + "`" + ` environment variable.

* ` + "`" + `timeout` + "`" + ` - (Optional) API request timeout in seconds. Defaults to ` + "`" + `300` + "`" + `.

* ` + "`" + `retry_count` + "`" + ` - (Optional) API request retry count. Defaults to ` + "`" + `3` + "`" + `.

## Resources and Data Sources

The {{.cloud_title}} provider supports the following resource types:

{{range .Products}}
{{if ne .Name "Provider Data Sources"}}
### {{.Name}}

{{if .Resources}}
#### Resources

{{range .Resources}}
* [` + "`" + `{{.}}` + "`" + `]({{getResourceLink .}}) - Manage {{getResourceDesc .}}
{{end}}
{{end}}

{{if .DataSources}}
#### Data Sources

{{range .DataSources}}
* [` + "`" + `{{.}}` + "`" + `]({{getDataSourceLink .}}) - Query {{getDataSourceDesc .}}
{{end}}
{{end}}
{{end}}
{{end}}
`
)
