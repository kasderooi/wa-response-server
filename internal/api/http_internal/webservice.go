package http_internal

import (
	"context"
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/sirupsen/logrus"
	"net/http"
	"time"
	"whatsapp-bot/internal/config"
	"whatsapp-bot/internal/controllers"
	"whatsapp-bot/internal/usecase"

	"github.com/go-chi/chi/v5/middleware"
)

func StartWebservice(ctx context.Context, cfg *config.Config, waClient *usecase.WhatsAppClient, logger *logrus.Logger) error {
	logger.WithField("bind_addr", cfg.HttpAddr).Debug("Binding to address")

	h := controllers.NewWhatsAppController(waClient, logger)
	mainServer := createMainServer(cfg, h)

	serverErrChan := make(chan error, 2)

	go func() {
		logger.Info("Starting HTTP server")
		serverErrChan <- mainServer.ListenAndServe()
	}()

	select {
	case <-ctx.Done():
		logger.Info("Shutdown signal received, initiating graceful shutdown")
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
		defer cancel()

		if err := mainServer.Shutdown(shutdownCtx); err != nil {
			logger.WithError(err).Error("Graceful shutdown of main server failed")
			return err
		}
		logger.Info("Both servers shut down gracefully")
	case err := <-serverErrChan:
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.WithError(err).Error("HTTP server error")
			return err
		}
	}

	return nil
}

func createMainServer(cfg *config.Config, h *controllers.WhatsAppController) *http.Server {
	router := chi.NewRouter()
	router.Use(
		middleware.Recoverer,
		middleware.RealIP,
	)
	SetupRoutes(router, h)

	return &http.Server{
		Addr:           cfg.HttpAddr,
		Handler:        router,
		ReadTimeout:    cfg.HttpIOTimeout,
		WriteTimeout:   cfg.HttpIOTimeout,
		IdleTimeout:    cfg.HttpIOTimeout * 2,
		MaxHeaderBytes: 1 << 16,
	}
}
