package domain_group

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/edgenextapisdk/terraform-provider-edgenext/edgenext/connectivity"
	"github.com/stretchr/testify/assert"
)

// Re-implementing test helpers locally since they are not exported from parent package

type TestConfig struct {
	AccessKey              string `json:"access_key"`
	SecretKey              string `json:"secret_key"`
	Endpoint               string `json:"endpoint"`
	EnableIntegrationTests bool   `json:"enable_integration_tests"`
}

func getTestConfig() *TestConfig {
	config := &TestConfig{
		Endpoint:               "https://api.edgenextscdn.com",
		EnableIntegrationTests: false,
	}

	// Try to find config file by traversing up
	paths := []string{
		"test_config.json",
		"../test_config.json",
		"../../test_config.json",
		"../../../test_config.json",
	}

	for _, path := range paths {
		if _, err := os.Stat(path); err == nil {
			data, err := ioutil.ReadFile(path)
			if err == nil {
				json.Unmarshal(data, config)
				return config
			}
		}
	}

	// Fallback to Env if config file not found or incomplete ?
	// The existing pattern uses config file. Let's stick to it or envs.
	if os.Getenv("EDGENEXT_ACCESS_KEY") != "" {
		config.AccessKey = os.Getenv("EDGENEXT_ACCESS_KEY")
	}
	if os.Getenv("EDGENEXT_SECRET_KEY") != "" {
		config.SecretKey = os.Getenv("EDGENEXT_SECRET_KEY")
	}
	if os.Getenv("EDGENEXT_API_ENDPOINT") != "" {
		config.Endpoint = os.Getenv("EDGENEXT_API_ENDPOINT")
	}
	if config.AccessKey != "" && config.SecretKey != "" {
		config.EnableIntegrationTests = true
	}

	return config
}

func createTestClient(t *testing.T) *connectivity.EdgeNextClient {
	config := getTestConfig()
	if !config.EnableIntegrationTests {
		t.Skip("Skipping integration test")
	}

	connectivityConfig := &connectivity.Config{
		AccessKey: config.AccessKey,
		SecretKey: config.SecretKey,
		Endpoint:  config.Endpoint,
	}

	client, err := connectivityConfig.Client()
	if err != nil {
		t.Fatalf("Failed to create EdgeNextClient: %v", err)
	}
	return client
}

func TestDomainGroupService_CRUD(t *testing.T) {
	client := createTestClient(t)
	service := NewDomainGroupService(client)

	// 1. Create Group
	groupName := "tf-test-group-" + randomString(6)
	addReq := DomainGroupSaveRequest{
		GroupName: groupName,
		Remark:    "Created by Unit Test",
	}
	addResp, err := service.AddDomainGroup(addReq)
	if err != nil {
		t.Fatalf("AddDomainGroup failed: %v", err)
	}
	assert.Equal(t, 1, addResp.Status.Code)
	groupID := addResp.Data.ID
	t.Logf("Created Domain Group ID: %s", groupID)

	// CLEANUP
	defer func() {
		groupIDInt, _ := strconv.Atoi(groupID)
		delReq := DomainGroupDelRequest{GroupID: groupIDInt}
		_, _ = service.DeleteDomainGroup(delReq)
	}()

	// 2. Info
	groupIDInt, _ := strconv.Atoi(groupID)
	info, err := service.GetDomainGroupInfo(groupIDInt)
	if err != nil {
		t.Errorf("GetDomainGroupInfo failed: %v", err)
	} else {
		assert.Equal(t, groupName, info.GroupName)
	}

	// 3. Update
	updateReq := DomainGroupSaveRequest{
		GroupID:   groupIDInt,
		GroupName: groupName + "_updated",
		Remark:    "Updated by Unit Test",
	}
	updateResp, err := service.UpdateDomainGroup(updateReq)
	if err != nil {
		t.Errorf("UpdateDomainGroup failed: %v", err)
	}
	assert.Equal(t, 1, updateResp.Status.Code)

	// 4. List
	listReq := DomainGroupListRequest{
		GroupName: groupName, // Filter by original name might fail if updated?
		Page:      1,
		PerPage:   10,
	}
	listResp, err := service.ListDomainGroups(listReq)
	if err != nil {
		t.Errorf("ListDomainGroups failed: %v", err)
	}
	if listResp != nil {
		// Just check if we get a response
		t.Logf("List count: %s", listResp.Data.Total)
	}

	// 5. Delete (handled by defer, but let's test explicit fail or success logic if needed?
	// defer will run at end. We can explicitly delete here.)
	delReq := DomainGroupDelRequest{GroupID: groupIDInt}
	delResp, err := service.DeleteDomainGroup(delReq)
	if err != nil {
		t.Errorf("DeleteDomainGroup failed: %v", err)
	}
	assert.Equal(t, 1, delResp.Status.Code)
}

func randomString(n int) string {
	const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[time.Now().UnixNano()%int64(len(letterBytes))]
	}
	return string(b)
}
