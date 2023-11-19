package main

import (
	"log"
	"net/http"

	dashboard_server "github.com/ivanov-slk/tma-dashboard/adapters/httpserver"
	"github.com/ivanov-slk/tma-dashboard/adapters/natsclient"
)

func main() {
	natsClient, err := natsclient.NewDashboardNATSClient()
	if err != nil {
		log.Fatal(err)
	}
	if err := http.ListenAndServe(":1337", &dashboard_server.DashboardServer{NATSClient: natsClient}); err != nil {
		log.Fatal(err)
	}
}
