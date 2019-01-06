package utils

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCheckStringsDeviation(t *testing.T) {
	tests := []struct {
		a, b     string
		expected bool
	}{
		{"", "hello", false},
		{"hello", "hello", true},
		{"ab", "aa", true},
		{"ab", "aaa", true},
		{"kitten", "sitting", false},
	}
	for _, d := range tests {
		assert.Equal(t, d.expected, CheckStringsDeviation(2, d.a, d.b))
	}
}

func TestLevenshteinDistance(t *testing.T) {
	tests := []struct {
		a, b     string
		expected int
	}{
		{"", "hello", 5},
		{"hello", "hello", 0},
		{"ab", "aa", 1},
		{"ab", "aaa", 2},
		{"kitten", "sitting", 3},
	}
	for _, d := range tests {
		assert.Equal(t, d.expected, levenshteinDistance(d.a, len(d.a), d.b, len(d.b)))
	}
}

func TestMinimum(t *testing.T) {
	assert.Equal(t, 1, minimum(1, 2, 3))
	assert.Equal(t, 2, minimum(4, 2, 3))
	assert.Equal(t, 3, minimum(4, 5, 3))
}
