package app

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log/slog"
	"net/http"
)

type errorBag struct {
	baseBag
	Message error
}

func errorResponse(w http.ResponseWriter, r *http.Request, code int, err error) {
	log := slog.With("code", code, "path", r.URL.Path)

	data := errorBag{
		Message: err,
	}

	log.Error(err.Error())
	if r.Header.Get("HX-Request") != "" {
		w.WriteHeader(code)
		fmt.Fprint(w, "Error: "+err.Error())
		return
	}

	renderHtml(w, code, "error.gohtml", data)
}

func renderHtml(w http.ResponseWriter, code int, file string, data any) {
	log := slog.With("name", file, "code", code)

	var t *template.Template
	var err error
	if dev {
		t, err = template.ParseGlob("app/templates/**")
	} else {
		t, err = template.ParseFS(templates, "templates/**")
	}
	if err != nil {
		log.Error("Could not parse templates", "error", err)
		return
	}

	log.Debug("Rendering file", "dev", dev)
	w.WriteHeader(code)
	err = t.ExecuteTemplate(w, file, data)
	if err != nil {
		log.Error("Could not render template", "error", err)
	}
}

func renderJson(w http.ResponseWriter, code int, data any) {
	log := slog.With("code", code)

	log.Debug("Rendering json", "dev", dev)
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(code)

	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		log.Error("Could not render template", "error", err)
	}
}

func renderRedirect(w http.ResponseWriter, code int, location string) {
	log := slog.With("code", code)

	log.Debug("Rendering redirect", "dev", dev)
	w.Header().Add("Location", location)
	w.WriteHeader(code)
}
