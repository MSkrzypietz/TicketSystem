package utils

// Checks the similarity between two strings using the allowedDeviation with the levenshtein algorithm
func CheckStringsDeviation(allowedDeviation int, a, b string) bool {
	return levenshteinDistance(a, len(a), b, len(b)) <= allowedDeviation
}

// This is a recursive implementation of the Levenshtein distance
// For more information see https://en.wikipedia.org/wiki/Levenshtein_distance
func levenshteinDistance(a string, len_a int, b string, len_b int) int {
	var cost int

	// base case: empty strings
	if len_a == 0 {
		return len_b
	}
	if len_b == 0 {
		return len_a
	}

	// test if last characters of the strings match
	if a[len_a-1] == b[len_b-1] {
		cost = 0
	} else {
		cost = 1
	}

	// return minimum of delete char from s, delete char from t, and delete char from both
	return minimum(
		levenshteinDistance(a, len_a-1, b, len_b)+1,
		levenshteinDistance(a, len_a, b, len_b-1)+1,
		levenshteinDistance(a, len_a-1, b, len_b-1)+cost)
}

// Returns the minimum from three numbers
func minimum(a, b, c int) int {
	if a <= b && a <= c {
		return a
	}

	if b <= c && b <= a {
		return b
	}

	return c
}
