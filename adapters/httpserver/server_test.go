package httpserver

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandler(t *testing.T) {
	t.Run("temperature is 15", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "/metrics", nil)
		resp := httptest.NewRecorder()
		inputChan := make(chan []byte)
		inputMsg := "{\"temperature\":15,\"humidity\":0.6,\"pressure\":1000,\"datetime\":\"2024-01-04T16:27:40Z\",\"id\":\"1\"}"
		expectedResp := "Temperature is 15 degrees Celsius!"
		go func() { inputChan <- []byte(inputMsg) }()

		server := &DashboardServer{InputChan: inputChan}
		server.ServeHTTP(resp, req)

		if resp.Body.String() != expectedResp {
			t.Errorf("incorrect response from handler: got %s, want %s", resp.Body.String(), expectedResp)
		}
	})

	t.Run("temperature is 20", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "/metrics", nil)
		resp := httptest.NewRecorder()
		inputChan := make(chan []byte)
		inputMsg := "{\"temperature\":20,\"humidity\":0.6,\"pressure\":1000,\"datetime\":\"2024-01-04T16:27:40Z\",\"id\":\"1\"}"
		expectedResp := "Temperature is 20 degrees Celsius!"
		go func() { inputChan <- []byte(inputMsg) }()

		server := &DashboardServer{InputChan: inputChan}
		server.ServeHTTP(resp, req)

		if resp.Body.String() != expectedResp {
			t.Errorf("incorrect response from handler: got %s, want %s", resp.Body.String(), expectedResp)
		}
	})

	t.Run("server stashes last non-none message and returns it if the most recent is none", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "/metrics", nil)
		resp := httptest.NewRecorder()
		inputChan := make(chan []byte)
		inputMsg := "{\"temperature\":20,\"humidity\":0.6,\"pressure\":1000,\"datetime\":\"2024-01-04T16:27:40Z\",\"id\":\"1\"}"
		expectedResp := "Temperature is 20 degrees Celsius!"
		go func() { inputChan <- []byte(inputMsg); inputChan <- nil }()

		server := &DashboardServer{InputChan: inputChan}
		server.ServeHTTP(resp, req)

		if resp.Body.String() != expectedResp {
			t.Fatalf("incorrect response from handler: got %s, want %s", resp.Body.String(), expectedResp)
		}

		req, _ = http.NewRequest(http.MethodGet, "/metrics", nil)
		resp = httptest.NewRecorder()
		server.ServeHTTP(resp, req)

		if resp.Body.String() != expectedResp {
			t.Errorf("incorrect response from handler: got %s, want %s", resp.Body.String(), expectedResp)
		}
	})

	t.Run("server stashes last non-none message and returns it if the most recent produces an error", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "/metrics", nil)
		resp := httptest.NewRecorder()
		inputChan := make(chan []byte)
		inputMsg := "{\"temperature\":20,\"humidity\":0.6,\"pressure\":1000,\"datetime\":\"2024-01-04T16:27:40Z\",\"id\":\"1\"}"
		inputMsgErr := "{this-is-unmarshalling-error}"
		expectedResp := "Temperature is 20 degrees Celsius!"
		go func() { inputChan <- []byte(inputMsg); inputChan <- []byte(inputMsgErr) }()

		server := &DashboardServer{InputChan: inputChan}
		server.ServeHTTP(resp, req)

		if resp.Body.String() != expectedResp {
			t.Fatalf("incorrect response from handler: got %s, want %s", resp.Body.String(), expectedResp)
		}

		req, _ = http.NewRequest(http.MethodGet, "/metrics", nil)
		resp = httptest.NewRecorder()
		server.ServeHTTP(resp, req)

		if resp.Body.String() != expectedResp {
			t.Errorf("incorrect response from handler: got %s, want %s", resp.Body.String(), expectedResp)
		}
	})
}
