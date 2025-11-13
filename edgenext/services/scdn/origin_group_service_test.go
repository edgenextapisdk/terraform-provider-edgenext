package scdn

import (
	"strings"
	"testing"
)

// ============================================================================
// Origin Group Tests
// ============================================================================

func TestScdnService_ListOriginGroups(t *testing.T) {
	tests := []struct {
		name string
		req  OriginGroupListRequest
	}{
		{
			name: "Test ListOriginGroups",
			req: OriginGroupListRequest{
				Page:     1,
				PageSize: 20,
			},
		},
		{
			name: "Test ListOriginGroups with name filter",
			req: OriginGroupListRequest{
				Page:     1,
				PageSize: 20,
				Name:     "origin",
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
			got, err := service.ListOriginGroups(tt.req)
			if err != nil {
				t.Errorf("ScdnService.ListOriginGroups() error = %v", err)
				return
			}
			t.Logf("Response Status: Code=%d, Message=%s", got.Status.Code, got.Status.Message)
			t.Logf("Total Origin Groups: %d", got.Data.Total)
			t.Logf("Origin Groups Count: %d", len(got.Data.List))
			if len(got.Data.List) > 0 {
				t.Logf("First Origin Group: ID=%d, Name=%s",
					got.Data.List[0].ID,
					got.Data.List[0].Name)
			}
		})
	}
}

func TestScdnService_GetOriginGroupDetail(t *testing.T) {
	tests := []struct {
		name string
		req  OriginGroupDetailRequest
	}{
		{
			name: "Test GetOriginGroupDetail",
			req: OriginGroupDetailRequest{
				ID: 82, // Replace with actual origin group ID
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
			got, err := service.GetOriginGroupDetail(tt.req)
			if err != nil {
				if strings.Contains(err.Error(), "code:") {
					t.Logf("Get failed (may not exist): %s", err.Error())
					return
				}
				t.Errorf("ScdnService.GetOriginGroupDetail() error = %v", err)
				return
			}
			t.Logf("Response Status: Code=%d, Message=%s", got.Status.Code, got.Status.Message)
			if got.Data.OriginGroup.ID > 0 {
				t.Logf("Origin Group: ID=%d, Name=%s, Remark=%s",
					got.Data.OriginGroup.ID,
					got.Data.OriginGroup.Name,
					got.Data.OriginGroup.Remark)
			}
		})
	}
}

func TestScdnService_CreateOriginGroup(t *testing.T) {
	tests := []struct {
		name string
		req  OriginGroupCreateRequest
	}{
		{
			name: "Test CreateOriginGroup",
			req: OriginGroupCreateRequest{
				Name:   "test-origin",
				Remark: "Test origin group",
				Origins: []OriginGroupOrigin{
					{
						OriginType: 0, // IP
						Records: []OriginGroupRecord{
							{
								Value:    "34.236.82.201",
								Port:     80,
								Priority: 10,
								View:     "primary",
								Host:     "httpbin.org",
							},
						},
						ProtocolPorts: []OriginGroupProtocolPort{
							{
								Protocol:    0, // http
								ListenPorts: []int{80, 8080},
							},
						},
						OriginProtocol: 0, // http
						LoadBalance:    1, // round_robin
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
			got, err := service.CreateOriginGroup(tt.req)
			if err != nil {
				t.Errorf("ScdnService.CreateOriginGroup() error = %v", err)
				return
			}
			t.Logf("Response Status: Code=%d, Message=%s", got.Status.Code, got.Status.Message)
			t.Logf("Created Origin Group ID: %d", got.Data.ID)
		})
	}
}

func TestScdnService_UpdateOriginGroup(t *testing.T) {
	tests := []struct {
		name string
		req  OriginGroupUpdateRequest
	}{
		{
			name: "Test UpdateOriginGroup",
			req: OriginGroupUpdateRequest{
				ID:     82, // Replace with actual origin group ID
				Name:   "updated-origin",
				Remark: "Updated remark",
				Origins: []OriginGroupOrigin{
					{
						OriginType: 0, // IP
						Records: []OriginGroupRecord{
							{
								Value:    "3.211.176.227",
								Port:     80,
								Priority: 20,
								View:     "primary",
							},
						},
						ProtocolPorts: []OriginGroupProtocolPort{
							{
								Protocol:    0, // http
								ListenPorts: []int{80},
							},
						},
						OriginProtocol: 0, // http
						LoadBalance:    1, // round_robin
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
			got, err := service.UpdateOriginGroup(tt.req)
			if err != nil {
				if strings.Contains(err.Error(), "code:") {
					t.Logf("Update failed (may not exist): %s", err.Error())
					return
				}
				t.Errorf("ScdnService.UpdateOriginGroup() error = %v", err)
				return
			}
			t.Logf("Response Status: Code=%d, Message=%s", got.Status.Code, got.Status.Message)
		})
	}
}

func TestScdnService_DeleteOriginGroups(t *testing.T) {
	tests := []struct {
		name string
		req  OriginGroupDeleteRequest
	}{
		{
			name: "Test DeleteOriginGroups",
			req: OriginGroupDeleteRequest{
				IDs: []int{83}, // Replace with actual origin group ID for testing
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
			got, err := service.DeleteOriginGroups(tt.req)
			if err != nil {
				if strings.Contains(err.Error(), "code:") {
					t.Logf("Delete failed (may not exist): %s", err.Error())
					return
				}
				t.Errorf("ScdnService.DeleteOriginGroups() error = %v", err)
				return
			}
			t.Logf("Response Status: Code=%d, Message=%s", got.Status.Code, got.Status.Message)
		})
	}
}

func TestScdnService_BindOriginGroupToDomains(t *testing.T) {
	tests := []struct {
		name string
		req  OriginGroupBindDomainsRequest
	}{
		{
			name: "Test BindOriginGroupToDomains with domain_ids",
			req: OriginGroupBindDomainsRequest{
				OriginGroupID: 82,            // Replace with actual origin group ID
				DomainIDs:     []int{116040}, // Replace with actual domain IDs
			},
		},
		{
			name: "Test BindOriginGroupToDomains with domains",
			req: OriginGroupBindDomainsRequest{
				OriginGroupID: 81, // Replace with actual origin group ID
				Domains:       []string{"terraform.example.com"},
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
			got, err := service.BindOriginGroupToDomains(tt.req)
			if err != nil {
				if strings.Contains(err.Error(), "code:") {
					t.Logf("Bind failed (may not exist): %s", err.Error())
					return
				}
				t.Errorf("ScdnService.BindOriginGroupToDomains() error = %v", err)
				return
			}
			t.Logf("Response Status: Code=%d, Message=%s", got.Status.Code, got.Status.Message)
			t.Logf("Job ID: %s", got.Data.JobID)
		})
	}
}

func TestScdnService_GetAllOriginGroups(t *testing.T) {
	tests := []struct {
		name string
		req  OriginGroupAllRequest
	}{
		{
			name: "Test GetAllOriginGroups with scdn",
			req: OriginGroupAllRequest{
				ProtectStatus: "scdn",
			},
		},
		{
			name: "Test GetAllOriginGroups with exclusive",
			req: OriginGroupAllRequest{
				ProtectStatus: "exclusive",
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
			got, err := service.GetAllOriginGroups(tt.req)
			if err != nil {
				t.Errorf("ScdnService.GetAllOriginGroups() error = %v", err)
				return
			}
			t.Logf("Response Status: Code=%d, Message=%s", got.Status.Code, got.Status.Message)
			t.Logf("Total Origin Groups: %d", got.Data.Total)
			t.Logf("Origin Groups Count: %d", len(got.Data.List))
			if len(got.Data.List) > 0 {
				t.Logf("First Origin Group: ID=%d, Name=%s",
					got.Data.List[0].ID,
					got.Data.List[0].Name)
			}
		})
	}
}

func TestScdnService_CopyOriginGroupToDomain(t *testing.T) {
	tests := []struct {
		name string
		req  OriginGroupCopyRequest
	}{
		{
			name: "Test CopyOriginGroupToDomain",
			req: OriginGroupCopyRequest{
				OriginGroupID: 82,     // Replace with actual origin group ID
				DomainID:      102021, // Replace with actual domain ID
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
			got, err := service.CopyOriginGroupToDomain(tt.req)
			if err != nil {
				if strings.Contains(err.Error(), "code:") {
					t.Logf("Copy failed (may not exist): %s", err.Error())
					return
				}
				t.Errorf("ScdnService.CopyOriginGroupToDomain() error = %v", err)
				return
			}
			t.Logf("Response Status: Code=%d, Message=%s", got.Status.Code, got.Status.Message)
		})
	}
}

func TestScdnService_GetOriginGroupBindHistory(t *testing.T) {
	tests := []struct {
		name string
		req  OriginGroupBindHistoryRequest
	}{
		{
			name: "Test GetOriginGroupBindHistory",
			req: OriginGroupBindHistoryRequest{
				OriginGroupID: 81, // Replace with actual origin group ID
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
			got, err := service.GetOriginGroupBindHistory(tt.req)
			if err != nil {
				if strings.Contains(err.Error(), "code:") {
					t.Logf("Get history failed (may not exist): %s", err.Error())
					return
				}
				t.Errorf("ScdnService.GetOriginGroupBindHistory() error = %v", err)
				return
			}
			t.Logf("Response Status: Code=%d, Message=%s", got.Status.Code, got.Status.Message)
			if got.Data.History.ID > 0 {
				t.Logf("History: ID=%d, OriginGroupID=%d, DomainCount=%d",
					got.Data.History.ID,
					got.Data.History.OriginGroupID,
					len(got.Data.History.Domains))
			}
		})
	}
}
