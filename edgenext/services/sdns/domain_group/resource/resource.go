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

func ResourceEdgenextDnsGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceDnsGroupCreate,
		Read:   resourceDnsGroupRead,
		Update: resourceDnsGroupUpdate,
		Delete: resourceDnsGroupDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"group_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the DNS domain group",
			},
			"remark": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Remark for the domain group",
			},
			"domain_ids": {
				Type:        schema.TypeSet,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "List of domain IDs to bind to the group",
			},
		},
	}
}

func resourceDnsGroupCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*connectivity.EdgeNextClient)
	service := sdns.NewSdnsService(client)

	req := sdns.DnsGroupAddRequest{
		GroupName: d.Get("group_name").(string),
		Remark:    d.Get("remark").(string),
	}

	if v, ok := d.GetOk("domain_ids"); ok {
		req.DomainIDs = expandIntList(v.(*schema.Set).List())
	}

	log.Printf("[INFO] Creating DNS group: %s", req.GroupName)
	resp, err := service.AddDnsGroup(req)
	if err != nil {
		return fmt.Errorf("failed to create DNS group: %w", err)
	}

	d.SetId(strconv.Itoa(resp.Data.ID))

	return resourceDnsGroupRead(d, m)
}

func resourceDnsGroupRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*connectivity.EdgeNextClient)
	service := sdns.NewSdnsService(client)

	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return fmt.Errorf("invalid group ID: %s", d.Id())
	}

	log.Printf("[DEBUG] Reading DNS group: %d", id)
	info, err := service.GetDnsGroupInfo(id)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("failed to get DNS group info: %w", err)
	}

	d.Set("group_name", info.GroupName)
	d.Set("remark", info.Remark)

	return nil
}

func resourceDnsGroupUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*connectivity.EdgeNextClient)
	service := sdns.NewSdnsService(client)

	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return fmt.Errorf("invalid group ID: %s", d.Id())
	}

	req := sdns.DnsGroupSaveRequest{
		GroupID:   id,
		GroupName: d.Get("group_name").(string),
		Remark:    d.Get("remark").(string),
	}

	if v, ok := d.GetOk("domain_ids"); ok {
		req.DomainIDs = expandIntList(v.(*schema.Set).List())
	}

	log.Printf("[INFO] Updating DNS group: %d", id)
	if err := service.UpdateDnsGroup(req); err != nil {
		return fmt.Errorf("failed to update DNS group: %w", err)
	}

	return resourceDnsGroupRead(d, m)
}

func resourceDnsGroupDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*connectivity.EdgeNextClient)
	service := sdns.NewSdnsService(client)

	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return fmt.Errorf("invalid group ID: %s", d.Id())
	}

	log.Printf("[INFO] Deleting DNS group: %d", id)
	err = service.DeleteDnsGroup(id)
	if err != nil {
		return fmt.Errorf("failed to delete DNS group: %w", err)
	}

	d.SetId("")
	return nil
}

func expandIntList(list []interface{}) []int {
	vs := make([]int, 0, len(list))
	for _, v := range list {
		if s, ok := v.(string); ok {
			i, _ := strconv.Atoi(s)
			vs = append(vs, i)
		} else if i, ok := v.(int); ok {
			vs = append(vs, i)
		}
	}
	return vs
}
