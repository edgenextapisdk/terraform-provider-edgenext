package helper

import (
	"testing"
)

func TestParseAPIResponseList_topLevelArray(t *testing.T) {
	resp := map[string]interface{}{
		"code": float64(0),
		"data": []interface{}{
			map[string]interface{}{"id": "a", "name": "n1"},
		},
	}
	out, err := ParseAPIResponseList(resp)
	if err != nil {
		t.Fatal(err)
	}
	if len(out) != 1 {
		t.Fatalf("len=%d", len(out))
	}
}

func TestParseAPIResponseList_envelopeList(t *testing.T) {
	resp := map[string]interface{}{
		"code": float64(0),
		"data": map[string]interface{}{
			"list": []interface{}{
				map[string]interface{}{"id": "1", "name": "kp"},
			},
			"total": float64(1),
		},
	}
	out, err := ParseAPIResponseList(resp)
	if err != nil {
		t.Fatal(err)
	}
	if len(out) != 1 {
		t.Fatalf("len=%d", len(out))
	}
}

func TestParseAPIResponseList_keypairsField(t *testing.T) {
	resp := map[string]interface{}{
		"code": float64(0),
		"data": map[string]interface{}{
			"keypairs": []interface{}{
				map[string]interface{}{"name": "only-name"},
			},
		},
	}
	out, err := ParseAPIResponseList(resp)
	if err != nil {
		t.Fatal(err)
	}
	if len(out) != 1 {
		t.Fatalf("len=%d", len(out))
	}
}

func TestParseAPIResponseList_keypairsNovaEnvelope(t *testing.T) {
	resp := map[string]interface{}{
		"code": float64(0),
		"msg":  "success",
		"data": map[string]interface{}{
			"keypairs": []interface{}{
				map[string]interface{}{
					"keypair": map[string]interface{}{
						"name": "example-key",
						"type": "ssh",
					},
				},
				map[string]interface{}{
					"keypair": map[string]interface{}{
						"name": "other",
					},
				},
			},
		},
	}
	out, err := ParseAPIResponseList(resp)
	if err != nil {
		t.Fatal(err)
	}
	if len(out) != 2 {
		t.Fatalf("len=%d", len(out))
	}
}

func TestResponseHelpers_MapExtraction(t *testing.T) {
	row := map[string]interface{}{
		"name":   "demo",
		"count":  float64(3),
		"nested": map[string]interface{}{"id": "n-1"},
		"tags":   []interface{}{"a", 1, "b"},
	}

	if got := StringFromMap(row, "name"); got != "demo" {
		t.Fatalf("name=%s", got)
	}
	if got := IntFromMap(row, "count"); got != 3 {
		t.Fatalf("count=%d", got)
	}
	if got := StringFromMap(MapFromMap(row, "nested"), "id"); got != "n-1" {
		t.Fatalf("nested.id=%s", got)
	}
	tags := InterfaceToStringSlice(row["tags"])
	if len(tags) != 2 {
		t.Fatalf("tags len=%d", len(tags))
	}
}
