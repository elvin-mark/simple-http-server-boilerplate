package main

import (
	"fmt"
	"http-server/config"
	"http-server/handlers"
	"http-server/services"
	"http-server/storage"
	"http-server/utils"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		panic(fmt.Sprintf("failed to load configuration: %v", err))
	}

	// Initialize logger
	utils.Logger, err = zap.NewProduction()
	if err != nil {
		panic(err)
	}
	defer utils.Logger.Sync()

	// Initialize database
	db, err := storage.InitDB(&cfg.Database)
	if err != nil {
		utils.Logger.Fatal("Failed to initialize database", zap.Error(err))
	}
	defer db.Close()

	// Create user repository, service, and handler
	userRepo := storage.NewUserRepository(db)
	userService := services.NewUserService(userRepo)
	userHandler := handlers.NewUserHandler(userService)

	// Create router
	r := chi.NewRouter()

	// ===== Middleware =====
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.AllowContentType("application/json", "text/plain"))
	r.Use(middleware.Timeout(60 * time.Second))

	// ===== Routes =====
	r.Get("/", homeHandler)
	r.Get("/health", handlers.HealthCheckHandler)

	r.Route("/users", func(r chi.Router) {
		r.Get("/", userHandler.GetUsersHandler)
		r.Post("/", userHandler.CreateUserHandler)
		r.Get("/{id}", userHandler.GetUserHandler)
		r.Delete("/{id}", userHandler.DeleteUserHandler)
	})

	utils.PrintStartupInfo()

	// Start server
	serverAddr := fmt.Sprintf(":%d", cfg.Server.Port)
	fmt.Printf("✅ Server running on http://localhost%s\n", serverAddr)
	if err := http.ListenAndServe(serverAddr, r); err != nil {
		utils.Logger.Fatal("Server failed to start", zap.Error(err))
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
<p>Built with Chi, Resty, Go-Cache, Zap</p>
<p><a href="/health">Health Check</a></p>
</body>
</html>`
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(html))
}
