package utils

import (
	"strings"
)

// EscapeCharacters for a bash string
func EscapeCharacters(str string) string {
	str = strings.ReplaceAll(str, `"`, `\"`)
	str = strings.ReplaceAll(str, "'", "\\'")
	return str
}
