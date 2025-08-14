package utils

import "math/rand/v2"

// Gets a slice's string value with a fallback.
//
// Returns string if found, provided def string if not.
func GetSliceStr(slice []string, i int, def string) string {
	if len(slice) > i {
		return slice[i]
	}
	return def
}

// Return a ranged float64.
func RandFloat(min, max float64) float64 {
	return min + rand.Float64()*(max-min)
}

func RandomStringFromList(l []string) string {
	return l[rand.IntN(len(l))]
}
