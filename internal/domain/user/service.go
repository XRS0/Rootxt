package user

import "strings"

func NormalizeEmail(value string) string {
	return strings.ToLower(strings.TrimSpace(value))
}
