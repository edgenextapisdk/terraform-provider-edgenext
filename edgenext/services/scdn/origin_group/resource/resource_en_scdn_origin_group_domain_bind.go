package resource

import (
	"fmt"
	"log"

	"github.com/edgenextapisdk/terraform-provider-edgenext/edgenext/connectivity"
	"github.com/edgenextapisdk/terraform-provider-edgenext/edgenext/services/scdn"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// ResourceEdgenextScdnOriginGroupDomainBind returns the SCDN origin group domain bind resource
func ResourceEdgenextScdnOriginGroupDomainBind() *schema.Resource {
	return &schema.Resource{
		Create: resourceScdnOriginGroupDomainBindCreate,
		Read:   resourceScdnOriginGroupDomainBindRead,
		Delete: resourceScdnOriginGroupDomainBindDelete,

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
			"domain_ids": {
				Type:        schema.TypeList,
				Optional:    true,
				ForceNew:    true,
				Description: "Domain ID array",
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
			},
			"domain_group_ids": {
				Type:        schema.TypeList,
				Optional:    true,
				ForceNew:    true,
				Description: "Domain group ID array",
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
			},
			"domains": {
				Type:        schema.TypeList,
				Optional:    true,
				ForceNew:    true,
				Description: "Domain array",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"job_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Batch job ID",
			},
		},
	}
}

func resourceScdnOriginGroupDomainBindCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*connectivity.EdgeNextClient)
	service := scdn.NewScdnService(client)

	originGroupID := d.Get("origin_group_id").(int)

	req := scdn.OriginGroupBindDomainsRequest{
		OriginGroupID: originGroupID,
	}

	if domainIDs, ok := d.GetOk("domain_ids"); ok {
		domainIDsList := domainIDs.([]interface{})
		req.DomainIDs = make([]int, len(domainIDsList))
		for i, v := range domainIDsList {
			req.DomainIDs[i] = v.(int)
		}
	}

	if domainGroupIDs, ok := d.GetOk("domain_group_ids"); ok {
		domainGroupIDsList := domainGroupIDs.([]interface{})
		req.DomainGroupIDs = make([]int, len(domainGroupIDsList))
		for i, v := range domainGroupIDsList {
			req.DomainGroupIDs[i] = v.(int)
		}
	}

	if domains, ok := d.GetOk("domains"); ok {
		domainsList := domains.([]interface{})
		req.Domains = make([]string, len(domainsList))
		for i, v := range domainsList {
			req.Domains[i] = v.(string)
		}
	}

	log.Printf("[INFO] Binding SCDN origin group to domains: origin_group_id=%d", originGroupID)
	response, err := service.BindOriginGroupToDomains(req)
	if err != nil {
		return fmt.Errorf("failed to bind origin group to domains: %w", err)
	}

	// Set resource ID
	d.SetId(fmt.Sprintf("origin-group-domain-bind-%d", originGroupID))

	// Set job_id
	if err := d.Set("job_id", response.Data.JobID); err != nil {
		log.Printf("[WARN] Failed to set job_id: %v", err)
	}

	return resourceScdnOriginGroupDomainBindRead(d, m)
}

func resourceScdnOriginGroupDomainBindRead(d *schema.ResourceData, m interface{}) error {
	// Binding is a one-time operation, we can't really "read" the binding state from API
	// Just verify the resource exists by ensuring origin_group_id is set
	originGroupID := d.Get("origin_group_id").(int)
	if originGroupID == 0 {
		// Try to parse from resource ID
		id := d.Id()
		if id != "" {
			var parsedID int
			_, err := fmt.Sscanf(id, "origin-group-domain-bind-%d", &parsedID)
			if err == nil {
				originGroupID = parsedID
				if err := d.Set("origin_group_id", originGroupID); err != nil {
					log.Printf("[WARN] Failed to set origin_group_id from ID: %v", err)
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

func resourceScdnOriginGroupDomainBindDelete(d *schema.ResourceData, m interface{}) error {
	// Unbinding is not directly supported by the API
	// The binding can be removed by updating the domain configuration
	// For now, we just remove the resource from state
	log.Printf("[INFO] Unbinding SCDN origin group from domains: origin_group_id=%d", d.Get("origin_group_id").(int))
	d.SetId("")
	return nil
}
