package cert

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Resources returns all certificate-related resources
func Resources() map[string]*schema.Resource {
	return map[string]*schema.Resource{
		"edgenext_scdn_certificate":       ResourceEdgenextScdnCertificate(),
		"edgenext_scdn_certificate_apply": ResourceEdgenextScdnCertificateApply(),
	}
}

// DataSources returns all certificate-related data sources
func DataSources() map[string]*schema.Resource {
	return map[string]*schema.Resource{
		"edgenext_scdn_certificate":             DataSourceEdgenextScdnCertificate(),
		"edgenext_scdn_certificates":            DataSourceEdgenextScdnCertificates(),
		"edgenext_scdn_certificates_by_domains": DataSourceEdgenextScdnCertificatesByDomains(),
		"edgenext_scdn_certificate_export":      DataSourceEdgenextScdnCertificateExport(),
	}
}
