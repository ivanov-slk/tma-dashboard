package httpserver

import (
	"fmt"
	"net/http"

	"github.com/ivanov-slk/tma-dashboard/adapters/natsclient"
)

type DashboardServer struct {
	// move outside at a later stage of refactoring, client and server should be completely uncoupled.
	NATSClient natsclient.DashboardClient
}

func (d *DashboardServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	messageData := d.NATSClient.FetchMessage()
	fmt.Fprint(w, messageData)
}
