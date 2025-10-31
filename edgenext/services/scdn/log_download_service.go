package scdn

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strconv"

	"github.com/edgenextapisdk/terraform-provider-edgenext/edgenext/connectivity"
)

// ============================================================================
// Log Download Management Methods
// ============================================================================

// ListLogDownloadTasks lists log download tasks
func (s *ScdnService) ListLogDownloadTasks(req LogDownloadTaskListRequest) (*LogDownloadTaskListResponse, error) {
	ctx := context.Background()

	// Get SCDN client from EdgeNextClient
	scdnClient, err := s.client.ScdnClient()
	if err != nil {
		return nil, fmt.Errorf("failed to get SCDN client: %w", err)
	}

	// Convert request to ScdnRequest format with query parameters
	scdnReq := &connectivity.ScdnRequest{
		Query: make(map[string]interface{}),
	}

	// Set default values for page and per_page if not specified
	page := req.Page
	if page <= 0 {
		page = 1
	}
	scdnReq.Query["page"] = page

	perPage := req.PerPage
	if perPage <= 0 {
		perPage = 20
	}
	scdnReq.Query["per_page"] = perPage

	// Optional parameters - only add if they have meaningful values
	// Status: -1 means don't filter by status, >= 0 means filter by specific status
	if req.Status >= 0 {
		scdnReq.Query["status"] = req.Status
	}
	// Note: If Status is -1 or not set, we don't add status parameter to query all statuses
	if req.TaskName != "" {
		scdnReq.Query["task_name"] = req.TaskName
	}
	if req.FileType != "" {
		scdnReq.Query["file_type"] = req.FileType
	}
	if req.DataSource != "" {
		scdnReq.Query["data_source"] = req.DataSource
	}

	// Call SCDN API (supports both GET and POST)
	scdnResp, err := scdnClient.Get(ctx, EndpointLogDownloadTaskList, scdnReq)
	if err != nil {
		return nil, fmt.Errorf("failed to list log download tasks: %w", err)
	}

	// Check business status code
	if scdnResp.Status.Code != 1 {
		return nil, fmt.Errorf("API error: %s (code: %d)", scdnResp.Status.Message, scdnResp.Status.Code)
	}

	// Convert response
	response := &LogDownloadTaskListResponse{
		Status: Status{
			Code:    scdnResp.Status.Code,
			Message: scdnResp.Status.Message,
		},
	}

	if scdnResp.Data != nil {
		dataBytes, err := json.Marshal(scdnResp.Data)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal response data: %w", err)
		}

		var taskListData LogDownloadTaskListData
		if err := json.Unmarshal(dataBytes, &taskListData); err != nil {
			return nil, fmt.Errorf("failed to unmarshal task list data: %w", err)
		}
		response.Data = taskListData
	}

	return response, nil
}

// AddLogDownloadTask adds a log download task
func (s *ScdnService) AddLogDownloadTask(req LogDownloadTaskAddRequest) (*LogDownloadTaskAddResponse, error) {
	ctx := context.Background()

	scdnClient, err := s.client.ScdnClient()
	if err != nil {
		return nil, fmt.Errorf("failed to get SCDN client: %w", err)
	}

	// Convert request to ScdnRequest format
	scdnReq := &connectivity.ScdnRequest{}
	if req.TaskName != "" {
		scdnReq.Data = make(map[string]interface{})
		reqBytes, err := json.Marshal(req)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request: %w", err)
		}

		if err := json.Unmarshal(reqBytes, &scdnReq.Data); err != nil {
			return nil, fmt.Errorf("failed to unmarshal request data: %w", err)
		}
	}

	scdnResp, err := scdnClient.Post(ctx, EndpointLogDownloadTaskAdd, scdnReq)
	if err != nil {
		return nil, fmt.Errorf("failed to call SCDN API: %w", err)
	}

	if scdnResp.Status.Code != 1 {
		return nil, fmt.Errorf("API error: %s (code: %d)", scdnResp.Status.Message, scdnResp.Status.Code)
	}

	// Convert response
	response := &LogDownloadTaskAddResponse{
		Status: Status{
			Code:    scdnResp.Status.Code,
			Message: scdnResp.Status.Message,
		},
	}

	if scdnResp.Data != nil {
		dataBytes, err := json.Marshal(scdnResp.Data)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal response data: %w", err)
		}

		var taskAddData LogDownloadTaskAddData
		if err := json.Unmarshal(dataBytes, &taskAddData); err != nil {
			return nil, fmt.Errorf("failed to unmarshal task add data: %w", err)
		}
		response.Data = taskAddData
	}

	return response, nil
}

// CancelLogDownloadTask cancels a log download task
func (s *ScdnService) CancelLogDownloadTask(req LogDownloadTaskCancelRequest) (*LogDownloadTaskCancelResponse, error) {
	ctx := context.Background()

	var response LogDownloadTaskCancelResponse
	err := s.callSCDNAPI(ctx, MethodPOST, EndpointLogDownloadTaskCancel, req, &response)
	if err != nil {
		return nil, fmt.Errorf("failed to cancel log download task: %w", err)
	}

	return &response, nil
}

// BatchCancelLogDownloadTasks batch cancels log download tasks
func (s *ScdnService) BatchCancelLogDownloadTasks(req LogDownloadTaskBatchCancelRequest) (*LogDownloadTaskBatchCancelResponse, error) {
	ctx := context.Background()

	var response LogDownloadTaskBatchCancelResponse
	err := s.callSCDNAPI(ctx, MethodDELETE, EndpointLogDownloadTaskBatchCancel, req, &response)
	if err != nil {
		return nil, fmt.Errorf("failed to batch cancel log download tasks: %w", err)
	}

	return &response, nil
}

// DeleteLogDownloadTask deletes a log download task
func (s *ScdnService) DeleteLogDownloadTask(req LogDownloadTaskDeleteRequest) (*LogDownloadTaskDeleteResponse, error) {
	ctx := context.Background()

	var response LogDownloadTaskDeleteResponse
	err := s.callSCDNAPI(ctx, MethodDELETE, EndpointLogDownloadTaskDelete, req, &response)
	if err != nil {
		return nil, fmt.Errorf("failed to delete log download task: %w", err)
	}

	return &response, nil
}

// BatchDeleteLogDownloadTasks batch deletes log download tasks
func (s *ScdnService) BatchDeleteLogDownloadTasks(req LogDownloadTaskBatchDeleteRequest) (*LogDownloadTaskBatchDeleteResponse, error) {
	ctx := context.Background()

	var response LogDownloadTaskBatchDeleteResponse
	err := s.callSCDNAPI(ctx, MethodDELETE, EndpointLogDownloadTaskBatchDelete, req, &response)
	if err != nil {
		return nil, fmt.Errorf("failed to batch delete log download tasks: %w", err)
	}

	return &response, nil
}

// RegenerateLogDownloadTask regenerates a log download task
func (s *ScdnService) RegenerateLogDownloadTask(req LogDownloadTaskRegenerateRequest) (*LogDownloadTaskRegenerateResponse, error) {
	ctx := context.Background()

	var response LogDownloadTaskRegenerateResponse
	err := s.callSCDNAPI(ctx, MethodPOST, EndpointLogDownloadTaskRegenerate, req, &response)
	if err != nil {
		return nil, fmt.Errorf("failed to regenerate log download task: %w", err)
	}

	return &response, nil
}

// GetLogDownloadFields gets log download fields
func (s *ScdnService) GetLogDownloadFields() (*LogDownloadFieldsResponse, error) {
	ctx := context.Background()

	// Get SCDN client from EdgeNextClient
	scdnClient, err := s.client.ScdnClient()
	if err != nil {
		return nil, fmt.Errorf("failed to get SCDN client: %w", err)
	}

	// Call SCDN API
	scdnReq := &connectivity.ScdnRequest{}
	scdnResp, err := scdnClient.Get(ctx, EndpointLogDownloadFields, scdnReq)
	if err != nil {
		return nil, fmt.Errorf("failed to get log download fields: %w", err)
	}

	// Check business status code
	if scdnResp.Status.Code != 1 {
		return nil, fmt.Errorf("API error: %s (code: %d)", scdnResp.Status.Message, scdnResp.Status.Code)
	}

	// Convert response
	response := &LogDownloadFieldsResponse{
		Status: Status{
			Code:    scdnResp.Status.Code,
			Message: scdnResp.Status.Message,
		},
	}

	if scdnResp.Data != nil {
		dataBytes, err := json.Marshal(scdnResp.Data)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal response data: %w", err)
		}

		var fieldConfigs map[string]LogDownloadFieldConfig
		if err := json.Unmarshal(dataBytes, &fieldConfigs); err != nil {
			return nil, fmt.Errorf("failed to unmarshal fields data: %w", err)
		}
		response.Data = fieldConfigs
	}

	return response, nil
}

// ListLogDownloadTemplates lists log download templates
func (s *ScdnService) ListLogDownloadTemplates(req LogDownloadTemplateListRequest) (*LogDownloadTemplateListResponse, error) {
	ctx := context.Background()

	// Get SCDN client from EdgeNextClient
	scdnClient, err := s.client.ScdnClient()
	if err != nil {
		return nil, fmt.Errorf("failed to get SCDN client: %w", err)
	}

	// Convert request to ScdnRequest format with query parameters
	scdnReq := &connectivity.ScdnRequest{
		Query: make(map[string]interface{}),
	}

	// Set default values for page and per_page if not specified
	page := req.Page
	if page <= 0 {
		page = 1
	}
	scdnReq.Query["page"] = page

	perPage := req.PerPage
	if perPage <= 0 {
		perPage = 20
	}
	scdnReq.Query["per_page"] = perPage

	// Optional parameters - only add if they have meaningful values
	// Status: 0=disabled, 1=enabled, -1 or unset means "all"
	// We use -1 as a sentinel value to indicate "not set"
	if req.Status >= 0 {
		scdnReq.Query["status"] = req.Status
	}
	if req.GroupID > 0 {
		scdnReq.Query["group_id"] = req.GroupID
	}
	if req.TemplateName != "" {
		scdnReq.Query["template_name"] = req.TemplateName
	}
	if req.DataSource != "" {
		scdnReq.Query["data_source"] = req.DataSource
	}

	// Call SCDN API (supports both GET and POST)
	scdnResp, err := scdnClient.Get(ctx, EndpointLogDownloadTemplateList, scdnReq)
	if err != nil {
		return nil, fmt.Errorf("failed to list log download templates: %w", err)
	}

	// Check business status code
	if scdnResp.Status.Code != 1 {
		return nil, fmt.Errorf("API error: %s (code: %d)", scdnResp.Status.Message, scdnResp.Status.Code)
	}

	// Convert response
	response := &LogDownloadTemplateListResponse{
		Status: Status{
			Code:    scdnResp.Status.Code,
			Message: scdnResp.Status.Message,
		},
	}

	if scdnResp.Data != nil {
		dataBytes, err := json.Marshal(scdnResp.Data)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal response data: %w", err)
		}

		var templateListData LogDownloadTemplateListData
		if err := json.Unmarshal(dataBytes, &templateListData); err != nil {
			return nil, fmt.Errorf("failed to unmarshal template list data: %w", err)
		}
		response.Data = templateListData
	}

	return response, nil
}

// ListLogDownloadTemplateDomains lists template domains
func (s *ScdnService) ListLogDownloadTemplateDomains(req LogDownloadTemplateDomainListRequest) (*LogDownloadTemplateDomainListResponse, error) {
	ctx := context.Background()

	// Get SCDN client from EdgeNextClient
	scdnClient, err := s.client.ScdnClient()
	if err != nil {
		return nil, fmt.Errorf("failed to get SCDN client: %w", err)
	}

	// Convert request to ScdnRequest format with query parameters
	scdnReq := &connectivity.ScdnRequest{
		Query: make(map[string]interface{}),
	}

	if req.Domain != "" {
		scdnReq.Query["domain"] = req.Domain
	}

	// Call SCDN API
	scdnResp, err := scdnClient.Get(ctx, EndpointLogDownloadTemplateDomainList, scdnReq)
	if err != nil {
		return nil, fmt.Errorf("failed to list template domains: %w", err)
	}

	// Check business status code
	if scdnResp.Status.Code != 1 {
		return nil, fmt.Errorf("API error: %s (code: %d)", scdnResp.Status.Message, scdnResp.Status.Code)
	}

	// Convert response
	response := &LogDownloadTemplateDomainListResponse{
		Status: Status{
			Code:    scdnResp.Status.Code,
			Message: scdnResp.Status.Message,
		},
	}

	if scdnResp.Data != nil {
		dataBytes, err := json.Marshal(scdnResp.Data)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal response data: %w", err)
		}

		// Extract total and data from response
		var respData map[string]interface{}
		if err := json.Unmarshal(dataBytes, &respData); err == nil {
			// Extract total - handle both number and string types
			if total, ok := respData["total"]; ok {
				switch v := total.(type) {
				case float64:
					response.Total = int(v)
				case int:
					response.Total = v
				case int64:
					response.Total = int(v)
				case string:
					// Try to parse string as int
					if parsed, err := strconv.Atoi(v); err == nil {
						response.Total = parsed
					}
				}
			}

			// Extract data array
			if data, ok := respData["data"].([]interface{}); ok {
				domains := make([]string, 0, len(data))
				for _, d := range data {
					if domain, ok := d.(string); ok {
						domains = append(domains, domain)
					}
				}
				response.Data = domains
				// If total was not set, use the length of data array
				if response.Total == 0 && len(domains) > 0 {
					response.Total = len(domains)
				}
			} else if dataStr, ok := respData["data"].([]string); ok {
				// Direct string array
				response.Data = dataStr
				if response.Total == 0 && len(dataStr) > 0 {
					response.Total = len(dataStr)
				}
			}
		} else {
			// Try to unmarshal as array directly
			var domains []string
			if err := json.Unmarshal(dataBytes, &domains); err != nil {
				return nil, fmt.Errorf("failed to unmarshal domain list data: %w", err)
			}
			response.Data = domains
			response.Total = len(domains)
		}
	}

	return response, nil
}

// AddLogDownloadTemplate adds a log download template
func (s *ScdnService) AddLogDownloadTemplate(req LogDownloadTemplateAddRequest) (*LogDownloadTemplateAddResponse, error) {
	ctx := context.Background()

	scdnClient, err := s.client.ScdnClient()
	if err != nil {
		return nil, fmt.Errorf("failed to get SCDN client: %w", err)
	}

	// Convert request to ScdnRequest format
	scdnReq := &connectivity.ScdnRequest{}
	scdnReq.Data = make(map[string]interface{})

	// Manually build the request data to ensure all fields are included, even if 0
	reqBytes, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	if err := json.Unmarshal(reqBytes, &scdnReq.Data); err != nil {
		return nil, fmt.Errorf("failed to unmarshal request data: %w", err)
	}

	// Ensure status is always included, even if 0
	scdnReq.Data["status"] = req.Status

	log.Printf("[DEBUG] Request data: %+v", scdnReq.Data)

	scdnResp, err := scdnClient.Post(ctx, EndpointLogDownloadTemplateAdd, scdnReq)
	if err != nil {
		return nil, fmt.Errorf("failed to call SCDN API: %w", err)
	}

	if scdnResp.Status.Code != 1 {
		return nil, fmt.Errorf("API error: %s (code: %d)", scdnResp.Status.Message, scdnResp.Status.Code)
	}

	// Convert response
	response := &LogDownloadTemplateAddResponse{
		Status: Status{
			Code:    scdnResp.Status.Code,
			Message: scdnResp.Status.Message,
		},
		Data: scdnResp.Data,
	}

	return response, nil
}

// SaveLogDownloadTemplate saves (updates) a log download template
func (s *ScdnService) SaveLogDownloadTemplate(req LogDownloadTemplateSaveRequest) (*LogDownloadTemplateSaveResponse, error) {
	ctx := context.Background()

	scdnClient, err := s.client.ScdnClient()
	if err != nil {
		return nil, fmt.Errorf("failed to get SCDN client: %w", err)
	}

	// Convert request to ScdnRequest format
	scdnReq := &connectivity.ScdnRequest{}
	scdnReq.Data = make(map[string]interface{})

	// Manually build the request data to ensure all fields are included, even if 0
	reqBytes, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	if err := json.Unmarshal(reqBytes, &scdnReq.Data); err != nil {
		return nil, fmt.Errorf("failed to unmarshal request data: %w", err)
	}

	// Ensure status is always included, even if 0
	scdnReq.Data["status"] = req.Status
	// Ensure domain_select_type is always included, even if 0
	scdnReq.Data["domain_select_type"] = req.DomainSelectType

	log.Printf("[DEBUG] Save template request data: %+v", scdnReq.Data)

	scdnResp, err := scdnClient.Post(ctx, EndpointLogDownloadTemplateSave, scdnReq)
	if err != nil {
		return nil, fmt.Errorf("failed to call SCDN API: %w", err)
	}

	if scdnResp.Status.Code != 1 {
		return nil, fmt.Errorf("API error: %s (code: %d)", scdnResp.Status.Message, scdnResp.Status.Code)
	}

	// Convert response
	response := &LogDownloadTemplateSaveResponse{
		Status: Status{
			Code:    scdnResp.Status.Code,
			Message: scdnResp.Status.Message,
		},
		Data: scdnResp.Data,
	}

	return response, nil
}

// DeleteLogDownloadTemplate deletes a log download template
func (s *ScdnService) DeleteLogDownloadTemplate(req LogDownloadTemplateDeleteRequest) (*LogDownloadTemplateDeleteResponse, error) {
	ctx := context.Background()

	var response LogDownloadTemplateDeleteResponse
	err := s.callSCDNAPI(ctx, MethodDELETE, EndpointLogDownloadTemplateDelete, req, &response)
	if err != nil {
		return nil, fmt.Errorf("failed to delete log download template: %w", err)
	}

	return &response, nil
}

// BatchDeleteLogDownloadTemplates batch deletes log download templates
func (s *ScdnService) BatchDeleteLogDownloadTemplates(req LogDownloadTemplateBatchDeleteRequest) (*LogDownloadTemplateBatchDeleteResponse, error) {
	ctx := context.Background()

	var response LogDownloadTemplateBatchDeleteResponse
	err := s.callSCDNAPI(ctx, MethodDELETE, EndpointLogDownloadTemplateBatchDelete, req, &response)
	if err != nil {
		return nil, fmt.Errorf("failed to batch delete log download templates: %w", err)
	}

	return &response, nil
}

// ChangeLogDownloadTemplateStatus changes template status
func (s *ScdnService) ChangeLogDownloadTemplateStatus(req LogDownloadTemplateChangeStatusRequest) (*LogDownloadTemplateChangeStatusResponse, error) {
	ctx := context.Background()

	var response LogDownloadTemplateChangeStatusResponse
	err := s.callSCDNAPI(ctx, MethodPOST, EndpointLogDownloadTemplateChangeStatus, req, &response)
	if err != nil {
		return nil, fmt.Errorf("failed to change template status: %w", err)
	}

	return &response, nil
}

// BatchChangeLogDownloadTemplateStatus batch changes template status
func (s *ScdnService) BatchChangeLogDownloadTemplateStatus(req LogDownloadTemplateBatchChangeStatusRequest) (*LogDownloadTemplateBatchChangeStatusResponse, error) {
	ctx := context.Background()

	var response LogDownloadTemplateBatchChangeStatusResponse
	err := s.callSCDNAPI(ctx, MethodPOST, EndpointLogDownloadTemplateBatchChangeStatus, req, &response)
	if err != nil {
		return nil, fmt.Errorf("failed to batch change template status: %w", err)
	}

	return &response, nil
}

// GetAllLogDownloadTemplates gets all templates (for adding tasks)
func (s *ScdnService) GetAllLogDownloadTemplates() (*LogDownloadTemplateAllResponse, error) {
	ctx := context.Background()

	// Get SCDN client from EdgeNextClient
	scdnClient, err := s.client.ScdnClient()
	if err != nil {
		return nil, fmt.Errorf("failed to get SCDN client: %w", err)
	}

	// Call SCDN API (supports both GET and POST)
	scdnReq := &connectivity.ScdnRequest{}
	scdnResp, err := scdnClient.Get(ctx, EndpointLogDownloadTemplateAll, scdnReq)
	if err != nil {
		return nil, fmt.Errorf("failed to get all log download templates: %w", err)
	}

	// Check business status code
	if scdnResp.Status.Code != 1 {
		return nil, fmt.Errorf("API error: %s (code: %d)", scdnResp.Status.Message, scdnResp.Status.Code)
	}

	// Convert response
	response := &LogDownloadTemplateAllResponse{
		Status: Status{
			Code:    scdnResp.Status.Code,
			Message: scdnResp.Status.Message,
		},
	}

	if scdnResp.Data != nil {
		dataBytes, err := json.Marshal(scdnResp.Data)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal response data: %w", err)
		}

		var templateGroups map[string]LogDownloadTemplateGroup
		if err := json.Unmarshal(dataBytes, &templateGroups); err != nil {
			return nil, fmt.Errorf("failed to unmarshal template groups data: %w", err)
		}
		response.Data = templateGroups
	}

	return response, nil
}

// GetAllLogDownloadTemplateGroups gets all template groups
func (s *ScdnService) GetAllLogDownloadTemplateGroups() (*LogDownloadTemplateGroupAllResponse, error) {
	ctx := context.Background()

	// Get SCDN client from EdgeNextClient
	scdnClient, err := s.client.ScdnClient()
	if err != nil {
		return nil, fmt.Errorf("failed to get SCDN client: %w", err)
	}

	// Call SCDN API (supports both GET and POST)
	scdnReq := &connectivity.ScdnRequest{}
	scdnResp, err := scdnClient.Get(ctx, EndpointLogDownloadTemplateGroupAll, scdnReq)
	if err != nil {
		return nil, fmt.Errorf("failed to get all template groups: %w", err)
	}

	// Check business status code
	if scdnResp.Status.Code != 1 {
		return nil, fmt.Errorf("API error: %s (code: %d)", scdnResp.Status.Message, scdnResp.Status.Code)
	}

	// Convert response
	response := &LogDownloadTemplateGroupAllResponse{
		Status: Status{
			Code:    scdnResp.Status.Code,
			Message: scdnResp.Status.Message,
		},
	}

	// Extract array data from various formats
	if err := extractArrayFromData(scdnResp.Data, &response.Data); err != nil {
		return nil, fmt.Errorf("failed to extract template groups data: %w", err)
	}

	return response, nil
}
