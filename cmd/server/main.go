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
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	// Load environment variables from .env file
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found or error loading .env file:", err)
	}

	// Load configuration
	cfg := config.Load()

	// Initialize databasee
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
	openaiService := services.NewOpenAIService(cfg.OpenAIKey)

	// Learning platform services
	topicService := services.NewTopicService(db)
	methodService := services.NewLearningMethodService(db)
	progressService := services.NewProgressService(db)
	learningService := services.NewLearningService(db, topicService, methodService, progressService)

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(authService)
	folderHandler := handlers.NewFolderHandler(folderService, contentService, userService, openaiService)
	contentHandler := handlers.NewContentHandler(contentService, folderService, userService, openaiService)
	learningHandler := handlers.NewLearningHandler(topicService, methodService, progressService, learningService, userService)
	pageHandler := handlers.NewPageHandler(userService)

	// Setup router
	r := setupRouter(cfg, authService, authHandler, folderHandler, contentHandler, learningHandler, pageHandler)

	log.Printf("Server starting on port %s", cfg.Port)
	log.Printf("Server failed to start: %v", http.ListenAndServe(":"+cfg.Port, r))
}

func setupRouter(cfg *config.Config, authService *services.AuthService, authHandler *handlers.AuthHandler, folderHandler *handlers.FolderHandler, contentHandler *handlers.ContentHandler, learningHandler *handlers.LearningHandler, pageHandler *handlers.PageHandler) *chi.Mux {
	r := chi.NewRouter()

	// Global middleware
	r.Use(chiMiddleware.Logger)
	r.Use(chiMiddleware.Recoverer)
	r.Use(middleware.CORS)
	r.Use(middleware.ContentType)

	// Static files
	r.Handle("/static/*", http.StripPrefix("/static/", http.FileServer(http.Dir("static/"))))

	// Health check endpoint (public - no auth)
	r.Get("/api/health", handlers.HealthCheckHandler)

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

		// Overview (replaces dashboard)
		r.Get("/", folderHandler.OverviewPage)
		r.Get("/overview", folderHandler.OverviewPage)

		// Learning Platform
		r.Get("/learning", learningHandler.LearningPlatformPage)
		r.Post("/learning/start", learningHandler.StartLearning)
		r.Get("/learning/session/{sessionID}", learningHandler.LearningSessionPage)

		// New learning-focused routes matching architect template
		r.Get("/remember", pageHandler.RememberPage)
		r.Get("/optimize", pageHandler.OptimizePage)

		// Keep existing folder functionality
		r.Get("/folders", folderHandler.FoldersPage)
		r.Get("/folders/{folderID}", folderHandler.FolderPage)

		// Other pages
		r.Get("/settings", pageHandler.SettingsPage)

		// API routes
		r.Route("/api", func(r chi.Router) {
			r.Use(middleware.RequireAuth(authService))

				r.Route("/folders", func(r chi.Router) {
					// Folder routes
					r.Get("/", folderHandler.GetFolders)
					r.Post("/", folderHandler.CreateFolder)
					r.Get("/{folderID}", folderHandler.GetFolder)
					r.Put("/{folderID}", folderHandler.UpdateFolder)
					r.Delete("/{folderID}", folderHandler.DeleteFolder)
					r.Get("/{folderID}/path", folderHandler.GetFolderPath)

					// Content routes
					r.Get("/{folderID}/content", folderHandler.GetContent)
					r.Post("/{folderID}/content", folderHandler.CreateContent)
					r.Put("/{folderID}/content/{contentID}", folderHandler.UpdateContent)
					r.Delete("/{folderID}/content/{contentID}", folderHandler.DeleteContent)
				})

				// Content routes
				r.Route("/content", func(r chi.Router) {
					r.Get("/search", folderHandler.SearchContent)
					r.Post("/", contentHandler.CreateContent)
					r.Get("/{contentID}", contentHandler.GetContent)
					r.Get("/{contentID}/edit", contentHandler.GetEditContentModal)
					r.Put("/{contentID}", contentHandler.UpdateContent)
					r.Delete("/{contentID}", contentHandler.DeleteContent)
				})

				// Learning platform API routes
				r.Route("/learning", func(r chi.Router) {
					r.Get("/topics", learningHandler.GetTopics)
					r.Get("/topics/search", learningHandler.SearchTopics)
					r.Get("/pathways", learningHandler.GetLearningPathways)
					r.Get("/stats", learningHandler.GetDashboardStats)
					r.Get("/sessions/active", learningHandler.GetActiveLearningSessions)
					r.Get("/prompts", learningHandler.GetUserPrompts)
					r.Get("/preferences", learningHandler.GetUserPreferences)
					r.Post("/preferences", learningHandler.UpdateUserPreferences)
					r.Post("/prompts", learningHandler.CreateLearningPrompt)
					r.Post("/progress", learningHandler.UpdateProgress)
					r.Post("/sessions", learningHandler.StartLearningSession)
				r.Put("/sessions/{sessionID}", learningHandler.EndLearningSession)
			})
		})
	})

	return r
}
