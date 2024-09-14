package app

import (
	"net/http"
)

type errorBag struct {
	baseBag
	Message error
}

func errorResponse(w http.ResponseWriter, code int, err error) {
	data := errorBag{
		Message: err,
	}

	renderHtml(w, code, "error.gohtml", data)
}
