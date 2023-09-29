package httpserver

import (
	"io"
	"net/http"
)

type Driver struct {
	BaseURL string
	Client  *http.Client
}

func (d Driver) ExtractMetrics() (string, error) {
	res, err := d.Client.Get(d.BaseURL + "/metrics")
	if err != nil {
		return "", err
	}

	defer res.Body.Close()

	metrics, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	return string(metrics), nil
}
