// Package httpserver provides functionalities for serving content to the user.
package httpserver

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/ivanov-slk/tma-data-generator/pkg/generator"
)

// DashboardServer is the HTTP server serving the frontend-related content.
type DashboardServer struct {
	InputChan chan []byte
}

// ServeHTTP fetches the most recent message from the input channel of DashboardServer.
func (d *DashboardServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	messageData := <-d.InputChan
	slog.Info("Message fetched from channel:", "messageData", messageData) // TODO align logging - slog or log or fmt ...
	temperatureStats := &generator.TemperatureStats{}
	err := json.Unmarshal([]byte(messageData), temperatureStats)
	slog.Info("JSON parsing done.")
	if err != nil {
		fmt.Fprintf(w, "Error parsing the output: %s. Full message: %s", err, messageData)
	} else {
		fmt.Fprintf(w, "Temperature is %d degrees Celsius!", temperatureStats.Temperature)
	}
}
