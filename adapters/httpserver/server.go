// Package httpserver provides functionalities for serving content to the user.
package httpserver

import (
	"embed"
	"encoding/json"
	"html/template"
	"log/slog"
	"net/http"
	"time"

	"github.com/ivanov-slk/tma-data-generator/pkg/generator"
)

var (
	//go:embed "templates/*"
	dashboardTemplates embed.FS
)

// DashboardServer is the HTTP server serving the frontend-related content.
type DashboardServer struct {
	InputChan       chan []byte
	lastFetchedData []byte
}

// ServeHTTP fetches the most recent message from the input channel of DashboardServer.
func (d *DashboardServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	messageData := d.fetchMessage()
	temperatureStats := d.parseMessage(messageData)

	templ, _ := template.ParseFS(dashboardTemplates, "templates/*.gohtml")
	err := templ.ExecuteTemplate(w, "main.gohtml", temperatureStats)
	if err != nil {
		slog.Error("Error parsing the templates.", "error", err)
	}
}

func (d *DashboardServer) fetchMessage() []byte {
	messageData := d.lastFetchedData
	select {
	case mes := <-d.InputChan:
		messageData = mes
		slog.Info("Message fetched from channel:", "messageData", messageData) // TODO align logging - slog or log or fmt ...
	case <-time.After(1 * time.Second):
		slog.Info("Timed out waiting for message. Using the last valid one.")
	}
	return messageData
}

func (d *DashboardServer) parseMessage(messageData []byte) *generator.TemperatureStats {
	temperatureStats := &generator.TemperatureStats{}
	err := json.Unmarshal([]byte(messageData), temperatureStats)
	slog.Info("JSON parsing done.")
	if err != nil {
		slog.Warn("Message parsing resulted in an error:", "error", err, "message", messageData)
		temperatureStats = &generator.TemperatureStats{}
		json.Unmarshal([]byte(d.lastFetchedData), temperatureStats)
	} else {
		d.lastFetchedData = messageData
	}
	return temperatureStats
}
