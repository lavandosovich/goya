package tests

import (
	"fmt"
	"github.com/lavandosovich/goya/internal"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestMemStorage_ReduceMetricsToHTML(t *testing.T) {
	tests := []struct {
		name                string
		metricName          string
		counterStorageValue internal.Counter
		gaugeStorageValue   internal.Gauge
	}{
		{
			name:                "positive case #1",
			metricName:          "totalGc",
			counterStorageValue: internal.Counter(1),
			gaugeStorageValue:   internal.Gauge(1.1),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			memStorage := internal.NewMemStorage()

			memStorage.SetCounterMetric(tt.metricName, tt.counterStorageValue)
			memStorage.SetGaugeMetric(tt.metricName, tt.gaugeStorageValue)
			body := memStorage.ReduceMetricsToHTML().String()

			assert.True(
				t,
				strings.Contains(
					body, fmt.Sprintf("<div>%s: %d</div>", tt.metricName, tt.counterStorageValue)))
			assert.True(
				t,
				strings.Contains(
					body, fmt.Sprintf("<div>%s: %f</div>", tt.metricName, tt.gaugeStorageValue)))
		})
	}
}
