package template

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/edgenextapisdk/terraform-provider-edgenext/edgenext/connectivity"
	"github.com/edgenextapisdk/terraform-provider-edgenext/edgenext/services/scdn"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// ResourceEdgenextScdnRuleTemplateDomainBind returns the SCDN rule template domain bind resource
func ResourceEdgenextScdnRuleTemplateDomainBind() *schema.Resource {
	return &schema.Resource{
		Create: resourceScdnRuleTemplateDomainBindCreate,
		Read:   resourceScdnRuleTemplateDomainBindRead,
		Delete: resourceScdnRuleTemplateDomainBindDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"template_id": {
				Type:        schema.TypeInt,
				Required:    true,
				ForceNew:    true,
				Description: "The ID of the rule template",
			},
			"domain_ids": {
				Type:        schema.TypeList,
				Required:    true,
				ForceNew:    true,
				Description: "List of domain IDs to bind to the template",
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
			},
			// Computed fields
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The unique identifier for this bind operation",
			},
		},
	}
}

func resourceScdnRuleTemplateDomainBindCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*connectivity.EdgeNextClient)
	service := scdn.NewScdnService(client)

	templateID := d.Get("template_id").(int)
	domainIDsInterface := d.Get("domain_ids").([]interface{})
	domainIDs := make([]int, len(domainIDsInterface))
	for i, v := range domainIDsInterface {
		domainIDs[i] = v.(int)
	}

	req := scdn.RuleTemplateBindDomainRequest{
		ID:        templateID,
		DomainIDs: domainIDs,
	}

	log.Printf("[INFO] Binding domains to SCDN rule template: %+v", req)
	_, err := service.BindRuleTemplateDomains(req)
	if err != nil {
		return fmt.Errorf("failed to bind domains to rule template: %w", err)
	}

	// Create a unique ID for this bind operation
	// Format: template_id:domain_id1,domain_id2,...
	domainIDsStr := make([]string, len(domainIDs))
	for i, id := range domainIDs {
		domainIDsStr[i] = strconv.Itoa(id)
	}
	bindID := fmt.Sprintf("%d:%s", templateID, strings.Join(domainIDsStr, ","))
	d.SetId(bindID)

	log.Printf("[INFO] SCDN rule template domain bind completed successfully: %s", d.Id())
	return resourceScdnRuleTemplateDomainBindRead(d, m)
}

func resourceScdnRuleTemplateDomainBindRead(d *schema.ResourceData, m interface{}) error {
	// For bind operations, we verify the template exists and domains are bound
	client := m.(*connectivity.EdgeNextClient)
	service := scdn.NewScdnService(client)

	var templateID int
	var domainIDs []int
	var err error

	// If we have an ID, parse it (for import scenarios)
	if d.Id() != "" {
		templateID, domainIDs, err = parseBindID(d.Id())
		if err != nil {
			return fmt.Errorf("invalid bind ID: %w", err)
		}
	} else {
		// Otherwise, get from schema
		templateID = d.Get("template_id").(int)
		domainIDsInterface := d.Get("domain_ids").([]interface{})
		domainIDs = make([]int, len(domainIDsInterface))
		for i, v := range domainIDsInterface {
			domainIDs[i] = v.(int)
		}
	}

	// Verify template exists by listing templates
	req := scdn.RuleTemplateListRequest{
		Page:     1,
		PageSize: 1000,
		AppType:  "network_speed", // Default, may need to be configurable
	}

	response, err := service.ListRuleTemplates(req)
	if err != nil {
		return fmt.Errorf("failed to read SCDN rule template: %w", err)
	}

	// Find the template
	var found bool
	for _, tpl := range response.Data.List {
		if tpl.ID == templateID {
			found = true
			break
		}
	}

	if !found {
		log.Printf("[WARN] Rule template not found: %d", templateID)
		d.SetId("")
		return nil
	}

	// Set fields
	if err := d.Set("template_id", templateID); err != nil {
		log.Printf("[WARN] Failed to set template_id: %v", err)
	}

	domainIDsInterface := make([]interface{}, len(domainIDs))
	for i, id := range domainIDs {
		domainIDsInterface[i] = id
	}
	if err := d.Set("domain_ids", domainIDsInterface); err != nil {
		log.Printf("[WARN] Failed to set domain_ids: %v", err)
	}

	return nil
}

func resourceScdnRuleTemplateDomainBindDelete(d *schema.ResourceData, m interface{}) error {
	// For bind operations, delete means unbind the domains
	client := m.(*connectivity.EdgeNextClient)
	service := scdn.NewScdnService(client)

	templateID, domainIDs, err := parseBindID(d.Id())
	if err != nil {
		return fmt.Errorf("invalid bind ID: %w", err)
	}

	req := scdn.RuleTemplateUnbindDomainRequest{
		ID:        templateID,
		DomainIDs: domainIDs,
	}

	log.Printf("[INFO] Unbinding domains from SCDN rule template (delete): %+v", req)
	_, err = service.UnbindRuleTemplateDomains(req)
	if err != nil {
		return fmt.Errorf("failed to unbind domains from rule template: %w", err)
	}

	log.Printf("[INFO] SCDN rule template domain unbind completed (delete): %s", d.Id())
	return nil
}

// parseBindID parses the bind ID format: template_id:domain_id1,domain_id2,...
func parseBindID(id string) (templateID int, domainIDs []int, err error) {
	parts := strings.Split(id, ":")
	if len(parts) != 2 {
		return 0, nil, fmt.Errorf("invalid bind ID format: %s", id)
	}

	templateID, err = strconv.Atoi(parts[0])
	if err != nil {
		return 0, nil, fmt.Errorf("invalid template ID in bind ID: %w", err)
	}

	domainIDsStr := strings.Split(parts[1], ",")
	domainIDs = make([]int, len(domainIDsStr))
	for i, idStr := range domainIDsStr {
		domainID, err := strconv.Atoi(idStr)
		if err != nil {
			return 0, nil, fmt.Errorf("invalid domain ID in bind ID: %w", err)
		}
		domainIDs[i] = domainID
	}

	return templateID, domainIDs, nil
}
