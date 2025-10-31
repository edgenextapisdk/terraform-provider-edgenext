package domain

import (
	"fmt"
	"log"
	"strconv"

	"github.com/edgenextapisdk/terraform-provider-edgenext/edgenext/connectivity"
	"github.com/edgenextapisdk/terraform-provider-edgenext/edgenext/services/scdn"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// ResourceEdgenextScdnDomainNodeSwitch returns the SCDN domain node switch resource
func ResourceEdgenextScdnDomainNodeSwitch() *schema.Resource {
	return &schema.Resource{
		Create: resourceScdnDomainNodeSwitchCreate,
		Read:   resourceScdnDomainNodeSwitchRead,
		Update: resourceScdnDomainNodeSwitchUpdate,
		Delete: resourceScdnDomainNodeSwitchDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"domain_id": {
				Type:        schema.TypeInt,
				Required:    true,
				ForceNew:    true,
				Description: "The ID of the domain to switch nodes",
			},
			"protect_status": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The edge node type. Valid values: back_source, scdn, exclusive",
			},
			"exclusive_resource_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "The ID of the exclusive resource package (required if protect_status is exclusive)",
			},
		},
	}
}

func resourceScdnDomainNodeSwitchCreate(d *schema.ResourceData, m interface{}) error {
	return resourceScdnDomainNodeSwitchUpdate(d, m)
}

func resourceScdnDomainNodeSwitchRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*connectivity.EdgeNextClient)
	service := scdn.NewScdnService(client)

	domainID := d.Get("domain_id").(int)

	// Query domain to get current protect_status
	req := scdn.DomainListRequest{
		Page:     1,
		PageSize: 100,
	}

	response, err := service.ListDomains(req)
	if err != nil {
		return fmt.Errorf("failed to read SCDN domain node switch: %w", err)
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

	if err := d.Set("protect_status", domainInfo.ProtectStatus); err != nil {
		return fmt.Errorf("error setting protect_status: %w", err)
	}
	if err := d.Set("exclusive_resource_id", domainInfo.ExclusiveResourceID); err != nil {
		return fmt.Errorf("error setting exclusive_resource_id: %w", err)
	}

	log.Printf("[INFO] SCDN domain node switch read successfully: %d", domainID)
	return nil
}

func resourceScdnDomainNodeSwitchUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*connectivity.EdgeNextClient)
	service := scdn.NewScdnService(client)

	domainID := d.Get("domain_id").(int)
	protectStatus := d.Get("protect_status").(string)
	exclusiveResourceID := 0
	if v, ok := d.GetOk("exclusive_resource_id"); ok {
		exclusiveResourceID = v.(int)
	}

	req := scdn.DomainNodeSwitchRequest{
		DomainID:            domainID,
		ProtectStatus:       protectStatus,
		ExclusiveResourceID: exclusiveResourceID,
	}

	log.Printf("[INFO] Switching SCDN domain nodes for domain: %d to %s", domainID, protectStatus)
	_, err := service.SwitchDomainNodes(req)
	if err != nil {
		return fmt.Errorf("failed to switch SCDN domain nodes: %w", err)
	}

	d.SetId(strconv.Itoa(domainID))
	log.Printf("[INFO] SCDN domain nodes switched successfully: %d", domainID)
	return resourceScdnDomainNodeSwitchRead(d, m)
}

func resourceScdnDomainNodeSwitchDelete(d *schema.ResourceData, m interface{}) error {
	// Node switch cannot be reverted, just remove from state
	log.Printf("[WARN] Domain node switch cannot be reverted, removing from state")
	d.SetId("")
	return nil
}
