package webapi

import (
	"regexp"
	"strings"
)

var idForbiddenChars = regexp.MustCompile(`[^a-zA-Z0-9-]`)

// ID sanitizes entity ID, removing forbidden characters.
// ID moze pochadzat z URL alebo ineho uzivatelskeho vstupu.
func ID(entityID string) string {
	return idForbiddenChars.ReplaceAllString(entityID, "")
}

func Or(conds ...string) string {
	return "(" + strings.Join(conds, " or ") + ")"
}

func And(conds ...string) string {
	return "(" + strings.Join(conds, " and ") + ")"
}
