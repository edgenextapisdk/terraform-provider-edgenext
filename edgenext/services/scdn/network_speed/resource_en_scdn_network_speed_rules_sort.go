package networkspeed

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/edgenextapisdk/terraform-provider-edgenext/edgenext/connectivity"
	"github.com/edgenextapisdk/terraform-provider-edgenext/edgenext/services/scdn"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// ResourceEdgenextScdnNetworkSpeedRulesSort returns the SCDN network speed rules sort resource
func ResourceEdgenextScdnNetworkSpeedRulesSort() *schema.Resource {
	return &schema.Resource{
		Create: resourceScdnNetworkSpeedRulesSortCreate,
		Read:   resourceScdnNetworkSpeedRulesSortRead,
		Update: resourceScdnNetworkSpeedRulesSortUpdate,
		Delete: resourceScdnNetworkSpeedRulesSortDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"business_id": {
				Type:        schema.TypeInt,
				Required:    true,
				ForceNew:    true,
				Description: "Business ID (template ID for 'tpl' type, user ID for 'global' type)",
			},
			"business_type": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Business type: 'tpl' (template) or 'global'",
			},
			"config_group": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Rule group: 'custom_page', 'upstream_uri_change_rule', 'resp_headers_rule', or 'customized_req_headers_rule'",
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

func resourceScdnNetworkSpeedRulesSortCreate(d *schema.ResourceData, m interface{}) error {
	return resourceScdnNetworkSpeedRulesSortUpdate(d, m)
}

func resourceScdnNetworkSpeedRulesSortRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*connectivity.EdgeNextClient)
	service := scdn.NewScdnService(client)

	// Parse ID (format: business_id-business_type-config_group)
	id := d.Id()
	var businessID int
	var businessType string
	var configGroup string

	if id != "" {
		// Parse composite ID
		parts := strings.Split(id, "-")
		if len(parts) >= 3 {
			var err error
			businessID, err = strconv.Atoi(parts[0])
			if err != nil {
				log.Printf("[WARN] Failed to parse business_id from ID: %v", err)
			}
			businessType = parts[1]
			configGroup = parts[2]
		}
	}

	// Fallback to state if ID parsing failed
	if businessID == 0 {
		businessID = d.Get("business_id").(int)
	}
	if businessType == "" {
		businessType = d.Get("business_type").(string)
	}
	if configGroup == "" {
		configGroup = d.Get("config_group").(string)
	}

	// Get current rules to verify order
	req := scdn.NetworkSpeedGetRulesRequest{
		BusinessID:   businessID,
		BusinessType: businessType,
		ConfigGroup:  configGroup,
	}

	log.Printf("[INFO] Reading SCDN network speed rules to verify sort order: business_id=%d, business_type=%s, config_group=%s", businessID, businessType, configGroup)
	response, err := service.GetNetworkSpeedRules(req)
	if err != nil {
		return fmt.Errorf("failed to read network speed rules: %w", err)
	}

	// Extract current rule IDs in order
	currentIDs := make([]int, 0, len(response.Data.List))
	for _, rule := range response.Data.List {
		currentIDs = append(currentIDs, rule.ID)
	}

	// Set fields
	if err := d.Set("business_id", businessID); err != nil {
		log.Printf("[WARN] Failed to set business_id: %v", err)
	}
	if err := d.Set("business_type", businessType); err != nil {
		log.Printf("[WARN] Failed to set business_type: %v", err)
	}
	if err := d.Set("config_group", configGroup); err != nil {
		log.Printf("[WARN] Failed to set config_group: %v", err)
	}
	if err := d.Set("ids", currentIDs); err != nil {
		log.Printf("[WARN] Failed to set ids: %v", err)
	}
	if err := d.Set("sorted_ids", currentIDs); err != nil {
		log.Printf("[WARN] Failed to set sorted_ids: %v", err)
	}

	log.Printf("[INFO] Network speed rules sort read successfully: current_ids=%v", currentIDs)
	return nil
}

func resourceScdnNetworkSpeedRulesSortUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*connectivity.EdgeNextClient)
	service := scdn.NewScdnService(client)

	businessID := d.Get("business_id").(int)
	businessType := d.Get("business_type").(string)
	configGroup := d.Get("config_group").(string)

	// Get IDs list
	idsList := d.Get("ids").([]interface{})
	ids := make([]int, len(idsList))
	for i, v := range idsList {
		ids[i] = v.(int)
	}

	req := scdn.NetworkSpeedSortRulesRequest{
		BusinessID:   businessID,
		BusinessType: businessType,
		ConfigGroup:  configGroup,
		IDs:          ids,
	}

	log.Printf("[INFO] Sorting SCDN network speed rules: business_id=%d, business_type=%s, config_group=%s, ids=%v", businessID, businessType, configGroup, ids)
	response, err := service.SortNetworkSpeedRules(req)
	if err != nil {
		return fmt.Errorf("failed to sort network speed rules: %w", err)
	}

	// Set resource ID (composite format: business_id-business_type-config_group)
	d.SetId(fmt.Sprintf("%d-%s-%s", businessID, businessType, configGroup))

	// Set sorted IDs
	if err := d.Set("sorted_ids", response.Data.IDs); err != nil {
		log.Printf("[WARN] Failed to set sorted_ids: %v", err)
	}

	log.Printf("[INFO] Network speed rules sorted successfully: sorted_ids=%v", response.Data.IDs)
	return resourceScdnNetworkSpeedRulesSortRead(d, m)
}

func resourceScdnNetworkSpeedRulesSortDelete(d *schema.ResourceData, m interface{}) error {
	// Sort operation is idempotent - no need to delete anything
	// Just remove from state
	log.Printf("[INFO] Deleting network speed rules sort from state (no API call needed)")
	return nil
}
