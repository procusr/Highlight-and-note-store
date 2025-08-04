package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"mh-go/internal/config"
	"mh-go/internal/database"
	"mh-go/internal/handlers"
	"mh-go/internal/middleware"
	"mh-go/internal/services"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Initialize database
	db, err := database.New(cfg.DatabaseURL)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	// Run migrations
	if err := database.RunMigrations(db, cfg.DatabaseURL); err != nil {
		log.Printf("Failed to run migrations: %v", err)
		// Don't exit here, continue with the application
	}

	// Initialize services
	authService := services.NewAuthService(db, cfg.JWTSecret)
	userService := services.NewUserService(db)
	folderService := services.NewFolderService(db)
	contentService := services.NewContentService(db)

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(authService)
	userHandler := handlers.NewUserHandler(userService)
	folderHandler := handlers.NewFolderHandler(folderService)
	contentHandler := handlers.NewContentHandler(contentService)
	staticHandler := handlers.NewStaticHandler()

	// Setup router
	r := chi.NewRouter()

	// Middleware
	r.Use(chimiddleware.Logger)
	r.Use(chimiddleware.Recoverer)
	r.Use(chimiddleware.RealIP)
	r.Use(chimiddleware.RequestID)
	r.Use(chimiddleware.Timeout(60 * time.Second))
	r.Use(middleware.CORS())

	// Static files
	r.Handle("/static/*", http.StripPrefix("/static/", http.FileServer(http.Dir("static/"))))

	// Public routes
	r.Group(func(r chi.Router) {
		r.Get("/", staticHandler.Index)
		r.Get("/login", authHandler.LoginPage)
		r.Post("/login", authHandler.Login)
		r.Get("/register", authHandler.RegisterPage)
		r.Post("/register", authHandler.Register)
		r.Get("/forgot-password", authHandler.ForgotPasswordPage)
		r.Post("/forgot-password", authHandler.ForgotPassword)
		r.Get("/reset-password", authHandler.ResetPasswordPage)
		r.Post("/reset-password", authHandler.ResetPassword)
	})

	// Protected routes
	r.Group(func(r chi.Router) {
		r.Use(middleware.AuthRequired(authService))

		r.Get("/dashboard", userHandler.Dashboard)
		r.Post("/logout", authHandler.Logout)

		// Folder routes
		r.Route("/folders", func(r chi.Router) {
			r.Get("/", folderHandler.List)
			r.Post("/", folderHandler.Create)
			r.Get("/{id}", folderHandler.Get)
			r.Put("/{id}", folderHandler.Update)
			r.Delete("/{id}", folderHandler.Delete)
		})

		// Content routes
		r.Route("/content", func(r chi.Router) {
			r.Get("/", contentHandler.List)
			r.Post("/", contentHandler.Create)
			r.Get("/{id}", contentHandler.Get)
			r.Put("/{id}", contentHandler.Update)
			r.Delete("/{id}", contentHandler.Delete)
			r.Post("/{id}/archive", contentHandler.Archive)
			r.Post("/{id}/favorite", contentHandler.ToggleFavorite)
		})
	})

	// Start server
	srv := &http.Server{
		Addr:    cfg.Port,
		Handler: r,
	}

	// Graceful shutdown
	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
		<-sigChan

		log.Println("Shutting down server...")
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		if err := srv.Shutdown(ctx); err != nil {
			log.Fatal("Server forced to shutdown:", err)
		}
	}()

	log.Printf("Server starting on %s", cfg.Port)
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Printf("Server failed to start: %v", err)
	}
}
