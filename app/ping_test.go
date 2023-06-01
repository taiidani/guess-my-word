package app

import (
	"testing"

	"net/http"
	"net/http/httptest"
)

func Test_PingHandler(t *testing.T) {
	router := setupRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/ping", nil)
	router.ServeHTTP(w, req)

	if w.Body.String() != "pong" {
		t.Error("Output did not contain pong")
	}
}
