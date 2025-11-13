package scdn

import (
	"strings"
	"testing"
)

func TestScdnService_GetNetworkSpeedConfig(t *testing.T) {
	tests := []struct {
		name string
		req  NetworkSpeedGetConfigRequest
	}{
		{
			name: "Test GetNetworkSpeedConfig",
			req: NetworkSpeedGetConfigRequest{
				BusinessID:   1246,
				BusinessType: "tpl",
				ConfigGroups: []string{
					"domain_proxy_conf",
					"upstream_redirect",
					"https",
					"page_gzip",
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
			got, err := service.GetNetworkSpeedConfig(tt.req)
			if err != nil {
				t.Errorf("ScdnService.GetNetworkSpeedConfig() error = %v", err)
				return
			}
			t.Logf("Response: %+v", got)
		})
	}
}

func TestScdnService_UpdateNetworkSpeedConfig(t *testing.T) {
	tests := []struct {
		name string
		req  NetworkSpeedUpdateConfigRequest
	}{
		{
			name: "Test UpdateNetworkSpeedConfig",
			req: NetworkSpeedUpdateConfigRequest{
				BusinessID:   1246,
				BusinessType: "tpl",
				UpstreamURIChange: &UpstreamURIChange{
					Status: "on",
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
			got, err := service.UpdateNetworkSpeedConfig(tt.req)
			if err != nil {
				t.Errorf("ScdnService.UpdateNetworkSpeedConfig() error = %v", err)
				return
			}
			t.Logf("Response: %+v", got)
		})
	}
}

func TestScdnService_GetNetworkSpeedRules(t *testing.T) {
	tests := []struct {
		name string
		req  NetworkSpeedGetRulesRequest
	}{
		{
			name: "Test GetNetworkSpeedRules",
			req: NetworkSpeedGetRulesRequest{
				BusinessID:   1246,
				BusinessType: "tpl",
				ConfigGroup:  "upstream_uri_change_rule",
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
			got, err := service.GetNetworkSpeedRules(tt.req)
			if err != nil {
				t.Errorf("ScdnService.GetNetworkSpeedRules() error = %v", err)
				return
			}
			t.Logf("Response: %+v", got)
		})
	}
}

func TestScdnService_CreateNetworkSpeedRule(t *testing.T) {
	tests := []struct {
		name string
		req  NetworkSpeedCreateRuleRequest
	}{
		{
			name: "Test CreateNetworkSpeedRule - CustomizedReqHeadersRule",
			req: NetworkSpeedCreateRuleRequest{
				BusinessID:   1246,
				BusinessType: "tpl",
				ConfigGroup:  "customized_req_headers_rule",
				CustomizedReqHeadersRule: &CustomizedReqHeadersRule{
					Type:    "User-Agent",
					Content: "test-content",
					Remark:  "test-remark",
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
			got, err := service.CreateNetworkSpeedRule(tt.req)
			if err != nil {
				t.Errorf("ScdnService.CreateNetworkSpeedRule() error = %v", err)
				return
			}
			t.Logf("Response: %+v", got)
		})
	}
}

func TestScdnService_DeleteNetworkSpeedRule(t *testing.T) {
	tests := []struct {
		name string
		req  NetworkSpeedDeleteRuleRequest
	}{
		{
			name: "Test DeleteNetworkSpeedRule",
			req: NetworkSpeedDeleteRuleRequest{
				BusinessID:   1246,
				BusinessType: "tpl",
				ConfigGroup:  "customized_req_headers_rule",
				IDs:          []int{601}, // Replace with actual rule ID for testing
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
			got, err := service.DeleteNetworkSpeedRule(tt.req)
			if err != nil {
				if strings.Contains(err.Error(), "code: 103404") {
					t.Logf("Rule not found (expected for test): %s", err.Error())
					return
				}
				t.Errorf("ScdnService.DeleteNetworkSpeedRule() error = %v", err)
				return
			}
			t.Logf("Response: %+v", got)
		})
	}
}

func TestScdnService_SortNetworkSpeedRules(t *testing.T) {
	tests := []struct {
		name string
		req  NetworkSpeedSortRulesRequest
	}{
		{
			name: "Test SortNetworkSpeedRules",
			req: NetworkSpeedSortRulesRequest{
				BusinessID:   1246,
				BusinessType: "tpl",
				ConfigGroup:  "customized_req_headers_rule",
				IDs:          []int{603, 602}, // Replace with actual rule IDs for testing
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
			got, err := service.SortNetworkSpeedRules(tt.req)
			if err != nil {
				t.Errorf("ScdnService.SortNetworkSpeedRules() error = %v", err)
				return
			}
			t.Logf("Response: %+v", got)
		})
	}
}

func TestScdnService_UpdateNetworkSpeedRule(t *testing.T) {
	tests := []struct {
		name string
		req  NetworkSpeedUpdateRuleRequest
	}{
		{
			name: "Test UpdateNetworkSpeedRule - CustomizedReqHeadersRule",
			req: NetworkSpeedUpdateRuleRequest{
				ID:          603, // Replace with actual rule ID for testing
				ConfigGroup: "customized_req_headers_rule",
				CustomizedReqHeadersRule: &CustomizedReqHeadersRule{
					Type:    "User-Agent-Updated",
					Content: "test-content-updated",
					Remark:  "test-remark-updated",
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
			got, err := service.UpdateNetworkSpeedRule(tt.req)
			if err != nil {
				if strings.Contains(err.Error(), "code: 103404") {
					t.Logf("Rule not found (expected for test): %s", err.Error())
					return
				}
				t.Errorf("ScdnService.UpdateNetworkSpeedRule() error = %v", err)
				return
			}
			t.Logf("Response: %+v", got)
		})
	}
}
