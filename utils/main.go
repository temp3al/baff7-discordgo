package utils

// Gets a slice's string value with a fallback.
//
// Returns string if found, provided def string if not.
func GetSliceStr(slice []string, i int, def string) string {
	if len(slice) > i {
		return slice[i]
	}
	return def
}
