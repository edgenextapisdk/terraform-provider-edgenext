package helper

import (
	"fmt"
	"strconv"
)

// ParseAPIResponsePayload validates common code/msg envelope and returns payload data.
// If "data" does not exist, it returns the full response map as payload.
func ParseAPIResponsePayload(resp map[string]interface{}) (interface{}, error) {
	if codeRaw, ok := resp["code"]; ok {
		code, ok := toInt(codeRaw)
		if !ok {
			return nil, fmt.Errorf("invalid response code type: %T", codeRaw)
		}
		if code != 0 {
			msg := "unknown error"
			if v, ok := resp["msg"].(string); ok && v != "" {
				msg = v
			}
			return nil, fmt.Errorf("api error: code=%d msg=%s", code, msg)
		}
	}

	if data, ok := resp["data"]; ok {
		return data, nil
	}
	return resp, nil
}

func ParseAPIResponseMap(resp map[string]interface{}) (map[string]interface{}, error) {
	payload, err := ParseAPIResponsePayload(resp)
	if err != nil {
		return nil, err
	}
	if payload == nil {
		return map[string]interface{}{}, nil
	}
	m, ok := payload.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("expected object payload, got %T", payload)
	}
	return m, nil
}

func ParseAPIResponseList(resp map[string]interface{}) ([]interface{}, error) {
	payload, err := ParseAPIResponsePayload(resp)
	if err != nil {
		return nil, err
	}
	return listFromPayloadContent(payload)
}

// listFromPayloadContent extracts a JSON array from either a top-level array or a wrapper object
// (e.g. {"list":[...], "total":0}) returned as the API "data" field.
func listFromPayloadContent(payload interface{}) ([]interface{}, error) {
	if payload == nil {
		return []interface{}{}, nil
	}
	if l, ok := payload.([]interface{}); ok {
		return l, nil
	}
	m, ok := payload.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("expected list payload, got %T", payload)
	}
	for _, key := range []string{"list", "items", "records", "keypairs", "ports"} {
		if v, ok := m[key]; ok {
			if l, ok := v.([]interface{}); ok {
				return l, nil
			}
		}
	}
	// Nested envelope: data -> { list: [...] } or data -> [...]
	if v, ok := m["data"]; ok {
		return listFromPayloadContent(v)
	}
	return nil, fmt.Errorf("expected list payload, got map[string]interface{} without a recognized list field")
}

func ExtractIDFromPayload(payload map[string]interface{}, fallback string) string {
	if id, ok := payload["id"]; ok {
		if str := toString(id); str != "" {
			return str
		}
	}
	return fallback
}

// StringFromMap returns a string value from map by key.
func StringFromMap(m map[string]interface{}, key string) string {
	v, ok := m[key]
	if !ok || v == nil {
		return ""
	}
	if s, ok := v.(string); ok {
		return s
	}
	return ""
}

// IntFromMap returns an int value from map by key.
func IntFromMap(m map[string]interface{}, key string) int {
	v, ok := m[key]
	if !ok || v == nil {
		return 0
	}
	if n, ok := toInt(v); ok {
		return n
	}
	return 0
}

// MapFromMap returns a map object from map by key.
func MapFromMap(m map[string]interface{}, key string) map[string]interface{} {
	v, ok := m[key]
	if !ok || v == nil {
		return nil
	}
	obj, ok := v.(map[string]interface{})
	if !ok {
		return nil
	}
	return obj
}

// ListFromMap returns a list value from map by key.
func ListFromMap(m map[string]interface{}, key string) []interface{} {
	v, ok := m[key]
	if !ok || v == nil {
		return []interface{}{}
	}
	return InterfaceToList(v)
}

// InterfaceToList converts interface{} to list if possible.
func InterfaceToList(v interface{}) []interface{} {
	rawList, ok := v.([]interface{})
	if !ok {
		return []interface{}{}
	}
	return rawList
}

// InterfaceToStringSlice converts interface list to string list.
func InterfaceToStringSlice(v interface{}) []interface{} {
	rawList := InterfaceToList(v)
	out := make([]interface{}, 0, len(rawList))
	for _, item := range rawList {
		if s, ok := item.(string); ok {
			out = append(out, s)
		}
	}
	return out
}

func toString(v interface{}) string {
	switch t := v.(type) {
	case string:
		return t
	case int:
		return strconv.Itoa(t)
	case int32:
		return strconv.FormatInt(int64(t), 10)
	case int64:
		return strconv.FormatInt(t, 10)
	case float64:
		return strconv.FormatInt(int64(t), 10)
	default:
		return fmt.Sprintf("%v", v)
	}
}

func toInt(v interface{}) (int, bool) {
	switch t := v.(type) {
	case int:
		return t, true
	case int32:
		return int(t), true
	case int64:
		return int(t), true
	case float64:
		return int(t), true
	case string:
		n, err := strconv.Atoi(t)
		if err != nil {
			return 0, false
		}
		return n, true
	default:
		return 0, false
	}
}
