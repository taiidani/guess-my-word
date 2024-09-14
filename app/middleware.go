package app

import (
	"net/http"
)

func standardHeadersMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if dev {
			w.Header().Set("Access-Control-Allow-Origin", "*") // Let everyone in in dev mode only!
		} else {
			w.Header().Set("Strict-Transport-Security", "max-age=63072000") // 2 years
		}

		next.ServeHTTP(w, r)
	})
}
