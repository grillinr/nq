package db

func safeStringValue(value *string) string {
	if value == nil {
		return ""
	}
	return *value
}
