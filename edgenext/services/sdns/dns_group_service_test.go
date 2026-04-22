package sdns

import (
	"testing"
)

func TestSdnsService_DnsGroupCRUD(t *testing.T) {
	if !isIntegrationTest() {
		t.Skip("Skipping integration test")
	}

	client := createTestClient(t)
	service := NewSdnsService(client)

	groupName := "test-dns-group"

	// 1. Add
	var groupID int
	t.Run("AddGroup", func(t *testing.T) {
		resp, err := service.AddDnsGroup(DnsGroupAddRequest{
			GroupName: groupName,
			Remark:    "Test group",
		})
		if err != nil {
			t.Fatalf("AddDnsGroup failed: %v", err)
		}
		groupID = resp.Data.ID
	})

	// 2. List & Info
	t.Run("ListAndInfo", func(t *testing.T) {
		resp, err := service.ListDnsGroups(DnsGroupListRequest{GroupName: groupName})
		if err != nil {
			t.Fatalf("ListDnsGroups failed: %v", err)
		}
		if len(resp.List) == 0 {
			t.Fatalf("Group not found in list")
		}

		info, err := service.GetDnsGroupInfo(groupID)
		if err != nil {
			t.Fatalf("GetDnsGroupInfo failed: %v", err)
		}
		if info.GroupName != groupName {
			t.Errorf("Group name mismatch: expected %s, got %s", groupName, info.GroupName)
		}
	})

	// 3. Update
	t.Run("UpdateGroup", func(t *testing.T) {
		err := service.UpdateDnsGroup(DnsGroupSaveRequest{
			GroupID:   groupID,
			GroupName: groupName + "-updated",
			Remark:    "Updated group",
			DomainIDs: []int{304, 582},
		})
		if err != nil {
			t.Fatalf("UpdateDnsGroup failed: %v", err)
		}
	})

	// 4. Delete
	t.Run("DeleteGroup", func(t *testing.T) {
		err := service.DeleteDnsGroup(groupID)
		if err != nil {
			t.Fatalf("DeleteDnsGroup failed: %v", err)
		}
	})
}
