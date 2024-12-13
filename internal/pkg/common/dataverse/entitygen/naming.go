package entitygen

import (
	"fmt"
	"strings"
	"unicode"

	"github.com/michaljurecko/ch-demo/internal/pkg/common/dataverse/metadata"
)

func entityFileName(entity metadata.Entity) (string, error) {
	v := normalizeName(entity.DisplayName.LocalizedLabels[0].Label, false)
	if v == "" {
		return "", fmt.Errorf(`empty entity "%s" struct name`, entity.LogicalName)
	}
	return v + ".go", nil
}

func entityGoName(entity metadata.Entity) (string, error) {
	v := normalizeName(entity.DisplayName.LocalizedLabels[0].Label, true)
	if v == "" {
		return "", fmt.Errorf(`empty entity "%s" struct name`, entity.LogicalName)
	}
	return v, nil
}

func attributeGoName(entity *entitySpec, attr metadata.Attribute, structName string) (string, error) {
	// There may be one primary identifier attribute
	if attr.IsPrimaryID {
		return "ID", nil
	}

	var name string
	if len(attr.DisplayName.LocalizedLabels) > 0 {
		name = attr.DisplayName.LocalizedLabels[0].Label
	} else {
		// Fallback
		name = attr.SchemaName
	}

	// Remove struct name from attribute name, if present
	name = strings.TrimPrefix(name, structName)

	v := normalizeName(name, true)
	if v == "" {
		return "", fmt.Errorf(`empty attribute "%s.%s" name`, entity.LogicalName, attr.LogicalName)
	}
	return v, nil
}

func entityDesc(entity metadata.Entity) string {
	if len(entity.Description.LocalizedLabels) > 0 {
		return entity.Description.LocalizedLabels[0].Label
	}
	return ""
}

func attributeDesc(attr metadata.Attribute) string {
	if !attr.IsPrimaryID && len(attr.Description.LocalizedLabels) > 0 {
		return attr.Description.LocalizedLabels[0].Label
	}
	return ""
}

// normalizeName converts a string by removing separators.
// Separators are whitespace characters and special characters.
// If camelCase is true, the first letter of each word is capitalized.
func normalizeName(input string, camelCase bool) string {
	var result string
	upperNext := true

	for _, char := range input {
		if unicode.IsLetter(char) {
			switch {
			case upperNext && camelCase:
				result += string(unicode.ToUpper(char))
				upperNext = false
			case !camelCase:
				result += string(unicode.ToLower(char))
			default:
				result += string(char)
			}
		} else {
			upperNext = true
		}
	}

	return result
}
