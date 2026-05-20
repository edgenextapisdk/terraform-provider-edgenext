package connectivity

import (
	"testing"
)

func TestEdgeNextClientServiceClientsRequireConfig(t *testing.T) {
	cfg := &Config{
		AccessKey: "",
		SecretKey: "",
		Endpoint:  "",
	}
	edge, err := cfg.Client()
	if err != nil {
		t.Fatalf("Config.Client() unexpected error: %v", err)
	}

	cases := []struct {
		name string
		fn   func() (interface{}, error)
	}{
		{"APIClient", func() (interface{}, error) { return edge.APIClient() }},
		{"OSSClient", func() (interface{}, error) { return edge.OSSClient() }},
		{"ScdnClient", func() (interface{}, error) { return edge.ScdnClient() }},
		{"ECSClient", func() (interface{}, error) { return edge.ECSClient() }},
		{"RDSClient", func() (interface{}, error) { return edge.RDSClient() }},
		{"ELBClient", func() (interface{}, error) { return edge.ELBClient() }},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := tc.fn()
			if err == nil {
				t.Fatalf("%s: expected initialization error, got nil error", tc.name)
			}
		})
	}
}

func TestEdgeNextClientECSClientSuccess(t *testing.T) {
	cfg := &Config{
		AccessKey: "ak",
		SecretKey: "sk",
		Endpoint:  "https://api.example.com",
		Region:    "tokyo-a",
	}
	edge, err := cfg.Client()
	if err != nil {
		t.Fatalf("Config.Client() unexpected error: %v", err)
	}

	client, err := edge.ECSClient()
	if err != nil {
		t.Fatalf("ECSClient() unexpected error: %v", err)
	}
	if client == nil {
		t.Fatal("ECSClient() returned nil client without error")
	}

	// sync.Once: second call returns same client and no error
	client2, err := edge.ECSClient()
	if err != nil {
		t.Fatalf("ECSClient() second call unexpected error: %v", err)
	}
	if client2 != client {
		t.Fatal("ECSClient() second call returned different instance")
	}
}
