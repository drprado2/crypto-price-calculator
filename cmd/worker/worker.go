package main

import (
	"context"
	"crypto-price-calculator/internal/bootstrap/worker"
	"crypto-price-calculator/internal/observability/applog"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	shutdownCh := make(chan os.Signal, 1)
	signal.Notify(shutdownCh, os.Interrupt, syscall.SIGTERM, syscall.SIGHUP, syscall.SIGINT, syscall.SIGQUIT)

	mainCtx, cancelCtx := context.WithCancel(context.Background())

	worker.Setup(mainCtx)

	errCh := make(chan error, 1)

	go func() {
		errCh <- worker.Start(mainCtx)
	}()

	select {
	case err := <-errCh:
		applog.Logger(mainCtx).WithError(err).Error("error received from server")
	case sig := <-shutdownCh:
		applog.Logger(mainCtx).Infof("shutdown signal received, sig: %v", sig)
	}

	gracefulTimeout := time.Second * 15
	ctx, cancel := context.WithTimeout(mainCtx, gracefulTimeout)
	defer cancel()
	cancelCtx()
	worker.Close(ctx)
	applog.Logger(ctx).Info("shutting down worker")
	os.Exit(0)
}
