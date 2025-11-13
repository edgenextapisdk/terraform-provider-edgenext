package resource

import (
	"fmt"
	"log"

	"github.com/edgenextapisdk/terraform-provider-edgenext/edgenext/connectivity"
	"github.com/edgenextapisdk/terraform-provider-edgenext/edgenext/services/scdn"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// ResourceEdgenextScdnSecurityProtectionDdosConfig returns the SCDN security protection DDoS config resource
func ResourceEdgenextScdnSecurityProtectionDdosConfig() *schema.Resource {
	return &schema.Resource{
		Create: resourceScdnSecurityProtectionDdosConfigCreate,
		Read:   resourceScdnSecurityProtectionDdosConfigRead,
		Update: resourceScdnSecurityProtectionDdosConfigUpdate,
		Delete: resourceScdnSecurityProtectionDdosConfigDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"business_id": {
				Type:        schema.TypeInt,
				Required:    true,
				ForceNew:    true,
				Description: "Business ID",
			},
			"application_ddos_protection": {
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				Description: "Application layer DDoS protection configuration",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"status": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Status: on, off, keep",
						},
						"ai_cc_status": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "AI protection status: on, off",
						},
						"type": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Protection type: default, normal, strict, captcha, keep",
						},
						"need_attack_detection": {
							Type:        schema.TypeInt,
							Optional:    true,
							Default:     1,
							Description: "Attack detection switch: 0 or 1",
						},
						"ai_status": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "AI status: on, off",
						},
					},
				},
			},
			"visitor_authentication": {
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				Description: "Visitor authentication configuration",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"status": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Status: on, off",
						},
						"auth_token": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Authentication token",
						},
						"pass_still_check": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Pass still check: 0 or 1",
						},
					},
				},
			},
		},
	}
}

func resourceScdnSecurityProtectionDdosConfigCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*connectivity.EdgeNextClient)
	service := scdn.NewScdnService(client)

	businessID := d.Get("business_id").(int)

	req := scdn.DdosProtectionUpdateConfigRequest{
		BusinessID: businessID,
	}

	// Build application_ddos_protection
	if v, ok := d.GetOk("application_ddos_protection"); ok {
		appDdosList := v.([]interface{})
		if len(appDdosList) > 0 {
			appDdosMap := appDdosList[0].(map[string]interface{})
			appDdos := &scdn.ApplicationDdosProtection{}
			if val, ok := appDdosMap["status"].(string); ok && val != "" {
				appDdos.Status = val
			}
			if val, ok := appDdosMap["ai_cc_status"].(string); ok && val != "" {
				appDdos.AICCStatus = val
			}
			if val, ok := appDdosMap["type"].(string); ok && val != "" {
				appDdos.Type = val
			}
			if val, ok := appDdosMap["need_attack_detection"].(int); ok {
				appDdos.NeedAttackDetection = val
			}
			if val, ok := appDdosMap["ai_status"].(string); ok && val != "" {
				appDdos.AIStatus = val
			}
			req.ApplicationDdosProtection = appDdos
		}
	}

	// Build visitor_authentication
	if v, ok := d.GetOk("visitor_authentication"); ok {
		visitorAuthList := v.([]interface{})
		if len(visitorAuthList) > 0 {
			visitorAuthMap := visitorAuthList[0].(map[string]interface{})
			visitorAuth := &scdn.VisitorAuthentication{}
			if val, ok := visitorAuthMap["status"].(string); ok && val != "" {
				visitorAuth.Status = val
			}
			if val, ok := visitorAuthMap["auth_token"].(string); ok && val != "" {
				visitorAuth.AuthToken = val
			}
			if val, ok := visitorAuthMap["pass_still_check"].(int); ok {
				visitorAuth.PassStillCheck = val
			}
			req.VisitorAuthentication = visitorAuth
		}
	}

	log.Printf("[INFO] Creating/Updating SCDN security protection DDoS config: business_id=%d", businessID)
	_, err := service.UpdateDdosProtectionConfig(req)
	if err != nil {
		return fmt.Errorf("failed to create/update DDoS protection config: %w", err)
	}

	d.SetId(fmt.Sprintf("ddos-config-%d", businessID))
	return resourceScdnSecurityProtectionDdosConfigRead(d, m)
}

func resourceScdnSecurityProtectionDdosConfigRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*connectivity.EdgeNextClient)
	service := scdn.NewScdnService(client)

	businessID := d.Get("business_id").(int)

	req := scdn.DdosProtectionGetConfigRequest{
		BusinessID: businessID,
	}

	log.Printf("[INFO] Reading SCDN security protection DDoS config: business_id=%d", businessID)
	response, err := service.GetDdosProtectionConfig(req)
	if err != nil {
		return fmt.Errorf("failed to read DDoS protection config: %w", err)
	}

	// Set application_ddos_protection
	if response.Data.ApplicationDdosProtection != nil {
		appDdos := []map[string]interface{}{
			{
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
				"status":           response.Data.VisitorAuthentication.Status,
				"auth_token":       response.Data.VisitorAuthentication.AuthToken,
				"pass_still_check": response.Data.VisitorAuthentication.PassStillCheck,
			},
		}
		if err := d.Set("visitor_authentication", visitorAuth); err != nil {
			return fmt.Errorf("error setting visitor_authentication: %w", err)
		}
	}

	return nil
}

func resourceScdnSecurityProtectionDdosConfigUpdate(d *schema.ResourceData, m interface{}) error {
	return resourceScdnSecurityProtectionDdosConfigCreate(d, m)
}

func resourceScdnSecurityProtectionDdosConfigDelete(d *schema.ResourceData, m interface{}) error {
	// DDoS protection config cannot be deleted, only reset
	// Set all values to default/off
	client := m.(*connectivity.EdgeNextClient)
	service := scdn.NewScdnService(client)

	businessID := d.Get("business_id").(int)

	req := scdn.DdosProtectionUpdateConfigRequest{
		BusinessID: businessID,
		ApplicationDdosProtection: &scdn.ApplicationDdosProtection{
			Status:              "off",
			AICCStatus:          "off",
			Type:                "default",
			NeedAttackDetection: 0,
			AIStatus:            "off",
		},
		VisitorAuthentication: &scdn.VisitorAuthentication{
			Status:         "off",
			AuthToken:      "",
			PassStillCheck: 0,
		},
	}

	log.Printf("[INFO] Resetting SCDN security protection DDoS config: business_id=%d", businessID)
	_, err := service.UpdateDdosProtectionConfig(req)
	if err != nil {
		return fmt.Errorf("failed to reset DDoS protection config: %w", err)
	}

	d.SetId("")
	return nil
}
