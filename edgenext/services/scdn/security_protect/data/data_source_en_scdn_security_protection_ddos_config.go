package data

import (
	"fmt"
	"log"

	"github.com/edgenextapisdk/terraform-provider-edgenext/edgenext/connectivity"
	"github.com/edgenextapisdk/terraform-provider-edgenext/edgenext/helper"
	"github.com/edgenextapisdk/terraform-provider-edgenext/edgenext/services/scdn"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// DataSourceEdgenextScdnSecurityProtectionDdosConfig returns the SCDN security protection DDoS config data source
func DataSourceEdgenextScdnSecurityProtectionDdosConfig() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceScdnSecurityProtectionDdosConfigRead,

		Schema: map[string]*schema.Schema{
			"business_id": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Business ID",
			},
			"keys": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Specify config keys",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results to a file",
			},
			"application_ddos_protection": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Application layer DDoS protection configuration",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "ID",
						},
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Status: on, off, keep",
						},
						"ai_cc_status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "AI protection status: on, off",
						},
						"type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Protection type: default, normal, strict, captcha, keep",
						},
						"need_attack_detection": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Attack detection switch: 0 or 1",
						},
						"ai_status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "AI status: on, off",
						},
					},
				},
			},
			"visitor_authentication": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Visitor authentication configuration",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "ID",
						},
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Status: on, off",
						},
						"auth_token": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Authentication token",
						},
						"pass_still_check": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Pass still check: 0 or 1",
						},
					},
				},
			},
		},
	}
}

func dataSourceScdnSecurityProtectionDdosConfigRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*connectivity.EdgeNextClient)
	service := scdn.NewScdnService(client)

	businessID := d.Get("business_id").(int)

	req := scdn.DdosProtectionGetConfigRequest{
		BusinessID: businessID,
	}

	// Get keys if provided
	if v, ok := d.GetOk("keys"); ok {
		keysList := v.([]interface{})
		keys := make([]string, len(keysList))
		for i, v := range keysList {
			keys[i] = v.(string)
		}
		req.Keys = keys
	}

	log.Printf("[INFO] Querying SCDN security protection DDoS config: business_id=%d", businessID)
	response, err := service.GetDdosProtectionConfig(req)
	if err != nil {
		return fmt.Errorf("failed to query DDoS protection config: %w", err)
	}

	// Set resource ID
	d.SetId(fmt.Sprintf("ddos-config-%d", businessID))

	// Set application_ddos_protection
	if response.Data.ApplicationDdosProtection != nil {
		appDdos := []map[string]interface{}{
			{
				"id":                    response.Data.ApplicationDdosProtection.ID,
				"status":                response.Data.ApplicationDdosProtection.Status,
				"ai_cc_status":          response.Data.ApplicationDdosProtection.AICCStatus,
				"type":                  response.Data.ApplicationDdosProtection.Type,
				"need_attack_detection": response.Data.ApplicationDdosProtection.NeedAttackDetection,
				"ai_status":             response.Data.ApplicationDdosProtection.AIStatus,
			},
		}
		if err := d.Set("application_ddos_protection", appDdos); err != nil {
			return fmt.Errorf("error setting application_ddos_protection: %w", err)
		}
	}

	// Set visitor_authentication
	if response.Data.VisitorAuthentication != nil {
		visitorAuth := []map[string]interface{}{
			{
				"id":               response.Data.VisitorAuthentication.ID,
				"status":           response.Data.VisitorAuthentication.Status,
				"auth_token":       response.Data.VisitorAuthentication.AuthToken,
				"pass_still_check": response.Data.VisitorAuthentication.PassStillCheck,
			},
		}
		if err := d.Set("visitor_authentication", visitorAuth); err != nil {
			return fmt.Errorf("error setting visitor_authentication: %w", err)
		}
	}

	// Write result to output file if specified
	if outputFile := d.Get("result_output_file").(string); outputFile != "" {
		outputData := map[string]interface{}{
			"application_ddos_protection": response.Data.ApplicationDdosProtection,
			"visitor_authentication":      response.Data.VisitorAuthentication,
		}
		if err := helper.WriteToFile(d, outputData); err != nil {
			return fmt.Errorf("failed to write output file: %w", err)
		}
	}

	log.Printf("[INFO] SCDN security protection DDoS config queried successfully: business_id=%d", businessID)
	return nil
}
