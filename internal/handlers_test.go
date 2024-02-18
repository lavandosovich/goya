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
			name:    "negative test #1",
			request: "/update/counter/124124",
			want: want{
				statusCode:  http.StatusMethodNotAllowed,
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
