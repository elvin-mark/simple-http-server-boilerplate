package middleware

import (
	"net"
	"net/http"
	"time"

	"github.com/go-chi/cors"
	"github.com/go-chi/httprate"
)

// CorsMiddleware sets up the CORS middleware.
func CorsMiddleware() func(http.Handler) http.Handler {
	return cors.New(cors.Options{
		AllowedOrigins:   []string{"*"}, // You might want to restrict this in production
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300, // Maximum value not ignored by any major browsers
	}).Handler
}

// RateLimiterMiddleware sets up the rate limiting middleware.
func RateLimiterMiddleware() func(http.Handler) http.Handler {
	// Limit requests to 100 per minute.
	return httprate.Limit(
		100,
		1*time.Minute,
		httprate.WithKeyFuncs(func(r *http.Request) (string, error) {
			ip, _, err := net.SplitHostPort(r.RemoteAddr)
			if err != nil {
				return "", err
			}
			return ip, nil
		}),
	)
}
