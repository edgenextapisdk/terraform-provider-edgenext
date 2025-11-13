package resource

import (
	"fmt"
	"log"

	"github.com/edgenextapisdk/terraform-provider-edgenext/edgenext/connectivity"
	"github.com/edgenextapisdk/terraform-provider-edgenext/edgenext/services/scdn"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// ResourceEdgenextScdnOriginGroupDomainCopy returns the SCDN origin group domain copy resource
func ResourceEdgenextScdnOriginGroupDomainCopy() *schema.Resource {
	return &schema.Resource{
		Create: resourceScdnOriginGroupDomainCopyCreate,
		Read:   resourceScdnOriginGroupDomainCopyRead,
		Delete: resourceScdnOriginGroupDomainCopyDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"origin_group_id": {
				Type:        schema.TypeInt,
				Required:    true,
				ForceNew:    true,
				Description: "Origin group ID",
			},
			"domain_id": {
				Type:        schema.TypeInt,
				Required:    true,
				ForceNew:    true,
				Description: "Domain ID",
			},
		},
	}
}

func resourceScdnOriginGroupDomainCopyCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*connectivity.EdgeNextClient)
	service := scdn.NewScdnService(client)

	req := scdn.OriginGroupCopyRequest{
		OriginGroupID: d.Get("origin_group_id").(int),
		DomainID:      d.Get("domain_id").(int),
	}

	log.Printf("[INFO] Copying SCDN origin group to domain: origin_group_id=%d, domain_id=%d", req.OriginGroupID, req.DomainID)
	_, err := service.CopyOriginGroupToDomain(req)
	if err != nil {
		return fmt.Errorf("failed to copy origin group to domain: %w", err)
	}

	// Set resource ID
	d.SetId(fmt.Sprintf("origin-group-domain-copy-%d-%d", req.OriginGroupID, req.DomainID))

	return resourceScdnOriginGroupDomainCopyRead(d, m)
}

func resourceScdnOriginGroupDomainCopyRead(d *schema.ResourceData, m interface{}) error {
	// Copy is a one-time operation, we can't really "read" the copy state from API
	// Just verify the resource exists
	originGroupID := d.Get("origin_group_id").(int)
	domainID := d.Get("domain_id").(int)
	if originGroupID == 0 || domainID == 0 {
		// Try to parse from resource ID
		id := d.Id()
		if id != "" {
			var parsedOriginGroupID, parsedDomainID int
			_, err := fmt.Sscanf(id, "origin-group-domain-copy-%d-%d", &parsedOriginGroupID, &parsedDomainID)
			if err == nil {
				originGroupID = parsedOriginGroupID
				domainID = parsedDomainID
				if err := d.Set("origin_group_id", originGroupID); err != nil {
					log.Printf("[WARN] Failed to set origin_group_id from ID: %v", err)
				}
				if err := d.Set("domain_id", domainID); err != nil {
					log.Printf("[WARN] Failed to set domain_id from ID: %v", err)
				}
			} else {
				d.SetId("")
				return nil
			}
		} else {
			d.SetId("")
			return nil
		}
	}

	return nil
}

func resourceScdnOriginGroupDomainCopyDelete(d *schema.ResourceData, m interface{}) error {
	// Copy operation cannot be undone via API
	// For now, we just remove the resource from state
	log.Printf("[INFO] Removing SCDN origin group domain copy from state: origin_group_id=%d, domain_id=%d",
		d.Get("origin_group_id").(int), d.Get("domain_id").(int))
	d.SetId("")
	return nil
}
