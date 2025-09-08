package main

import (
	"database/sql"
	"log"
	"mh-go/internal/handlers"
	"mh-go/internal/middleware"
	"mh-go/internal/services"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	_ "github.com/lib/pq"
)

func main() {
	// Database connection
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = "postgres://user:password@localhost:5433/mhgo?sslmode=disable"
	}

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Printf("Failed to connect to database: %v", err)
		return
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Println("Database is not reachable:", err)
		return
	}

	log.Println("Connected to database successfully")

	// JWT secret
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "your-super-secret-jwt-key-change-in-production"
	}

	// Initialize services
	authService := services.NewAuthService(db, jwtSecret)
	folderService := services.NewFolderService(db)
	contentService := services.NewContentService(db)
	userService := services.NewUserService(db)

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(authService)
	folderHandler := handlers.NewFolderHandler(folderService, contentService, userService)
	pageHandler := handlers.NewPageHandler()

	// Setup router
	r := setupRouter(authService, authHandler, folderHandler, pageHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s", port)
	if err := http.ListenAndServe(":"+port, r); err != nil {
		log.Printf("Server error: %v", err)
	}
}

func setupRouter(authService *services.AuthService, authHandler *handlers.AuthHandler, folderHandler *handlers.FolderHandler, pageHandler *handlers.PageHandler) *chi.Mux {
	r := chi.NewRouter()

	// Global middleware
	r.Use(chiMiddleware.Logger)
	r.Use(chiMiddleware.Recoverer)
	r.Use(middleware.CORS)
	r.Use(middleware.ContentType)

	// Static files
	r.Handle("/static/*", http.StripPrefix("/static/", http.FileServer(http.Dir("static/"))))

	// Auth routes (public)
	r.Get("/login", authHandler.LoginPage)
	r.Post("/login", authHandler.Login)
	r.Get("/register", authHandler.RegisterPage)
	r.Post("/register", authHandler.Register)
	r.Post("/logout", authHandler.Logout)

	// Protected routes
	r.Group(func(r chi.Router) {
		r.Use(middleware.RequireAuth(authService))

		// Overview (replaces dashboard)
		r.Get("/", folderHandler.OverviewPage)
		r.Get("/overview", folderHandler.OverviewPage)

		// Page routes
		r.Get("/flashcards", pageHandler.FlashcardsPage)
		r.Get("/settings", pageHandler.SettingsPage)

		// Register folder routes
		folderHandler.RegisterRoutes(r, authService)
	})

	return r
}
