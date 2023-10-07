package main

import (
	"log"
	"net/http"

	"github.com/ivanov-slk/tma-dashboard/adapters/httpserver"
)

func main() {
	if err := http.ListenAndServe(":1337", http.HandlerFunc(httpserver.MetricsHandler)); err != nil {
		log.Fatal(err)
	}
}
