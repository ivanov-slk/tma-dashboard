package specifications

import (
	"strings"
	"testing"
)

type MetricsExtractor interface {
	ExtractMetrics() (string, error)
}

func ExtractMetricsSpecification(t testing.TB, extractor MetricsExtractor) {
	got, err := extractor.ExtractMetrics()
	expected := "Temperature is 15 degrees Celsius."

	if err != nil {
		t.Fatalf("failed to extract metrics: %s", err)
	}

	if !strings.Contains(got, expected) {
		t.Errorf("did not get correct response from server: %s should contain %s, but it did not.", got, expected)
	}
}
