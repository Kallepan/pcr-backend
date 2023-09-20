package test

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"

	"github.com/gin-gonic/gin"
)

// Get Test Gin context
func GetTestGinContext(w *httptest.ResponseRecorder) *gin.Context {
	gin.SetMode(gin.TestMode)

	ctx, _ := gin.CreateTestContext(w)
	ctx.Request = &http.Request{
		Header: make(http.Header),
		URL:    &url.URL{},
	}

	return ctx
}

// Mock GET request with JSON
func MockJsonGet(ctx *gin.Context, params gin.Params, u url.Values) {
	ctx.Request.Method = "GET"
	ctx.Request.Header.Set("Content-Type", "application/json")

	// Set the params
	ctx.Params = params

	// Set the query params
	ctx.Request.URL.RawQuery = u.Encode()
}

// Mock POST request with JSON
func MockJsonPost(ctx *gin.Context, params gin.Params, u url.Values, content interface{}) {
	ctx.Request.Method = "POST"
	ctx.Request.Header.Set("Content-Type", "application/json")

	// Set the params
	ctx.Params = params

	// Set the query params
	ctx.Request.URL.RawQuery = u.Encode()

	// json
	jsonbytes, err := json.Marshal(content)
	if err != nil {
		panic(err)
	}

	ctx.Request.Body = io.NopCloser(bytes.NewReader(jsonbytes))
}
