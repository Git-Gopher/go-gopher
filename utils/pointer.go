package utils

// String returns a pointer to string.
func String(s string) *string {
	return &s
}

// Int returns a pointer to int.
func Int(i int) *int {
	return &i
}
