package main

import (
	"context"
	"encoding/base64"
	"log"
	"net/http"
	"time"

	"backend/internal/config"
	"backend/internal/database"
	"backend/internal/handler"
	"backend/internal/middleware"
	"backend/internal/repository"
	"backend/internal/service"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

func main() {
	cfg := config.Load()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	db, err := database.NewPostgresConnection(cfg.DBURL)
	if err != nil {
		log.Fatalf("postgres connection failed: %v", err)
	}
	defer db.Close()

	rdb, err := database.NewRedisClient(cfg.REDIS_URL)
	if err != nil {
		log.Fatalf("redis connection failed: %v", err)
	}
	defer rdb.Close()

	key, err := base64.StdEncoding.DecodeString(cfg.MasterKey)
	if err != nil || len(key) != 32 {
		log.Fatal("MASTER_KEY must be valid base64-encoded 32 bytes")
	}

	userRepo := repository.NewUserRepository(db)
	apiKeyRepo := repository.NewApiKeyRepository(db)
	walletRepo := repository.NewWalletRepository(db)

	// Initialize Email Service - use SMTP if configured, otherwise use SendGrid or console fallback
	var emailService *service.EmailService
	if cfg.SMTPHost != "" && cfg.SMTPUsername != "" && cfg.SMTPPassword != "" {
		// Use SMTP for real email sending
		emailService = service.NewEmailServiceWithSMTP(cfg.SMTPHost, cfg.SMTPPort, cfg.SMTPUsername, cfg.SMTPPassword, cfg.FromEmail)
		log.Println("📧 Email Service initialized with SMTP:", cfg.SMTPHost)
	} else if cfg.SendGridAPIKey != "" {
		// Use SendGrid as fallback
		emailService = service.NewEmailService(cfg.SendGridAPIKey, cfg.SendGridEmail)
		log.Println("📧 Email Service initialized with SendGrid")
	} else {
		// Use console fallback for development
		emailService = service.NewEmailService("", cfg.FromEmail)
		log.Println("⚠️  Email Service in FALLBACK mode (console logging) - configure SMTP_HOST or SENDGRID_API_KEY for production")
	}

	authService := service.NewAuthService(userRepo, rdb, cfg.JWTSecret, emailService)
	apiKeyService := service.NewApiKeyService(apiKeyRepo)

	walletService, err := service.NewWalletService(walletRepo, key, cfg.EthRPCURL)
	if err != nil {
		log.Fatalf("wallet service init failed: %v", err)
	}
	defer walletService.Close()

	authHandler := handler.NewAuthHandler(authService)
	apiKeyHandler := handler.NewApiKeyHandler(apiKeyService)
	walletHandler := handler.NewWalletHandler(walletService)

	txWatcher := service.NewTxWatcher(walletRepo, walletService.EthClient())
	go func() {
		log.Println("tx watcher started")
		if err := txWatcher.Start(ctx); err != nil && err != context.Canceled {
			log.Printf("tx watcher stopped: %v", err)
		}
	}()

	r := chi.NewRouter()

	// CORS Configuration - Allow frontend domain
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{
			"http://localhost:3000",
			"https://mastergo-pr.onrender.com",
			"https://mastergo-pr-1.onrender.com",
		},
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
	}))

	r.Use(chimiddleware.RequestID)
	r.Use(chimiddleware.RealIP)
	r.Use(chimiddleware.Logger)
	r.Use(chimiddleware.Recoverer)
	r.Use(chimiddleware.Timeout(60 * time.Second))

	// Explicit OPTIONS handler for preflight requests
	r.Method("OPTIONS", "/*", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})

	r.Route("/auth", func(r chi.Router) {
		r.Post("/signup", authHandler.Signup)
		r.Post("/verify-email", authHandler.VerifyEmail)
		r.Post("/login", authHandler.Login)
		r.Post("/resend-otp", authHandler.ResendOTP)
		r.Get("/password-strength", authHandler.GetPasswordStrength)
		r.Get("/email-available", authHandler.CheckEmailAvailability)
		r.Group(func(r chi.Router) {
			r.Use(middleware.AuthMiddleware(cfg.JWTSecret))
			r.Get("/me", authHandler.Me)
			r.Patch("/profile", authHandler.UpdateProfile)
			r.Post("/logout", authHandler.Logout)
		})
	})
	r.Route("/api/v1", func(r chi.Router) {
		r.Group(func(r chi.Router) {
			r.Use(middleware.AuthMiddleware(cfg.JWTSecret))

			// Wallet routes with auth
			r.Route("/wallet", func(r chi.Router) {
				r.Post("/create", walletHandler.Create)
				r.Get("/balance", walletHandler.GetBalance)
				r.Get("/transactions", walletHandler.GetTransactions)
				r.Post("/withdraw", walletHandler.Withdraw)
			})

			// API key routes with auth
			r.Route("/apikeys", func(r chi.Router) {
				r.Post("/", apiKeyHandler.Create)
				r.Get("/", apiKeyHandler.List)
				r.Delete("/{id}", apiKeyHandler.Delete)
				r.Post("/{id}/regenerate", apiKeyHandler.Regenerate)
			})
		})
	})

	server := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: r,
	}

	log.Printf("server running on :%s", cfg.Port)
	log.Fatal(server.ListenAndServe())
}
