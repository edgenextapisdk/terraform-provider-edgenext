package data

import (
	"fmt"
	"log"

	"github.com/edgenextapisdk/terraform-provider-edgenext/edgenext/connectivity"
	"github.com/edgenextapisdk/terraform-provider-edgenext/edgenext/helper"
	"github.com/edgenextapisdk/terraform-provider-edgenext/edgenext/services/scdn"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// DataSourceEdgenextScdnOriginGroups returns the SCDN origin groups data source
func DataSourceEdgenextScdnOriginGroups() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceScdnOriginGroupsRead,

		Schema: map[string]*schema.Schema{
			"page": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     1,
				Description: "Page number",
			},
			"page_size": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     20,
				Description: "Page size",
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Origin group name filter",
			},
			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results to a file",
			},
			"total": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Total number of origin groups",
			},
			"origin_groups": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Origin group list",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Origin group ID",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Origin group name",
						},
						"remark": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Remark",
						},
						"member_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Member ID",
						},
						"username": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Username",
						},
						"created_at": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Creation time",
						},
						"updated_at": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Update time",
						},
					},
				},
			},
		},
	}
}

func dataSourceScdnOriginGroupsRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*connectivity.EdgeNextClient)
	service := scdn.NewScdnService(client)

	req := scdn.OriginGroupListRequest{
		Page:     d.Get("page").(int),
		PageSize: d.Get("page_size").(int),
	}

	if name, ok := d.GetOk("name"); ok {
		req.Name = name.(string)
	}

	log.Printf("[INFO] Querying SCDN origin groups")
	response, err := service.ListOriginGroups(req)
	if err != nil {
		return fmt.Errorf("failed to query origin groups: %w", err)
	}

	// Set resource ID
	d.SetId("origin-groups")

	// Set total
	if err := d.Set("total", response.Data.Total); err != nil {
		return fmt.Errorf("error setting total: %w", err)
	}

	// Convert origin groups to schema format
	originGroupList := make([]map[string]interface{}, 0, len(response.Data.List))
	for _, og := range response.Data.List {
		ogMap := map[string]interface{}{
			"id":         og.ID,
			"name":       og.Name,
			"remark":     og.Remark,
			"member_id":  og.MemberID,
			"username":   og.Username,
			"created_at": og.CreatedAt,
			"updated_at": og.UpdatedAt,
		}
		originGroupList = append(originGroupList, ogMap)
	}

	if err := d.Set("origin_groups", originGroupList); err != nil {
		return fmt.Errorf("error setting origin_groups: %w", err)
	}

	// Write result to output file if specified
	if outputFile := d.Get("result_output_file").(string); outputFile != "" {
		outputData := map[string]interface{}{
			"total":         response.Data.Total,
			"origin_groups": originGroupList,
		}
		if err := helper.WriteToFile(d, outputData); err != nil {
			return fmt.Errorf("failed to write output file: %w", err)
		}
	}

	log.Printf("[INFO] SCDN origin groups queried successfully: total=%d", response.Data.Total)
	return nil
}
