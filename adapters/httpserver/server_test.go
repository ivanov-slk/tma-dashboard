package httpserver

import (
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	approvals "github.com/approvals/go-approval-tests"
)

func TestMetricsHandler(t *testing.T) {
	t.Run("temperature is 15", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "/metrics", nil)
		resp := httptest.NewRecorder()
		inputChan := make(chan []byte)
		inputMsg := "{\"temperature\":15,\"humidity\":0.6,\"pressure\":1000,\"datetime\":\"2024-01-04T16:27:40Z\",\"id\":\"1\"}"

		go func() { inputChan <- []byte(inputMsg) }()

		server := NewDashboardServer(inputChan)
		server.ServeHTTP(resp, req)

		approvals.VerifyString(t, resp.Body.String())
	})

	t.Run("temperature is 20", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "/metrics", nil)
		resp := httptest.NewRecorder()
		inputChan := make(chan []byte)
		inputMsg := "{\"temperature\":20,\"humidity\":0.6,\"pressure\":1000,\"datetime\":\"2024-01-04T16:27:40Z\",\"id\":\"1\"}"
		go func() { inputChan <- []byte(inputMsg) }()

		server := NewDashboardServer(inputChan)
		server.ServeHTTP(resp, req)

		approvals.VerifyString(t, resp.Body.String())
	})

	t.Run("server stashes last non-none message and returns it if the most recent is none", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "/metrics", nil)
		resp := httptest.NewRecorder()
		inputChan := make(chan []byte)
		inputMsg := "{\"temperature\":20,\"humidity\":0.6,\"pressure\":1000,\"datetime\":\"2024-01-04T16:27:40Z\",\"id\":\"1\"}"
		go func() { inputChan <- []byte(inputMsg); inputChan <- nil }()

		server := NewDashboardServer(inputChan)
		server.ServeHTTP(resp, req)

		approvals.VerifyString(t, resp.Body.String())

		req, _ = http.NewRequest(http.MethodGet, "/metrics", nil)
		resp = httptest.NewRecorder()
		server.ServeHTTP(resp, req)

		approvals.VerifyString(t, resp.Body.String())
	})

	t.Run("server stashes last non-none message and returns it if the most recent produces an error", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "/metrics", nil)
		resp := httptest.NewRecorder()
		inputChan := make(chan []byte)
		inputMsg := "{\"temperature\":20,\"humidity\":0.6,\"pressure\":1000,\"datetime\":\"2024-01-04T16:27:40Z\",\"id\":\"1\"}"
		inputMsgErr := "{this-is-unmarshalling-error}"
		go func() { inputChan <- []byte(inputMsg); inputChan <- []byte(inputMsgErr) }()

		server := NewDashboardServer(inputChan)
		server.ServeHTTP(resp, req)

		approvals.VerifyString(t, resp.Body.String())

		req, _ = http.NewRequest(http.MethodGet, "/metrics", nil)
		resp = httptest.NewRecorder()
		server.ServeHTTP(resp, req)

		approvals.VerifyString(t, resp.Body.String())
	})

	t.Run("server timeouts after 1 second and returns last message if no message is receved from channel", func(t *testing.T) {
		timeout := time.After(3 * time.Second)
		testOutput := make(chan httptest.ResponseRecorder)
		go func(t testing.TB) {
			req, _ := http.NewRequest(http.MethodGet, "/metrics", nil)
			resp := httptest.NewRecorder()
			inputChan := make(chan []byte)
			inputMsg := "{\"temperature\":25,\"humidity\":0.6,\"pressure\":1000,\"datetime\":\"2024-01-04T16:27:40Z\",\"id\":\"1\"}"
			go func() { inputChan <- []byte(inputMsg) }()

			server := NewDashboardServer(inputChan)
			server.ServeHTTP(resp, req)

			testOutput <- *resp

			req, _ = http.NewRequest(http.MethodGet, "/metrics", nil)
			resp = httptest.NewRecorder()
			server.ServeHTTP(resp, req)

			testOutput <- *resp
		}(t)

		for i := 0; i < 2; i++ {
			// Ugly, but need to fetch output twice before the timeout.
			select {
			case <-timeout:
				t.Fatal("The test didn't finish in time.")
			case resp := <-testOutput:

				approvals.VerifyString(t, resp.Body.String())
				slog.Info("here")
			}
		}
	})

	t.Run("should render zeroes when no messages found on channel", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "/metrics", nil)
		resp := httptest.NewRecorder()
		inputChan := make(chan []byte)

		server := NewDashboardServer(inputChan)
		server.ServeHTTP(resp, req)

		approvals.VerifyString(t, resp.Body.String())
	})

	t.Run("should render zeroes when no valid messages found on channel", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "/metrics", nil)
		resp := httptest.NewRecorder()
		inputChan := make(chan []byte)
		inputMsg := "{this-is-unmarshalling-error}"

		go func() { inputChan <- []byte(inputMsg) }()

		server := NewDashboardServer(inputChan)
		server.ServeHTTP(resp, req)

		approvals.VerifyString(t, resp.Body.String())
	})
}

func TestWelcomeHandler(t *testing.T) {
	t.Run("should display welcome page", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "/welcome", nil)
		resp := httptest.NewRecorder()
		inputChan := make(chan []byte)

		server := NewDashboardServer(inputChan)
		server.ServeHTTP(resp, req)

		approvals.VerifyString(t, resp.Body.String())
	})
}

func TestStaticHandler(t *testing.T) {
	t.Run("should load htmx", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "/static/htmx.js", nil)
		resp := httptest.NewRecorder()
		inputChan := make(chan []byte)

		server := NewDashboardServer(inputChan)
		server.ServeHTTP(resp, req)

		if !strings.Contains(resp.Body.String(), "(function(e,t)") {
			t.Errorf("Expected HTMX to be served, but got %s", resp.Body.String())
		}
	})
}
