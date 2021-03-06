package actions

import (
	"strings"
	"testing"

	"net/http"
	"net/http/httptest"
)

func Test_HomeHandler(t *testing.T) {
	router := setupRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/", nil)
	router.ServeHTTP(w, req)

	if !strings.Contains(w.Body.String(), "Guess My Word") {
		t.Error("Output did not contain expected phrase")
	}
}
