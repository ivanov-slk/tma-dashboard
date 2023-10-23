package httpserver

import (
	"fmt"
	"net/http"

	"github.com/ivanov-slk/tma-dashboard/adapters/natsclient"
)

func MetricsHandler(w http.ResponseWriter, r *http.Request) {
	messageData := natsclient.FetchMessage()
	fmt.Fprint(w, messageData)
}
