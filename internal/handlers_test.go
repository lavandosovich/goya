package internal

import (
	"github.com/stretchr/testify/assert"
	"io"
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
	// создаём массив тестов: имя и желаемый результат
	tests := []struct {
		name    string
		want    want
		request string
	}{
		// определяем все тесты
		{
			name:    "positive test #1",
			request: "/update/gauge/nextgc/124124",
			want: want{
				statusCode:  http.StatusOK,
				response:    "ok",
				contentType: "application/text",
				metricValue: Gauge(124124),
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
				metricValue: Counter(12),
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
				metricValue: Counter(100),
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
				metricValue: Gauge(101),
				metricName:  "testGauge",
			},
		},
		{
			name:    "negative test #1",
			request: "/update/counter/124124",
			want: want{
				statusCode:  http.StatusNotFound,
				response:    "fail",
				contentType: "application/text",
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
				contentType: "application/text",
				metricValue: nil,
			},
		},
	}
	for _, tt := range tests {
		// запускаем каждый тест
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodPost, tt.request, nil)
			memStorage := NewMemStorage()
			// создаём новый Recorder
			w := httptest.NewRecorder()
			// определяем хендлер
			h := HandlerWrapper(memStorage, PostHandler)

			// запускаем сервер
			h.ServeHTTP(w, request)
			res := w.Result()

			// проверяем код ответа
			assert.Equal(t, tt.want.statusCode, res.StatusCode)

			// получаем и проверяем тело запроса
			defer res.Body.Close()
			resBody, err := io.ReadAll(res.Body)
			if err != nil {
				t.Fatal(err)
			}

			// заголовок ответа
			assert.Equal(t, tt.want.response, string(resBody))
			assert.Equal(t, res.Header.Get("Content-Type"), tt.want.contentType)
			assert.Equal(t, tt.want.metricValue, memStorage.GetMetric(tt.want.metricName))
		})
	}
}
