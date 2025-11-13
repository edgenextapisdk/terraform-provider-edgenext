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

// ResourceEdgenextScdnRuleTemplateDomainUnbind returns the SCDN rule template domain unbind resource
func ResourceEdgenextScdnRuleTemplateDomainUnbind() *schema.Resource {
	return &schema.Resource{
		Create: resourceScdnRuleTemplateDomainUnbindCreate,
		Read:   resourceScdnRuleTemplateDomainUnbindRead,
		Delete: resourceScdnRuleTemplateDomainUnbindDelete,

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
				Description: "List of domain IDs to unbind from the template",
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
			},
			// Computed fields
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The unique identifier for this unbind operation",
			},
		},
	}
}

func resourceScdnRuleTemplateDomainUnbindCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*connectivity.EdgeNextClient)
	service := scdn.NewScdnService(client)

	templateID := d.Get("template_id").(int)
	domainIDsInterface := d.Get("domain_ids").([]interface{})
	domainIDs := make([]int, len(domainIDsInterface))
	for i, v := range domainIDsInterface {
		domainIDs[i] = v.(int)
	}

	req := scdn.RuleTemplateUnbindDomainRequest{
		ID:        templateID,
		DomainIDs: domainIDs,
	}

	log.Printf("[INFO] Unbinding domains from SCDN rule template: %+v", req)
	_, err := service.UnbindRuleTemplateDomains(req)
	if err != nil {
		return fmt.Errorf("failed to unbind domains from rule template: %w", err)
	}

	// Create a unique ID for this unbind operation
	// Format: template_id:domain_id1,domain_id2,...
	domainIDsStr := make([]string, len(domainIDs))
	for i, id := range domainIDs {
		domainIDsStr[i] = strconv.Itoa(id)
	}
	unbindID := fmt.Sprintf("%d:%s", templateID, strings.Join(domainIDsStr, ","))
	d.SetId(unbindID)

	log.Printf("[INFO] SCDN rule template domain unbind completed successfully: %s", d.Id())
	return resourceScdnRuleTemplateDomainUnbindRead(d, m)
}

func resourceScdnRuleTemplateDomainUnbindRead(d *schema.ResourceData, m interface{}) error {
	// For unbind operations, we just verify the template exists
	// The domains are already unbound, so we don't need to check their binding status
	client := m.(*connectivity.EdgeNextClient)
	service := scdn.NewScdnService(client)

	var templateID int
	var domainIDs []int
	var err error

	// If we have an ID, parse it (for import scenarios)
	if d.Id() != "" {
		templateID, domainIDs, err = parseUnbindID(d.Id())
		if err != nil {
			return fmt.Errorf("invalid unbind ID: %w", err)
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

func resourceScdnRuleTemplateDomainUnbindDelete(d *schema.ResourceData, m interface{}) error {
	// For unbind operations, delete is a no-op since unbinding is already done
	// The resource is removed from state, which is the desired behavior
	log.Printf("[INFO] Unbind operation resource deleted from state: %s", d.Id())
	return nil
}

// parseUnbindID parses the unbind ID format: template_id:domain_id1,domain_id2,...
func parseUnbindID(id string) (templateID int, domainIDs []int, err error) {
	parts := strings.Split(id, ":")
	if len(parts) != 2 {
		return 0, nil, fmt.Errorf("invalid unbind ID format: %s", id)
	}

	templateID, err = strconv.Atoi(parts[0])
	if err != nil {
		return 0, nil, fmt.Errorf("invalid template ID in unbind ID: %w", err)
	}

	domainIDsStr := strings.Split(parts[1], ",")
	domainIDs = make([]int, len(domainIDsStr))
	for i, idStr := range domainIDsStr {
		domainID, err := strconv.Atoi(idStr)
		if err != nil {
			return 0, nil, fmt.Errorf("invalid domain ID in unbind ID: %w", err)
		}
		domainIDs[i] = domainID
	}

	return templateID, domainIDs, nil
}
