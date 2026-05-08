package rds

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
)

//go:embed rds_mysql_charset_collations.json
var rdsMySQLCharsetCollationsJSON []byte

var (
	rdsMySQLCharsetCollations map[string][]string
	rdsMySQLCharsetNames      []string
)

func init() {
	if err := json.Unmarshal(rdsMySQLCharsetCollationsJSON, &rdsMySQLCharsetCollations); err != nil {
		panic("edgenext rds: rds_mysql_charset_collations.json: " + err.Error())
	}
	if len(rdsMySQLCharsetCollations) == 0 {
		panic("edgenext rds: empty rds_mysql_charset_collations.json")
	}
	for k := range rdsMySQLCharsetCollations {
		rdsMySQLCharsetNames = append(rdsMySQLCharsetNames, k)
	}
	sort.Strings(rdsMySQLCharsetNames)
}

func rdsMySQLCollateValid(charset, collate string) bool {
	charset = strings.ToLower(strings.TrimSpace(charset))
	collate = strings.TrimSpace(collate)
	list, ok := rdsMySQLCharsetCollations[charset]
	if !ok {
		return false
	}
	for _, c := range list {
		if strings.EqualFold(c, collate) {
			return true
		}
	}
	return false
}

func validateRDSMySQLCharacterSet(v interface{}, _ cty.Path) diag.Diagnostics {
	if v == nil {
		return nil
	}
	raw, ok := v.(string)
	if !ok {
		return diag.Diagnostics{{
			Severity: diag.Error,
			Summary:  "Invalid character_set",
			Detail:   fmt.Sprintf("Expected string, got %T", v),
		}}
	}
	s := strings.ToLower(strings.TrimSpace(raw))
	if s == "" {
		return nil
	}
	if _, exists := rdsMySQLCharsetCollations[s]; !exists {
		return diag.Diagnostics{{
			Severity: diag.Error,
			Summary:  "Invalid character_set",
			Detail: fmt.Sprintf(
				"%q is not a supported MySQL character set. Examples: ascii, binary, gbk, latin1, utf8, utf8mb4.",
				raw,
			),
		}}
	}
	return nil
}
