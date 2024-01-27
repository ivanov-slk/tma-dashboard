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
	InputChan       chan []byte
	lastFetchedData []byte
}

// ServeHTTP fetches the most recent message from the input channel of DashboardServer.
func (d *DashboardServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	messageData := <-d.InputChan
	slog.Info("Message fetched from channel:", "messageData", messageData) // TODO align logging - slog or log or fmt ...

	if messageData == nil {
		messageData = d.lastFetchedData
	}

	temperatureStats := &generator.TemperatureStats{}
	err := json.Unmarshal([]byte(messageData), temperatureStats)
	slog.Info("JSON parsing done.")
	if err != nil {
		slog.Warn("Message parsing resulted in an error:", "error", err, "message", messageData)
		temperatureStats = &generator.TemperatureStats{}
		json.Unmarshal([]byte(d.lastFetchedData), temperatureStats)
		fmt.Fprintf(w, "Temperature is %d degrees Celsius!", temperatureStats.Temperature)
	} else {
		d.lastFetchedData = messageData
		fmt.Fprintf(w, "Temperature is %d degrees Celsius!", temperatureStats.Temperature)
	}
}
