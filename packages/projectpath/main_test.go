package projectpath

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRootPath(t *testing.T) {
	// Test that the root path is correct
	// We could test against the full path but that would be brittle
	assert.NotContains(t, Root, "main_test.go")
	assert.NotContains(t, Root, "packages")
	assert.Contains(t, Root, "pcr-backend")
}
