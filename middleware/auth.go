package middleware

import (
	"net/http"
)

// BasicAuth is a middleware that provides basic authentication.
func BasicAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, pass, ok := r.BasicAuth()
		if !ok || !checkCredentials(user, pass) {
			w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func checkCredentials(user, pass string) bool {
	// Sample logic for checking credentials
	return user == "admin" && pass == "password"
}
