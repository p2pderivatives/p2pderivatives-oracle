package decompose

import "strconv"

func DecomposeValue(value int, base int, length int) []string {
	var result []string
	for value > 0 {
		digit := strconv.Itoa(value % base)
		result = append(result, digit)
		value = value / base
	}

	for len(result) < length {
		result = append(result, "0")
	}

	return result
}
