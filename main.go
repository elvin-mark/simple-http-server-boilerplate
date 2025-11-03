package main

import (
	"fmt"
	"http-server/config"
	"http-server/handlers"
	"http-server/middleware"
	"http-server/services"
	"http-server/storage"
	"http-server/utils"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		panic(fmt.Sprintf("failed to load configuration: %v", err))
	}

	// Initialize logger
	utils.InitLogger(cfg.LogLevel)

	// Initialize database
	db, err := storage.InitDB(&cfg.Database)
	if err != nil {
		utils.Logger.Error("Failed to initialize database", "error", err)
		os.Exit(1)
	}
	defer db.Close()

	// Create user repository, service, and handler
	userRepo := storage.NewUserRepository(db)
	userService := services.NewUserService(userRepo)
	userHandler := handlers.NewUserHandler(userService)

	// Create router
	r := chi.NewRouter()

	// ===== Middleware =====
	r.Use(chiMiddleware.RequestID)
	r.Use(middleware.LoggerMiddleware)
	r.Use(chiMiddleware.Recoverer)
	r.Use(middleware.CorsMiddleware())
	r.Use(middleware.RateLimiterMiddleware())
	r.Use(chiMiddleware.AllowContentType("application/json", "text/plain"))
	r.Use(chiMiddleware.Timeout(60 * time.Second))

	// ===== Routes =====
	r.Get("/", homeHandler)
	r.Get("/health", handlers.HealthCheckHandler)
	r.Handle("/metrics", handlers.MetricsHandler())

	r.Route("/users", func(r chi.Router) {
		r.Use(middleware.BasicAuth)
		r.Get("/", userHandler.GetUsersHandler)
		r.Post("/", userHandler.CreateUserHandler)
		r.Get("/{id}", userHandler.GetUserHandler)
		r.Delete("/{id}", userHandler.DeleteUserHandler)
	})

	// Start server
	serverAddr := fmt.Sprintf(":%d", cfg.Server.Port)
	fmt.Printf("✅ Server running on http://localhost%s\n", serverAddr)
	if err := http.ListenAndServe(serverAddr, r); err != nil {
		utils.Logger.Error("Server failed to start", "error", err)
		os.Exit(1)
	}
}

// ============== HANDLERS ==============

func homeHandler(w http.ResponseWriter, r *http.Request) {
	html := `
<!DOCTYPE html>
<html>
<head><title>Go-Chi HTTP Server</title></head>
<body>
<h1>🐹 Go-Chi Server</h1>
<p>Built with Chi, Resty, Go-Cache</p>
<p><a href="/health">Health Check</a></p>
</body>
</html>`
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(html))
}
