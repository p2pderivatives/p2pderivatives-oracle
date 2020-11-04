package decompose

import "strconv"

// Value takes an integer value and returns its decomposition under
// the provided base.
func Value(value int, base int, length int) []string {
	var result []string
	for value > 0 {
		digit := strconv.Itoa(value % base)
		result = append([]string{digit}, result...)
		value = value / base
	}

	for len(result) < length {
		result = append([]string{"0"}, result...)
	}

	return result
}
