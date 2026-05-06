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

	authService := service.NewAuthService(userRepo, rdb, cfg.JWTSecret)
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

	//  CORS Configuration
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{
			"http://localhost:3001",
			"http://localhost:3000",
			"http://127.0.0.1:3001",
			"http://127.0.0.1:3000",
			"https://your-frontend-url.onrender.com",
		},
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-Requested-With"},
		AllowCredentials: true,
	}))

	r.Use(chimiddleware.RequestID)
	r.Use(chimiddleware.RealIP)
	r.Use(chimiddleware.Logger)
	r.Use(chimiddleware.Recoverer)
	r.Use(chimiddleware.Timeout(60 * time.Second))

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
