package resource

import (
	"fmt"
	"log"
	"strconv"

	"github.com/edgenextapisdk/terraform-provider-edgenext/edgenext/connectivity"
	"github.com/edgenextapisdk/terraform-provider-edgenext/edgenext/services/scdn/domain_group"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceEdgenextScdnDomainGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceScdnDomainGroupCreate,
		Read:   resourceScdnDomainGroupRead,
		Update: resourceScdnDomainGroupUpdate,
		Delete: resourceScdnDomainGroupDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"group_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the domain group",
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
			"domains": {
				Type:        schema.TypeSet,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "List of domains to bind to the group",
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Creation time",
			},
			"updated_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Last update time",
			},
		},
	}
}

func resourceScdnDomainGroupCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*connectivity.EdgeNextClient)
	service := domain_group.NewDomainGroupService(client)

	req := domain_group.DomainGroupSaveRequest{
		GroupName: d.Get("group_name").(string),
		Remark:    d.Get("remark").(string),
	}

	log.Printf("[INFO] Creating SCDN Domain Group: %s", req.GroupName)
	resp, err := service.AddDomainGroup(req)
	if err != nil {
		return fmt.Errorf("failed to create Domain Group: %w", err)
	}

	groupID := resp.Data.ID
	d.SetId(groupID)

	// Bind domains if provided
	if err := bindDomains(d, service, groupID); err != nil {
		// If binding fails, we still created the group, so we should probably return the id but also error?
		// Or try to cleanup? Terraform will mark as tainted if we error.
		return fmt.Errorf("failed to bind domains: %w", err)
	}

	return resourceScdnDomainGroupRead(d, m)
}

func resourceScdnDomainGroupRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*connectivity.EdgeNextClient)
	service := domain_group.NewDomainGroupService(client)

	groupID, err := strconv.Atoi(d.Id())
	if err != nil {
		return fmt.Errorf("invalid group ID: %s", d.Id())
	}

	log.Printf("[DEBUG] Reading SCDN Domain Group: %s", d.Id())
	groupInfo, err := service.GetDomainGroupInfo(groupID)
	if err != nil {
		// TODO: detecting 404/not found to unset ID like d.SetId("")
		// Assuming API returns error if not found, let's treat it as gone if specific error code?
		// For now simple error return.
		return fmt.Errorf("failed to get domain group info: %w", err)
	}

	d.Set("group_name", groupInfo.GroupName)
	d.Set("remark", groupInfo.Remark)
	d.Set("created_at", groupInfo.CreatedAt)
	d.Set("updated_at", groupInfo.UpdatedAt)

	// Read bound domains
	domainListReq := domain_group.DomainGroupDomainListRequest{
		GroupID: groupID,
		Page:    1,
		PerPage: 1000,
	}
	domainResp, err := service.ListGroupDomains(domainListReq)
	if err == nil && domainResp != nil {
		var domainIDs []string
		var domains []string
		for _, item := range domainResp.Data.List {
			domainIDs = append(domainIDs, item.DomainID)
			domains = append(domains, item.Domain)
		}
		// Only set if we manage them?
		// If user configured domain_ids, we update state.
		// If user configured domains, we update state.
		if _, ok := d.GetOk("domain_ids"); ok {
			d.Set("domain_ids", domainIDs)
		}
		if _, ok := d.GetOk("domains"); ok {
			d.Set("domains", domains)
		}
	}

	return nil
}

func resourceScdnDomainGroupUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*connectivity.EdgeNextClient)
	service := domain_group.NewDomainGroupService(client)

	groupID, _ := strconv.Atoi(d.Id())

	if d.HasChange("group_name") || d.HasChange("remark") {
		req := domain_group.DomainGroupSaveRequest{
			GroupID:   groupID,
			GroupName: d.Get("group_name").(string),
			Remark:    d.Get("remark").(string),
		}
		log.Printf("[INFO] Updating SCDN Domain Group: %s", d.Id())
		_, err := service.UpdateDomainGroup(req)
		if err != nil {
			return fmt.Errorf("failed to update Domain Group: %w", err)
		}
	}

	if d.HasChange("domain_ids") || d.HasChange("domains") {
		// Re-bind (add/remove) logic could be complex (API supports action 'add' or 'del')
		// The `BindDomainsToGroup` API seems to add to existing? Or replace?
		// The API doc says: action: add/del.
		// So to sync with Terraform state (which is absolute), we might need to calc diff or unbind all then bind?
		// Simplest for now:
		// 1. Get current bound domains
		// 2. Diff with new config
		// 3. Add new, Del old

		// For simplicity in this implementation, if we assume declarative:
		// We probably need to implement full diff logic.
		// Or, simpler, if the API supports 'set' (replace all)... but it seems to support add/del only.

		if err := syncDomains(d, service, d.Id()); err != nil {
			return err
		}
	}

	return resourceScdnDomainGroupRead(d, m)
}

func resourceScdnDomainGroupDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*connectivity.EdgeNextClient)
	service := domain_group.NewDomainGroupService(client)

	groupID, _ := strconv.Atoi(d.Id())

	// First, unbind all domains from the group
	log.Printf("[INFO] Unbinding all domains from SCDN Domain Group: %s before deletion", d.Id())

	listReq := domain_group.DomainGroupDomainListRequest{
		GroupID: groupID,
		Page:    1,
		PerPage: 1000,
	}
	domainResp, err := service.ListGroupDomains(listReq)
	if err != nil {
		log.Printf("[WARN] Failed to list domains in group %s: %v, attempting delete anyway", d.Id(), err)
	} else if len(domainResp.Data.List) > 0 {
		// Collect all domain IDs to unbind
		var domainIDs []string
		for _, item := range domainResp.Data.List {
			domainIDs = append(domainIDs, item.DomainID)
		}

		log.Printf("[INFO] Unbinding %d domains from group %s", len(domainIDs), d.Id())
		unbindReq := domain_group.DomainGroupDomainSaveRequest{
			GroupID:   groupID,
			DomainIDs: domainIDs,
			Action:    "del",
		}
		if _, err := service.BindDomainsToGroup(unbindReq); err != nil {
			return fmt.Errorf("failed to unbind domains before deleting group: %w", err)
		}
	}

	// Now delete the group
	req := domain_group.DomainGroupDelRequest{
		GroupID: groupID,
	}

	log.Printf("[INFO] Deleting SCDN Domain Group: %s", d.Id())
	_, err = service.DeleteDomainGroup(req)
	if err != nil {
		return fmt.Errorf("failed to delete Domain Group: %w", err)
	}

	d.SetId("")
	return nil
}

func bindDomains(d *schema.ResourceData, service *domain_group.DomainGroupService, groupID string) error {
	groupIDInt, _ := strconv.Atoi(groupID)

	// Handle domain_ids
	if v, ok := d.GetOk("domain_ids"); ok {
		ids := v.(*schema.Set).List()
		var strIDs []string
		for _, id := range ids {
			strIDs = append(strIDs, id.(string))
		}
		if len(strIDs) > 0 {
			req := domain_group.DomainGroupDomainSaveRequest{
				GroupID:   groupIDInt,
				DomainIDs: strIDs,
				Action:    "add",
			}
			if _, err := service.BindDomainsToGroup(req); err != nil {
				return err
			}
		}
	}

	// Handle domains
	if v, ok := d.GetOk("domains"); ok {
		doms := v.(*schema.Set).List()
		var strDoms []string
		for _, dom := range doms {
			strDoms = append(strDoms, dom.(string))
		}
		if len(strDoms) > 0 {
			req := domain_group.DomainGroupDomainSaveRequest{
				GroupID: groupIDInt,
				Domains: strDoms,
				Action:  "add",
			}
			if _, err := service.BindDomainsToGroup(req); err != nil {
				return err
			}
		}
	}
	return nil
}

func syncDomains(d *schema.ResourceData, service *domain_group.DomainGroupService, groupID string) error {
	// 1. Get current domains
	groupIDInt, _ := strconv.Atoi(groupID)
	listReq := domain_group.DomainGroupDomainListRequest{GroupID: groupIDInt, Page: 1, PerPage: 1000}
	currentResp, err := service.ListGroupDomains(listReq)
	if err != nil {
		return fmt.Errorf("failed to list current domains: %w", err)
	}

	currentDomainIDs := make(map[string]bool)
	currentDomains := make(map[string]bool)
	for _, item := range currentResp.Data.List {
		currentDomainIDs[item.DomainID] = true
		currentDomains[item.Domain] = true
	}

	// 2. Handle domain_ids set
	if d.HasChange("domain_ids") {
		o, n := d.GetChange("domain_ids")
		oldSet := o.(*schema.Set)
		newSet := n.(*schema.Set)

		toAdd := newSet.Difference(oldSet).List()
		toDel := oldSet.Difference(newSet).List()

		if len(toAdd) > 0 {
			var ids []string
			for _, id := range toAdd {
				ids = append(ids, id.(string))
			}
			req := domain_group.DomainGroupDomainSaveRequest{GroupID: groupIDInt, DomainIDs: ids, Action: "add"}
			if _, err := service.BindDomainsToGroup(req); err != nil {
				return err
			}
		}
		if len(toDel) > 0 {
			var ids []string
			for _, id := range toDel {
				ids = append(ids, id.(string))
			}
			req := domain_group.DomainGroupDomainSaveRequest{GroupID: groupIDInt, DomainIDs: ids, Action: "del"}
			if _, err := service.BindDomainsToGroup(req); err != nil {
				return err
			}
		}
	}

	// 3. Handle domains set
	if d.HasChange("domains") {
		o, n := d.GetChange("domains")
		oldSet := o.(*schema.Set)
		newSet := n.(*schema.Set)

		toAdd := newSet.Difference(oldSet).List()
		toDel := oldSet.Difference(newSet).List()

		if len(toAdd) > 0 {
			var doms []string
			for _, d := range toAdd {
				doms = append(doms, d.(string))
			}
			req := domain_group.DomainGroupDomainSaveRequest{GroupID: groupIDInt, Domains: doms, Action: "add"}
			if _, err := service.BindDomainsToGroup(req); err != nil {
				return err
			}
		}
		if len(toDel) > 0 {
			var doms []string
			for _, d := range toDel {
				doms = append(doms, d.(string))
			}
			req := domain_group.DomainGroupDomainSaveRequest{GroupID: groupIDInt, Domains: doms, Action: "del"}
			if _, err := service.BindDomainsToGroup(req); err != nil {
				return err
			}
		}
	}

	return nil
}
