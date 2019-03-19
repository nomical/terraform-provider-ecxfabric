package apiclient

// String returns a pointer to the string value passed in.
func StringPtr(v string) *string {
	return &v
}

// StringValue returns the value of the string pointer passed in or
// "" if the pointer is nil.
func StringPtrValue(v *string) string {
	if v != nil {
		return *v
	}
	return ""
}

// Int returns a pointer to the int value passed in.
func IntPtr(v int) *int {
	return &v
}

// IntValue returns the value of the int pointer passed in or
// 0 if the pointer is nil.
func IntPtrValue(v *int) int {
	if v != nil {
		return *v
	}
	return 0
}
