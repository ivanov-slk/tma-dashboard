package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	if err := http.ListenAndServe(":1337", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "hello message")
	})); err != nil {
		log.Fatal(err)
	}
}
