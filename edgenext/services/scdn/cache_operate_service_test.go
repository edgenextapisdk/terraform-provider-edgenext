package scdn

import (
	"testing"
)

func TestScdnService_GetCacheCleanConfig(t *testing.T) {
	tests := []struct {
		name string
		req  CacheCleanGetConfigRequest
	}{
		{
			name: "Test GetCacheCleanConfig",
			req:  CacheCleanGetConfigRequest{},
		},
	}
	if !isIntegrationTest() {
		t.Skip("Skipping integration test: set EDGENEXT_ACCESS_KEY and EDGENEXT_SECRET_KEY to run")
	}

	client := createTestClient(t)
	service := NewScdnService(client)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := service.GetCacheCleanConfig(tt.req)
			if err != nil {
				t.Errorf("ScdnService.GetCacheCleanConfig() error = %v", err)
				return
			}
			t.Logf("Response: %+v", got)
		})
	}
}

func TestScdnService_SaveCacheCleanTask(t *testing.T) {
	tests := []struct {
		name string
		req  CacheCleanSaveRequest
	}{
		{
			name: "Test SaveCacheCleanTask - Whole Site",
			req: CacheCleanSaveRequest{
				Wholesite: []string{"terraform.example.com", "terraform.example.com"},
			},
		},
		{
			name: "Test SaveCacheCleanTask - Special URL",
			req: CacheCleanSaveRequest{
				Specialurl: []string{"http://terraform.example.com/a", "http://terraform.example.com/a1"},
			},
		},
		{
			name: "Test SaveCacheCleanTask - Special Directory",
			req: CacheCleanSaveRequest{
				Specialdir: []string{"http://terraform.example.com/a/", "http://terraform.example.com/a1/a2/"},
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
			got, err := service.SaveCacheCleanTask(tt.req)
			if err != nil {
				t.Errorf("ScdnService.SaveCacheCleanTask() error = %v", err)
				return
			}
			t.Logf("Response: %+v", got)
		})
	}
}

func TestScdnService_GetCacheCleanTaskList(t *testing.T) {
	tests := []struct {
		name string
		req  CacheCleanTaskListRequest
	}{
		{
			name: "Test GetCacheCleanTaskList",
			req: CacheCleanTaskListRequest{
				Page:      1,
				PerPage:   20,
				StartTime: "2025-11-01 00:00:00",
				EndTime:   "2025-11-10 23:59:59",
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
			got, err := service.GetCacheCleanTaskList(tt.req)
			if err != nil {
				t.Errorf("ScdnService.GetCacheCleanTaskList() error = %v", err)
				return
			}
			t.Logf("Response: %+v", got)
		})
	}
}

func TestScdnService_GetCacheCleanTaskDetail(t *testing.T) {
	tests := []struct {
		name string
		req  CacheCleanTaskDetailRequest
	}{
		{
			name: "Test GetCacheCleanTaskDetail",
			req: CacheCleanTaskDetailRequest{
				TaskID:  2247,
				Page:    1,
				PerPage: 20,
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
			got, err := service.GetCacheCleanTaskDetail(tt.req)
			if err != nil {
				t.Errorf("ScdnService.GetCacheCleanTaskDetail() error = %v", err)
				return
			}
			t.Logf("Response: %+v", got)
		})
	}
}

func TestScdnService_GetCachePreheatTaskList(t *testing.T) {
	tests := []struct {
		name string
		req  CachePreheatTaskListRequest
	}{
		{
			name: "Test GetCachePreheatTaskList",
			req: CachePreheatTaskListRequest{
				Page:    1,
				PerPage: 10,
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
			got, err := service.GetCachePreheatTaskList(tt.req)
			if err != nil {
				t.Errorf("ScdnService.GetCachePreheatTaskList() error = %v", err)
				return
			}
			t.Logf("Response: %+v", got)
		})
	}
}

func TestScdnService_SaveCachePreheatTask(t *testing.T) {
	tests := []struct {
		name string
		req  CachePreheatSaveRequest
	}{
		{
			name: "Test SaveCachePreheatTask",
			req: CachePreheatSaveRequest{
				PreheatURL: []string{"http://terraform.example.com/a.jpg"},
			},
		},
		{
			name: "Test SaveCachePreheatTask - Multiple URLs",
			req: CachePreheatSaveRequest{
				PreheatURL: []string{
					"http://terraform.example.com/a.jpg",
					"http://terraform.example.com/b.jpg",
					"http://terraform.example.com/c.jpg",
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
			got, err := service.SaveCachePreheatTask(tt.req)
			if err != nil {
				t.Errorf("ScdnService.SaveCachePreheatTask() error = %v", err)
				return
			}
			t.Logf("Response: %+v", got)
		})
	}
}
