package resource

import (
	"fmt"
	"log"
	"strconv"

	"github.com/edgenextapisdk/terraform-provider-edgenext/edgenext/connectivity"
	"github.com/edgenextapisdk/terraform-provider-edgenext/edgenext/services/sdns"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceEdgenextDnsRecord() *schema.Resource {
	return &schema.Resource{
		Create: resourceDnsRecordCreate,
		Read:   resourceDnsRecordRead,
		Update: resourceDnsRecordUpdate,
		Delete: resourceDnsRecordDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"domain_id": {
				Type:        schema.TypeInt,
				Required:    true,
				ForceNew:    true,
				Description: "The ID of the domain",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the record (e.g., www)",
			},
			"type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The type of the record (A, CNAME, etc.)",
			},
			"view": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The view/line for the record",
			},
			"value": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The value of the record",
			},
			"mx": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     0,
				Description: "MX priority",
			},
			"ttl": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     600,
				Description: "TTL in seconds",
			},
			"remark": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Remark for the record",
			},
			"status": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     1,
				Description: "Status of the record (1 for enabled, 2 for paused)",
			},
		},
	}
}

func resourceDnsRecordCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*connectivity.EdgeNextClient)
	service := sdns.NewSdnsService(client)

	req := sdns.DnsRecordAddRequest{
		DomainID:     d.Get("domain_id").(int),
		RecordName:   d.Get("name").(string),
		RecordType:   d.Get("type").(string),
		RecordView:   d.Get("view").(string),
		RecordValue:  d.Get("value").(string),
		RecordMX:     d.Get("mx").(int),
		RecordTTL:    d.Get("ttl").(int),
		RecordRemark: d.Get("remark").(string),
	}

	log.Printf("[INFO] Creating DNS record: %s in domain %d", req.RecordName, req.DomainID)
	id, err := service.AddDnsRecord(req)
	if err != nil {
		return fmt.Errorf("failed to create DNS record: %w", err)
	}

	d.SetId(strconv.Itoa(id))
	return resourceDnsRecordRead(d, m)
}

func resourceDnsRecordRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*connectivity.EdgeNextClient)
	service := sdns.NewSdnsService(client)

	domainID := d.Get("domain_id").(int)
	recordID, err := strconv.Atoi(d.Id())
	if err != nil {
		return fmt.Errorf("invalid record ID: %s", d.Id())
	}

	log.Printf("[DEBUG] Reading DNS record: %d in domain %d", recordID, domainID)
	listReq := sdns.DnsRecordListRequest{
		DomainID: domainID,
		PerPage:  1000,
	}
	resp, err := service.ListDnsRecords(listReq)
	if err != nil {
		return fmt.Errorf("failed to list DNS records: %w", err)
	}

	var foundRecord *sdns.DnsRecord
	for _, r := range resp.List {
		if r.ID == recordID {
			foundRecord = &r
			break
		}
	}

	if foundRecord == nil {
		log.Printf("[WARN] DNS record %d not found in domain %d", recordID, domainID)
		d.SetId("")
		return nil
	}

	d.Set("name", foundRecord.Name)
	d.Set("type", foundRecord.Type)
	d.Set("view", foundRecord.View)
	d.Set("value", foundRecord.Value)
	d.Set("mx", foundRecord.MX)
	d.Set("ttl", foundRecord.TTL)
	d.Set("status", foundRecord.Status)
	d.Set("remark", foundRecord.Remark)

	return nil
}

func resourceDnsRecordUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*connectivity.EdgeNextClient)
	service := sdns.NewSdnsService(client)

	recordID, err := strconv.Atoi(d.Id())
	if err != nil {
		return fmt.Errorf("invalid record ID: %s", d.Id())
	}
	domainID := d.Get("domain_id").(int)

	req := sdns.DnsRecordEditRequest{
		RecordID:     recordID,
		DomainID:     domainID,
		RecordName:   d.Get("name").(string),
		RecordType:   d.Get("type").(string),
		RecordView:   d.Get("view").(string),
		RecordValue:  d.Get("value").(string),
		RecordMX:     d.Get("mx").(int),
		RecordTTL:    d.Get("ttl").(int),
		RecordRemark: d.Get("remark").(string),
	}

	log.Printf("[INFO] Updating DNS record: %d in domain %d", recordID, domainID)
	err = service.UpdateDnsRecord(req)
	if err != nil {
		return fmt.Errorf("failed to update DNS record: %w", err)
	}

	return resourceDnsRecordRead(d, m)
}

func resourceDnsRecordDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*connectivity.EdgeNextClient)
	service := sdns.NewSdnsService(client)

	recordID, err := strconv.Atoi(d.Id())
	if err != nil {
		return fmt.Errorf("invalid record ID: %s", d.Id())
	}
	domainID := d.Get("domain_id").(int)

	log.Printf("[INFO] Deleting DNS record: %d in domain %d", recordID, domainID)
	err = service.DeleteDnsRecord(recordID, domainID)
	if err != nil {
		return fmt.Errorf("failed to delete DNS record: %w", err)
	}

	d.SetId("")
	return nil
}
