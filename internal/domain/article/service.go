package article

import (
	"regexp"
	"strings"
)

var slugCleanup = regexp.MustCompile(`[^a-z0-9]+`)

func Slugify(value string) string {
	clean := strings.ToLower(strings.TrimSpace(value))
	clean = slugCleanup.ReplaceAllString(clean, "-")
	clean = strings.Trim(clean, "-")
	if clean == "" {
		return "article"
	}
	return clean
}
