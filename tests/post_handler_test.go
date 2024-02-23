package tests

import (
	"github.com/lavandosovich/goya/internal"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestPostHandler(t *testing.T) {
	type want struct {
		statusCode  int
		response    string
		contentType string
		metricValue any
		metricName  string
	}
	tests := []struct {
		name    string
		want    want
		request string
	}{
		{
			name:    "positive test #1",
			request: "/update/gauge/nextgc/124124",
			want: want{
				statusCode:  http.StatusOK,
				response:    "ok",
				contentType: "application/text",
				metricValue: internal.Gauge(124124),
				metricName:  "nextgc",
			},
		},
		{
			name:    "positive test #2",
			request: "/update/counter/pollcount/12",
			want: want{
				statusCode:  http.StatusOK,
				response:    "ok",
				contentType: "application/text",
				metricValue: internal.Counter(12),
				metricName:  "pollcount",
			},
		},
		{
			name:    "positive test from CI #3",
			request: "/update/counter/testCounter/100",
			want: want{
				statusCode:  http.StatusOK,
				response:    "ok",
				contentType: "application/text",
				metricValue: internal.Counter(100),
				metricName:  "testCounter",
			},
		},
		{
			name:    "positive test from CI #4",
			request: "/update/gauge/testGauge/101",
			want: want{
				statusCode:  http.StatusOK,
				response:    "ok",
				contentType: "application/text",
				metricValue: internal.Gauge(101),
				metricName:  "testGauge",
			},
		},
		{
			name:    "negative test #1",
			request: "/update/counter/124124",
			want: want{
				statusCode:  http.StatusNotFound,
				response:    "404 page not found\n",
				contentType: "text/plain; charset=utf-8",
				metricValue: nil,
			},
		},
		{
			name:    "negative test from CI #2",
			request: "/update/counter/testCounter/none",
			want: want{
				statusCode:  http.StatusBadRequest,
				response:    "wrong type",
				contentType: "application/text",
				metricValue: nil,
			},
		},
		{
			name:    "negative test from CI #3",
			request: "/update/gauge/testGauge/none",
			want: want{
				statusCode:  http.StatusBadRequest,
				response:    "wrong type",
				contentType: "application/text",
				metricValue: nil,
			},
		},
		{
			name:    "negative test from CI #4",
			request: "/update/unknown/testCounter/100",
			want: want{
				statusCode:  http.StatusNotImplemented,
				response:    "fail",
				contentType: "text/plain; charset=utf-8",
				metricValue: nil,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			memStorage := internal.NewMemStorage()

			ts := httptest.NewServer(internal.InitChiRouter(memStorage))
			defer ts.Close()

			statusCode, body := testRequest(t, ts, http.MethodPost, tt.request, func(response *http.Response) {
				assert.Equal(t, response.Header.Get("Content-Type"), tt.want.contentType)
			})
			assert.Equal(t, tt.want.statusCode, statusCode)
			assert.Equal(t, tt.want.response, body)
			//request := httptest.NewRequest(http.MethodPost, tt.request, nil)
			// создаём новый Recorder
			//w := httptest.NewRecorder()
			switch tt.want.metricValue.(type) {
			case internal.Gauge:
				metric, _ := memStorage.GetGaugeMetric(tt.want.metricName)
				assert.Equal(t, tt.want.metricValue, metric)
			case internal.Counter:
				metric, _ := memStorage.GetCounterMetric(tt.want.metricName)
				assert.Equal(t, tt.want.metricValue, metric)
			default:
				if tt.want.metricValue != nil {
					assert.NotNil(t, nil)
				}
			}
		})
	}
}
