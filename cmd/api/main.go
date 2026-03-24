package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/FelixWinchester/ssubench/internal/config"
	"github.com/FelixWinchester/ssubench/internal/handler"
	"github.com/FelixWinchester/ssubench/internal/repo"
	"github.com/FelixWinchester/ssubench/internal/service"
)

func main() {
	// Config
	cfg := config.MustLoad()

	// Logger
	log := setupLogger(cfg.App.LogLevel)
	log.Info("starting ssubench", slog.String("env", cfg.App.LogLevel))

	// Database
	ctx := context.Background()
	db, err := repo.NewPostgresPool(ctx, cfg.Database.DSN())
	if err != nil {
		log.Error("failed to connect to database", slog.Any("error", err))
		os.Exit(1)
	}
	defer db.Close()
	log.Info("connected to database")

	// Repos
	userRepo := repo.NewUserRepo(db)
	taskRepo := repo.NewTaskRepo(db)
	bidRepo := repo.NewBidRepo(db)
	paymentRepo := repo.NewPaymentRepo(db)

	// Services
	authSvc := service.NewAuthService(userRepo, cfg.JWT.Secret, cfg.JWT.TTL)
	userSvc := service.NewUserService(userRepo)
	taskSvc := service.NewTxTaskService(db, taskRepo, userRepo, bidRepo, paymentRepo)
	bidSvc := service.NewBidService(bidRepo, taskRepo)
	paymentSvc := service.NewPaymentService(paymentRepo)

	// Handler
	h := handler.New(authSvc, userSvc, taskSvc, bidSvc, paymentSvc, log, cfg.JWT.Secret)
	router := h.InitRoutes()

	// Server
	srv := &http.Server{
		Addr:         ":" + cfg.Server.Port,
		Handler:      router,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
		IdleTimeout:  cfg.Server.IdleTimeout,
	}

	// Запуск сервера в горутине
	go func() {
		log.Info("server started", slog.String("addr", srv.Addr))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Error("server error", slog.Any("error", err))
			os.Exit(1)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info("shutting down server...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Error("server shutdown error", slog.Any("error", err))
	}

	log.Info("server stopped")
}

func setupLogger(level string) *slog.Logger {
	var l slog.Level
	switch level {
	case "debug":
		l = slog.LevelDebug
	case "info":
		l = slog.LevelInfo
	case "warn":
		l = slog.LevelWarn
	case "error":
		l = slog.LevelError
	default:
		l = slog.LevelInfo
	}

	return slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: l}))
}
