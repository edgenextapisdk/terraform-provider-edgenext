package scdn

import (
	"testing"
)

// ============================================================================
// Task Management Tests
// ============================================================================

func TestScdnService_ListLogDownloadTasks(t *testing.T) {
	tests := []struct {
		name string
		req  LogDownloadTaskListRequest
	}{
		{
			name: "Test ListLogDownloadTasks",
			req: LogDownloadTaskListRequest{
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
			got, err := service.ListLogDownloadTasks(tt.req)
			if err != nil {
				t.Errorf("ScdnService.ListLogDownloadTasks() error = %v", err)
				return
			}
			t.Logf("Response: %+v", got)
		})
	}
}

func TestScdnService_AddLogDownloadTask(t *testing.T) {
	tests := []struct {
		name string
		req  LogDownloadTaskAddRequest
	}{
		{
			name: "Test AddLogDownloadTask",
			req: LogDownloadTaskAddRequest{
				TaskName:       "test-task1",
				IsUseTemplate:  0,
				DataSource:     "ng",
				DownloadFields: []string{"http_host", "city", "country"},
				SearchTerms: map[string]string{
					"http_host": "terraform.example.com",
				},
				FileType:  "xls",
				StartTime: "2025-11-01 00:00:00",
				EndTime:   "2025-11-02 00:00:00",
				Lang:      "zh_CN",
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
			got, err := service.AddLogDownloadTask(tt.req)
			if err != nil {
				t.Errorf("ScdnService.AddLogDownloadTask() error = %v", err)
				return
			}
			t.Logf("Response: %+v", got)
		})
	}
}

func TestScdnService_CancelLogDownloadTask(t *testing.T) {
	tests := []struct {
		name string
		req  LogDownloadTaskCancelRequest
	}{
		{
			name: "Test CancelLogDownloadTask",
			req: LogDownloadTaskCancelRequest{
				TaskID: 109, // Note: Set a valid task_id for actual testing
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
			// Skip if task_id is 0 (not set)
			if tt.req.TaskID == 0 {
				t.Skip("Skipping test: task_id not set (set a valid task_id to test)")
			}
			got, err := service.CancelLogDownloadTask(tt.req)
			if err != nil {
				t.Errorf("ScdnService.CancelLogDownloadTask() error = %v", err)
				return
			}
			t.Logf("Response: %+v", got)
		})
	}
}

func TestScdnService_BatchCancelLogDownloadTasks(t *testing.T) {
	tests := []struct {
		name string
		req  LogDownloadTaskBatchCancelRequest
	}{
		{
			name: "Test BatchCancelLogDownloadTasks",
			req: LogDownloadTaskBatchCancelRequest{
				TaskIDs: []int{109, 110}, // Note: Set valid task_ids for actual testing
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
			// Skip if task_ids is empty
			if len(tt.req.TaskIDs) == 0 {
				t.Skip("Skipping test: task_ids not set (set valid task_ids to test)")
			}
			got, err := service.BatchCancelLogDownloadTasks(tt.req)
			if err != nil {
				t.Errorf("ScdnService.BatchCancelLogDownloadTasks() error = %v", err)
				return
			}
			t.Logf("Response: %+v", got)
		})
	}
}

func TestScdnService_RegenerateLogDownloadTask(t *testing.T) {
	tests := []struct {
		name string
		req  LogDownloadTaskRegenerateRequest
	}{
		{
			name: "Test RegenerateLogDownloadTask",
			req: LogDownloadTaskRegenerateRequest{
				TaskID: 117, // Note: Set a valid task_id for actual testing
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
			// Skip if task_id is 0 (not set)
			if tt.req.TaskID == 0 {
				t.Skip("Skipping test: task_id not set (set a valid task_id to test)")
			}
			got, err := service.RegenerateLogDownloadTask(tt.req)
			if err != nil {
				t.Errorf("ScdnService.RegenerateLogDownloadTask() error = %v", err)
				return
			}
			t.Logf("Response: %+v", got)
		})
	}
}

func TestScdnService_DeleteLogDownloadTask(t *testing.T) {
	tests := []struct {
		name string
		req  LogDownloadTaskDeleteRequest
	}{
		{
			name: "Test DeleteLogDownloadTask",
			req: LogDownloadTaskDeleteRequest{
				TaskID: 109, // Note: Set a valid task_id for actual testing
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
			// Skip if task_id is 0 (not set)
			if tt.req.TaskID == 0 {
				t.Skip("Skipping test: task_id not set (set a valid task_id to test)")
			}
			got, err := service.DeleteLogDownloadTask(tt.req)
			if err != nil {
				t.Errorf("ScdnService.DeleteLogDownloadTask() error = %v", err)
				return
			}
			t.Logf("Response: %+v", got)
		})
	}
}

func TestScdnService_BatchDeleteLogDownloadTasks(t *testing.T) {
	tests := []struct {
		name string
		req  LogDownloadTaskBatchDeleteRequest
	}{
		{
			name: "Test BatchDeleteLogDownloadTasks",
			req: LogDownloadTaskBatchDeleteRequest{
				TaskIDs: []int{53}, // Note: Set valid task_ids for actual testing
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
			// Skip if task_ids is empty
			if len(tt.req.TaskIDs) == 0 {
				t.Skip("Skipping test: task_ids not set (set valid task_ids to test)")
			}
			got, err := service.BatchDeleteLogDownloadTasks(tt.req)
			if err != nil {
				t.Errorf("ScdnService.BatchDeleteLogDownloadTasks() error = %v", err)
				return
			}
			t.Logf("Response: %+v", got)
		})
	}
}

// ============================================================================
// Template Management Tests
// ============================================================================

func TestScdnService_ListLogDownloadTemplates(t *testing.T) {
	tests := []struct {
		name string
		req  LogDownloadTemplateListRequest
	}{
		{
			name: "Test ListLogDownloadTemplates",
			req: LogDownloadTemplateListRequest{
				Page:    1,
				PerPage: 20,
				Status:  1,
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
			got, err := service.ListLogDownloadTemplates(tt.req)
			if err != nil {
				t.Errorf("ScdnService.ListLogDownloadTemplates() error = %v", err)
				return
			}
			t.Logf("Response: %+v", got)
		})
	}
}

func TestScdnService_ListLogDownloadTemplateDomains(t *testing.T) {
	tests := []struct {
		name string
		req  LogDownloadTemplateDomainListRequest
	}{
		{
			name: "Test ListLogDownloadTemplateDomains without domain",
			req:  LogDownloadTemplateDomainListRequest{},
		},
		{
			name: "Test ListLogDownloadTemplateDomains with domain",
			req: LogDownloadTemplateDomainListRequest{
				Domain: "terraform.example.com",
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
			got, err := service.ListLogDownloadTemplateDomains(tt.req)
			if err != nil {
				t.Errorf("ScdnService.ListLogDownloadTemplateDomains() error = %v", err)
				return
			}
			t.Logf("Response: %+v", got)
		})
	}
}

func TestScdnService_GetAllLogDownloadTemplates(t *testing.T) {
	if !isIntegrationTest() {
		t.Skip("Skipping integration test: set EDGENEXT_ACCESS_KEY and EDGENEXT_SECRET_KEY to run")
	}

	client := createTestClient(t)
	service := NewScdnService(client)
	got, err := service.GetAllLogDownloadTemplates()
	if err != nil {
		t.Errorf("ScdnService.GetAllLogDownloadTemplates() error = %v", err)
		return
	}
	t.Logf("Response: %+v", got)
}

func TestScdnService_GetAllLogDownloadTemplateGroups(t *testing.T) {
	if !isIntegrationTest() {
		t.Skip("Skipping integration test: set EDGENEXT_ACCESS_KEY and EDGENEXT_SECRET_KEY to run")
	}

	client := createTestClient(t)
	service := NewScdnService(client)
	got, err := service.GetAllLogDownloadTemplateGroups()
	if err != nil {
		t.Errorf("ScdnService.GetAllLogDownloadTemplateGroups() error = %v", err)
		return
	}
	t.Logf("Response: %+v", got)
}

func TestScdnService_AddLogDownloadTemplate(t *testing.T) {
	tests := []struct {
		name string
		req  LogDownloadTemplateAddRequest
	}{
		{
			name: "Test AddLogDownloadTemplate",
			req: LogDownloadTemplateAddRequest{
				TemplateName:   "test-template",
				GroupName:      "test-group",
				DataSource:     "ng",
				Status:         1,
				DownloadFields: []string{"http_host", "city", "country"},
				SearchTerms: map[string]string{
					"http_host": "example.com",
				},
				DomainSelectType: 0,
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
			got, err := service.AddLogDownloadTemplate(tt.req)
			if err != nil {
				t.Errorf("ScdnService.AddLogDownloadTemplate() error = %v", err)
				return
			}
			t.Logf("Response: %+v", got)
		})
	}
}

func TestScdnService_SaveLogDownloadTemplate(t *testing.T) {
	tests := []struct {
		name string
		req  LogDownloadTemplateSaveRequest
	}{
		{
			name: "Test SaveLogDownloadTemplate",
			req: LogDownloadTemplateSaveRequest{
				TemplateID:     26, // Note: Set a valid template_id for actual testing
				TemplateName:   "test-template-updated",
				GroupName:      "test-group",
				GroupID:        0, // Note: Set a valid group_id for actual testing
				DataSource:     "ng",
				Status:         1,
				DownloadFields: []string{"http_host", "city", "country", "province"},
				SearchTerms: map[string]string{
					"http_host": "terraform.example.com",
				},
				DomainSelectType: 0,
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
			// Skip if template_id is 0 (not set)
			if tt.req.TemplateID == 0 {
				t.Skip("Skipping test: template_id not set (set a valid template_id to test)")
			}
			got, err := service.SaveLogDownloadTemplate(tt.req)
			if err != nil {
				t.Errorf("ScdnService.SaveLogDownloadTemplate() error = %v", err)
				return
			}
			t.Logf("Response: %+v", got)
		})
	}
}

func TestScdnService_ChangeLogDownloadTemplateStatus(t *testing.T) {
	tests := []struct {
		name string
		req  LogDownloadTemplateChangeStatusRequest
	}{
		{
			name: "Test ChangeLogDownloadTemplateStatus enable",
			req: LogDownloadTemplateChangeStatusRequest{
				TemplateID: 26, // Note: Set a valid template_id for actual testing
				Status:     1,  // Enable
			},
		},
		{
			name: "Test ChangeLogDownloadTemplateStatus disable",
			req: LogDownloadTemplateChangeStatusRequest{
				TemplateID: 26, // Note: Set a valid template_id for actual testing
				Status:     0,  // Disable
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
			// Skip if template_id is 0 (not set)
			if tt.req.TemplateID == 0 {
				t.Skip("Skipping test: template_id not set (set a valid template_id to test)")
			}
			got, err := service.ChangeLogDownloadTemplateStatus(tt.req)
			if err != nil {
				t.Errorf("ScdnService.ChangeLogDownloadTemplateStatus() error = %v", err)
				return
			}
			t.Logf("Response: %+v", got)
		})
	}
}

func TestScdnService_BatchChangeLogDownloadTemplateStatus(t *testing.T) {
	tests := []struct {
		name string
		req  LogDownloadTemplateBatchChangeStatusRequest
	}{
		{
			name: "Test BatchChangeLogDownloadTemplateStatus enable",
			req: LogDownloadTemplateBatchChangeStatusRequest{
				TemplateIDs: []int{26}, // Note: Set valid template_ids for actual testing
				Status:      1,         // Enable
			},
		},
		{
			name: "Test BatchChangeLogDownloadTemplateStatus disable",
			req: LogDownloadTemplateBatchChangeStatusRequest{
				TemplateIDs: []int{26}, // Note: Set valid template_ids for actual testing
				Status:      0,         // Disable
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
			// Skip if template_ids is empty
			if len(tt.req.TemplateIDs) == 0 {
				t.Skip("Skipping test: template_ids not set (set valid template_ids to test)")
			}
			got, err := service.BatchChangeLogDownloadTemplateStatus(tt.req)
			if err != nil {
				t.Errorf("ScdnService.BatchChangeLogDownloadTemplateStatus() error = %v", err)
				return
			}
			t.Logf("Response: %+v", got)
		})
	}
}

func TestScdnService_DeleteLogDownloadTemplate(t *testing.T) {
	tests := []struct {
		name string
		req  LogDownloadTemplateDeleteRequest
	}{
		{
			name: "Test DeleteLogDownloadTemplate",
			req: LogDownloadTemplateDeleteRequest{
				TemplateID: 27, // Note: Set a valid template_id for actual testing
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
			// Skip if template_id is 0 (not set)
			if tt.req.TemplateID == 0 {
				t.Skip("Skipping test: template_id not set (set a valid template_id to test)")
			}
			got, err := service.DeleteLogDownloadTemplate(tt.req)
			if err != nil {
				t.Errorf("ScdnService.DeleteLogDownloadTemplate() error = %v", err)
				return
			}
			t.Logf("Response: %+v", got)
		})
	}
}

func TestScdnService_BatchDeleteLogDownloadTemplates(t *testing.T) {
	tests := []struct {
		name string
		req  LogDownloadTemplateBatchDeleteRequest
	}{
		{
			name: "Test BatchDeleteLogDownloadTemplates",
			req: LogDownloadTemplateBatchDeleteRequest{
				TemplateIDs: []int{28}, // Note: Set valid template_ids for actual testing
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
			// Skip if template_ids is empty
			if len(tt.req.TemplateIDs) == 0 {
				t.Skip("Skipping test: template_ids not set (set valid template_ids to test)")
			}
			got, err := service.BatchDeleteLogDownloadTemplates(tt.req)
			if err != nil {
				t.Errorf("ScdnService.BatchDeleteLogDownloadTemplates() error = %v", err)
				return
			}
			t.Logf("Response: %+v", got)
		})
	}
}

// ============================================================================
// Fields and Configuration Tests
// ============================================================================

func TestScdnService_GetLogDownloadFields(t *testing.T) {
	if !isIntegrationTest() {
		t.Skip("Skipping integration test: set EDGENEXT_ACCESS_KEY and EDGENEXT_SECRET_KEY to run")
	}

	client := createTestClient(t)
	service := NewScdnService(client)
	got, err := service.GetLogDownloadFields()
	if err != nil {
		t.Errorf("ScdnService.GetLogDownloadFields() error = %v", err)
		return
	}
	t.Logf("Response: %+v", got)
}
