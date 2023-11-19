package httpserver

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

type StubNATSClient struct {
}

func (c *StubNATSClient) FetchMessage() string {
	return "stub message"
}

func TestHandler(t *testing.T) {
	req, _ := http.NewRequest(http.MethodGet, "/metrics", nil)
	resp := httptest.NewRecorder()

	server := &DashboardServer{NATSClient: &StubNATSClient{}}
	server.ServeHTTP(resp, req)

	if resp.Body.String() != "stub message" {
		t.Errorf("incorrect response from handler: got %s, want %s", resp.Body.String(), "stub message")
	}

}