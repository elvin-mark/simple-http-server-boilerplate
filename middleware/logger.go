package middleware

import (
	"http-server/utils"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5/middleware"
)

func LoggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		writer := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
		t0 := time.Now()

		defer func() {
			utils.Logger.Info("Request",
				"method", r.Method,
				"path", r.URL.Path,
				"status", writer.Status(),
				"latency", time.Since(t0),
				"request_id", middleware.GetReqID(r.Context()),
			)
		}()

		next.ServeHTTP(writer, r)
	})
}
