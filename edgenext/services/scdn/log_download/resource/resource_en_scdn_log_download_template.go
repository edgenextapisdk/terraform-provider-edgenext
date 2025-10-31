package resource

import (
	"fmt"
	"log"
	"strconv"

	"github.com/edgenextapisdk/terraform-provider-edgenext/edgenext/connectivity"
	"github.com/edgenextapisdk/terraform-provider-edgenext/edgenext/services/scdn"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// convertSearchTerms converts search_terms from map[string]string or map[string][]string to []map[string]string
func convertSearchTermsForTemplate(searchTerms interface{}) []map[string]interface{} {
	if searchTerms == nil {
		return nil
	}

	// Try to convert from map format (API response)
	if termMap, ok := searchTerms.(map[string]interface{}); ok {
		result := make([]map[string]interface{}, 0)
		for key, value := range termMap {
			// First try map[string]string format (API document format)
			if strValue, ok := value.(string); ok {
				result = append(result, map[string]interface{}{
					"key":   key,
					"value": strValue,
				})
				continue
			}

			// Then try map[string][]string format (array format, for backward compatibility)
			var values []string
			if strSlice, ok := value.([]string); ok {
				values = strSlice
			} else if ifaceSlice, ok := value.([]interface{}); ok {
				values = make([]string, len(ifaceSlice))
				for i, v := range ifaceSlice {
					if str, ok := v.(string); ok {
						values[i] = str
					}
				}
			}
			// Create one entry per value in the array
			for _, val := range values {
				result = append(result, map[string]interface{}{
					"key":   key,
					"value": val,
				})
			}
		}
		return result
	}

	// Try to convert from []LogDownloadSearchTerm format (if already converted)
	if termSlice, ok := searchTerms.([]scdn.LogDownloadSearchTerm); ok {
		result := make([]map[string]interface{}, len(termSlice))
		for i, term := range termSlice {
			result[i] = map[string]interface{}{
				"key":   term.Key,
				"value": term.Value,
			}
		}
		return result
	}

	return nil
}

// ResourceEdgenextScdnLogDownloadTemplate returns the SCDN log download template resource
func ResourceEdgenextScdnLogDownloadTemplate() *schema.Resource {
	return &schema.Resource{
		Create: resourceScdnLogDownloadTemplateCreate,
		Read:   resourceScdnLogDownloadTemplateRead,
		Update: resourceScdnLogDownloadTemplateUpdate,
		Delete: resourceScdnLogDownloadTemplateDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"template_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Template name",
			},
			"group_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Group name",
			},
			"group_id": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Group ID",
			},
			"data_source": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Data source: ng, cc, waf",
			},
			"status": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     1,
				Description: "Status: 1-enabled, 0-disabled, default: 1",
			},
			"download_fields": {
				Type:        schema.TypeList,
				Required:    true,
				Description: "Download fields",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"search_terms": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Search conditions",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Search key",
						},
						"value": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Search value",
						},
					},
				},
			},
			"domain_select_type": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     0,
				Description: "Domain select type: 0-partial, 1-all, default: 0",
			},
			// Computed fields
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The ID of the log download template",
			},
			"template_id": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The template ID",
			},
			"member_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The member ID",
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Creation timestamp",
			},
			"updated_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Last update timestamp",
			},
		},
	}
}

func resourceScdnLogDownloadTemplateCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*connectivity.EdgeNextClient)
	service := scdn.NewScdnService(client)

	// Build request
	req := scdn.LogDownloadTemplateAddRequest{
		TemplateName:     d.Get("template_name").(string),
		GroupName:        d.Get("group_name").(string),
		GroupID:          d.Get("group_id").(int),
		DataSource:       d.Get("data_source").(string),
		DomainSelectType: d.Get("domain_select_type").(int),
	}

	// Status - always set, even if 0
	if v, ok := d.GetOk("status"); ok {
		req.Status = v.(int)
	} else {
		// Use default value if not set
		req.Status = d.Get("status").(int)
	}
	log.Printf("[DEBUG] Setting status to: %d", req.Status)

	// Download fields
	if v, ok := d.GetOk("download_fields"); ok {
		fieldsList := v.([]interface{})
		req.DownloadFields = make([]string, len(fieldsList))
		for i, item := range fieldsList {
			req.DownloadFields[i] = item.(string)
		}
	}

	// Search terms - convert from array format to map format
	if v, ok := d.GetOk("search_terms"); ok {
		termsList := v.([]interface{})
		req.SearchTerms = make(map[string]string)
		for _, item := range termsList {
			termMap := item.(map[string]interface{})
			key := termMap["key"].(string)
			value := termMap["value"].(string)
			// If multiple entries have the same key, the last one will be used
			req.SearchTerms[key] = value
		}
	}

	log.Printf("[INFO] Creating SCDN log download template: %s", req.TemplateName)
	response, err := service.AddLogDownloadTemplate(req)
	if err != nil {
		return fmt.Errorf("failed to create SCDN log download template: %w", err)
	}

	log.Printf("[DEBUG] Log download template creation response: %+v", response)
	log.Printf("[DEBUG] Response Data: %+v", response.Data)

	// Try to extract template_id from response if available
	var templateID int
	if response.Data != nil {
		// Try to parse response.Data as a map
		if dataMap, ok := response.Data.(map[string]interface{}); ok {
			if idVal, ok := dataMap["template_id"]; ok {
				switch v := idVal.(type) {
				case float64:
					templateID = int(v)
				case int:
					templateID = v
				case int64:
					templateID = int(v)
				case string:
					if parsed, err := strconv.Atoi(v); err == nil {
						templateID = parsed
					}
				}
				log.Printf("[DEBUG] Found template_id in response: %d", templateID)
			}
		}
	}

	// If template_id is not in response, query by template_name
	if templateID == 0 {
		log.Printf("[WARN] Template ID not found in creation response, querying by template_name: %s", req.TemplateName)

		// Query the created template to get its ID
		listReq := scdn.LogDownloadTemplateListRequest{
			TemplateName: req.TemplateName,
			Page:         1,
			PerPage:      100,
			Status:       -1, // Query all statuses
		}

		// If group_id is available, use it for more precise query
		if req.GroupID > 0 {
			listReq.GroupID = req.GroupID
		}

		listResponse, err := service.ListLogDownloadTemplates(listReq)
		if err != nil {
			return fmt.Errorf("failed to list log download templates: %w", err)
		}

		log.Printf("[DEBUG] Template list response: total=%v, list_count=%d", listResponse.Data.Total, len(listResponse.Data.List))

		// Find the template by name and group_id
		for _, tpl := range listResponse.Data.List {
			if tpl.TemplateName == req.TemplateName {
				// Try to match group_id
				tplGroupID, err := strconv.Atoi(tpl.GroupID)
				if err == nil && tplGroupID == req.GroupID {
					templateID = tpl.TemplateID
					log.Printf("[INFO] Found template by name and group_id: template_id=%d, template_name=%s", templateID, req.TemplateName)
					break
				} else if tplGroupID == 0 || req.GroupID == 0 {
					// If group_id is 0 or doesn't match, use the first matching template name
					templateID = tpl.TemplateID
					log.Printf("[INFO] Found template by name: template_id=%d, template_name=%s", templateID, req.TemplateName)
					break
				}
			}
		}

		// If still not found, search across multiple pages
		if templateID == 0 {
			log.Printf("[DEBUG] Template not found on first page, searching across multiple pages")
			for page := 1; page <= 5 && templateID == 0; page++ {
				listReq.Page = page
				listResponse, err = service.ListLogDownloadTemplates(listReq)
				if err != nil {
					log.Printf("[WARN] Failed to list templates on page %d: %v", page, err)
					break
				}

				log.Printf("[DEBUG] Searching page %d: found %d templates", page, len(listResponse.Data.List))

				for _, tpl := range listResponse.Data.List {
					if tpl.TemplateName == req.TemplateName {
						tplGroupID, err := strconv.Atoi(tpl.GroupID)
						if err == nil && tplGroupID == req.GroupID {
							templateID = tpl.TemplateID
							log.Printf("[INFO] Found template by name and group_id on page %d: template_id=%d", page, templateID)
							break
						} else if tplGroupID == 0 || req.GroupID == 0 {
							templateID = tpl.TemplateID
							log.Printf("[INFO] Found template by name on page %d: template_id=%d", page, templateID)
							break
						}
					}
				}

				if len(listResponse.Data.List) == 0 {
					break
				}
			}
		}

		if templateID == 0 {
			return fmt.Errorf("failed to find created template by name: %s", req.TemplateName)
		}
	}

	// Set ID
	id := strconv.Itoa(templateID)
	d.SetId(id)
	log.Printf("[DEBUG] Set resource ID to: %d", templateID)

	// Set basic fields from creation request
	if err := d.Set("template_id", templateID); err != nil {
		log.Printf("[WARN] Failed to set template_id: %v", err)
	}
	if err := d.Set("template_name", req.TemplateName); err != nil {
		log.Printf("[WARN] Failed to set template_name: %v", err)
	}
	if err := d.Set("group_id", req.GroupID); err != nil {
		log.Printf("[WARN] Failed to set group_id: %v", err)
	}
	if err := d.Set("group_name", req.GroupName); err != nil {
		log.Printf("[WARN] Failed to set group_name: %v", err)
	}
	if err := d.Set("data_source", req.DataSource); err != nil {
		log.Printf("[WARN] Failed to set data_source: %v", err)
	}
	if err := d.Set("status", req.Status); err != nil {
		log.Printf("[WARN] Failed to set status: %v", err)
	}
	if err := d.Set("domain_select_type", req.DomainSelectType); err != nil {
		log.Printf("[WARN] Failed to set domain_select_type: %v", err)
	}
	if err := d.Set("download_fields", req.DownloadFields); err != nil {
		log.Printf("[WARN] Failed to set download_fields: %v", err)
	}
	if len(req.SearchTerms) > 0 {
		searchTermsList := make([]map[string]interface{}, 0, len(req.SearchTerms))
		for key, value := range req.SearchTerms {
			searchTermsList = append(searchTermsList, map[string]interface{}{
				"key":   key,
				"value": value,
			})
		}
		if err := d.Set("search_terms", searchTermsList); err != nil {
			log.Printf("[WARN] Failed to set search_terms: %v", err)
		}
	} else {
		if err := d.Set("search_terms", []interface{}{}); err != nil {
			log.Printf("[WARN] Failed to set empty search_terms: %v", err)
		}
	}

	log.Printf("[INFO] SCDN log download template created successfully: %s", d.Id())

	// Always call read to get full details from API
	return resourceScdnLogDownloadTemplateRead(d, m)
}

func resourceScdnLogDownloadTemplateRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*connectivity.EdgeNextClient)
	service := scdn.NewScdnService(client)

	templateID, err := strconv.Atoi(d.Id())
	if err != nil {
		return fmt.Errorf("invalid template ID: %w", err)
	}

	// Query template list to find the template
	req := scdn.LogDownloadTemplateListRequest{
		Page:    1,
		PerPage: 100,
		Status:  -1, // Query all statuses
	}

	// If template_name is available in state, use it for more precise query
	templateName := ""
	if templateNameVal, ok := d.GetOk("template_name"); ok && templateNameVal.(string) != "" {
		templateName = templateNameVal.(string)
		req.TemplateName = templateName
		log.Printf("[DEBUG] Querying log download template by name: %s (template_id=%d)", req.TemplateName, templateID)
	}

	// If group_id is available, use it for more precise query
	if groupIDVal, ok := d.GetOk("group_id"); ok {
		req.GroupID = groupIDVal.(int)
	}

	log.Printf("[DEBUG] Querying log download templates: page=%d, per_page=%d, template_name=%s, group_id=%d", req.Page, req.PerPage, req.TemplateName, req.GroupID)
	response, err := service.ListLogDownloadTemplates(req)
	if err != nil {
		return fmt.Errorf("failed to list log download templates: %w", err)
	}

	log.Printf("[DEBUG] Template list response: total=%v, list_count=%d", response.Data.Total, len(response.Data.List))

	// Find the template by ID
	var template *scdn.LogDownloadTemplateInfo
	for i := range response.Data.List {
		log.Printf("[DEBUG] Checking template: template_id=%d, template_name=%s", response.Data.List[i].TemplateID, response.Data.List[i].TemplateName)
		if response.Data.List[i].TemplateID == templateID {
			template = &response.Data.List[i]
			log.Printf("[DEBUG] Found matching template: template_id=%d", templateID)
			break
		}
	}

	// If not found by ID, try searching by template_name if available
	if template == nil && templateName != "" {
		log.Printf("[DEBUG] Template not found by ID, trying to find by name: %s", templateName)
		for i := range response.Data.List {
			if response.Data.List[i].TemplateName == templateName {
				template = &response.Data.List[i]
				log.Printf("[DEBUG] Found template by name: template_id=%d, template_name=%s", template.TemplateID, template.TemplateName)
				// Update the resource ID to the found template's ID
				if template.TemplateID != templateID {
					log.Printf("[INFO] Updating resource ID from %d to %d based on template name match", templateID, template.TemplateID)
					d.SetId(strconv.Itoa(template.TemplateID))
					templateID = template.TemplateID
				}
				break
			}
		}
	}

	// If still not found, try searching without template_name filter and across multiple pages
	if template == nil {
		log.Printf("[DEBUG] Template not found with current filters, trying without template_name filter")
		req.TemplateName = ""
		req.Page = 1
		req.PerPage = 100

		// Search up to 5 pages
		for page := 1; page <= 5 && template == nil; page++ {
			req.Page = page
			response, err = service.ListLogDownloadTemplates(req)
			if err != nil {
				log.Printf("[WARN] Failed to list templates on page %d: %v", page, err)
				break
			}

			log.Printf("[DEBUG] Searching page %d: found %d templates", page, len(response.Data.List))

			// Search by ID first
			for i := range response.Data.List {
				if response.Data.List[i].TemplateID == templateID {
					template = &response.Data.List[i]
					log.Printf("[DEBUG] Found template by ID on page %d: template_id=%d", page, templateID)
					break
				}
			}

			// If not found by ID and we have template_name, search by name
			if template == nil && templateName != "" {
				for i := range response.Data.List {
					if response.Data.List[i].TemplateName == templateName {
						template = &response.Data.List[i]
						log.Printf("[DEBUG] Found template by name on page %d: template_id=%d, template_name=%s", page, template.TemplateID, template.TemplateName)
						// Update the resource ID to the found template's ID
						if template.TemplateID != templateID {
							log.Printf("[INFO] Updating resource ID from %d to %d based on template name match", templateID, template.TemplateID)
							d.SetId(strconv.Itoa(template.TemplateID))
							templateID = template.TemplateID
						}
						break
					}
				}
			}

			// If no more templates on this page, stop searching
			if len(response.Data.List) == 0 {
				break
			}
		}
	}

	if template == nil {
		// Template not found after exhaustive search
		log.Printf("[WARN] Log download template %d not found after searching multiple pages", templateID)

		// If templateID is 0, it means the resource was never properly created
		// Clear the ID to mark the resource for creation
		if templateID == 0 {
			log.Printf("[INFO] Template ID is 0, clearing resource ID to mark for creation")
			d.SetId("")
			return nil
		}

		// If templateID is not 0 but template not found, it might have been deleted
		// However, we should not clear the ID immediately - it might be a timing issue
		// Return an error so Terraform knows the resource state is inconsistent
		return fmt.Errorf("log download template %d not found - template may have been deleted or there is an API issue", templateID)
	}

	// Set all fields
	if err := d.Set("template_id", template.TemplateID); err != nil {
		return fmt.Errorf("error setting template_id: %w", err)
	}
	if err := d.Set("template_name", template.TemplateName); err != nil {
		return fmt.Errorf("error setting template_name: %w", err)
	}

	// group_id is returned as string from API, convert to int
	groupID := 0
	if template.GroupID != "" {
		if parsed, err := strconv.Atoi(template.GroupID); err == nil {
			groupID = parsed
		}
	}
	if err := d.Set("group_id", groupID); err != nil {
		return fmt.Errorf("error setting group_id: %w", err)
	}
	if err := d.Set("group_name", template.GroupName); err != nil {
		return fmt.Errorf("error setting group_name: %w", err)
	}
	if err := d.Set("data_source", template.DataSource); err != nil {
		return fmt.Errorf("error setting data_source: %w", err)
	}

	// status is returned as int from API
	if err := d.Set("status", template.Status); err != nil {
		return fmt.Errorf("error setting status: %w", err)
	}
	if template.MemberID != "" {
		if err := d.Set("member_id", template.MemberID); err != nil {
			log.Printf("[WARN] Failed to set member_id: %v", err)
		}
	}
	if err := d.Set("created_at", template.CreatedAt); err != nil {
		log.Printf("[WARN] Failed to set created_at: %v", err)
	}
	if err := d.Set("updated_at", template.UpdatedAt); err != nil {
		log.Printf("[WARN] Failed to set updated_at: %v", err)
	}

	// Download fields
	if err := d.Set("download_fields", template.DownloadFields); err != nil {
		return fmt.Errorf("error setting download_fields: %w", err)
	}

	// Search terms - convert from map format to array format
	searchTerms := convertSearchTermsForTemplate(template.SearchTerms)
	if len(searchTerms) > 0 {
		if err := d.Set("search_terms", searchTerms); err != nil {
			log.Printf("[WARN] Failed to set search_terms: %v", err)
		}
	} else {
		// Set empty list if search_terms is nil or empty
		if err := d.Set("search_terms", []interface{}{}); err != nil {
			log.Printf("[WARN] Failed to set empty search_terms: %v", err)
		}
	}

	log.Printf("[INFO] Log download template read successfully: template_id=%d", template.TemplateID)
	return nil
}

func resourceScdnLogDownloadTemplateUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*connectivity.EdgeNextClient)
	service := scdn.NewScdnService(client)

	templateID, err := strconv.Atoi(d.Id())
	if err != nil {
		return fmt.Errorf("invalid template ID: %w", err)
	}

	// Build request
	req := scdn.LogDownloadTemplateSaveRequest{
		TemplateID:       templateID,
		TemplateName:     d.Get("template_name").(string),
		GroupName:        d.Get("group_name").(string),
		GroupID:          d.Get("group_id").(int),
		DataSource:       d.Get("data_source").(string),
		Status:           d.Get("status").(int), // Always set status, even if 0
		DomainSelectType: d.Get("domain_select_type").(int),
	}

	log.Printf("[DEBUG] Updating template with status: %d", req.Status)

	// Download fields
	if v, ok := d.GetOk("download_fields"); ok {
		fieldsList := v.([]interface{})
		req.DownloadFields = make([]string, len(fieldsList))
		for i, item := range fieldsList {
			req.DownloadFields[i] = item.(string)
		}
	}

	// Search terms - convert from array format to map format
	if v, ok := d.GetOk("search_terms"); ok {
		termsList := v.([]interface{})
		req.SearchTerms = make(map[string]string)
		for _, item := range termsList {
			termMap := item.(map[string]interface{})
			key := termMap["key"].(string)
			value := termMap["value"].(string)
			// If multiple entries have the same key, the last one will be used
			req.SearchTerms[key] = value
		}
	}

	log.Printf("[INFO] Updating SCDN log download template: %d", templateID)
	_, err = service.SaveLogDownloadTemplate(req)
	if err != nil {
		return fmt.Errorf("failed to update SCDN log download template: %w", err)
	}

	log.Printf("[INFO] SCDN log download template updated successfully: %d", templateID)
	return resourceScdnLogDownloadTemplateRead(d, m)
}

func resourceScdnLogDownloadTemplateDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*connectivity.EdgeNextClient)
	service := scdn.NewScdnService(client)

	templateID, err := strconv.Atoi(d.Id())
	if err != nil {
		return fmt.Errorf("invalid template ID: %w", err)
	}

	req := scdn.LogDownloadTemplateDeleteRequest{
		TemplateID: templateID,
	}

	log.Printf("[INFO] Deleting SCDN log download template: %d", templateID)
	_, err = service.DeleteLogDownloadTemplate(req)
	if err != nil {
		return fmt.Errorf("failed to delete SCDN log download template: %w", err)
	}

	log.Printf("[INFO] SCDN log download template deleted successfully: %d", templateID)
	d.SetId("")
	return nil
}
