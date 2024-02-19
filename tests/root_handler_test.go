package tests

import (
	"fmt"
	"github.com/lavandosovich/goya/internal"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestRootHandler(t *testing.T) {
	type want struct {
		statusCode  int
		contentType string
		metricName  string
	}
	tests := []struct {
		name               string
		want               want
		request            string
		gaugeMetricValue   internal.Gauge
		counterMetricValue internal.Counter
	}{
		{
			name:               "positive test #1",
			request:            "/",
			gaugeMetricValue:   internal.Gauge(6.1),
			counterMetricValue: internal.Counter(12),
			want: want{
				statusCode:  http.StatusOK,
				contentType: "text/html; charset=utf-8",
				metricName:  "nextgc",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			memStorage := internal.NewMemStorage()

			memStorage.SetGaugeMetric(tt.want.metricName, tt.gaugeMetricValue)
			memStorage.SetCounterMetric(tt.want.metricName, tt.counterMetricValue)

			ts := httptest.NewServer(internal.InitChiRouter(memStorage))
			defer ts.Close()

			statusCode, body := testRequest(t, ts, http.MethodGet, tt.request, func(response *http.Response) {
				assert.Equal(t, response.Header.Get("Content-Type"), tt.want.contentType)
			})
			assert.Equal(t, tt.want.statusCode, statusCode)

			assert.True(
				t,
				strings.Contains(
					body, fmt.Sprintf("<div>%s: %d</div>", tt.want.metricName, tt.counterMetricValue)))
			assert.True(
				t,
				strings.Contains(
					body, fmt.Sprintf("<div>%s: %f</div>", tt.want.metricName, tt.gaugeMetricValue)))
		})
	}
}

func TestEmptyStorageRootHandler(t *testing.T) {
	type want struct {
		statusCode  int
		contentType string
		metricName  string
	}
	tests := []struct {
		name    string
		want    want
		request string
	}{
		{
			name:    "positive test #1",
			request: "/",
			want: want{
				statusCode:  http.StatusOK,
				contentType: "text/html; charset=utf-8",
				metricName:  "nextgc",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			memStorage := internal.NewMemStorage()

			ts := httptest.NewServer(internal.InitChiRouter(memStorage))
			defer ts.Close()

			statusCode, body := testRequest(t, ts, http.MethodGet, tt.request, func(response *http.Response) {
				assert.Equal(t, response.Header.Get("Content-Type"), tt.want.contentType)
			})
			assert.Equal(t, tt.want.statusCode, statusCode)

			assert.Equal(t, "", body)
		})
	}
}
