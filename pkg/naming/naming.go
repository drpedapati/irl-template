package naming

import (
	"regexp"
	"strings"
	"time"
)

// Slugify converts a purpose/description into a clean folder name
func Slugify(input string) string {
	// Lowercase
	s := strings.ToLower(input)

	// Remove common filler words
	fillers := []string{" the ", " a ", " an ", " for ", " of ", " in ", " on ", " to ", " and ", " with "}
	for _, f := range fillers {
		s = strings.ReplaceAll(s, f, " ")
	}

	// Replace non-alphanumeric with hyphens
	reg := regexp.MustCompile(`[^a-z0-9]+`)
	s = reg.ReplaceAllString(s, "-")

	// Trim hyphens from ends
	s = strings.Trim(s, "-")

	// Collapse multiple hyphens
	reg = regexp.MustCompile(`-+`)
	s = reg.ReplaceAllString(s, "-")

	// Truncate to reasonable length
	if len(s) > 40 {
		// Try to cut at a hyphen
		if idx := strings.LastIndex(s[:40], "-"); idx > 20 {
			s = s[:idx]
		} else {
			s = s[:40]
		}
	}

	return s
}

// GenerateName creates a timestamped project name from a purpose
// Format: YYMMDD-slug (e.g., 260129-erp-analysis)
func GenerateName(purpose string) string {
	timestamp := time.Now().Format("060102")
	slug := Slugify(purpose)

	if slug == "" {
		slug = "project"
	}

	return timestamp + "-" + slug
}

// Timestamp returns current timestamp prefix
func Timestamp() string {
	return time.Now().Format("060102")
}
