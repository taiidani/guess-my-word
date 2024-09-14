package app

import (
	"context"
	"guess_my_word/internal/sessions"
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

func sessionMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if dev {
			w.Header().Set("Access-Control-Allow-Origin", "*") // Let everyone in in dev mode only!
		} else {
			w.Header().Set("Strict-Transport-Security", "max-age=63072000") // 2 years
		}

		next.ServeHTTP(w, r)
	})
}

func getSession(ctx context.Context) *sessions.Session {
	sess := ctx.Value("session")
	if sess == nil {
		return nil
	}

	return sess.(*sessions.Session)
}
