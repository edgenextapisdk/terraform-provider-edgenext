package domain

import (
	"fmt"
	"log"
	"strconv"

	"github.com/edgenextapisdk/terraform-provider-edgenext/edgenext/connectivity"
	"github.com/edgenextapisdk/terraform-provider-edgenext/edgenext/services/scdn"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// ResourceEdgenextScdnDomainAccessMode returns the SCDN domain access mode switch resource
func ResourceEdgenextScdnDomainAccessMode() *schema.Resource {
	return &schema.Resource{
		Create: resourceScdnDomainAccessModeCreate,
		Read:   resourceScdnDomainAccessModeRead,
		Update: resourceScdnDomainAccessModeUpdate,
		Delete: resourceScdnDomainAccessModeDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"domain_id": {
				Type:        schema.TypeInt,
				Required:    true,
				ForceNew:    true,
				Description: "The ID of the domain to switch access mode",
			},
			"access_mode": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The access mode. Valid values: ns, cname",
			},
		},
	}
}

func resourceScdnDomainAccessModeCreate(d *schema.ResourceData, m interface{}) error {
	return resourceScdnDomainAccessModeUpdate(d, m)
}

func resourceScdnDomainAccessModeRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*connectivity.EdgeNextClient)
	service := scdn.NewScdnService(client)

	domainID := d.Get("domain_id").(int)

	// Query domain to get current access_mode
	req := scdn.DomainListRequest{
		Page:     1,
		PageSize: 100,
	}

	response, err := service.ListDomains(req)
	if err != nil {
		return fmt.Errorf("failed to read SCDN domain access mode: %w", err)
	}

	var domainInfo *scdn.DomainInfo
	for _, domain := range response.Data.List {
		if domain.ID == domainID {
			domainInfo = &domain
			break
		}
	}

	if domainInfo == nil {
		log.Printf("[WARN] SCDN domain not found: %d", domainID)
		d.SetId("")
		return nil
	}

	if err := d.Set("access_mode", domainInfo.AccessMode); err != nil {
		return fmt.Errorf("error setting access_mode: %w", err)
	}

	log.Printf("[INFO] SCDN domain access mode read successfully: %d", domainID)
	return nil
}

func resourceScdnDomainAccessModeUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*connectivity.EdgeNextClient)
	service := scdn.NewScdnService(client)

	domainID := d.Get("domain_id").(int)
	accessMode := d.Get("access_mode").(string)

	req := scdn.DomainAccessModeSwitchRequest{
		DomainID:   domainID,
		AccessMode: accessMode,
	}

	log.Printf("[INFO] Switching SCDN domain access mode for domain: %d to %s", domainID, accessMode)
	_, err := service.SwitchDomainAccessMode(req)
	if err != nil {
		return fmt.Errorf("failed to switch SCDN domain access mode: %w", err)
	}

	d.SetId(strconv.Itoa(domainID))
	log.Printf("[INFO] SCDN domain access mode switched successfully: %d", domainID)
	return resourceScdnDomainAccessModeRead(d, m)
}

func resourceScdnDomainAccessModeDelete(d *schema.ResourceData, m interface{}) error {
	// Access mode switch cannot be reverted, just remove from state
	log.Printf("[WARN] Domain access mode switch cannot be reverted, removing from state")
	d.SetId("")
	return nil
}
