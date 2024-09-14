package app

import (
	"fmt"
	"net/http"
)

// PingHandler interacts with service healthchecks
func PingHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "pong")
}
