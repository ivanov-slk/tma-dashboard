package specifications

import "testing"

type MetricsExtractor interface {
	ExtractMetrics() (string, error)
}

func ExtractMetricsSpecification(t testing.TB, extractor MetricsExtractor) {
	got, err := extractor.ExtractMetrics()

	if err != nil {
		t.Fatalf("failed to extract metrics: %s", err)
	}

	if got != "Temperature is 15 degrees Celsius!" {
		t.Errorf("did not get correct response from server: got %s, want %s", got, "Temperature is 15 degrees Celsius!")
	}
}
