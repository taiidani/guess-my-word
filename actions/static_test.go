package actions

import (
	"strings"
	"testing"

	"net/http"
	"net/http/httptest"
)

func Test_StaticHandler(t *testing.T) {
	router := setupRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/assets/css/application.css", nil)
	router.ServeHTTP(w, req)

	if !strings.Contains(w.Body.String(), "text-decoration: underline;") {
		t.Error("Output did not contain expected phrase")
	}
}
