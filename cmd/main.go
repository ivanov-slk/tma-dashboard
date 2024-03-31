package main

import (
	"log"
	"net/http"

	"github.com/ivanov-slk/tma-dashboard/adapters/httpserver"
)

func main() {
	input_data := make(chan []byte)
	go func(c chan []byte) {
		// natsClient, err := natsclient.NewDashboardNATSClient()
		// if err != nil {
		// 	log.Fatalf("failed to initialize NATS client: %s", err)
		// }
		// natsClient.Consume(c)
		c <- []byte("{\"temperature\":20,\"humidity\":0.6,\"pressure\":1000,\"datetime\":\"2024-01-04T16:27:40Z\",\"id\":\"1\"}")
	}(input_data)

	go func(c chan []byte) {
		// TODO: listen and serve with the raw input data channel for now. Later switch to a
		// presentation layer format as part of other functionality or microservice.
		if err := http.ListenAndServe(":1337", httpserver.NewDashboardServer(input_data)); err != nil {
			log.Fatal(err)
		}
	}(input_data)

	select {}
}
