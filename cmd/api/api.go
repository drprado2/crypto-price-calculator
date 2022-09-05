package main

import (
	"context"
	"crypto-price-calculator/internal/bootstrap/api"
	"crypto-price-calculator/internal/observability/applog"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	shutdownCh := make(chan os.Signal, 1)
	signal.Notify(shutdownCh, os.Interrupt)
	signal.Notify(shutdownCh, syscall.SIGTERM)

	mainCtx, cancelCtx := context.WithCancel(context.Background())

	api.Setup(mainCtx)

	errCh := make(chan error, 1)

	go func() {
		errCh <- api.Start(mainCtx)
	}()

	select {
	case err := <-errCh:
		applog.Logger(mainCtx).WithError(err).Error("error received from server")
	case sig := <-shutdownCh:
		applog.Logger(mainCtx).Infof("shutdown signal received, sig: %v", sig)
	}

	gracefulTimeout := time.Second * 5
	ctx, cancel := context.WithTimeout(mainCtx, gracefulTimeout)
	defer cancel()
	cancelCtx()
	api.Close(ctx)
	applog.Logger(ctx).Info("shutting down API")
	os.Exit(0)
}
