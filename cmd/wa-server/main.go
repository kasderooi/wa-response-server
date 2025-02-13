package main

import (
	"context"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"syscall"
	"whatsapp-bot/internal/api/http_internal"
	"whatsapp-bot/internal/config"
	"whatsapp-bot/internal/infrastructure"
	"whatsapp-bot/internal/usecase"
)

func main() {
	logger := logrus.New()
	logger.Info("Starting service")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cfg, err := config.ParseFromEnvVars()
	if err != nil {
		logger.Fatalf("config parsing error: %v", err)
	}

	errChan := make(chan error, 1)
	closers, err := run(ctx, cfg, logger, errChan)
	if err != nil {
		logger.Fatalf("startup error: %v", err)
	}

	onExit(errChan, closers...)
}

func run(ctx context.Context, cfg config.Config, logger *logrus.Logger, errChan chan error) ([]func() error, error) {
	var closers []func() error

	logger.Info("starting db")
	gormDb, sqlDb, err := infrastructure.InitDb(infrastructure.Database{
		MaxConnLifetime: cfg.MaxConnLifetime,
		MaxIdleConns:    cfg.MaxIdleConns,
		MaxOpenConns:    cfg.MaxOpenConns,
		DataSourceName:  cfg.DbDsn,
		Dialect:         infrastructure.PostgresSQL,
	})
	if err != nil {
		return closers, fmt.Errorf("could not connect to DB: %v", err)
	}
	closers = append(closers, sqlDb.Close)

	logger.Info("migrate db")
	if err = infrastructure.MigrateDB(gormDb); err != nil {
		return closers, fmt.Errorf("could not migrate DB: %v", err)
	}

	logger.Info("starting whatsapp client")
	waClient := usecase.NewWhatsAppClient(sqlDb, cfg.DbDialect, cfg.DbLogLevel)
	closers = append(closers, waClient.Stop)
	err = waClient.Start()
	if err != nil {
		return closers, fmt.Errorf("could not start whats app client: %v", err)
	}

	logger.Info("starting webservice")
	go func() {
		errChan <- http_internal.StartWebservice(ctx, &cfg, waClient, logger)
	}()

	return closers, nil
}

func onExit(errChan chan error, run ...func() error) {
	killSignal := make(chan os.Signal, 1)
	signal.Notify(killSignal, syscall.SIGINT, syscall.SIGTERM)

	select {
	case s := <-killSignal:
		fmt.Printf("exited with signal: %s", s)
	case err := <-errChan:
		fmt.Printf("exited with error: %s", err)
	}

	for _, r := range run {
		err := r()
		if err != nil {
			fmt.Printf("error while closing: %v", err)
		}
	}
}
