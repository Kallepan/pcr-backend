package main

import (
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.com/kallepan/pcr-backend/test"
)

// Test Ping
func TestPing(t *testing.T) {
	w := httptest.NewRecorder()
	ctx := test.GetTestGinContext(w)

	test.MockJsonGet(ctx, nil, nil)

	assert.Equal(t, 200, w.Code)
}
