package cache

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/edgenextapisdk/terraform-provider-edgenext/edgenext/connectivity"
	"github.com/edgenextapisdk/terraform-provider-edgenext/edgenext/services/scdn"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// ResourceEdgenextScdnCacheRulesSort returns the SCDN cache rules sort resource
func ResourceEdgenextScdnCacheRulesSort() *schema.Resource {
	return &schema.Resource{
		Create: resourceScdnCacheRulesSortCreate,
		Read:   resourceScdnCacheRulesSortRead,
		Update: resourceScdnCacheRulesSortUpdate,
		Delete: resourceScdnCacheRulesSortDelete,

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
			"ids": {
				Type:        schema.TypeList,
				Required:    true,
				Description: "Sorted rule IDs array (order matters)",
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
			},
			// Computed fields
			"sorted_ids": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Sorted rule IDs after sorting",
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
			},
		},
	}
}

func resourceScdnCacheRulesSortCreate(d *schema.ResourceData, m interface{}) error {
	return resourceScdnCacheRulesSortUpdate(d, m)
}

func resourceScdnCacheRulesSortRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*connectivity.EdgeNextClient)
	service := scdn.NewScdnService(client)

	// Parse ID (format: business_id-business_type)
	id := d.Id()
	var businessID int
	var businessType string

	if id != "" {
		// Parse composite ID
		parts := strings.Split(id, "-")
		if len(parts) >= 2 {
			var err error
			businessID, err = strconv.Atoi(parts[0])
			if err != nil {
				log.Printf("[WARN] Failed to parse business_id from ID: %v", err)
			}
			businessType = parts[1]
		}
	}

	// Fallback to state if ID parsing failed
	if businessID == 0 {
		businessID = d.Get("business_id").(int)
	}
	if businessType == "" {
		businessType = d.Get("business_type").(string)
	}

	// Get current rules to verify order
	req := scdn.CacheRuleGetRulesRequest{
		BusinessID:   businessID,
		BusinessType: businessType,
	}

	log.Printf("[INFO] Reading SCDN cache rules to verify sort order: business_id=%d, business_type=%s", businessID, businessType)
	response, err := service.GetCacheRules(req)
	if err != nil {
		return fmt.Errorf("failed to read cache rules: %w", err)
	}

	// Extract current rule IDs in order
	currentIDs := make([]int, 0, len(response.Data.List))
	for _, rule := range response.Data.List {
		currentIDs = append(currentIDs, rule.ID)
	}

	// Set resource ID
	d.SetId(fmt.Sprintf("%d-%s", businessID, businessType))

	// Set fields
	if err := d.Set("business_id", businessID); err != nil {
		log.Printf("[WARN] Failed to set business_id: %v", err)
	}
	if err := d.Set("business_type", businessType); err != nil {
		log.Printf("[WARN] Failed to set business_type: %v", err)
	}
	if err := d.Set("ids", currentIDs); err != nil {
		log.Printf("[WARN] Failed to set ids: %v", err)
	}
	if err := d.Set("sorted_ids", currentIDs); err != nil {
		log.Printf("[WARN] Failed to set sorted_ids: %v", err)
	}

	log.Printf("[INFO] Cache rules sort read successfully: current_ids=%v", currentIDs)
	return nil
}

func resourceScdnCacheRulesSortUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*connectivity.EdgeNextClient)
	service := scdn.NewScdnService(client)

	businessID := d.Get("business_id").(int)
	businessType := d.Get("business_type").(string)

	// Get IDs list
	idsList := d.Get("ids").([]interface{})
	ids := make([]int, len(idsList))
	for i, v := range idsList {
		ids[i] = v.(int)
	}

	req := scdn.CacheRuleSortRequest{
		BusinessID:   businessID,
		BusinessType: businessType,
		IDs:          ids,
	}

	log.Printf("[INFO] Sorting SCDN cache rules: business_id=%d, business_type=%s, ids=%v", businessID, businessType, ids)
	response, err := service.SortCacheRules(req)
	if err != nil {
		return fmt.Errorf("failed to sort cache rules: %w", err)
	}

	// Set resource ID
	d.SetId(fmt.Sprintf("%d-%s", businessID, businessType))

	// Set computed sorted_ids
	if err := d.Set("sorted_ids", response.Data.IDs); err != nil {
		log.Printf("[WARN] Failed to set sorted_ids: %v", err)
	}

	return resourceScdnCacheRulesSortRead(d, m)
}

func resourceScdnCacheRulesSortDelete(d *schema.ResourceData, m interface{}) error {
	// Sorting is an idempotent operation, so deletion is a no-op.
	// The resource is simply removed from the state.
	log.Printf("[INFO] Deleting cache rules sort resource from state: %s", d.Id())
	d.SetId("") // Mark as deleted
	return nil
}
