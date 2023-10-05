package httpserver

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandler(t *testing.T) {
	req, _ := http.NewRequest(http.MethodGet, fmt.Sprint("/metrics"), nil)
	resp := httptest.NewRecorder()

	MetricsHandler(resp, req)

	if resp.Body.String() != "hello message" {
		t.Errorf("incorrect response from handler: got %s, want %s", resp.Body.String(), "hello message")
	}

}
