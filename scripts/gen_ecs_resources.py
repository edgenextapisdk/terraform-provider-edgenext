import os

def to_camel_case(snake_str):
    components = snake_str.split('_')
    return ''.join(x.title() for x in components)

def get_paths(res_type):
    if res_type == "floating_ip":
        return "/floatingIp"
    elif res_type == "network_interface":
        return "/port"
    elif res_type == "disk":
        return "/volume"
    elif res_type == "tag":
        return "/tags"
    elif res_type == "instance":
        return "/instance"
    elif res_type == "key_pair":
        return "/keypair"
    elif res_type == "image":
        return "/image"
    elif res_type == "security_group":
        return "/security_group"
    else:
        return f"/{res_type}"

resources = [
    "instance",
    "key_pair",
    "image",
    "vpc",
    "router",
    "floating_ip",
    "network_interface",
    "security_group",
    "disk",
    "tag"
]

base_dir = "/Users/tianyidong/go/terraform-provider-edgenext/edgenext/services/ecs"
os.makedirs(base_dir, exist_ok=True)

for res in resources:
    camel_res = to_camel_case(res)
    api_path = get_paths(res)
    
    # 1. resource file
    res_file_content = f"""package ecs

import (
	"context"

	"github.com/edgenextapisdk/terraform-provider-edgenext/edgenext/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// ResourceENECS{camel_res} returns the resource schema for ECS {res}.
func ResourceENECS{camel_res}() *schema.Resource {{
	return &schema.Resource{{
		CreateContext: resourceENECS{camel_res}Create,
		ReadContext:   resourceENECS{camel_res}Read,
		UpdateContext: resourceENECS{camel_res}Update,
		DeleteContext: resourceENECS{camel_res}Delete,
		Importer: &schema.ResourceImporter{{
			StateContext: schema.ImportStatePassthroughContext,
		}},
		Description: "Provides an EdgeNext ECS {res} resource.",
		Schema: map[string]*schema.Schema{{
			"name": {{
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the {res}.",
			}},
			// Add other necessary fields here based on API request parameters
		}},
	}}
}}

func resourceENECS{camel_res}Create(ctx context.Context, d *schema.ResourceData, m interface{{}}) diag.Diagnostics {{
	client := m.(*connectivity.EdgeNextClient)
	ecsClient, err := client.ECSClient()
	if err != nil {{
		return diag.FromErr(err)
	}}

	name := d.Get("name").(string)
	req := map[string]interface{{}}{{
		"name": name,
	}}
	var resp map[string]interface{{}}

	// Usually instance uses create_order, others use create
	path := "/openapi/v2{api_path}/create"
	if "{res}" == "instance" {{
		path = "/openapi/v2{api_path}/create_order"
	}}

	err = ecsClient.Post(ctx, path, req, &resp)
	if err != nil {{
		return diag.Errorf("failed to create ECS {res}: %s", err)
	}}

	// Assuming the response contains an ID field
	if id, ok := resp["id"].(string); ok {{
		d.SetId(id)
	}} else {{
		d.SetId(name) // fallback
	}}

	return resourceENECS{camel_res}Read(ctx, d, m)
}}

func resourceENECS{camel_res}Read(ctx context.Context, d *schema.ResourceData, m interface{{}}) diag.Diagnostics {{
	client := m.(*connectivity.EdgeNextClient)
	ecsClient, err := client.ECSClient()
	if err != nil {{
		return diag.FromErr(err)
	}}

	req := map[string]interface{{}}{{
		"id": d.Id(),
	}}
	var resp map[string]interface{{}}

	err = ecsClient.Post(ctx, "/openapi/v2{api_path}/detail", req, &resp)
	if err != nil {{
		d.SetId("") // If not found, mark as destroyed
		return nil
	}}

	if name, ok := resp["name"].(string); ok {{
		d.Set("name", name)
	}}

	return nil
}}

func resourceENECS{camel_res}Update(ctx context.Context, d *schema.ResourceData, m interface{{}}) diag.Diagnostics {{
	client := m.(*connectivity.EdgeNextClient)
	ecsClient, err := client.ECSClient()
	if err != nil {{
		return diag.FromErr(err)
	}}

	if d.HasChange("name") {{
		req := map[string]interface{{}}{{
			"id":   d.Id(),
			"name": d.Get("name"),
		}}
		var resp map[string]interface{{}}
		
		err = ecsClient.Post(ctx, "/openapi/v2{api_path}/update", req, &resp)
		if err != nil {{
			return diag.Errorf("failed to update ECS {res}: %s", err)
		}}
	}}

	return resourceENECS{camel_res}Read(ctx, d, m)
}}

func resourceENECS{camel_res}Delete(ctx context.Context, d *schema.ResourceData, m interface{{}}) diag.Diagnostics {{
	client := m.(*connectivity.EdgeNextClient)
	ecsClient, err := client.ECSClient()
	if err != nil {{
		return diag.FromErr(err)
	}}

	req := map[string]interface{{}}{{
		"id": d.Id(),
	}}
	var resp map[string]interface{{}}

	err = ecsClient.Post(ctx, "/openapi/v2{api_path}/delete", req, &resp)
	if err != nil {{
		return diag.Errorf("failed to delete ECS {res}: %s", err)
	}}

	return nil
}}
"""

    with open(os.path.join(base_dir, f"resource_en_ecs_{res}.go"), "w") as f:
        f.write(res_file_content)
        
    # 2. data file
    data_file_content = f"""package ecs

import (
	"context"

	"github.com/edgenextapisdk/terraform-provider-edgenext/edgenext/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// DataSourceENECS{camel_res}s returns the data source schema for ECS {res}s.
func DataSourceENECS{camel_res}s() *schema.Resource {{
	return &schema.Resource{{
		ReadContext: dataSourceENECS{camel_res}sRead,
		Description: "Data source to query EdgeNext ECS {res}s.",
		Schema: map[string]*schema.Schema{{
			"name_regex": {{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "A regex string to filter results by {res} name.",
			}},
			"ids": {{
				Type:        schema.TypeList,
				Optional:    true,
				Elem:        &schema.Schema{{Type: schema.TypeString}},
				Description: "A list of {res} IDs.",
			}},
			"{res}s": {{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "A list of ECS {res}s.",
				Elem: &schema.Resource{{
					Schema: map[string]*schema.Schema{{
						"id": {{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The ID of the {res}.",
						}},
						"name": {{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The name of the {res}.",
						}},
					}},
				}},
			}},
		}},
	}}
}}

func dataSourceENECS{camel_res}sRead(ctx context.Context, d *schema.ResourceData, m interface{{}}) diag.Diagnostics {{
	client := m.(*connectivity.EdgeNextClient)
	ecsClient, err := client.ECSClient()
	if err != nil {{
		return diag.FromErr(err)
	}}

	req := map[string]interface{{}}{{}}
	var resp map[string]interface{{}}

	// List action
	err = ecsClient.Post(ctx, "/openapi/v2{api_path}/list", req, &resp)
	if err != nil {{
		return diag.Errorf("failed to read ECS {res}s: %s", err)
	}}

	// Extract and filter logic should be implemented below based on the actual API struct
	// This is a minimal placeholder
	if dataList, ok := resp["data"].([]interface{{}}); ok {{
		// Filter logic here
		d.SetId("en-ecs-{res}s-id")
		d.Set("{res}s", dataList)
	}}

	return nil
}}
"""
    with open(os.path.join(base_dir, f"data_source_en_ecs_{res}s.go"), "w") as f:
        f.write(data_file_content)

print("Regenerated all files successfully with error handling.")
