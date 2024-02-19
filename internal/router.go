package internal

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func InitChiRouter(memStorage *MemStorage) chi.Router {
	r := chi.NewRouter()
	r.Use(middleware.Recoverer)
	r.Use(middleware.Logger)

	r.Get("/", HandlerWrapper(memStorage, RootHandler))

	r.Route(
		"/update/{metricType}/{metricName}",
		func(r chi.Router) {
			r.Use(MetricTypeMiddleware)
			r.Post("/{metricValue}", HandlerWrapper(memStorage, PostHandler))
		},
	)
	r.Route(
		"/value/{metricType}/{metricName}",
		func(r chi.Router) {
			r.Use(MetricTypeMiddleware)
			r.Get("/", HandlerWrapper(memStorage, GetHandler))
		},
	)
	return r
}
