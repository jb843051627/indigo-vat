package validation

import "strings"

func CleanText(value string) string {
	return strings.TrimSpace(strings.Join(strings.Fields(value), " "))
}
func FirstNonEmpty(values ...string) string {
	for _, value := range values {
		if CleanText(value) != "" {
			return CleanText(value)
		}
	}
	return ""
}
func IsID(value string) bool {
	value = CleanText(value)
	return len(value) >= 6 && !strings.ContainsAny(value, "\r\n\t ")
}
