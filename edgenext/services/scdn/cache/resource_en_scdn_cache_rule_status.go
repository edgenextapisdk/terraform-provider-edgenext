package cache

import (
	"fmt"
	"log"

	"github.com/edgenextapisdk/terraform-provider-edgenext/edgenext/connectivity"
	"github.com/edgenextapisdk/terraform-provider-edgenext/edgenext/services/scdn"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// ResourceEdgenextScdnCacheRuleStatus returns the SCDN cache rule status resource
func ResourceEdgenextScdnCacheRuleStatus() *schema.Resource {
	return &schema.Resource{
		Create: resourceScdnCacheRuleStatusCreate,
		Read:   resourceScdnCacheRuleStatusRead,
		Update: resourceScdnCacheRuleStatusUpdate,
		Delete: resourceScdnCacheRuleStatusDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"business_id": {
				Type:        schema.TypeInt,
				Required:    true,
				ForceNew:    true,
				Description: "Business ID (template ID for 'tpl' type, domain ID for 'domain' type)",
			},
			"business_type": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Business type: 'tpl' (template) or 'domain'",
			},
			"rule_ids": {
				Type:        schema.TypeList,
				Required:    true,
				Description: "Rule IDs array to update status",
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
			},
			"status": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Status: 1 (enabled) or 2 (disabled)",
			},
			// Computed fields
			"updated_ids": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Rule IDs that were updated",
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
			},
		},
	}
}

func resourceScdnCacheRuleStatusCreate(d *schema.ResourceData, m interface{}) error {
	return resourceScdnCacheRuleStatusUpdate(d, m)
}

func resourceScdnCacheRuleStatusRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*connectivity.EdgeNextClient)
	service := scdn.NewScdnService(client)

	businessID := d.Get("business_id").(int)
	businessType := d.Get("business_type").(string)

	// Get rules to verify status
	req := scdn.CacheRuleGetRulesRequest{
		BusinessID:   businessID,
		BusinessType: businessType,
	}

	log.Printf("[INFO] Reading SCDN cache rules to verify status: business_id=%d, business_type=%s", businessID, businessType)
	response, err := service.GetCacheRules(req)
	if err != nil {
		return fmt.Errorf("failed to read cache rules: %w", err)
	}

	// Get rule IDs from state
	ruleIDsList := d.Get("rule_ids").([]interface{})
	ruleIDs := make([]int, len(ruleIDsList))
	for i, v := range ruleIDsList {
		ruleIDs[i] = v.(int)
	}

	// Verify rules exist and get their status
	updatedIDs := make([]int, 0)
	for _, ruleID := range ruleIDs {
		for _, rule := range response.Data.List {
			if rule.ID == ruleID {
				updatedIDs = append(updatedIDs, ruleID)
				break
			}
		}
	}

	// Set resource ID
	d.SetId(fmt.Sprintf("%d-%s", businessID, businessType))

	// Set computed fields
	if err := d.Set("updated_ids", updatedIDs); err != nil {
		log.Printf("[WARN] Failed to set updated_ids: %v", err)
	}

	log.Printf("[INFO] Cache rule status read successfully: updated_ids=%v", updatedIDs)
	return nil
}

func resourceScdnCacheRuleStatusUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*connectivity.EdgeNextClient)
	service := scdn.NewScdnService(client)

	businessID := d.Get("business_id").(int)
	businessType := d.Get("business_type").(string)
	status := d.Get("status").(int)

	// Get rule IDs list
	ruleIDsList := d.Get("rule_ids").([]interface{})
	ruleIDs := make([]int, len(ruleIDsList))
	for i, v := range ruleIDsList {
		ruleIDs[i] = v.(int)
	}

	req := scdn.CacheRuleUpdateStatusRequest{
		BusinessID:   businessID,
		BusinessType: businessType,
		IDs:          ruleIDs,
		Status:       status,
	}

	log.Printf("[INFO] Updating SCDN cache rule status: business_id=%d, business_type=%s, rule_ids=%v, status=%d", businessID, businessType, ruleIDs, status)
	response, err := service.UpdateCacheRuleStatus(req)
	if err != nil {
		return fmt.Errorf("failed to update cache rule status: %w", err)
	}

	// Set resource ID
	d.SetId(fmt.Sprintf("%d-%s", businessID, businessType))

	// Set computed fields
	if err := d.Set("updated_ids", response.Data.IDs); err != nil {
		log.Printf("[WARN] Failed to set updated_ids: %v", err)
	}

	log.Printf("[INFO] Cache rule status updated successfully: updated_ids=%v", response.Data.IDs)
	return nil
}

func resourceScdnCacheRuleStatusDelete(d *schema.ResourceData, m interface{}) error {
	// Status update is idempotent, so deletion is a no-op.
	// The resource is simply removed from the state.
	log.Printf("[INFO] Deleting cache rule status resource from state: %s", d.Id())
	d.SetId("") // Mark as deleted
	return nil
}
