package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/SerzhLimon/GopherMartSolo/internal/config"
	"github.com/SerzhLimon/GopherMartSolo/internal/handler"
	"github.com/SerzhLimon/GopherMartSolo/internal/middleware"
	"github.com/SerzhLimon/GopherMartSolo/pkg/accrual"
	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog/log"
)

func main() {

	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to load config")
	}

	router := chi.NewRouter()
	h := handler.NewHandler(cfg)

	router.Use(middleware.LoggingMiddleware)

	router.Post("/api/user/register", h.RegisterUserHandler)
	router.Post("/api/user/login", h.LoginHandler)

	router.Group(func(router chi.Router) {
		router.Use(middleware.AuthMiddleware(cfg.SecretKey))

		router.Post("/api/user/orders", h.UploadOrderHandler)
		router.Get("/api/user/orders", h.GetOrdersHandler)

		router.Get("/api/user/balance", h.GetBalanceHandler)

		router.Post("/api/user/balance/withdraw", h.WithdrawHandler)
		router.Get("/api/user/withdrawals", h.GetWithdrawalsHandler)
	})

	accrual := accrual.NewAccrual(cfg.AccrualAddress, h.Storage)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go accrual.Start(ctx)
	log.Printf("Accrual service available on: %s...", cfg.AccrualAddress)

	srv := &http.Server{
		Addr:    cfg.Host,
		Handler: router,
	}
	serverErr := make(chan error, 1)
	go func() {
		log.Printf("Starting server on: %s...", cfg.Host)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			serverErr <- err
		}
	}()

	exit := make(chan os.Signal, 1)
	signal.Notify(exit, os.Interrupt, syscall.SIGTERM)

	select {
	case <-exit:
		log.Info().Msg("Shutdown signal received")
	case err := <-serverErr:
		log.Error().Err(err).Msg("Server error occurred")
	}

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Error().Err(err).Msg("Server shutdown failed")
	}

	cancel()
	time.Sleep(500 * time.Millisecond)

	log.Info().Msg("Server stopped gracefully")
}
