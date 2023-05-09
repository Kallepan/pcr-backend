package main

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gitlab.com/kaka/pcr-backend/common/controllers"
)

func SetUpRouter() *gin.Engine {
	router := gin.Default()
	return router
}

// TestPing test ping
func TestPing(t *testing.T) {
	r := SetUpRouter()
	r.GET("/api/v1", controllers.Ping)
	req, _ := http.NewRequest("GET", "/api/v1", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	expectedBody := `{"message":"pong"}`

	responseData, _ := ioutil.ReadAll(w.Body)
	assert.Equal(t, expectedBody, string(responseData))
	assert.Equal(t, 200, w.Code)
}
