package httpserver

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandler(t *testing.T) {
	req, _ := http.NewRequest(http.MethodGet, "/metrics", nil)
	resp := httptest.NewRecorder()
	inputChan := make(chan string)
	go func() { inputChan <- "stub message" }()

	server := &DashboardServer{InputChan: inputChan}
	server.ServeHTTP(resp, req)

	if resp.Body.String() != "stub message" {
		t.Errorf("incorrect response from handler: got %s, want %s", resp.Body.String(), "stub message")
	}

}
