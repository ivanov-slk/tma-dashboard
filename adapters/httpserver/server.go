// Package httpserver provides functionalities for serving content to the user.
package httpserver

import (
	"fmt"
	"net/http"
)

// DashboardServer is the HTTP server serving the frontend-related content.
type DashboardServer struct {
	// move outside at a later stage of refactoring, client and server should be completely uncoupled.
	InputChan chan string
}

// ServeHTTP fetches the most recent message from the input channel of DashboardServer.
func (d *DashboardServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	messageData := <-d.InputChan
	fmt.Fprint(w, messageData)
}
