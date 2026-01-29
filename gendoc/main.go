package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"runtime"
	"sort"
	"strings"
	"text/template"

	"github.com/fatih/color"
	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	edgenext "github.com/edgenextapisdk/terraform-provider-edgenext/edgenext"
)

// ensureDir creates directories if they don't exist
func ensureDir(filePath string) error {
	dir := filepath.Dir(filePath)
	return os.MkdirAll(dir, 0755)
}

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
	cloudMark      = "edgenext"
	cloudTitle     = "EdgeNext"
	cloudPrefix    = cloudMark + "_"
	cloudMarkShort = "en"
	docRoot        = "../website/docs"

	// Template for individual resource/data source documentation
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

	// Template for sidebar navigation (edgenext.erb)
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
)

var (
	hclMatch          = regexp.MustCompile("(?si)([^`]+)?```(hcl)?(.*?)```")
	usageMatch        = regexp.MustCompile(`(?s)(?m)^([^ \n].*?)(?:\n{2}|$)(.*)`)
	bigSymbol         = regexp.MustCompile("([\u007F-\uffff])")
	productNameRegexp = regexp.MustCompile(`^.*\((.*)\)$`)

	// Command line flags
	linkFormat = flag.String("link-format", "terraform", "Link format: 'terraform' for Registry or 'github' for GitHub")
)

func main() {
	flag.Parse()

	provider := edgenext.Provider()
	vProvider := runtime.FuncForPC(reflect.ValueOf(edgenext.Provider).Pointer())

	filename, _ := vProvider.FileLine(0)
	filePath := filepath.Dir(filename)
	message("generating doc from: %s\n", filePath)

	// document for Index
	products := genIdx(filePath)

	// document for Main Page
	genMainPage(filePath, products)

	for _, product := range products {
		// document for DataSources
		for _, dataSource := range product.DataSources {
			genDoc(product.Name, "data_source", filePath, dataSource, provider.DataSourcesMap[dataSource])
		}

		// document for Resources
		for _, resource := range product.Resources {
			genDoc(product.Name, "resource", filePath, resource, provider.ResourcesMap[resource])
		}
	}
}

// genIdx generating index for resource
func genIdx(filePath string) (prods []Product) {
	filename := "provider.md"

	message("[START]get description from file: %s\n", filename)

	raw, err := os.ReadFile(filepath.Join(filePath, filename))
	if err != nil {
		message("[SKIP!]get description failed, skip: %s", err)
		return
	}
	description := string(raw)

	description = strings.TrimSpace(description)
	if description == "" {
		message("[SKIP!]description empty, skip: %s\n", filename)
		return
	}

	pos := strings.Index(description, "\nResources List\n")
	if pos == -1 {
		message("[SKIP!]resource list missing, skip: %s\n", filename)
		return
	}

	doc := strings.TrimSpace(description[pos+16:])

	prods, err = GetIndex(doc)
	if err != nil {
		message("[FAIL!]: %s", err)
		os.Exit(1)
	}

	data := map[string]interface{}{
		"cloud_mark":  cloudMark,
		"cloud_title": cloudTitle,
		"cloudPrefix": cloudPrefix,
		"Products":    prods,
	}

	filename = filepath.Join(docRoot, "..", fmt.Sprintf("%s.erb", cloudMark))
	if err := ensureDir(filename); err != nil {
		message("[FAIL!]create directory for %s failed: %s", filename, err)
		os.Exit(1)
	}
	fd, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		message("[FAIL!]open file %s failed: %s", filename, err)
		os.Exit(1)
	}

	defer fd.Close()

	tmpl := template.Must(template.New("t").Funcs(template.FuncMap{"replace": replace}).Parse(idxTPL))

	if err := tmpl.Execute(fd, data); err != nil {
		message("[FAIL!]write file %s failed: %s", filename, err)
		os.Exit(1)
	}

	message("[SUCC.]write doc to file success: %s", filename)
	return
}

// genMainPage generating main page documentation from provider.md
func genMainPage(filePath string, products []Product) {
	// Read provider.md content
	providerMdPath := filepath.Join(filePath, "provider.md")
	raw, err := os.ReadFile(providerMdPath)
	if err != nil {
		message("[FAIL!]read provider.md failed: %s", err)
		os.Exit(1)
	}

	providerContent := string(raw)

	// Find the "Resources List" section
	pos := strings.Index(providerContent, "\nResources List\n")
	if pos == -1 {
		message("[FAIL!]Resources List section not found in provider.md")
		os.Exit(1)
	}

	// Get content before "Resources List"
	contentBeforeResourcesList := strings.TrimSpace(providerContent[:pos])

	// Generate resources and data sources section
	var resourcesSection strings.Builder
	resourcesSection.WriteString("\n\n## Resources and Data Sources\n\n")
	resourcesSection.WriteString("The EdgeNext provider supports the following resource types:\n\n")

	for _, product := range products {
		if product.Name == "Provider Data Sources" {
			continue
		}

		resourcesSection.WriteString(fmt.Sprintf("### %s\n\n", product.Name))

		if len(product.Resources) > 0 {
			resourcesSection.WriteString("#### Resources\n\n")
			for _, resource := range product.Resources {
				link := getResourceLink(resource)
				desc := getResourceDesc(resource)
				resourcesSection.WriteString(fmt.Sprintf("* [`%s`](%s) - Manage %s\n", resource, link, desc))
			}
			resourcesSection.WriteString("\n")
		}

		if len(product.DataSources) > 0 {
			resourcesSection.WriteString("#### Data Sources\n\n")
			for _, dataSource := range product.DataSources {
				link := getDataSourceLink(dataSource)
				desc := getDataSourceDesc(dataSource)
				resourcesSection.WriteString(fmt.Sprintf("* [`%s`](%s) - Query %s\n", dataSource, link, desc))
			}
			resourcesSection.WriteString("\n")
		}
	}

	// Combine content
	finalContent := contentBeforeResourcesList + resourcesSection.String()

	// Write to index.html.markdown
	filename := filepath.Join(docRoot, "index.html.markdown")
	if err := ensureDir(filename); err != nil {
		message("[FAIL!]create directory for %s failed: %s", filename, err)
		os.Exit(1)
	}

	err = os.WriteFile(filename, []byte(finalContent), 0644)
	if err != nil {
		message("[FAIL!]write file %s failed: %s", filename, err)
		os.Exit(1)
	}

	message("[SUCC.]write doc to file success: %s", filename)
}

// genDoc generating doc for data source and resource
func genDoc(product, dtype, fpath, name string, resource *schema.Resource) {
	data := map[string]string{
		"product":           product,
		"name":              name,
		"dtype":             strings.Replace(dtype, "_", "", -1),
		"resource":          name[len(cloudMark)+1:],
		"cloud_mark":        cloudMark,
		"cloud_title":       cloudTitle,
		"example":           "",
		"description":       "",
		"description_short": "",
		"import":            "",
	}

	productDir := strings.ToLower(product)
	groups := productNameRegexp.FindStringSubmatch(productDir)
	if groups != nil {
		productDir = groups[1]
	}
	if productDir == "provider data sources" {
		productDir = "common"
	}

	// Map product names to service directories
	switch productDir {
	case "content delivery network(cdn)", "content delivery network":
		productDir = "cdn"
	case "ssl certificate management(ssl)", "ssl certificate management":
		productDir = "ssl"
	case "object storage service(oss)", "object storage service":
		productDir = "oss"
	case "security cdn(scdn)", "security cdn", "scdn":
		productDir = "scdn"
	case "security dns(sdns)", "security dns", "sdns":
		productDir = "sdns"
	}

	// Try to find the file, first in the main directory, then in subdirectories
	filename := fmt.Sprintf("services/%s/%s_%s_%s.md", productDir, dtype, cloudMarkShort, data["resource"])
	message("[START]get description from file: %s\n", filename)

	var raw []byte
	var err error

	// First try the main directory
	filePath := filepath.Join(fpath, filename)
	raw, err = os.ReadFile(filePath)

	// If not found and it's SCDN or SDNS, try subdirectories
	if err != nil && (productDir == "scdn" || productDir == "sdns") {
		// List subdirectories in service path
		servicePath := filepath.Join(fpath, "services", productDir)
		if entries, err2 := os.ReadDir(servicePath); err2 == nil {
			for _, entry := range entries {
				if entry.IsDir() {
					subDirPath := filepath.Join(servicePath, entry.Name(), fmt.Sprintf("%s_%s_%s.md", dtype, cloudMarkShort, data["resource"]))
					if raw2, err2 := os.ReadFile(subDirPath); err2 == nil {
						raw = raw2
						err = nil
						filename = fmt.Sprintf("services/%s/%s/%s_%s_%s.md", productDir, entry.Name(), dtype, cloudMarkShort, data["resource"])
						message("[INFO]Found file in subdirectory: %s\n", filename)
						break
					}
				}
			}
		}
	}

	if err != nil {
		message("[FAIL!]get description failed: %s", err)
		os.Exit(1)
	}
	description := string(raw)

	description = strings.TrimSpace(description)
	if description == "" {
		message("[FAIL!]description empty: %s\n", filename)
		os.Exit(1)
	}

	importPos := strings.Index(description, "\nImport\n")
	if importPos != -1 {
		data["import"] = strings.TrimSpace(description[importPos+8:])
		description = strings.TrimSpace(description[:importPos])
	}

	pos := strings.Index(description, "\nExample Usage\n")
	if pos != -1 {
		data["example"] = formatHCL(description[pos+15:])
		description = strings.TrimSpace(description[:pos])
	} else {
		message("[FAIL!]example usage missing: %s\n", filename)
		os.Exit(1)
	}

	data["description"] = description
	pos = strings.Index(description, "\n\n")
	if pos != -1 {
		data["description_short"] = strings.TrimSpace(description[:pos])
	} else {
		data["description_short"] = description
	}

	var (
		requiredArgs []string
		optionalArgs []string
		attributes   []string
		subStruct    []string
	)

	if _, ok := resource.Schema["output_file"]; dtype == "data_source" && !ok {
		if resource.DeprecationMessage != "" {
			message("[SKIP!]argument 'output_file' is missing, skip: %s", filename)
		} else {
			message("[WARN!]argument 'output_file' is missing: %s", filename)
			// Don't exit, just warn for now
		}
	}

	for k, v := range resource.Schema {
		if v.Description == "" {
			message("[WARN!]description for '%s' is missing: %s\n", k, filename)
			// Don't exit, just warn for now
		} else {
			// Skip description format check for now to allow documentation generation
			// checkDescription(k, v.Description)
		}
		if dtype == "data_source" && v.ForceNew {
			message("[FAIL!]Don't set ForceNew on data source: '%s'", k)
			os.Exit(1)
		}
		if v.Required && v.Optional {
			message("[FAIL!]Don't set Required and Optional at the same time: '%s'", k)
			os.Exit(1)
		}
		if v.Required {
			opt := "Required"
			sub := getSubStruct(0, "", k, v)
			subStruct = append(subStruct, sub...)
			// get type
			res := parseSubtract(v, sub)
			valueType := parseType(v)
			if res == "" {
				opt += fmt.Sprintf(", %s", valueType)
			} else {
				opt += fmt.Sprintf(", %s: [`%s`]", valueType, res)
			}
			if v.ForceNew {
				opt += ", ForceNew"
			}
			if v.Deprecated != "" {
				opt += ", **Deprecated**"
				v.Description = fmt.Sprintf("%s %s", v.Deprecated, v.Description)
			}
			requiredArgs = append(requiredArgs, fmt.Sprintf("* `%s` - (%s) %s", k, opt, v.Description))
		} else if v.Optional {
			opt := "Optional"
			sub := getSubStruct(0, "", k, v)
			subStruct = append(subStruct, sub...)
			// get type
			res := parseSubtract(v, sub)
			valueType := parseType(v)
			if res == "" {
				opt += fmt.Sprintf(", %s", valueType)
			} else {
				opt += fmt.Sprintf(", %s: [`%s`]", valueType, res)
			}
			if v.ForceNew {
				opt += ", ForceNew"
			}
			if v.Deprecated != "" {
				opt += ", **Deprecated**"
				v.Description = fmt.Sprintf("%s %s", v.Deprecated, v.Description)
			}
			optionalArgs = append(optionalArgs, fmt.Sprintf("* `%s` - (%s) %s", k, opt, v.Description))
		} else {
			attrs := getAttributes(0, k, v)
			if len(attrs) > 0 {
				attributes = append(attributes, attrs...)
			}
		}
	}

	sort.Strings(requiredArgs)
	sort.Strings(optionalArgs)
	sort.Strings(attributes)

	// Sort subStruct with intelligent hierarchical ordering
	if len(subStruct) > 0 {
		var configMain []string
		var level1Objects []string // Direct children of config (e.g., add_response_head, head_control)
		var level2Objects []string // Children of level1 objects (e.g., list of add_response_head)
		var otherObjects []string  // Other objects

		for _, s := range subStruct {
			if strings.Contains(s, "The `config` object supports the following:") {
				configMain = append(configMain, s)
			} else if strings.Contains(s, "object of `config` supports the following:") {
				level1Objects = append(level1Objects, s)
			} else if strings.Contains(s, "object of `add_response_head` supports the following:") ||
				strings.Contains(s, "object of `head_control` supports the following:") ||
				strings.Contains(s, "object of `cache_rule` supports the following:") ||
				strings.Contains(s, "object of `cache_rule_list` supports the following:") {
				level2Objects = append(level2Objects, s)
			} else {
				otherObjects = append(otherObjects, s)
			}
		}

		// Sort each level separately
		sort.Strings(level1Objects)
		sort.Strings(level2Objects)
		sort.Strings(otherObjects)

		// Combine in hierarchical order: config main → level1 → level2 → others
		subStruct = append(configMain, level1Objects...)
		subStruct = append(subStruct, level2Objects...)
		subStruct = append(subStruct, otherObjects...)

		// remove duplicates
		uniqSubStruct := make([]string, 0, len(subStruct))
		var i int
		for i = 0; i < len(subStruct)-1; i++ {
			if subStruct[i] != subStruct[i+1] {
				uniqSubStruct = append(uniqSubStruct, subStruct[i])
			}
		}
		uniqSubStruct = append(uniqSubStruct, subStruct[i])
		subStruct = uniqSubStruct
	}

	requiredArgs = append(requiredArgs, optionalArgs...)
	data["arguments"] = strings.Join(requiredArgs, "\n")
	if len(subStruct) > 0 {
		data["arguments"] += "\n" + strings.Join(subStruct, "\n")
	}
	data["attributes"] = strings.Join(attributes, "\n")
	if dtype == "resource" {
		// Check if id is already defined in attributes
		hasID := false
		for _, attr := range attributes {
			if strings.HasPrefix(attr, "* `id` -") {
				hasID = true
				break
			}
		}
		if !hasID {
			idAttribute := "* `id` - ID of the resource.\n"
			data["attributes"] = idAttribute + data["attributes"]
		}
	}

	filename = filepath.Join(docRoot, dtype[:1], fmt.Sprintf("%s.html.markdown", data["resource"]))

	if err := ensureDir(filename); err != nil {
		message("[FAIL!]create directory for %s failed: %s", filename, err)
		os.Exit(1)
	}
	fd, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		message("[FAIL!]open file %s failed: %s", filename, err)
		os.Exit(1)
	}

	defer fd.Close()
	t := template.Must(template.New("t").Parse(docTPL))
	err = t.Execute(fd, data)
	if err != nil {
		message("[FAIL!]write file %s failed: %s", filename, err)
		os.Exit(1)
	}

	message("[SUCC.]write doc to file success: %s", filename)
}

// getAttributes get attributes from schema
func getAttributes(step int, k string, v *schema.Schema) []string {
	var attributes []string
	ident := strings.Repeat(" ", step*2)

	if v.Description == "" {
		return attributes
	} else {
		// Skip description format check for now
		// checkDescription(k, v.Description)
	}

	if v.Computed {
		if v.Deprecated != "" {
			v.Description = fmt.Sprintf("(**Deprecated**) %s %s", v.Deprecated, v.Description)
		}
		if _, ok := v.Elem.(*schema.Resource); ok {
			var listAttributes []string
			for kk, vv := range v.Elem.(*schema.Resource).Schema {
				attrs := getAttributes(step+1, kk, vv)
				if len(attrs) > 0 {
					listAttributes = append(listAttributes, attrs...)
				}
			}
			var slistAttributes string
			sort.Strings(listAttributes)
			if len(listAttributes) > 0 {
				slistAttributes = "\n" + strings.Join(listAttributes, "\n")
			}
			attributes = append(attributes, fmt.Sprintf("%s* `%s` - %s%s", ident, k, v.Description, slistAttributes))
		} else {
			attributes = append(attributes, fmt.Sprintf("%s* `%s` - %s", ident, k, v.Description))
		}
	}

	return attributes
}

// getSubStruct get sub structure from go file
func getSubStruct(step int, parentK, k string, v *schema.Schema) []string {
	var subStructs []string

	if v.Description == "" {
		return subStructs
	} else {
		// Skip description format check for now
		// checkDescription(k, v.Description)
	}

	var subStruct []string
	if v.Type == schema.TypeMap || v.Type == schema.TypeList || v.Type == schema.TypeSet {
		if _, ok := v.Elem.(*schema.Resource); ok {
			if step == 0 {
				subStruct = append(subStruct, fmt.Sprintf("\nThe `%s` object supports the following:\n", k))
			} else {
				subStruct = append(subStruct, fmt.Sprintf("\nThe `%s` object of `%s` supports the following:\n", k, parentK))
			}
			var (
				requiredArgs []string
				optionalArgs []string
			)
			for kk, vv := range v.Elem.(*schema.Resource).Schema {
				if vv.Required {
					opt := "Required"
					valueType := parseType(vv)
					opt += fmt.Sprintf(", %s", valueType)
					if vv.ForceNew {
						opt += ", ForceNew"
					}
					if vv.Deprecated != "" {
						opt += ", **Deprecated**"
						vv.Description = fmt.Sprintf("%s %s", vv.Deprecated, vv.Description)
					}
					requiredArgs = append(requiredArgs, fmt.Sprintf("* `%s` - (%s) %s", kk, opt, vv.Description))
				} else if vv.Optional {
					opt := "Optional"
					valueType := parseType(vv)
					opt += fmt.Sprintf(", %s", valueType)
					if vv.ForceNew {
						opt += ", ForceNew"
					}
					if vv.Deprecated != "" {
						opt += ", **Deprecated**"
						vv.Description = fmt.Sprintf("%s %s", vv.Deprecated, vv.Description)
					}
					optionalArgs = append(optionalArgs, fmt.Sprintf("* `%s` - (%s) %s", kk, opt, vv.Description))
				}
			}
			sort.Strings(requiredArgs)
			subStruct = append(subStruct, requiredArgs...)
			sort.Strings(optionalArgs)
			subStruct = append(subStruct, optionalArgs...)
			subStructs = append(subStructs, strings.Join(subStruct, "\n"))

			for kk, vv := range v.Elem.(*schema.Resource).Schema {
				subStructs = append(subStructs, getSubStruct(step+1, k, kk, vv)...)
			}
		}
	}

	return subStructs
}

// formatHCL format HLC code
func formatHCL(s string) string {
	var rr []string

	s = strings.TrimSpace(s)
	m := hclMatch.FindAllStringSubmatch(s, -1)
	if len(m) > 0 {
		for _, v := range m {
			p := strings.TrimSpace(v[1])
			if p != "" {
				p = formatUsageDesc(p)
			}
			b := hclwrite.Format([]byte(strings.TrimSpace(v[3])))
			rr = append(rr, fmt.Sprintf("\n%s\n\n```hcl\n%s\n```", p, string(b)))
		}
	}

	return strings.TrimSpace(strings.Join(rr, "\n"))
}

func formatUsageDesc(s string) string {
	var rr []string
	s = strings.TrimSpace(s)
	m := usageMatch.FindAllStringSubmatch(s, -1)

	for _, v := range m {
		title := strings.TrimSpace(v[1])
		descp := strings.TrimSpace(v[2])

		rr = append(rr, fmt.Sprintf("### %s\n\n%s", title, descp))
	}

	ret := strings.TrimSpace(strings.Join(rr, "\n\n"))
	return ret
}

// checkDescription check description format
func checkDescription(k, s string) {
	if s == "" {
		return
	}

	if strings.TrimLeft(s, " ") != s {
		message("[FAIL!]There is space on the left of description: '%s': '%s'", k, s)
		os.Exit(1)
	}

	if strings.TrimRight(s, " ") != s {
		message("[FAIL!]There is space on the right of description: '%s': '%s'", k, s)
		os.Exit(1)
	}

	if s[len(s)-1] != '.' && s[len(s)-1] != ':' {
		message("[FAIL!]There is no ending charset(. or :) on the description: '%s': '%s'", k, s)
		os.Exit(1)
	}

	if c := containsBigSymbol(s); c != "" {
		message("[FAIL!]There is unexcepted symbol '%s' on the description: '%s': '%s'", c, k, s)
		os.Exit(1)
	}

	for _, v := range []string{",", ".", ";", ":", "?", "!"} {
		if strings.Contains(s, " "+v) {
			message("[FAIL!]There is space before '%s' on the description: '%s': '%s'", v, k, s)
			os.Exit(1)
		}
	}
}

// containsBigSymbol returns the Big symbol if found
func containsBigSymbol(s string) string {
	m := bigSymbol.FindStringSubmatch(s)
	if len(m) > 0 {
		return m[0]
	}

	return ""
}

// message print color message
func message(msg string, v ...interface{}) {
	if strings.Contains(msg, "FAIL") {
		color.Red(fmt.Sprintf(msg, v...))
	} else if strings.Contains(msg, "SUCC") {
		color.Green(fmt.Sprintf(msg, v...))
	} else if strings.Contains(msg, "SKIP") {
		color.Yellow(fmt.Sprintf(msg, v...))
	} else {
		color.White(fmt.Sprintf(msg, v...))
	}
}

func parseType(v *schema.Schema) string {
	res := ""
	switch v.Type {
	case schema.TypeBool:
		res = "Bool"
	case schema.TypeInt:
		res = "Int"
	case schema.TypeFloat:
		res = "Float64"
	case schema.TypeString:
		res = "String"
	case schema.TypeList:
		res = "List"
	case schema.TypeMap:
		res = "Map"
	case schema.TypeSet:
		res = "Set"
	}
	return res
}

func parseSubtract(v *schema.Schema, subStruct []string) string {
	res := ""
	if v.Type == schema.TypeSet || v.Type == schema.TypeList {
		if len(subStruct) == 0 {
			vv := v.Elem.(*schema.Schema)
			res = parseType(vv)
		}
	}
	return res
}

func replace(input, from, to string) string {
	return strings.Replace(input, from, to, -1)
}

// getResourceDesc returns a friendly description for resources
func getResourceDesc(resourceName string) string {
	descriptions := map[string]string{
		"edgenext_cdn_domain":      "CDN domain configuration",
		"edgenext_cdn_purge":       "CDN cache purge tasks",
		"edgenext_cdn_prefetch":    "CDN cache prefetch tasks",
		"edgenext_ssl_certificate": "SSL certificates",
		"edgenext_oss_bucket":      "OSS buckets",
		"edgenext_oss_object":      "OSS objects",
		"edgenext_oss_object_copy": "OSS object copy",
		// SCDN resources
		"edgenext_scdn_domain":                                    "SCDN domain configuration",
		"edgenext_scdn_origin":                                    "SCDN origin servers",
		"edgenext_scdn_cert_binding":                              "SCDN certificate bindings",
		"edgenext_scdn_domain_base_settings":                      "SCDN domain base settings",
		"edgenext_scdn_domain_status":                             "SCDN domain status management",
		"edgenext_scdn_domain_node_switch":                        "SCDN domain node switching",
		"edgenext_scdn_domain_access_mode":                        "SCDN domain access mode",
		"edgenext_scdn_certificate":                               "SCDN certificates",
		"edgenext_scdn_certificate_apply":                         "SCDN certificate application",
		"edgenext_scdn_rule_template":                             "SCDN rule templates",
		"edgenext_scdn_rule_template_domain_bind":                 "SCDN rule template domain bindings",
		"edgenext_scdn_rule_template_domain_unbind":               "SCDN rule template domain unbindings",
		"edgenext_scdn_network_speed_config":                      "SCDN network speed configuration",
		"edgenext_scdn_network_speed_rule":                        "SCDN network speed rules",
		"edgenext_scdn_network_speed_rules_sort":                  "SCDN network speed rules sorting",
		"edgenext_scdn_cache_rule":                                "SCDN cache rules",
		"edgenext_scdn_cache_rule_status":                         "SCDN cache rule status",
		"edgenext_scdn_cache_rules_sort":                          "SCDN cache rules sorting",
		"edgenext_scdn_security_protection_ddos_config":           "SCDN DDoS protection configuration",
		"edgenext_scdn_security_protection_waf_config":            "SCDN WAF protection configuration",
		"edgenext_scdn_security_protection_template":              "SCDN security protection templates",
		"edgenext_scdn_security_protection_template_domain_bind":  "SCDN security protection template domain bindings",
		"edgenext_scdn_security_protection_template_batch_config": "SCDN security protection template batch configuration",
		"edgenext_scdn_origin_group":                              "SCDN origin groups",
		"edgenext_scdn_origin_group_domain_bind":                  "SCDN origin group domain bindings",
		"edgenext_scdn_origin_group_domain_copy":                  "SCDN origin group domain copying",
		"edgenext_scdn_cache_clean_task":                          "SCDN cache clean tasks",
		"edgenext_scdn_cache_preheat_task":                        "SCDN cache preheat tasks",
		"edgenext_scdn_log_download_task":                         "SCDN log download tasks",
		"edgenext_scdn_log_download_template":                     "SCDN log download templates",
		"edgenext_scdn_log_download_template_status":              "SCDN log download template status",
	}

	if desc, ok := descriptions[resourceName]; ok {
		return desc
	}

	// Fallback: convert resource name to description
	name := strings.TrimPrefix(resourceName, "edgenext_")
	name = strings.Replace(name, "_", " ", -1)
	return name
}

// getDataSourceDesc returns a friendly description for data sources
func getDataSourceDesc(dataSourceName string) string {
	descriptions := map[string]string{
		"edgenext_cdn_domain":       "CDN domain configuration",
		"edgenext_cdn_domains":      "CDN domains",
		"edgenext_cdn_purge":        "CDN purge task details",
		"edgenext_cdn_purges":       "CDN purge tasks",
		"edgenext_cdn_prefetch":     "CDN prefetch task details",
		"edgenext_cdn_prefetches":   "CDN prefetch tasks",
		"edgenext_ssl_certificate":  "SSL certificate details",
		"edgenext_ssl_certificates": "SSL certificates",
		"edgenext_oss_buckets":      "OSS buckets",
		"edgenext_oss_object":       "OSS object details",
		"edgenext_oss_objects":      "OSS objects",
		// SCDN data sources
		"edgenext_scdn_domain":                                       "SCDN domain details",
		"edgenext_scdn_domains":                                      "SCDN domains",
		"edgenext_scdn_origin":                                       "SCDN origin details",
		"edgenext_scdn_origins":                                      "SCDN origins",
		"edgenext_scdn_domain_base_settings":                         "SCDN domain base settings",
		"edgenext_scdn_access_progress":                              "SCDN access progress options",
		"edgenext_scdn_domain_templates":                             "SCDN domain templates",
		"edgenext_scdn_brief_domains":                                "SCDN brief domain information",
		"edgenext_scdn_certificate":                                  "SCDN certificate details",
		"edgenext_scdn_certificates":                                 "SCDN certificates",
		"edgenext_scdn_certificates_by_domains":                      "SCDN certificates by domains",
		"edgenext_scdn_certificate_export":                           "SCDN certificate export",
		"edgenext_scdn_rule_template":                                "SCDN rule template details",
		"edgenext_scdn_rule_templates":                               "SCDN rule templates",
		"edgenext_scdn_rule_template_domains":                        "SCDN rule template domains",
		"edgenext_scdn_network_speed_config":                         "SCDN network speed configuration",
		"edgenext_scdn_network_speed_rules":                          "SCDN network speed rules",
		"edgenext_scdn_cache_rules":                                  "SCDN cache rules",
		"edgenext_scdn_cache_global_config":                          "SCDN cache global configuration",
		"edgenext_scdn_security_protection_ddos_config":              "SCDN DDoS protection configuration",
		"edgenext_scdn_security_protection_waf_config":               "SCDN WAF protection configuration",
		"edgenext_scdn_security_protection_template":                 "SCDN security protection template details",
		"edgenext_scdn_security_protection_templates":                "SCDN security protection templates",
		"edgenext_scdn_security_protection_template_domains":         "SCDN security protection template domains",
		"edgenext_scdn_security_protection_template_unbound_domains": "SCDN security protection template unbound domains",
		"edgenext_scdn_security_protection_member_global_template":   "SCDN security protection member global template",
		"edgenext_scdn_security_protection_iota":                     "SCDN security protection IOTA",
		"edgenext_scdn_origin_group":                                 "SCDN origin group details",
		"edgenext_scdn_origin_groups":                                "SCDN origin groups",
		"edgenext_scdn_origin_groups_all":                            "SCDN all origin groups",
		"edgenext_scdn_cache_clean_config":                           "SCDN cache clean configuration",
		"edgenext_scdn_cache_clean_tasks":                            "SCDN cache clean tasks",
		"edgenext_scdn_cache_clean_task_detail":                      "SCDN cache clean task details",
		"edgenext_scdn_cache_preheat_tasks":                          "SCDN cache preheat tasks",
		"edgenext_scdn_log_download_tasks":                           "SCDN log download tasks",
		"edgenext_scdn_log_download_templates":                       "SCDN log download templates",
		"edgenext_scdn_log_download_fields":                          "SCDN log download fields",
	}

	if desc, ok := descriptions[dataSourceName]; ok {
		return desc
	}

	// Fallback: convert data source name to description
	name := strings.TrimPrefix(dataSourceName, "edgenext_")
	name = strings.Replace(name, "_", " ", -1)
	return name
}
