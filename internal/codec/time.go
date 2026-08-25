package codec

import "time"

func Location(name string) *time.Location {
	value, err := time.LoadLocation(name)
	if err != nil {
		return time.UTC
	}
	return value
}
func Parse(value string, location *time.Location) (time.Time, error) {
	if value == "" {
		return time.Time{}, nil
	}
	if location == nil {
		location = time.UTC
	}
	parsed, err := time.ParseInLocation(time.RFC3339, value, location)
	if err != nil {
		return time.Time{}, err
	}
	return parsed.UTC(), nil
}
func Format(value time.Time, location *time.Location) string {
	if value.IsZero() {
		return ""
	}
	if location == nil {
		location = time.UTC
	}
	return value.In(location).Format(time.RFC3339)
}
