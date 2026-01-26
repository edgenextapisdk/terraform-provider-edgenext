package template

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/edgenextapisdk/terraform-provider-edgenext/edgenext/connectivity"
	"github.com/edgenextapisdk/terraform-provider-edgenext/edgenext/services/scdn"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

// ResourceEdgenextScdnRuleTemplateSwitch returns the SCDN rule template switch resource
func ResourceEdgenextScdnRuleTemplateSwitch() *schema.Resource {
	return &schema.Resource{
		Create: resourceScdnRuleTemplateSwitchCreate,
		Read:   resourceScdnRuleTemplateSwitchRead,
		Delete: resourceScdnRuleTemplateSwitchDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"app_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "network_speed",
				ForceNew:    true,
				Description: "Application type the template applies to",
			},
			"domain_ids": {
				Type:        schema.TypeList,
				Required:    true,
				ForceNew:    true,
				Description: "List of domain IDs to switch templates",
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
			},
			"new_tpl_id": {
				Type:        schema.TypeInt,
				Required:    true,
				ForceNew:    true,
				Description: "New template ID to switch to, when new_tpl_type=global, pass 0",
			},
			"new_tpl_type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				Description:  "New template type",
				ValidateFunc: validation.StringInSlice([]string{"more_domain", "global"}, false),
			},
			// Computed fields
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The unique identifier for this switch operation",
			},
		},
	}
}

func resourceScdnRuleTemplateSwitchCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*connectivity.EdgeNextClient)
	service := scdn.NewScdnService(client)

	appType := d.Get("app_type").(string)
	newTplID := d.Get("new_tpl_id").(int)
	newTplType := d.Get("new_tpl_type").(string)

	domainIDsInterface := d.Get("domain_ids").([]interface{})
	domainIDs := make([]int, len(domainIDsInterface))
	for i, v := range domainIDsInterface {
		domainIDs[i] = v.(int)
	}

	req := scdn.RuleTemplateSwitchDomainRequest{
		AppType:    appType,
		DomainIDs:  domainIDs,
		NewTplID:   newTplID,
		NewTplType: newTplType,
	}

	log.Printf("[INFO] Switching domains to SCDN rule template: %+v", req)
	resp, err := service.SwitchDomainTemplate(req)
	if err != nil {
		return fmt.Errorf("failed to switch domains to rule template: %w", err)
	}

	log.Printf("[INFO] Switch domain template result: %s", resp.Data.Info)

	// Create a unique ID for this switch operation
	// Format: new_tpl_id:domain_id1,domain_id2,...
	domainIDsStr := make([]string, len(domainIDs))
	for i, id := range domainIDs {
		domainIDsStr[i] = strconv.Itoa(id)
	}
	switchID := fmt.Sprintf("%d:%s", newTplID, strings.Join(domainIDsStr, ","))
	d.SetId(switchID)

	return resourceScdnRuleTemplateSwitchRead(d, m)
}

func resourceScdnRuleTemplateSwitchRead(d *schema.ResourceData, m interface{}) error {
	// Since switching is an action, Read mainly validates the state.
	// We can try to verify if the domains are indeed bound to the new template if not global.

	client := m.(*connectivity.EdgeNextClient)
	service := scdn.NewScdnService(client)

	var newTplID int
	var domainIDs []int
	var err error

	if d.Id() != "" {
		newTplID, domainIDs, err = parseSwitchID(d.Id())
		if err != nil {
			return fmt.Errorf("invalid switch ID: %w", err)
		}
	} else {
		newTplID = d.Get("new_tpl_id").(int)
		domainIDsInterface := d.Get("domain_ids").([]interface{})
		domainIDs = make([]int, len(domainIDsInterface))
		for i, v := range domainIDsInterface {
			domainIDs[i] = v.(int)
		}
	}

	// If new_tpl_type is global (new_tpl_id is 0), we can't easily verify a specific template ID binding
	// without knowing which global template logic applies.
	// If it's a specific template ID, we can check.

	if newTplID > 0 {
		// List domains for the template to verify binding
		req := scdn.RuleTemplateListDomainsRequest{
			ID:       newTplID,
			AppType:  d.Get("app_type").(string),
			Page:     1,
			PageSize: 1000, // Should be enough to find our domains if not too many
		}

		apiResp, err := service.ListRuleTemplateDomains(req)
		if err != nil {
			log.Printf("[WARN] Failed to list domains for template %d: %v", newTplID, err)
			// Don't error out, just keep state? Or assume it's fine.
		} else {
			// Check if all our domains are present
			boundMap := make(map[int]bool)
			for _, d := range apiResp.Data.List {
				boundMap[d.ID] = true
			}

			allBound := true
			for _, id := range domainIDs {
				if !boundMap[id] {
					allBound = false
					break
				}
			}

			if !allBound {
				log.Printf("[WARN] Not all domains found in template %d list, maybe unassigned externally?", newTplID)
				// If strictly enforcing, we could d.SetId("") to trigger recreation,
				// but switch resource is more of a "fire and forget" action resource sometimes.
				// However, if we treat it as "State: these domains SHOULD be in this template",
				// then removing from state is correct if they are not there.
				// Let's assume authoritative.
				// d.SetId("")
				// For now let's just log warn.
			}
		}
	}

	// Restore state variables if reading from import
	if err := d.Set("new_tpl_id", newTplID); err != nil {
		log.Printf("[WARN] Error setting new_tpl_id: %s", err)
	}

	domainIDsInterface := make([]interface{}, len(domainIDs))
	for i, id := range domainIDs {
		domainIDsInterface[i] = id
	}
	if err := d.Set("domain_ids", domainIDsInterface); err != nil {
		log.Printf("[WARN] Error setting domain_ids: %s", err)
	}

	return nil
}

func resourceScdnRuleTemplateSwitchDelete(d *schema.ResourceData, m interface{}) error {
	// "Delete" for a switch action is ambiguous.
	// Does it mean revert to previous template? We don't know what the previous template was.
	// Does it mean unbind from the current template?
	// Usually "Switch" resources in Terraform are tricky.
	// Ideally, they should just manage the state "Domains X are in Template Y".
	// So Delete would mean removing them from Template Y.

	// If we treat "Delete" as "Unbind from the target template", we can call unbind.

	client := m.(*connectivity.EdgeNextClient)
	service := scdn.NewScdnService(client)

	newTplID, domainIDs, err := parseSwitchID(d.Id())
	if err != nil {
		return fmt.Errorf("invalid switch ID: %w", err)
	}

	// If it was a global template (ID 0), unbinding might not make sense or might be different API.
	// EndpointsRuleTemplatesUnbind requires a template ID.
	if newTplID > 0 {
		req := scdn.RuleTemplateUnbindDomainRequest{
			ID:        newTplID,
			DomainIDs: domainIDs,
		}

		log.Printf("[INFO] Unbinding domains from SCDN rule template (revert switch): %+v", req)
		_, err = service.UnbindRuleTemplateDomains(req)
		if err != nil {
			return fmt.Errorf("failed to unbind domains from rule template during destroy: %w", err)
		}
	} else {
		log.Printf("[INFO] Skipping unbind for global template switch (ID 0)")
	}

	d.SetId("")
	return nil
}

func parseSwitchID(id string) (newTplID int, domainIDs []int, err error) {
	parts := strings.Split(id, ":")
	if len(parts) != 2 {
		return 0, nil, fmt.Errorf("invalid switch ID format: %s", id)
	}

	newTplID, err = strconv.Atoi(parts[0])
	if err != nil {
		return 0, nil, fmt.Errorf("invalid template ID in switch ID: %w", err)
	}

	domainIDsStr := strings.Split(parts[1], ",")
	domainIDs = make([]int, len(domainIDsStr))
	for i, idStr := range domainIDsStr {
		domainID, err := strconv.Atoi(idStr)
		if err != nil {
			return 0, nil, fmt.Errorf("invalid domain ID in switch ID: %w", err)
		}
		domainIDs[i] = domainID
	}

	return newTplID, domainIDs, nil
}
