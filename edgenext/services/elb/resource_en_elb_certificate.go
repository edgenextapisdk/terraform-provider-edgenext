package elb

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/edgenextapisdk/terraform-provider-edgenext/edgenext/connectivity"
	"github.com/edgenextapisdk/terraform-provider-edgenext/edgenext/helper"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

const (
	elbCertCreatePath = "/elb/openapi/v2/certificates/create"
	elbCertDeletePath = "/elb/openapi/v2/certificates/delete"
	elbCertDetailPath = "/elb/openapi/v2/certificates/detail"
)

// ResourceENELBCertificate manages an EdgeNext ELB certificate (Barbican container). There is no update API.
func ResourceENELBCertificate() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceENELBCertificateCreate,
		ReadContext:   resourceENELBCertificateRead,
		DeleteContext: resourceENELBCertificateDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceENELBCertificateImport,
		},
		Description: "Manages an EdgeNext ELB TLS certificate (create and delete; refresh uses certificates/detail).",
		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				Description:  "Certificate name.",
				ValidateFunc: validation.StringIsNotWhiteSpace,
			},
			"certificate": {
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				Sensitive:        true,
				Description:      "PEM certificate body sent on create (chain supported).",
				ValidateFunc:     validation.StringIsNotWhiteSpace,
				DiffSuppressFunc: elbCertificatePEMDiffSuppress,
			},
			"certificate_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Certificate (Barbican container) ID returned by the API.",
			},
			"certificate_ref": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Certificate (Barbican container) reference URL.",
			},
			"type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Container type (for example certificate).",
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Certificate status.",
			},
			"created": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Creation time as Unix timestamp (seconds).",
			},
			"updated": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Last update time as Unix timestamp (seconds).",
			},
			"expiration": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Expiration time as Unix timestamp (seconds).",
			},
			"private_key": {
				Type:             schema.TypeString,
				Optional:         true,
				Computed:         true,
				ForceNew:         true,
				Sensitive:        true,
				Description:      "PEM private key sent on create (required for create). Detail API usually omits the key; when omitted, the prior value is kept and must not force replacement.",
				DiffSuppressFunc: elbCertificatePEMDiffSuppress,
			},
			"private_key_passphrase": {
				Type:        schema.TypeString,
				Computed:    true,
				Sensitive:   true,
				Description: "Private key passphrase from detail when present.",
			},
			"cert_detail": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Extra certificate detail as JSON when the API returns an object; empty when null.",
			},
		},
	}
}

func resourceENELBCertificateImport(ctx context.Context, d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	certificateID := strings.TrimSpace(d.Id())
	if certificateID == "" {
		return nil, fmt.Errorf("import id must be the certificate_id (Barbican container UUID, non-empty)")
	}
	if err := d.Set("certificate_id", certificateID); err != nil {
		return nil, err
	}
	d.SetId(certificateID)
	diags := resourceENELBCertificateRead(ctx, d, m)
	if diags.HasError() {
		return nil, elbCertificateDiagError(diags)
	}
	return []*schema.ResourceData{d}, nil
}

func elbCertificateDiagError(diags diag.Diagnostics) error {
	if len(diags) == 0 {
		return fmt.Errorf("import failed")
	}
	e := diags[0]
	if e.Detail != "" {
		return fmt.Errorf("%s: %s", e.Summary, e.Detail)
	}
	return fmt.Errorf("%s", e.Summary)
}

func resourceENELBCertificateCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*connectivity.EdgeNextClient)
	elbClient, err := client.ELBClient()
	if err != nil {
		return diag.FromErr(err)
	}

	privateKey := strings.TrimSpace(d.Get("private_key").(string))
	if privateKey == "" {
		return diag.Errorf("private_key is required for create")
	}
	req := map[string]interface{}{
		"name":        strings.TrimSpace(d.Get("name").(string)),
		"certificate": d.Get("certificate").(string),
		"private_key": privateKey,
	}
	var resp map[string]interface{}
	if err := elbClient.Post(ctx, elbCertCreatePath, req, &resp); err != nil {
		return diag.Errorf("failed to create ELB certificate: %s", err)
	}
	payload, err := helper.ParseAPIResponseMap(resp)
	if err != nil {
		return diag.Errorf("failed to parse ELB certificate create response: %s", err)
	}
	certificateID := helper.StringFromMap(payload, "container_id")
	if certificateID == "" {
		return diag.Errorf("create response missing container_id")
	}
	d.SetId(certificateID)
	_ = d.Set("certificate_id", certificateID)
	_ = d.Set("certificate_ref", helper.StringFromMap(payload, "container_ref"))
	return resourceENELBCertificateRead(ctx, d, m)
}

func resourceENELBCertificateRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*connectivity.EdgeNextClient)
	elbClient, err := client.ELBClient()
	if err != nil {
		return diag.FromErr(err)
	}

	certificateID := strings.TrimSpace(d.Get("certificate_id").(string))
	if certificateID == "" {
		certificateID = strings.TrimSpace(d.Id())
	}
	if certificateID == "" {
		d.SetId("")
		return nil
	}

	req := map[string]interface{}{
		"container_id": certificateID,
	}
	var resp map[string]interface{}
	if err := elbClient.Post(ctx, elbCertDetailPath, req, &resp); err != nil {
		if elbCertificateAPIGone(err) {
			d.SetId("")
			return nil
		}
		return diag.Errorf("failed to read ELB certificate: %s", err)
	}
	payload, err := helper.ParseAPIResponseMap(resp)
	if err != nil {
		if elbCertificateAPIGone(err) {
			d.SetId("")
			return nil
		}
		return diag.Errorf("failed to parse ELB certificate detail response: %s", err)
	}

	_ = d.Set("certificate_id", helper.StringFromMap(payload, "container_id"))
	_ = d.Set("certificate_ref", helper.StringFromMap(payload, "container_ref"))
	_ = d.Set("name", helper.StringFromMap(payload, "name"))
	_ = d.Set("type", helper.StringFromMap(payload, "type"))
	_ = d.Set("status", helper.StringFromMap(payload, "status"))
	_ = d.Set("created", helper.IntFromMap(payload, "created"))
	_ = d.Set("updated", helper.IntFromMap(payload, "updated"))
	_ = d.Set("expiration", helper.IntFromMap(payload, "expiration"))
	if cert := strings.TrimSpace(helper.StringFromMap(payload, "certificate")); cert != "" {
		_ = d.Set("certificate", cert)
	}
	if pk := strings.TrimSpace(helper.StringFromMap(payload, "private_key")); pk != "" {
		_ = d.Set("private_key", pk)
	}
	_ = d.Set("private_key_passphrase", helper.StringFromMap(payload, "private_key_passphrase"))
	_ = d.Set("cert_detail", certDetailToString(payload["cert_detail"]))

	id := helper.StringFromMap(payload, "container_id")
	if id == "" {
		id = certificateID
	}
	d.SetId(id)
	return nil
}

func resourceENELBCertificateDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*connectivity.EdgeNextClient)
	elbClient, err := client.ELBClient()
	if err != nil {
		return diag.FromErr(err)
	}

	certificateID := strings.TrimSpace(d.Get("certificate_id").(string))
	if certificateID == "" {
		certificateID = strings.TrimSpace(d.Id())
	}
	if certificateID == "" {
		return nil
	}

	req := map[string]interface{}{
		"container_id": certificateID,
	}
	var resp map[string]interface{}
	if err := elbClient.Post(ctx, elbCertDeletePath, req, &resp); err != nil {
		if elbCertificateAPIGone(err) {
			d.SetId("")
			return nil
		}
		return diag.Errorf("failed to delete ELB certificate: %s", err)
	}
	if _, err := helper.ParseAPIResponseMap(resp); err != nil {
		if elbCertificateAPIGone(err) {
			d.SetId("")
			return nil
		}
		return diag.Errorf("failed to parse ELB certificate delete response: %s", err)
	}
	d.SetId("")
	return nil
}

func elbCertificateAPIGone(err error) bool {
	if err == nil {
		return false
	}
	msg := strings.ToLower(err.Error())
	return strings.Contains(msg, "not found") ||
		strings.Contains(msg, "404") ||
		strings.Contains(msg, "does not exist") ||
		strings.Contains(msg, "not exist")
}

func certDetailToString(v interface{}) string {
	if v == nil {
		return ""
	}
	if m, ok := v.(map[string]interface{}); ok && len(m) == 0 {
		return ""
	}
	b, err := json.Marshal(v)
	if err != nil {
		return ""
	}
	return string(b)
}

// elbCertificatePEMDiffSuppress treats PEM bodies as equal when they differ only by whitespace or line endings.
// This avoids perpetual replacement when the detail API returns a normalized PEM that does not match the HCL heredoc.
func elbCertificatePEMDiffSuppress(_, old, new string, _ *schema.ResourceData) bool {
	return elbNormalizePEMString(old) == elbNormalizePEMString(new)
}

func elbNormalizePEMString(s string) string {
	s = strings.TrimSpace(s)
	s = strings.ReplaceAll(s, "\r\n", "\n")
	s = strings.ReplaceAll(s, "\r", "\n")
	return s
}
