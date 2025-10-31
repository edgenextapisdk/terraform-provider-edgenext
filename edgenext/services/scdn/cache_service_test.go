package scdn

import (
	"strings"
	"testing"
)

func TestScdnService_GetCacheRules(t *testing.T) {
	tests := []struct {
		name string
		req  CacheRuleGetRulesRequest
	}{
		{
			name: "Test GetCacheRules",
			req: CacheRuleGetRulesRequest{
				BusinessID:   1246,
				BusinessType: "tpl",
				Page:         1,
				PageSize:     20,
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
			got, err := service.GetCacheRules(tt.req)
			if err != nil {
				t.Errorf("ScdnService.GetCacheRules() error = %v", err)
				return
			}
			t.Logf("Response: %+v", got)
		})
	}
}

func TestScdnService_CreateCacheRule(t *testing.T) {
	tests := []struct {
		name string
		req  CacheRuleCreateRequest
	}{
		{
			name: "Test CreateCacheRule",
			req: CacheRuleCreateRequest{
				BusinessID:   1246,
				BusinessType: "tpl",
				Name:         "test-cache-rule",
				Expr:         "(http.request.uri.path eq \"/test\")",
				Remark:       "Test cache rule",
				Conf: &CacheRuleConf{
					NoCache: false,
					CacheRule: &CacheRule{
						CacheTime: 60,
						Action:    "cachetime",
					},
					CacheShare: &CacheShare{
						Scheme: "http",
					},
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
			got, err := service.CreateCacheRule(tt.req)
			if err != nil {
				t.Errorf("ScdnService.CreateCacheRule() error = %v", err)
				return
			}
			t.Logf("Response: %+v", got)
		})
	}
}

func TestScdnService_UpdateCacheRule(t *testing.T) {
	tests := []struct {
		name string
		req  CacheRuleUpdateRequest
	}{
		{
			name: "Test UpdateCacheRule",
			req: CacheRuleUpdateRequest{
				ID:     5742, // Replace with actual rule ID for testing
				Name:   "updated-cache-rule-name",
				Remark: "Updated remark",
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
			got, err := service.UpdateCacheRule(tt.req)
			if err != nil {
				t.Errorf("ScdnService.UpdateCacheRule() error = %v", err)
				return
			}
			t.Logf("Response: %+v", got)
		})
	}
}

func TestScdnService_UpdateCacheRuleConfig(t *testing.T) {
	tests := []struct {
		name string
		req  CacheRuleUpdateConfigRequest
	}{
		{
			name: "Test UpdateCacheRuleConfig",
			req: CacheRuleUpdateConfigRequest{
				ID:           5742, // Replace with actual rule ID for testing
				BusinessID:   1246,
				Name:         "updated-cache-rule-config",
				BusinessType: "tpl",
				Conf: &CacheRuleConf{
					NoCache: false,
					CacheRule: &CacheRule{
						CacheTime: 120,
						Action:    "force",
					},
					BrowserCacheRule: &BrowserCacheRule{
						CacheTime:       60,
						IgnoreCacheTime: true,
						NoCache:         false,
					},
					CacheShare: &CacheShare{
						Scheme: "https",
					},
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
			got, err := service.UpdateCacheRuleConfig(tt.req)
			if err != nil {
				t.Errorf("ScdnService.UpdateCacheRuleConfig() error = %v", err)
				return
			}
			t.Logf("Response: %+v", got)
		})
	}
}

func TestScdnService_UpdateCacheRuleStatus(t *testing.T) {
	tests := []struct {
		name string
		req  CacheRuleUpdateStatusRequest
	}{
		{
			name: "Test UpdateCacheRuleStatus",
			req: CacheRuleUpdateStatusRequest{
				BusinessID:   1246,
				BusinessType: "tpl",
				IDs:          []int{5742}, // Replace with actual rule IDs for testing
				Status:       1,           // 1: enabled, 2: disabled
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
			got, err := service.UpdateCacheRuleStatus(tt.req)
			if err != nil {
				t.Errorf("ScdnService.UpdateCacheRuleStatus() error = %v", err)
				return
			}
			t.Logf("Response: %+v", got)
		})
	}
}

func TestScdnService_SortCacheRules(t *testing.T) {
	tests := []struct {
		name string
		req  CacheRuleSortRequest
	}{
		{
			name: "Test SortCacheRules",
			req: CacheRuleSortRequest{
				BusinessID:   1246,
				BusinessType: "tpl",
				IDs:          []int{5746, 5742}, // Replace with actual rule IDs for testing
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
			got, err := service.SortCacheRules(tt.req)
			if err != nil {
				t.Errorf("ScdnService.SortCacheRules() error = %v", err)
				return
			}
			t.Logf("Response: %+v", got)
		})
	}
}

func TestScdnService_DeleteCacheRule(t *testing.T) {
	tests := []struct {
		name string
		req  CacheRuleDeleteRequest
	}{
		{
			name: "Test DeleteCacheRule",
			req: CacheRuleDeleteRequest{
				BusinessID:   1246,
				BusinessType: "tpl",
				IDs:          []int{5749}, // Replace with actual rule ID for testing
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
			got, err := service.DeleteCacheRule(tt.req)
			if err != nil {
				if strings.Contains(err.Error(), "code: 103404") {
					t.Logf("error: %s", err.Error())
					return
				}
				t.Errorf("ScdnService.DeleteCacheRule() error = %v", err)
				return
			}
			t.Logf("Response: %+v", got)
		})
	}
}

func TestScdnService_GetCacheGlobalConfig(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "Test GetCacheGlobalConfig",
		},
	}
	if !isIntegrationTest() {
		t.Skip("Skipping integration test: set EDGENEXT_ACCESS_KEY and EDGENEXT_SECRET_KEY to run")
	}

	client := createTestClient(t)
	service := NewScdnService(client)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := service.GetCacheGlobalConfig()
			if err != nil {
				t.Errorf("ScdnService.GetCacheGlobalConfig() error = %v", err)
				return
			}
			t.Logf("Response: %+v", got)
		})
	}
}
