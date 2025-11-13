package domain

import (
	"fmt"
	"log"

	"github.com/edgenextapisdk/terraform-provider-edgenext/edgenext/connectivity"
	"github.com/edgenextapisdk/terraform-provider-edgenext/edgenext/helper"
	"github.com/edgenextapisdk/terraform-provider-edgenext/edgenext/services/scdn"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// DataSourceEdgenextScdnDomainBaseSettings returns the SCDN domain base settings data source
func DataSourceEdgenextScdnDomainBaseSettings() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceScdnDomainBaseSettingsRead,

		Schema: map[string]*schema.Schema{
			"domain_id": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "The ID of the domain to query base settings",
			},
			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results to a file",
			},
			// Computed fields
			"proxy_host": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Proxy host configuration",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"proxy_host": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Proxy host value",
						},
						"proxy_host_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Proxy host type",
						},
					},
				},
			},
			"proxy_sni": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Proxy SNI configuration",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"proxy_sni": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Proxy SNI value",
						},
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Proxy SNI status",
						},
					},
				},
			},
			"domain_redirect": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Domain redirect configuration",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Redirect status",
						},
						"jump_to": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Redirect target URL",
						},
						"jump_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Redirect jump type",
						},
					},
				},
			},
		},
	}
}

func dataSourceScdnDomainBaseSettingsRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*connectivity.EdgeNextClient)
	service := scdn.NewScdnService(client)

	domainID := d.Get("domain_id").(int)

	req := scdn.DomainBaseSettingsGetRequest{
		DomainID: domainID,
	}

	log.Printf("[INFO] Querying SCDN domain base settings for domain: %d", domainID)
	response, err := service.GetDomainBaseSettings(req)
	if err != nil {
		return fmt.Errorf("failed to query SCDN domain base settings: %w", err)
	}

	// Set domain ID as resource ID
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

	// Write result to output file if specified
	if outputFile := d.Get("result_output_file").(string); outputFile != "" {
		outputData := map[string]interface{}{
			"domain_id":       domainID,
			"proxy_host":      response.Data.ProxyHost,
			"proxy_sni":       response.Data.ProxySNI,
			"domain_redirect": response.Data.DomainRedirect,
		}
		if err := helper.WriteToFile(d, outputData); err != nil {
			return fmt.Errorf("failed to write output file: %w", err)
		}
	}

	log.Printf("[INFO] SCDN domain base settings query successful for domain: %d", domainID)
	return nil
}
