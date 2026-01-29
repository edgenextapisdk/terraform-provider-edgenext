package scdn

import (
	"strings"
	"testing"
)

func TestScdnService_CreateRuleTemplate(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		req  RuleTemplateCreateRequest
	}{
		{
			name: "Test CreateRuleTemplate",
			req: RuleTemplateCreateRequest{
				Name:        "test-template",
				Description: "Test template description",
				AppType:     "network_speed",
				BindDomain: &RuleTemplateBindDomain{
					AllDomain: false,
					DomainIDs: []int{101753},
					IsBind:    true,
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
			got, err := service.CreateRuleTemplate(tt.req)
			if err != nil {
				t.Errorf("ScdnService.CreateRuleTemplate() error = %v", err)
				return
			}
			t.Logf("Response: %+v", got)
		})
	}
}

func TestScdnService_UpdateRuleTemplate(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		req  RuleTemplateUpdateRequest
	}{
		{
			name: "Test UpdateRuleTemplate",
			req: RuleTemplateUpdateRequest{
				ID:          1246, // Replace with actual template ID for testing
				Name:        "updated-template-name",
				Description: "Updated template description",
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
			got, err := service.UpdateRuleTemplate(tt.req)
			if err != nil {
				t.Errorf("ScdnService.UpdateRuleTemplate() error = %v", err)
				return
			}
			t.Logf("Response: %+v", got)
		})
	}
}

func TestScdnService_DeleteRuleTemplate(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		req  RuleTemplateDeleteRequest
	}{
		{
			name: "Test DeleteRuleTemplate",
			req: RuleTemplateDeleteRequest{
				ID: 1246, // Replace with actual template ID for testing
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
			got, err := service.DeleteRuleTemplate(tt.req)
			if err != nil {
				if strings.Contains(err.Error(), "code: 103404") {
					t.Logf("error: %s", err.Error())
					return
				}
				t.Errorf("ScdnService.DeleteRuleTemplate() error = %v", err)
				return
			}
			t.Logf("Response: %+v", got)
		})
	}
}

func TestScdnService_ListRuleTemplates(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		req  RuleTemplateListRequest
	}{
		{
			name: "Test ListRuleTemplates",
			req: RuleTemplateListRequest{
				Page:     1,
				PageSize: 10,
				AppType:  "network_speed",
				Name:     "template-name",
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
			got, err := service.ListRuleTemplates(tt.req)
			if err != nil {
				t.Errorf("ScdnService.ListRuleTemplates() error = %v", err)
				return
			}
			t.Logf("Response: %+v", got)
		})
	}
}

func TestScdnService_BindRuleTemplateDomains(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		req  RuleTemplateBindDomainRequest
	}{
		{
			name: "Test BindRuleTemplateDomains",
			req: RuleTemplateBindDomainRequest{
				ID:        1246, // Replace with actual template ID for testing
				DomainIDs: []int{101753},
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
			got, err := service.BindRuleTemplateDomains(tt.req)
			if err != nil {
				t.Errorf("ScdnService.BindRuleTemplateDomains() error = %v", err)
				return
			}
			t.Logf("Response: %+v", got)
		})
	}
}

func TestScdnService_UnbindRuleTemplateDomains(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		req  RuleTemplateUnbindDomainRequest
	}{
		{
			name: "Test UnbindRuleTemplateDomains",
			req: RuleTemplateUnbindDomainRequest{
				ID:        1246, // Replace with actual template ID for testing
				DomainIDs: []int{101753},
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
			got, err := service.UnbindRuleTemplateDomains(tt.req)
			if err != nil {
				t.Errorf("ScdnService.UnbindRuleTemplateDomains() error = %v", err)
				return
			}
			t.Logf("Response: %+v", got)
		})
	}
}

func TestScdnService_ListRuleTemplateDomains(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		req  RuleTemplateListDomainsRequest
	}{
		{
			name: "Test ListRuleTemplateDomains",
			req: RuleTemplateListDomainsRequest{
				ID:       1246, // Replace with actual template ID for testing
				AppType:  "network_speed",
				Page:     1,
				PageSize: 10,
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
			got, err := service.ListRuleTemplateDomains(tt.req)
			if err != nil {
				t.Errorf("ScdnService.ListRuleTemplateDomains() error = %v", err)
				return
			}
			t.Logf("Response: %+v", got)
		})
	}
}
func TestScdnService_SwitchDomainTemplate(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		req  RuleTemplateSwitchDomainRequest
	}{
		{
			name: "Test SwitchDomainTemplate",
			req: RuleTemplateSwitchDomainRequest{
				AppType:    "network_speed",
				DomainIDs:  []int{115188},
				NewTplID:   1636,
				NewTplType: "more_domain",
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
			got, err := service.SwitchDomainTemplate(tt.req)
			if err != nil {
				t.Errorf("ScdnService.SwitchDomainTemplate() error = %v", err)
				return
			}
			t.Logf("Response: %+v", got)
		})
	}
}
