package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/MartialM1nd/freefsm/internal/config"
	"github.com/MartialM1nd/freefsm/internal/database"
	"github.com/MartialM1nd/freefsm/internal/handlers"
	"github.com/MartialM1nd/freefsm/internal/middleware"
	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"
)

func main() {
	// Flags
	migrateFlag := flag.Bool("migrate", false, "Run database migrations and exit")
	flag.Parse()

	// Load config
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Connect to database
	db, err := database.New(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Run migrations if flag set
	if *migrateFlag {
		log.Println("Running migrations...")
		if err := db.Migrate(context.Background()); err != nil {
			log.Fatalf("Migration failed: %v", err)
		}
		log.Println("Migrations complete")
		return
	}

	// Setup router
	r := chi.NewRouter()

	// Middleware
	r.Use(chimw.Logger)
	r.Use(chimw.Recoverer)
	r.Use(chimw.RealIP)
	r.Use(chimw.Timeout(60 * time.Second))

	// Static files
	fileServer := http.FileServer(http.Dir("ui/static"))
	r.Handle("/static/*", http.StripPrefix("/static/", fileServer))

	// Initialize handlers
	h := handlers.New(db, cfg)

	// Public routes
	r.Get("/login", h.LoginPage)
	r.Post("/login", h.Login)
	r.Post("/logout", h.Logout)

	// Protected routes
	r.Group(func(r chi.Router) {
		r.Use(middleware.Auth(db))

		// Dashboard
		r.Get("/", h.Dashboard)

		// Customers
		r.Get("/customers", h.CustomersList)
		r.Get("/customers/new", h.CustomersNew)
		r.Post("/customers", h.CustomersCreate)
		r.Get("/customers/{id}", h.CustomersView)
		r.Get("/customers/{id}/edit", h.CustomersEdit)
		r.Put("/customers/{id}", h.CustomersUpdate)
		r.Delete("/customers/{id}", h.CustomersDelete)

		// Workers
		r.Get("/workers", h.WorkersList)
		r.Get("/workers/new", h.WorkersNew)
		r.Post("/workers", h.WorkersCreate)
		r.Get("/workers/{id}", h.WorkersView)
		r.Get("/workers/{id}/edit", h.WorkersEdit)
		r.Put("/workers/{id}", h.WorkersUpdate)
		r.Delete("/workers/{id}", h.WorkersDelete)

		// Jobs
		r.Get("/jobs", h.JobsList)
		r.Get("/jobs/new", h.JobsNew)
		r.Post("/jobs", h.JobsCreate)
		r.Get("/jobs/{id}", h.JobsView)
		r.Get("/jobs/{id}/edit", h.JobsEdit)
		r.Put("/jobs/{id}", h.JobsUpdate)
		r.Delete("/jobs/{id}", h.JobsDelete)
		r.Post("/jobs/{id}/notes", h.JobsAddNote)
		r.Put("/jobs/{id}/status", h.JobsUpdateStatus)
	})

	// Start server
	srv := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Graceful shutdown
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		log.Printf("Starting server on :%s", cfg.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()

	<-done
	log.Println("Shutting down...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Shutdown error: %v", err)
	}

	log.Println("Server stopped")
}
