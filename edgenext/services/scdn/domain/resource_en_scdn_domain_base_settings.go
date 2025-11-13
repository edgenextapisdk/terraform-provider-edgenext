package domain

import (
	"fmt"
	"log"

	"github.com/edgenextapisdk/terraform-provider-edgenext/edgenext/connectivity"
	"github.com/edgenextapisdk/terraform-provider-edgenext/edgenext/services/scdn"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// ResourceEdgenextScdnDomainBaseSettings returns the SCDN domain base settings resource
func ResourceEdgenextScdnDomainBaseSettings() *schema.Resource {
	return &schema.Resource{
		Create: resourceScdnDomainBaseSettingsCreate,
		Read:   resourceScdnDomainBaseSettingsRead,
		Update: resourceScdnDomainBaseSettingsUpdate,
		Delete: resourceScdnDomainBaseSettingsDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"domain_id": {
				Type:        schema.TypeInt,
				Required:    true,
				ForceNew:    true,
				Description: "The ID of the domain to update base settings",
			},
			"proxy_host": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Proxy host configuration",
				MaxItems:    1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"proxy_host": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Proxy host value",
						},
						"proxy_host_type": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Proxy host type",
						},
					},
				},
			},
			"proxy_sni": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Proxy SNI configuration",
				MaxItems:    1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"proxy_sni": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Proxy SNI value",
						},
						"status": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Proxy SNI status",
						},
					},
				},
			},
			"domain_redirect": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Domain redirect configuration",
				MaxItems:    1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"status": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Redirect status (on/off)",
						},
						"jump_to": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Redirect target URL",
						},
						"jump_type": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Redirect jump type",
						},
					},
				},
			},
		},
	}
}

func resourceScdnDomainBaseSettingsCreate(d *schema.ResourceData, m interface{}) error {
	return resourceScdnDomainBaseSettingsUpdate(d, m)
}

func resourceScdnDomainBaseSettingsRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*connectivity.EdgeNextClient)
	service := scdn.NewScdnService(client)

	domainID := d.Get("domain_id").(int)

	req := scdn.DomainBaseSettingsGetRequest{
		DomainID: domainID,
	}

	log.Printf("[INFO] Reading SCDN domain base settings for domain: %d", domainID)
	response, err := service.GetDomainBaseSettings(req)
	if err != nil {
		return fmt.Errorf("failed to read SCDN domain base settings: %w", err)
	}

	d.SetId(fmt.Sprintf("%d", domainID))

	// Set proxy host
	if response.Data.ProxyHost.ProxyHost != "" || response.Data.ProxyHost.ProxyHostType != "" {
		proxyHost := []map[string]interface{}{
			{
				"proxy_host":      response.Data.ProxyHost.ProxyHost,
				"proxy_host_type": response.Data.ProxyHost.ProxyHostType,
			},
		}
		if err := d.Set("proxy_host", proxyHost); err != nil {
			return fmt.Errorf("error setting proxy_host: %w", err)
		}
	}

	// Set proxy SNI
	proxySNI := []map[string]interface{}{
		{
			"proxy_sni": response.Data.ProxySNI.ProxySNI,
			"status":    response.Data.ProxySNI.Status,
		},
	}
	if err := d.Set("proxy_sni", proxySNI); err != nil {
		return fmt.Errorf("error setting proxy_sni: %w", err)
	}

	// Set domain redirect
	domainRedirect := []map[string]interface{}{
		{
			"status":    response.Data.DomainRedirect.Status,
			"jump_to":   response.Data.DomainRedirect.JumpTo,
			"jump_type": response.Data.DomainRedirect.JumpType,
		},
	}
	if err := d.Set("domain_redirect", domainRedirect); err != nil {
		return fmt.Errorf("error setting domain_redirect: %w", err)
	}

	log.Printf("[INFO] SCDN domain base settings read successfully for domain: %d", domainID)
	return nil
}

func resourceScdnDomainBaseSettingsUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*connectivity.EdgeNextClient)
	service := scdn.NewScdnService(client)

	domainID := d.Get("domain_id").(int)

	// Build settings value
	value := scdn.DomainBaseSettingsValue{}

	// Proxy host
	if v, ok := d.GetOk("proxy_host"); ok && len(v.([]interface{})) > 0 {
		proxyHostMap := v.([]interface{})[0].(map[string]interface{})
		value.ProxyHost = &scdn.ProxyHostConfig{
			ProxyHost:     proxyHostMap["proxy_host"].(string),
			ProxyHostType: proxyHostMap["proxy_host_type"].(string),
		}
	}

	// Proxy SNI
	if v, ok := d.GetOk("proxy_sni"); ok && len(v.([]interface{})) > 0 {
		proxySNIMap := v.([]interface{})[0].(map[string]interface{})
		value.ProxySNI = &scdn.ProxySNIConfig{
			ProxySNI: proxySNIMap["proxy_sni"].(string),
			Status:   proxySNIMap["status"].(string),
		}
	}

	// Domain redirect
	if v, ok := d.GetOk("domain_redirect"); ok && len(v.([]interface{})) > 0 {
		redirectMap := v.([]interface{})[0].(map[string]interface{})
		value.DomainRedirect = &scdn.DomainRedirectConfig{
			Status:   redirectMap["status"].(string),
			JumpTo:   redirectMap["jump_to"].(string),
			JumpType: redirectMap["jump_type"].(string),
		}
	}

	req := scdn.DomainBaseSettingsUpdateRequest{
		DomainID: domainID,
		Value:    value,
	}

	log.Printf("[INFO] Updating SCDN domain base settings for domain: %d", domainID)
	_, err := service.UpdateDomainBaseSettings(req)
	if err != nil {
		return fmt.Errorf("failed to update SCDN domain base settings: %w", err)
	}

	d.SetId(fmt.Sprintf("%d", domainID))
	log.Printf("[INFO] SCDN domain base settings updated successfully for domain: %d", domainID)
	return resourceScdnDomainBaseSettingsRead(d, m)
}

func resourceScdnDomainBaseSettingsDelete(d *schema.ResourceData, m interface{}) error {
	// Domain base settings cannot be deleted, only updated
	// Setting empty values would effectively "delete" them
	log.Printf("[WARN] Domain base settings cannot be deleted, only cleared")
	d.SetId("")
	return nil
}
