// Package httpserver provides functionalities for serving content to the user.
package httpserver

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/ivanov-slk/tma-data-generator/pkg/generator"
)

// DashboardServer is the HTTP server serving the frontend-related content.
type DashboardServer struct {
	// move outside at a later stage of refactoring, client and server should be completely uncoupled.
	InputChan chan []byte
}

// ServeHTTP fetches the most recent message from the input channel of DashboardServer.
func (d *DashboardServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	messageData := <-d.InputChan
	temperatureStats := &generator.TemperatureStats{}
	err := json.Unmarshal([]byte(messageData), temperatureStats)
	if err != nil {
		fmt.Fprintf(w, "Error parsing the output: %s. Full message: %s", err, messageData)
	} else {
		fmt.Fprintf(w, "Temperature is %d degrees Celsius!", temperatureStats.Temperature)
	}
}
