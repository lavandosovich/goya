package tests

import (
	"fmt"
	"github.com/lavandosovich/goya/internal"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetHandler(t *testing.T) {
	type want struct {
		statusCode         int
		response           string
		contentType        string
		gaugeMetricValue   internal.Gauge
		counterMetricValue internal.Counter
		metricName         string
	}
	tests := []struct {
		name    string
		want    want
		request string
	}{
		{
			name:    "positive test #1",
			request: "/update/gauge/nextgc",
			want: want{
				statusCode:       http.StatusOK,
				response:         fmt.Sprintf("%f", internal.Gauge(124124)),
				contentType:      "text/plain; charset=utf-8",
				gaugeMetricValue: internal.Gauge(124124),
				metricName:       "nextgc",
			},
		},
		{
			name:    "positive test #2",
			request: "/update/counter/pollcount",
			want: want{
				statusCode:         http.StatusOK,
				response:           fmt.Sprintf("%d", internal.Counter(12)),
				contentType:        "text/plain; charset=utf-8",
				counterMetricValue: internal.Counter(12),
				metricName:         "pollcount",
			},
		},
		{
			name:    "positive test from CI #3",
			request: "/update/counter/testCounter",
			want: want{
				statusCode:         http.StatusOK,
				response:           fmt.Sprintf("%d", internal.Counter(100)),
				contentType:        "text/plain; charset=utf-8",
				counterMetricValue: internal.Counter(100),
				metricName:         "testCounter",
			},
		},
		{
			name:    "positive empty test #1",
			request: "/update/celcius/124124",
			want: want{
				statusCode:  http.StatusNotImplemented,
				response:    "fail",
				contentType: "text/plain; charset=utf-8",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			memStorage := internal.NewMemStorage()
			switch {
			case tt.want.gaugeMetricValue != internal.Gauge(0):
				memStorage.SetGaugeMetric(tt.want.metricName, tt.want.gaugeMetricValue)
			case tt.want.counterMetricValue != internal.Counter(0):
				memStorage.SetCounterMetric(tt.want.metricName, tt.want.counterMetricValue)
			default:

			}

			ts := httptest.NewServer(internal.InitChiRouter(memStorage))
			defer ts.Close()

			statusCode, body := testRequest(t, ts, http.MethodGet, tt.request, func(response *http.Response) {
				assert.Equal(t, response.Header.Get("Content-Type"), tt.want.contentType)
			})
			assert.Equal(t, tt.want.statusCode, statusCode)
			assert.Equal(t, tt.want.response, body)

			//request := httptest.NewRequest(http.MethodPost, tt.request, nil)
			// создаём новый Recorder
			//w := httptest.NewRecorder()

			switch {
			case tt.want.gaugeMetricValue != internal.Gauge(0):
				assert.Equal(t, tt.want.gaugeMetricValue, memStorage.GetGaugeMetric(tt.want.metricName))
			case tt.want.counterMetricValue != internal.Counter(0):
				assert.Equal(t, tt.want.counterMetricValue, memStorage.GetCounterMetric(tt.want.metricName))
			default:

			}
		})
	}
}
