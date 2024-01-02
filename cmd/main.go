package main

import (
	"log"
	"net/http"

	"github.com/ivanov-slk/tma-dashboard/adapters/httpserver"
	"github.com/ivanov-slk/tma-dashboard/adapters/natsclient"
)

func main() {
	input_data := make(chan string)
	go func(c chan string) {
		natsClient, err := natsclient.NewDashboardNATSClient()
		if err != nil {
			log.Fatal(err)
		}
		c <- natsClient.FetchMessage()
	}(input_data)

	go func(c chan string) {
		// TODO: listen and serve with the raw input data channel for now. Later switch to a
		// presentation layer format as part of other functionality or service.
		if err := http.ListenAndServe(":1337", &httpserver.DashboardServer{InputChan: input_data}); err != nil {
			log.Fatal(err)
		}
	}(input_data)

	select {}
}
