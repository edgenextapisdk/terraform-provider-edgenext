package scdn

import (
	"strings"
	"testing"
)

// ============================================================================
// DDoS Protection Tests
// ============================================================================

func TestScdnService_GetDdosProtectionConfig(t *testing.T) {
	tests := []struct {
		name string
		req  DdosProtectionGetConfigRequest
	}{
		{
			name: "Test GetDdosProtectionConfig",
			req: DdosProtectionGetConfigRequest{
				BusinessID: 947,
			},
		},
	}
	if !isIntegrationTest() {
		t.Skip("Skipping integration test: set EDGENEXT_ACCESS_KEY and EDGENEXT_SECRET_KEY to run")
	}

	client := createTestClient(t)
	service := NewScdnService(client)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := service.GetDdosProtectionConfig(tt.req)
			if err != nil {
				t.Errorf("ScdnService.GetDdosProtectionConfig() error = %v", err)
				return
			}
			t.Logf("Response Status: Code=%d, Message=%s", got.Status.Code, got.Status.Message)
			if got.Data.ApplicationDdosProtection != nil {
				t.Logf("ApplicationDdosProtection: Status=%s, Type=%s, AIStatus=%s",
					got.Data.ApplicationDdosProtection.Status,
					got.Data.ApplicationDdosProtection.Type,
					got.Data.ApplicationDdosProtection.AIStatus)
			}
			if got.Data.VisitorAuthentication != nil {
				t.Logf("VisitorAuthentication: Status=%s", got.Data.VisitorAuthentication.Status)
			}
		})
	}
}

func TestScdnService_UpdateDdosProtectionConfig(t *testing.T) {
	tests := []struct {
		name string
		req  DdosProtectionUpdateConfigRequest
	}{
		{
			name: "Test UpdateDdosProtectionConfig",
			req: DdosProtectionUpdateConfigRequest{
				BusinessID: 947,
				ApplicationDdosProtection: &ApplicationDdosProtection{
					Status:              "on",
					AICCStatus:          "on",
					Type:                "strict",
					NeedAttackDetection: 1,
					AIStatus:            "on",
				},
				VisitorAuthentication: &VisitorAuthentication{
					Status:         "off",
					AuthToken:      "",
					PassStillCheck: 0,
				},
			},
		},
	}
	if !isIntegrationTest() {
		t.Skip("Skipping integration test: set EDGENEXT_ACCESS_KEY and EDGENEXT_SECRET_KEY to run")
	}

	client := createTestClient(t)
	service := NewScdnService(client)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := service.UpdateDdosProtectionConfig(tt.req)
			if err != nil {
				t.Errorf("ScdnService.UpdateDdosProtectionConfig() error = %v", err)
				return
			}
			t.Logf("Response Status: Code=%d, Message=%s", got.Status.Code, got.Status.Message)
		})
	}
}

// ============================================================================
// WAF Rule Config Tests
// ============================================================================

func TestScdnService_GetWafRuleConfig(t *testing.T) {
	tests := []struct {
		name string
		req  WafRuleConfigGetRequest
	}{
		{
			name: "Test GetWafRuleConfig",
			req: WafRuleConfigGetRequest{
				BusinessID: 947,
			},
		},
	}
	if !isIntegrationTest() {
		t.Skip("Skipping integration test: set EDGENEXT_ACCESS_KEY and EDGENEXT_SECRET_KEY to run")
	}

	client := createTestClient(t)
	service := NewScdnService(client)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := service.GetWafRuleConfig(tt.req)
			if err != nil {
				t.Errorf("ScdnService.GetWafRuleConfig() error = %v", err)
				return
			}
			t.Logf("Response Status: Code=%d, Message=%s", got.Status.Code, got.Status.Message)
			if got.Data.WafRuleConfig != nil {
				t.Logf("WafRuleConfig: Status=%s, WafLevel=%s, WafMode=%s, AIStatus=%s",
					got.Data.WafRuleConfig.Status,
					got.Data.WafRuleConfig.WafLevel,
					got.Data.WafRuleConfig.WafMode,
					got.Data.WafRuleConfig.AIStatus)
			}
			if got.Data.WafInterceptPage != nil {
				t.Logf("WafInterceptPage: Status=%s, Type=%s",
					got.Data.WafInterceptPage.Status,
					got.Data.WafInterceptPage.Type)
			}
		})
	}
}

func TestScdnService_UpdateWafRuleConfig(t *testing.T) {
	tests := []struct {
		name string
		req  WafRuleConfigUpdateRequest
	}{
		{
			name: "Test UpdateWafRuleConfig",
			req: WafRuleConfigUpdateRequest{
				BusinessID: 947,
				WafRuleConfig: &WafRuleConfig{
					Status:   "on",
					AIStatus: "on",
					WafLevel: "strict",
					WafMode:  "block",
				},
				WafInterceptPage: &WafInterceptPage{
					Status:  "on",
					Type:    "default",
					Content: "",
				},
			},
		},
	}
	if !isIntegrationTest() {
		t.Skip("Skipping integration test: set EDGENEXT_ACCESS_KEY and EDGENEXT_SECRET_KEY to run")
	}

	client := createTestClient(t)
	service := NewScdnService(client)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := service.UpdateWafRuleConfig(tt.req)
			if err != nil {
				t.Errorf("ScdnService.UpdateWafRuleConfig() error = %v", err)
				return
			}
			t.Logf("Response Status: Code=%d, Message=%s", got.Status.Code, got.Status.Message)
		})
	}
}

// ============================================================================
// Security Protection Template Tests
// ============================================================================

func TestScdnService_GetMemberGlobalTemplate(t *testing.T) {
	if !isIntegrationTest() {
		t.Skip("Skipping integration test: set EDGENEXT_ACCESS_KEY and EDGENEXT_SECRET_KEY to run")
	}

	client := createTestClient(t)
	service := NewScdnService(client)
	got, err := service.GetMemberGlobalTemplate()
	if err != nil {
		t.Errorf("ScdnService.GetMemberGlobalTemplate() error = %v", err)
		return
	}
	t.Logf("Response Status: Code=%d, Message=%s", got.Status.Code, got.Status.Message)
	if got.Data.Template != nil {
		t.Logf("Template: ID=%d, Name=%s, Type=%s, BindDomainCount=%d",
			got.Data.Template.ID,
			got.Data.Template.Name,
			got.Data.Template.Type,
			got.Data.BindDomainCount)
	}
}

func TestScdnService_CreateSecurityProtectionTemplate(t *testing.T) {
	tests := []struct {
		name string
		req  SecurityProtectionTemplateCreateRequest
	}{
		{
			name: "Test CreateSecurityProtectionTemplate",
			req: SecurityProtectionTemplateCreateRequest{
				Name:   "test-security-template",
				Remark: "Test security protection template",
			},
		},
	}
	if !isIntegrationTest() {
		t.Skip("Skipping integration test: set EDGENEXT_ACCESS_KEY and EDGENEXT_SECRET_KEY to run")
	}

	client := createTestClient(t)
	service := NewScdnService(client)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := service.CreateSecurityProtectionTemplate(tt.req)
			if err != nil {
				t.Errorf("ScdnService.CreateSecurityProtectionTemplate() error = %v", err)
				return
			}
			t.Logf("Response Status: Code=%d, Message=%s", got.Status.Code, got.Status.Message)
			t.Logf("Created Template ID: %d", got.Data.BusinessID)
			if len(got.Data.FailDomains) > 0 {
				t.Logf("FailDomains: %+v", got.Data.FailDomains)
			}
		})
	}
}

func TestScdnService_CreateSecurityProtectionTemplateDomain(t *testing.T) {
	tests := []struct {
		name string
		req  SecurityProtectionTemplateCreateDomainRequest
	}{
		{
			name: "Test CreateSecurityProtectionTemplateDomain",
			req: SecurityProtectionTemplateCreateDomainRequest{
				DomainIDs:        []int{102008}, // Replace with actual domain IDs
				TemplateSourceID: 947,           // Replace with actual template ID
			},
		},
	}
	if !isIntegrationTest() {
		t.Skip("Skipping integration test: set EDGENEXT_ACCESS_KEY and EDGENEXT_SECRET_KEY to run")
	}

	client := createTestClient(t)
	service := NewScdnService(client)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := service.CreateSecurityProtectionTemplateDomain(tt.req)
			if err != nil {
				t.Errorf("ScdnService.CreateSecurityProtectionTemplateDomain() error = %v", err)
				return
			}
			t.Logf("Response Status: Code=%d, Message=%s", got.Status.Code, got.Status.Message)
			if len(got.Data.FailDomains) > 0 {
				t.Logf("FailDomains: %+v", got.Data.FailDomains)
			}
		})
	}
}

func TestScdnService_SearchSecurityProtectionTemplates(t *testing.T) {
	tests := []struct {
		name string
		req  SecurityProtectionTemplateSearchRequest
	}{
		{
			name: "Test SearchSecurityProtectionTemplates",
			req: SecurityProtectionTemplateSearchRequest{
				TplType:  "global",
				Page:     1,
				PageSize: 20,
			},
		},
	}
	if !isIntegrationTest() {
		t.Skip("Skipping integration test: set EDGENEXT_ACCESS_KEY and EDGENEXT_SECRET_KEY to run")
	}

	client := createTestClient(t)
	service := NewScdnService(client)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := service.SearchSecurityProtectionTemplates(tt.req)
			if err != nil {
				t.Errorf("ScdnService.SearchSecurityProtectionTemplates() error = %v", err)
				return
			}
			t.Logf("Response Status: Code=%d, Message=%s", got.Status.Code, got.Status.Message)
			t.Logf("Total Templates: %d", got.Data.Total)
			t.Logf("Templates Count: %d", len(got.Data.Templates))
			if len(got.Data.Templates) > 0 {
				t.Logf("First Template: ID=%d, Name=%s, Type=%s",
					got.Data.Templates[0].ID,
					got.Data.Templates[0].Name,
					got.Data.Templates[0].Type)
			}
		})
	}
}

func TestScdnService_GetSecurityProtectionTemplateBindDomains(t *testing.T) {
	tests := []struct {
		name string
		req  SecurityProtectionTemplateBindDomainSearchRequest
	}{
		{
			name: "Test GetSecurityProtectionTemplateBindDomains",
			req: SecurityProtectionTemplateBindDomainSearchRequest{
				BusinessID: 947,
				Page:       1,
				PageSize:   20,
			},
		},
	}
	if !isIntegrationTest() {
		t.Skip("Skipping integration test: set EDGENEXT_ACCESS_KEY and EDGENEXT_SECRET_KEY to run")
	}

	client := createTestClient(t)
	service := NewScdnService(client)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := service.GetSecurityProtectionTemplateBindDomains(tt.req)
			if err != nil {
				t.Errorf("ScdnService.GetSecurityProtectionTemplateBindDomains() error = %v", err)
				return
			}
			t.Logf("Response Status: Code=%d, Message=%s", got.Status.Code, got.Status.Message)
			t.Logf("Total Domains: %d", got.Data.Total)
			t.Logf("Domains Count: %d", len(got.Data.Domains))
			if len(got.Data.Domains) > 0 {
				t.Logf("First Domain: ID=%d, Domain=%s",
					got.Data.Domains[0].ID,
					got.Data.Domains[0].Domain)
			}
		})
	}
}

func TestScdnService_BindSecurityProtectionTemplateDomain(t *testing.T) {
	tests := []struct {
		name string
		req  SecurityProtectionTemplateBindDomainRequest
	}{
		{
			name: "Test BindSecurityProtectionTemplateDomain",
			req: SecurityProtectionTemplateBindDomainRequest{
				BusinessID: 947,
				DomainIDs:  []int{102008}, // Replace with actual domain IDs
			},
		},
	}
	if !isIntegrationTest() {
		t.Skip("Skipping integration test: set EDGENEXT_ACCESS_KEY and EDGENEXT_SECRET_KEY to run")
	}

	client := createTestClient(t)
	service := NewScdnService(client)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := service.BindSecurityProtectionTemplateDomain(tt.req)
			if err != nil {
				t.Errorf("ScdnService.BindSecurityProtectionTemplateDomain() error = %v", err)
				return
			}
			t.Logf("Response Status: Code=%d, Message=%s", got.Status.Code, got.Status.Message)
			if len(got.Data.FailDomains) > 0 {
				t.Logf("FailDomains: %+v", got.Data.FailDomains)
			}
		})
	}
}

func TestScdnService_DeleteSecurityProtectionTemplate(t *testing.T) {
	tests := []struct {
		name string
		req  SecurityProtectionTemplateDeleteRequest
	}{
		{
			name: "Test DeleteSecurityProtectionTemplate",
			req: SecurityProtectionTemplateDeleteRequest{
				BusinessID: 1265, // Replace with actual template ID for testing
			},
		},
	}
	if !isIntegrationTest() {
		t.Skip("Skipping integration test: set EDGENEXT_ACCESS_KEY and EDGENEXT_SECRET_KEY to run")
	}

	client := createTestClient(t)
	service := NewScdnService(client)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := service.DeleteSecurityProtectionTemplate(tt.req)
			if err != nil {
				if strings.Contains(err.Error(), "code:") {
					t.Logf("Delete failed (may not exist): %s", err.Error())
					return
				}
				t.Errorf("ScdnService.DeleteSecurityProtectionTemplate() error = %v", err)
				return
			}
			t.Logf("Response Status: Code=%d, Message=%s", got.Status.Code, got.Status.Message)
		})
	}
}

func TestScdnService_BatchConfigSecurityProtectionTemplate(t *testing.T) {
	tests := []struct {
		name string
		req  SecurityProtectionTemplateBatchConfigRequest
	}{
		{
			name: "Test BatchConfigSecurityProtectionTemplate",
			req: SecurityProtectionTemplateBatchConfigRequest{
				TemplateIDs: []int{1194}, // Replace with actual template IDs
				DdosConfig: &DdosProtectionGetConfigData{
					ApplicationDdosProtection: &ApplicationDdosProtection{
						Status:              "on",
						Type:                "strict",
						AIStatus:            "on",
						AICCStatus:          "on",
						NeedAttackDetection: 1,
					},
				},
				WafRuleConfig: &BatchUpdateWafRuleConfigRequest{
					WafRuleConfig: &WafRuleConfig{
						Status:   "on",
						WafLevel: "strict",
						WafMode:  "block",
						AIStatus: "on",
					},
				},
				BotManagementConfig: &UpdateBotManagementConfigRequest{
					BusinessID: 947,
					IDs:        []int{},
				},
				PreciseAccessControlConfig: &UpdatePreciseAccessControlConfigRequest{
					Action:   "cover",
					Policies: []PreciseAccessControlPolicy{},
				},
			},
		},
	}
	if !isIntegrationTest() {
		t.Skip("Skipping integration test: set EDGENEXT_ACCESS_KEY and EDGENEXT_SECRET_KEY to run")
	}

	client := createTestClient(t)
	service := NewScdnService(client)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := service.BatchConfigSecurityProtectionTemplate(tt.req)
			if err != nil {
				t.Errorf("ScdnService.BatchConfigSecurityProtectionTemplate() error = %v", err)
				return
			}
			t.Logf("Response Status: Code=%d, Message=%s", got.Status.Code, got.Status.Message)
			if len(got.Data.FailTemplates) > 0 {
				t.Logf("FailTemplates: %+v", got.Data.FailTemplates)
			}
		})
	}
}

func TestScdnService_GetSecurityProtectionTemplateUnboundDomains(t *testing.T) {
	tests := []struct {
		name string
		req  SecurityProtectionTemplateUnboundDomainSearchRequest
	}{
		{
			name: "Test GetSecurityProtectionTemplateUnboundDomains",
			req: SecurityProtectionTemplateUnboundDomainSearchRequest{
				Page:     1,
				PageSize: 20,
			},
		},
	}
	if !isIntegrationTest() {
		t.Skip("Skipping integration test: set EDGENEXT_ACCESS_KEY and EDGENEXT_SECRET_KEY to run")
	}

	client := createTestClient(t)
	service := NewScdnService(client)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := service.GetSecurityProtectionTemplateUnboundDomains(tt.req)
			if err != nil {
				t.Errorf("ScdnService.GetSecurityProtectionTemplateUnboundDomains() error = %v", err)
				return
			}
			t.Logf("Response Status: Code=%d, Message=%s", got.Status.Code, got.Status.Message)
			t.Logf("Total Unbound Domains: %d", got.Data.Total)
			t.Logf("Unbound Domains Count: %d", len(got.Data.Domains))
			if len(got.Data.Domains) > 0 {
				t.Logf("First Unbound Domain: ID=%d, Domain=%s",
					got.Data.Domains[0].ID,
					got.Data.Domains[0].Domain)
			}
		})
	}
}

func TestScdnService_EditSecurityProtectionTemplate(t *testing.T) {
	tests := []struct {
		name string
		req  SecurityProtectionTemplateEditRequest
	}{
		{
			name: "Test EditSecurityProtectionTemplate",
			req: SecurityProtectionTemplateEditRequest{
				BusinessID: 947,
				Name:       "updated-template-name",
				Remark:     "Updated remark",
			},
		},
	}
	if !isIntegrationTest() {
		t.Skip("Skipping integration test: set EDGENEXT_ACCESS_KEY and EDGENEXT_SECRET_KEY to run")
	}

	client := createTestClient(t)
	service := NewScdnService(client)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := service.EditSecurityProtectionTemplate(tt.req)
			if err != nil {
				t.Errorf("ScdnService.EditSecurityProtectionTemplate() error = %v", err)
				return
			}
			t.Logf("Response Status: Code=%d, Message=%s", got.Status.Code, got.Status.Message)
		})
	}
}

// ============================================================================
// Security Protection Iota Tests
// ============================================================================

func TestScdnService_GetSecurityProtectionIota(t *testing.T) {
	if !isIntegrationTest() {
		t.Skip("Skipping integration test: set EDGENEXT_ACCESS_KEY and EDGENEXT_SECRET_KEY to run")
	}

	client := createTestClient(t)
	service := NewScdnService(client)
	got, err := service.GetSecurityProtectionIota()
	if err != nil {
		t.Errorf("ScdnService.GetSecurityProtectionIota() error = %v", err)
		return
	}
	t.Logf("Response Status: Code=%d, Message=%s", got.Status.Code, got.Status.Message)
	t.Logf("Iota Enum Values: %+v", got.Data.Iota)
	for k, v := range got.Data.Iota {
		t.Logf("  %s: %s", k, v)
	}
}
