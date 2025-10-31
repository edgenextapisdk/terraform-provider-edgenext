package template

import (
	"fmt"
	"log"
	"strconv"

	"github.com/edgenextapisdk/terraform-provider-edgenext/edgenext/connectivity"
	"github.com/edgenextapisdk/terraform-provider-edgenext/edgenext/helper"
	"github.com/edgenextapisdk/terraform-provider-edgenext/edgenext/services/scdn"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// DataSourceEdgenextScdnRuleTemplateDomains returns the SCDN rule template domains data source
func DataSourceEdgenextScdnRuleTemplateDomains() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceScdnRuleTemplateDomainsRead,

		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "The rule template ID",
			},
			"app_type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The application type (e.g., 'network_speed')",
			},
			"page": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     1,
				Description: "Page number for pagination",
			},
			"page_size": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     1000,
				Description: "Items per page",
			},
			"domain": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Filter by domain name",
			},
			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results to a file",
			},
			// Computed fields
			"total": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Total number of domains bound to template",
			},
			"list": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "List of domain information",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Domain ID",
						},
						"domain": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Domain name",
						},
						"created_at": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Domain binding timestamp",
						},
					},
				},
			},
		},
	}
}

func dataSourceScdnRuleTemplateDomainsRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*connectivity.EdgeNextClient)
	service := scdn.NewScdnService(client)

	templateID := d.Get("id").(int)
	appType := d.Get("app_type").(string)

	req := scdn.RuleTemplateListDomainsRequest{
		ID:       templateID,
		AppType:  appType,
		Page:     d.Get("page").(int),
		PageSize: d.Get("page_size").(int),
	}

	if domain, ok := d.GetOk("domain"); ok {
		req.Domain = domain.(string)
	}

	log.Printf("[INFO] Listing domains for rule template: %d", templateID)
	response, err := service.ListRuleTemplateDomains(req)
	if err != nil {
		return fmt.Errorf("failed to list rule template domains: %w", err)
	}

	// Set the resource ID
	d.SetId(strconv.Itoa(templateID))

	// Set computed fields
	if err := d.Set("total", response.Data.Total); err != nil {
		log.Printf("[WARN] Failed to set total: %v", err)
	}

	// Convert list to schema format
	domainList := make([]map[string]interface{}, 0, len(response.Data.List))
	for _, domain := range response.Data.List {
		domainList = append(domainList, map[string]interface{}{
			"id":         domain.ID,
			"domain":     domain.Domain,
			"created_at": domain.CreatedAt,
		})
	}

	if err := d.Set("list", domainList); err != nil {
		log.Printf("[WARN] Failed to set list: %v", err)
	}

	// Save to file if requested
	if _, ok := d.GetOk("result_output_file"); ok {
		output := map[string]interface{}{
			"total": response.Data.Total,
			"list":  domainList,
		}
		if err := helper.WriteToFile(d, output); err != nil {
			log.Printf("[WARN] Failed to write result to file: %v", err)
		}
	}

	log.Printf("[INFO] Rule template domains listed successfully: %d domains", response.Data.Total)
	return nil
}
