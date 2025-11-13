package domain

import (
	"fmt"
	"log"
	"strconv"

	"github.com/edgenextapisdk/terraform-provider-edgenext/edgenext/connectivity"
	"github.com/edgenextapisdk/terraform-provider-edgenext/edgenext/services/scdn"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// ResourceEdgenextScdnDomainStatus returns the SCDN domain status management resource
func ResourceEdgenextScdnDomainStatus() *schema.Resource {
	return &schema.Resource{
		Create: resourceScdnDomainStatusCreate,
		Read:   resourceScdnDomainStatusRead,
		Update: resourceScdnDomainStatusUpdate,
		Delete: resourceScdnDomainStatusDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"domain_id": {
				Type:        schema.TypeInt,
				Required:    true,
				ForceNew:    true,
				Description: "The ID of the domain to manage status",
			},
			"enabled": {
				Type:        schema.TypeBool,
				Required:    true,
				Description: "Whether the domain is enabled (true) or disabled (false)",
			},
		},
	}
}

func resourceScdnDomainStatusCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*connectivity.EdgeNextClient)
	service := scdn.NewScdnService(client)

	domainID := d.Get("domain_id").(int)
	enabled := d.Get("enabled").(bool)

	var err error
	if enabled {
		req := scdn.DomainEnableRequest{
			DomainIDs: []int{domainID},
		}
		log.Printf("[INFO] Enabling SCDN domain: %d", domainID)
		_, err = service.EnableDomain(req)
		if err != nil {
			return fmt.Errorf("failed to enable SCDN domain: %w", err)
		}
	} else {
		req := scdn.DomainDisableRequest{
			DomainIDs: []int{domainID},
		}
		log.Printf("[INFO] Disabling SCDN domain: %d", domainID)
		_, err = service.DisableDomain(req)
		if err != nil {
			return fmt.Errorf("failed to disable SCDN domain: %w", err)
		}
	}

	d.SetId(strconv.Itoa(domainID))
	log.Printf("[INFO] SCDN domain status updated successfully: %d (enabled=%v)", domainID, enabled)
	return resourceScdnDomainStatusRead(d, m)
}

func resourceScdnDomainStatusRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*connectivity.EdgeNextClient)
	service := scdn.NewScdnService(client)

	domainID := d.Get("domain_id").(int)

	// Query domain to get current status
	req := scdn.DomainListRequest{
		Page:     1,
		PageSize: 100,
		ID:       domainID,
	}

	response, err := service.ListDomains(req)
	if err != nil {
		return fmt.Errorf("failed to read SCDN domain status: %w", err)
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

	// Infer enabled status from access_progress
	// The API doesn't provide a direct enabled/disabled field, so we must infer from access_progress
	// Based on API behavior: disabled domains typically have access_progress="disabled"
	// Enabled domains typically have access_progress="completed" or "enabled"
	// IMPORTANT: This inference may not be 100% accurate, but it's the best we can do with available API fields
	enabled := domainInfo.AccessProgress != "disabled" &&
		(domainInfo.AccessProgress == "completed" || domainInfo.AccessProgress == "enabled")

	log.Printf("[DEBUG] Domain %d: access_progress=%s, inferred enabled=%v",
		domainID, domainInfo.AccessProgress, enabled)

	// Set the inferred value to state
	// Terraform will compare this with the desired value from configuration
	// If they differ, Terraform will call Update to synchronize
	if err := d.Set("enabled", enabled); err != nil {
		return fmt.Errorf("error setting enabled: %w", err)
	}

	log.Printf("[INFO] SCDN domain status read successfully: %d (enabled=%v, access_progress=%s)",
		domainID, enabled, domainInfo.AccessProgress)
	return nil
}

func resourceScdnDomainStatusUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*connectivity.EdgeNextClient)
	service := scdn.NewScdnService(client)

	domainID := d.Get("domain_id").(int)
	enabled := d.Get("enabled").(bool)

	// Check if enabled status actually changed
	oldEnabled, newEnabled := d.GetChange("enabled")
	log.Printf("[DEBUG] Update called for domain %d: old_enabled=%v, new_enabled=%v", domainID, oldEnabled, newEnabled)

	var err error
	if enabled {
		req := scdn.DomainEnableRequest{
			DomainIDs: []int{domainID},
		}
		log.Printf("[INFO] Enabling SCDN domain: %d", domainID)
		_, err = service.EnableDomain(req)
		if err != nil {
			return fmt.Errorf("failed to enable SCDN domain: %w", err)
		}
	} else {
		req := scdn.DomainDisableRequest{
			DomainIDs: []int{domainID},
		}
		log.Printf("[INFO] Disabling SCDN domain: %d", domainID)
		_, err = service.DisableDomain(req)
		if err != nil {
			return fmt.Errorf("failed to disable SCDN domain: %w", err)
		}
	}

	log.Printf("[INFO] SCDN domain status updated successfully: %d (enabled=%v)", domainID, enabled)
	return resourceScdnDomainStatusRead(d, m)
}

func resourceScdnDomainStatusDelete(d *schema.ResourceData, m interface{}) error {
	// Disable the domain when deleting the resource
	client := m.(*connectivity.EdgeNextClient)
	service := scdn.NewScdnService(client)

	domainID := d.Get("domain_id").(int)

	req := scdn.DomainDisableRequest{
		DomainIDs: []int{domainID},
	}

	log.Printf("[INFO] Disabling SCDN domain before delete: %d", domainID)
	_, err := service.DisableDomain(req)
	if err != nil {
		return fmt.Errorf("failed to disable SCDN domain: %w", err)
	}

	d.SetId("")
	log.Printf("[INFO] SCDN domain status deleted successfully: %d", domainID)
	return nil
}
