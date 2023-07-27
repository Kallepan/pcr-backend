package main

import (
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.com/kaka/pcr-backend/common/controllers"
	"gitlab.com/kaka/pcr-backend/test"
)

// Test Ping
func TestPing(t *testing.T) {
	w := httptest.NewRecorder()
	ctx := test.GetTestGinContext(w)

	test.MockJsonGet(ctx, nil, nil)

	controllers.Ping(ctx)

	assert.Equal(t, 200, w.Code)

	// Test the response body { "message": "pong" }
	assert.JSONEq(t, `{"message":"pong"}`, w.Body.String())
}
