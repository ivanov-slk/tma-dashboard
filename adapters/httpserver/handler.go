package httpserver

import (
	"fmt"
	"net/http"
)

func MetricsHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "hello message")
}
