package main

// Matrikelnummern: 6813128, 1665910, 7612558

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCheckPortBoundaries(t *testing.T) {
	tests := []struct {
		port     int
		expected bool
	}{
		{0, true},
		{1337, true},
		{65535, true},
		{-1, false},
		{65536, false},
	}
	for _, d := range tests {
		assert.Equal(t, d.expected, checkPortBoundaries(d.port))
	}
}

func TestExistsPath(t *testing.T) {
	tests := []struct {
		path     string
		expected bool
	}{
		{"test", false},
		{"ticketSystem.go", true},
		{"../ticketSystem/ticketSystem.go", true},
	}
	for _, d := range tests {
		assert.Equal(t, d.expected, existsPath(d.path))
	}
}
