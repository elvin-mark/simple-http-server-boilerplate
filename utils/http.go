package utils

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5/middleware"
)

func WriteJSON(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

func WriteJSONStatus(w http.ResponseWriter, data interface{}, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func LoggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		writer := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
		t0 := time.Now()

		defer func() {
			Logger.Info("Request",
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