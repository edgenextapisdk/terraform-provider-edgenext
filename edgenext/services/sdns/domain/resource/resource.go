package resource

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/edgenextapisdk/terraform-provider-edgenext/edgenext/connectivity"
	"github.com/edgenextapisdk/terraform-provider-edgenext/edgenext/services/sdns"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceEdgenextDnsDomain() *schema.Resource {
	return &schema.Resource{
		Create: resourceDnsDomainCreate,
		Read:   resourceDnsDomainRead,
		Delete: resourceDnsDomainDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"domain": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The domain name to be added to DNS",
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Status of the domain",
			},
		},
	}
}

func resourceDnsDomainCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*connectivity.EdgeNextClient)
	service := sdns.NewSdnsService(client)

	domain := d.Get("domain").(string)

	log.Printf("[INFO] Creating DNS domain: %s", domain)
	resp, err := service.AddDnsDomain(domain)
	if err != nil {
		return fmt.Errorf("failed to create DNS domain: %w", err)
	}

	d.SetId(strconv.Itoa(resp.ID))
	return resourceDnsDomainRead(d, m)
}

func resourceDnsDomainRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*connectivity.EdgeNextClient)
	service := sdns.NewSdnsService(client)

	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return fmt.Errorf("invalid domain ID: %s", d.Id())
	}

	log.Printf("[DEBUG] Reading DNS domain: %d", id)
	info, err := service.GetDnsDomainInfo(id)
	if err != nil {
		if strings.Contains(err.Error(), "not found") || strings.Contains(err.Error(), "exist") {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("failed to get DNS domain info: %w", err)
	}

	d.Set("domain", info.Domain)
	d.Set("status", strconv.Itoa(info.Status))

	return nil
}

func resourceDnsDomainDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*connectivity.EdgeNextClient)
	service := sdns.NewSdnsService(client)

	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return fmt.Errorf("invalid domain ID: %s", d.Id())
	}

	log.Printf("[INFO] Deleting DNS domain: %d", id)
	err = service.DeleteDnsDomain([]int{id})
	if err != nil {
		return fmt.Errorf("failed to delete DNS domain: %w", err)
	}

	d.SetId("")
	return nil
}
