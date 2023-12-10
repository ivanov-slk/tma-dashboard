// Package httpserver provides functionalities for serving content to the user.
package httpserver

import (
	"fmt"
	"net/http"

	"github.com/ivanov-slk/tma-dashboard/adapters/natsclient"
)

// DashboardServer is the HTTP server serving the frontend-related content.
type DashboardServer struct {
	// move outside at a later stage of refactoring, client and server should be completely uncoupled.
	NATSClient natsclient.DashboardClient
}

func (d *DashboardServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	messageData := d.NATSClient.FetchMessage()
	fmt.Fprint(w, messageData)
}
