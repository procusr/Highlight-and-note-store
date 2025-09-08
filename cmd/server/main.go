package main

import (
	"log"
	"net/http"

	"mh-go/internal/config"
	"mh-go/internal/database"
	"mh-go/internal/handlers"
	"mh-go/internal/middleware"
	"mh-go/internal/services"

	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	_ "github.com/lib/pq"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Initialize database
	db, err := database.Connect(cfg.DatabaseURL)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer func() {
		if cerr := database.Close(db); cerr != nil {
			log.Printf("Error closing database: %v", cerr)
		}
	}()

	// Run database migrations if needed
	if err := database.Migrate(cfg.DatabaseURL); err != nil {
		log.Println("Warning: Migration failed:", err)
	}

	log.Println("Connected to database successfully")

	// Initialize services
	authService := services.NewAuthService(db, cfg.JWTSecret)
	userService := services.NewUserService(db)
	folderService := services.NewFolderService(db)
	contentService := services.NewContentService(db)

	// Initialize services for learning
	topicService := services.NewTopicService(db)
	methodService := services.NewLearningMethodService(db)
	progressService := services.NewProgressService(db)
	learningService := services.NewLearningService(db, topicService, methodService, progressService)

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(authService)
	folderHandler := handlers.NewFolderHandler(folderService, contentService, userService)
	contentHandler := handlers.NewContentHandler(contentService, folderService)
	learningHandler := handlers.NewLearningHandler(topicService, methodService, progressService, learningService, userService)

	// Setup router
	r := setupRouter(cfg, authService, authHandler, folderHandler, contentHandler, learningHandler)

	log.Printf("Server starting on port %s", cfg.Port)
	log.Printf("Server failed to start: %v", http.ListenAndServe(":"+cfg.Port, r))
}

func setupRouter(cfg *config.Config, authService *services.AuthService, authHandler *handlers.AuthHandler, folderHandler *handlers.FolderHandler, contentHandler *handlers.ContentHandler, learningHandler *handlers.LearningHandler) *chi.Mux {
	r := chi.NewRouter()

	// Global middleware
	r.Use(chiMiddleware.Logger)
	r.Use(chiMiddleware.Recoverer)
	r.Use(middleware.CORS)
	r.Use(middleware.ContentType)

	// Static files
	r.Handle("/static/*", http.StripPrefix("/static/", http.FileServer(http.Dir("static/"))))

	// Auth routes (public)
	r.Get("/", authHandler.LoginPage)
	r.Get("/login", authHandler.LoginPage)
	r.Post("/login", authHandler.Login)
	r.Get("/register", authHandler.RegisterPage)
	r.Post("/register", authHandler.Register)
	r.Post("/logout", authHandler.Logout)

	// Protected routes
	r.Group(func(r chi.Router) {
		r.Use(middleware.RequireAuth(authService))

		// Dashboard
		r.Get("/dashboard", folderHandler.DashboardPage)

		// Learning page
		r.Get("/learning", learningHandler.LearningPlatformPage)

		// Folder pages
		r.Get("/folders", folderHandler.FoldersPage)
		r.Get("/folders/{folderID}", folderHandler.FolderPage)

		// API routes
		r.Route("/api", func(r chi.Router) {
			// Folder routes
			folderHandler.RegisterRoutes(r, authService)

			// Content routes
			contentHandler.RegisterRoutes(r, authService)

			// Learning routes
			r.Route("/learning", func(r chi.Router) {
				r.Get("/pathways", learningHandler.GetLearningPathways)
				r.Get("/topics", learningHandler.GetTopics)
				r.Get("/topics/search", learningHandler.SearchTopics)
				r.Post("/start", learningHandler.StartLearning)
				r.Post("/prompts", learningHandler.CreateLearningPrompt)
				r.Get("/prompts", learningHandler.GetUserPrompts)
				r.Get("/stats", learningHandler.GetDashboardStats)
				r.Post("/progress", learningHandler.UpdateProgress)
				r.Post("/sessions", learningHandler.StartLearningSession)
				r.Put("/sessions/{sessionID}", learningHandler.EndLearningSession)
				r.Get("/preferences", learningHandler.GetUserPreferences)
				r.Put("/preferences", learningHandler.UpdateUserPreferences)
			})
		})
	})

	return r
}
