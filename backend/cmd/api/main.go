package main

import (
	"context"
	"encoding/base64"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"backend/internal/config"
	"backend/internal/database"
	"backend/internal/handler"
	"backend/internal/middleware"
	"backend/internal/repository"
	"backend/internal/service"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
)

func main() {
	cfg := config.Load()

	// root context for background workers
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// postgres
	db, err := database.NewPostgresConnection(cfg.DBURL)
	if err != nil {
		log.Fatalf("postgres connection failed: %v", err)
	}
	defer db.Close()

	// redis
	rdb, err := database.NewRedisClient(cfg.REDIS_URL)
	if err != nil {
		log.Fatalf("redis connection failed: %v", err)
	}
	defer rdb.Close()

	// master key
	key, err := base64.StdEncoding.DecodeString(cfg.MasterKey)
	if err != nil || len(key) != 32 {
		log.Fatal("MASTER_KEY must be valid base64-encoded 32 bytes")
	}

	// repositories
	userRepo := repository.NewUserRepository(db)
	apiKeyRepo := repository.NewApiKeyRepository(db)
	walletRepo := repository.NewWalletRepository(db)

	// services
	authService := service.NewAuthService(userRepo, rdb, cfg.JWTSecret)
	apiKeyService := service.NewApiKeyService(apiKeyRepo)

	walletService, err := service.NewWalletService(walletRepo, key, cfg.EthRPCURL)
	if err != nil {
		log.Fatalf("wallet service init failed: %v", err)
	}
	defer walletService.Close()

	// handlers
	authHandler := handler.NewAuthHandler(authService)
	apiKeyHandler := handler.NewApiKeyHandler(apiKeyService)
	walletHandler := handler.NewWalletHandler(walletService)

	// tx watcher for withdrawal confirmation / refund flow
	txWatcher := service.NewTxWatcher(walletRepo, walletService.EthClient())
	go func() {
		log.Println("tx watcher started")
		if err := txWatcher.Start(ctx); err != nil && err != context.Canceled {
			log.Printf("tx watcher stopped: %v", err)
		}
	}()

	// router
	r := chi.NewRouter()

	r.Use(chimiddleware.RequestID)
	r.Use(chimiddleware.RealIP)
	r.Use(chimiddleware.Logger)
	r.Use(chimiddleware.Recoverer)
	r.Use(chimiddleware.Timeout(60 * time.Second))

	// limit request body size
	r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			r.Body = http.MaxBytesReader(w, r.Body, 1<<20) // 1 MB
			next.ServeHTTP(w, r)
		})
	})

	// health check
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("OK"))
	})

	// public auth routes
	r.Route("/auth", func(r chi.Router) {
		r.Post("/signup", authHandler.Signup)
		r.Post("/verify-email", authHandler.VerifyEmail)
		r.Post("/login", authHandler.Login)
		r.Post("/resend-otp", authHandler.ResendOTP)

		r.Group(func(r chi.Router) {
			r.Use(middleware.AuthMiddleware(cfg.JWTSecret))
			r.Get("/me", authHandler.Me)
			r.Patch("/profile", authHandler.UpdateProfile)
		})
	})

	// protected routes
	r.Route("/api/v1", func(r chi.Router) {
		r.Use(middleware.AuthMiddleware(cfg.JWTSecret))

		r.Route("/wallet", func(r chi.Router) {
			r.Post("/create", walletHandler.Create)
			r.Get("/balance", walletHandler.GetBalance)
			r.Get("/transactions", walletHandler.GetTransactions)
			r.Post("/withdraw", walletHandler.Withdraw)
		})

		r.Route("/apikeys", func(r chi.Router) {
			r.Post("/", apiKeyHandler.Create)
			r.Get("/", apiKeyHandler.List)
			r.Delete("/{id}", apiKeyHandler.Delete)
		})
	})

	// http server
	server := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 60 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	go func() {
		log.Printf("server running on :%s", cfg.Port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server error: %v", err)
		}
	}()

	// graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("shutdown signal received")

	cancel()

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer shutdownCancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("forced shutdown: %v", err)
	}

	log.Println("server stopped cleanly")
}
