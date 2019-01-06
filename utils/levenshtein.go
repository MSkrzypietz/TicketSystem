package utils

// This file provides a recursive implementation of the Levenshtein distance
// For more information see https://en.wikipedia.org/wiki/Levenshtein_distance

func CheckStringsDeviation(allowedDeviation int, a, b string) bool {
	return levenshteinDistance(a, len(a), b, len(b)) <= allowedDeviation
}

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

func minimum(a, b, c int) int {
	if a <= b && a <= c {
		return a
	}

	if b <= c && b <= a {
		return b
	}

	return c
}
