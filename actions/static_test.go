package actions

import (
	"strings"
	"testing"

	"net/http"
	"net/http/httptest"
)

func Test_StaticHandler(t *testing.T) {
	router := setupRouter()

	t.Run("CSS", func(tt *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/assets/css/application.css", nil)
		router.ServeHTTP(w, req)

		if !strings.Contains(w.Body.String(), "text-decoration: underline;") {
			t.Error("Output did not contain expected phrase")
		}
	})

	t.Run("JS", func(tt *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/assets/js/application.js", nil)
		router.ServeHTTP(w, req)

		if !strings.Contains(w.Body.String(), "return") {
			tt.Error("Output did not contain expected phrase")
		}
	})

	t.Run("SVG", func(tt *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/assets/images/between.svg", nil)
		router.ServeHTTP(w, req)

		if !strings.Contains(w.Body.String(), "xmlns") {
			tt.Error("Output did not contain expected phrase")
		}
	})

	t.Run("PNG", func(tt *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/assets/images/favicon.png", nil)
		router.ServeHTTP(w, req)

		if !strings.Contains(w.Body.String(), "PNG") {
			tt.Error("Output did not contain expected phrase")
		}
	})

	t.Run("Error on valid directory, not file", func(tt *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/assets/images/", nil)
		router.ServeHTTP(w, req)

		if w.Code != 400 {
			tt.Error("Request should have 400ed, instead received ", w.Code)
		}
	})

	t.Run("Error on missing file", func(tt *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/assets/images/blackhole.txt", nil)
		router.ServeHTTP(w, req)

		if w.Code != 404 {
			tt.Error("Request should have 404ed, instead received ", w.Code)
		}
	})
}
