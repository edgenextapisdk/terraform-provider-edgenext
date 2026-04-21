import os
import glob

base_dir = "/Users/tianyidong/go/terraform-provider-edgenext/edgenext/services/ecs"
os.makedirs(base_dir, exist_ok=True)

def to_camel_case(snake_str):
    components = snake_str.split('_')
    return ''.join(x.title() for x in components)

resources_spec = {
    "instance": {
        "path": "/instance",
        "fields": {
            "name": ("schema.TypeString", "Required", False),
            "region": ("schema.TypeString", "Optional", False),
            "flavor_ref": ("schema.TypeString", "Required", False),
            "image_ref": ("schema.TypeString", "Required", False),
            "admin_pass": ("schema.TypeString", "Required", False),
            "key_name": ("schema.TypeString", "Optional", False),
            "project_id": ("schema.TypeString", "Optional", False),
            "bandwidth": ("schema.TypeInt", "Optional", False),
            "status": ("schema.TypeString", "Computed", False),
            "networks": ("schema.TypeList", "Optional", True),
            "security_groups": ("schema.TypeList", "Optional", True)
        }
    },
    "key_pair": {
        "path": "/keypair",
        "fields": {
            "name": ("schema.TypeString", "Required", False),
            "public_key": ("schema.TypeString", "Optional", False),
            "private_key": ("schema.TypeString", "Computed", False)
        }
    },
    "image": {
        "path": "/image",
        "fields": {
            "name": ("schema.TypeString", "Required", False),
            "instance_id": ("schema.TypeString", "Optional", False),
            "description": ("schema.TypeString", "Optional", False),
            "os_distro": ("schema.TypeString", "Computed", False)
        }
    },
    "vpc": {
        "path": "/vpc",
        "fields": {
            "name": ("schema.TypeString", "Required", False),
            "cidr": ("schema.TypeString", "Required", False),
            "description": ("schema.TypeString", "Optional", False)
        }
    },
    "router": {
        "path": "/router",
        "fields": {
            "name": ("schema.TypeString", "Required", False),
            "vpc_id": ("schema.TypeString", "Required", False),
            "admin_state_up": ("schema.TypeBool", "Optional", False)
        }
    },
    "floating_ip": {
        "path": "/floatingIp",
        "fields": {
            "bandwidth": ("schema.TypeInt", "Required", False),
            "ip_address": ("schema.TypeString", "Computed", False)
        }
    },
    "network_interface": {
        "path": "/port",
        "fields": {
            "network_id": ("schema.TypeString", "Required", False),
            "subnet_id": ("schema.TypeString", "Optional", False),
            "mac_address": ("schema.TypeString", "Computed", False)
        }
    },
    "security_group": {
        "path": "/security_group",
        "fields": {
            "name": ("schema.TypeString", "Required", False),
            "description": ("schema.TypeString", "Optional", False)
        }
    },
    "disk": {
        "path": "/volume",
        "fields": {
            "name": ("schema.TypeString", "Required", False),
            "volume_type": ("schema.TypeString", "Required", False),
            "size": ("schema.TypeInt", "Required", False)
        }
    },
    "tag": {
        "path": "/tags",
        "fields": {
            "key": ("schema.TypeString", "Required", False),
            "value": ("schema.TypeString", "Required", False)
        }
    }
}

for res, spec in resources_spec.items():
    camel_res = to_camel_case(res)
    api_path = spec["path"]
    fields = spec["fields"]
    
    schema_block = []
    for f, f_meta in fields.items():
        elemType, constraints, is_list = f_meta
        if is_list:
            schema_block.append(f"""			"{f}": {{
				Type:        {elemType},
				{constraints}:    true,
				Elem: &schema.Schema{{Type: schema.TypeString}},
				Description: "{f} description",
			}},""")
        elif elemType == "schema.TypeBool":
            schema_block.append(f"""			"{f}": {{
				Type:        {elemType},
				{constraints}:    true,
                Default: false,
				Description: "{f} description",
			}},""")
        else:
            schema_block.append(f"""			"{f}": {{
				Type:        {elemType},
				{constraints}:    true,
				Description: "{f} description",
			}},""")
    schema_str = "\n".join(schema_block)
    
    req_block = []
    for f, f_meta in fields.items():
        if f_meta[1] == "Computed": continue
        if f == "name": continue
        
        if f_meta[0] == "schema.TypeInt":
            req_block.append(f'		"{f}": d.Get("{f}").(int),')
        elif f_meta[0] == "schema.TypeBool":
            req_block.append(f'		"{f}": d.Get("{f}").(bool),')
        elif f_meta[2]: 
            req_block.append(f"""		"{f}": func() []string {{
            lst := d.Get("{f}").([]interface{{}})
            res := make([]string, len(lst))
            for i, v := range lst {{ res[i] = v.(string) }}
            return res
        }}(),""")
        else:
            req_block.append(f'		"{f}": d.Get("{f}").(string),')
            
    req_str = "\n".join(req_block)
    
    read_block = []
    for f, f_meta in fields.items():
        if f == "name": continue
        read_block.append(f"""	if val, ok := resp["{f}"]; ok {{
		d.Set("{f}", val)
	}}""")
    read_str = "\n".join(read_block)

    res_file_content = f"""package ecs

import (
	"context"
    "fmt"

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
{schema_str}
		}},
	}}
}}

func resourceENECS{camel_res}Create(ctx context.Context, d *schema.ResourceData, m interface{{}}) diag.Diagnostics {{
	client := m.(*connectivity.EdgeNextClient)
	ecsClient, err := client.ECSClient()
	if err != nil {{
		return diag.FromErr(err)
	}}

	req := map[string]interface{{}}{{
{req_str}
	}}
    if n, ok := d.GetOk("name"); ok {{
        req["name"] = n.(string)
    }}
	var resp map[string]interface{{}}

	path := "/openapi/v2{api_path}/create"
	if "{res}" == "instance" {{
		path = "/openapi/v2{api_path}/create_order"
	}}

	err = ecsClient.Post(ctx, path, req, &resp)
	if err != nil {{
		return diag.Errorf("failed to create ECS {res}: %s", err)
	}}

	if id, ok := resp["id"].(string); ok {{
		d.SetId(id)
	}} else if idFloat, ok := resp["id"].(float64); ok {{
        d.SetId(fmt.Sprintf("%v", idFloat))
    }} else if n, ok := d.GetOk("name"); ok {{
        d.SetId(n.(string))
    }} else {{
		d.SetId("created-{res}")
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
		d.SetId("") // assume destroyed
		return nil
	}}

	if name, ok := resp["name"].(string); ok {{
		d.Set("name", name)
	}}
{read_str}

	return nil
}}

func resourceENECS{camel_res}Update(ctx context.Context, d *schema.ResourceData, m interface{{}}) diag.Diagnostics {{
	client := m.(*connectivity.EdgeNextClient)
	ecsClient, err := client.ECSClient()
	if err != nil {{
		return diag.FromErr(err)
	}}

    // Try name update
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

print("Regenerated resources without import errors.")
