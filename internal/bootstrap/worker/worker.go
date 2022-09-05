package worker

import (
	"context"
	"crypto-price-calculator/internal/adapters/repositories"
	"crypto-price-calculator/internal/adapters/vwapobservables"
	"crypto-price-calculator/internal/configs"
	"crypto-price-calculator/internal/core/usecases/registertradeorder"
	"crypto-price-calculator/internal/core/vwap"
	"crypto-price-calculator/internal/infra/database/postgres"
	"crypto-price-calculator/internal/infra/websocketserver"
	"crypto-price-calculator/internal/infra/websocketserver/coinbase"
	"crypto-price-calculator/internal/observability/applog"
	"crypto-price-calculator/internal/observability/apptracer"
)

var (
	worker         *websocketserver.Server
	vwapConsumer   *vwap.Consumer
	vwapCalculator vwap.CalculatorInterface
)

func Setup(ctx context.Context) {
	applog.Setup()
	logger := applog.Logger(ctx)

	config := configs.Get()

	dbPool, err := postgres.Setup(ctx, config)
	if err != nil {
		logger.WithError(err).Fatal("error on setup DB")
	}

	if err := apptracer.Setup(ctx, config); err != nil {
		logger.WithError(err).Fatal("error creating tracer provider")
	}

	tradesCh := make(chan *vwap.TradeEvent, 200)

	vwapLoggerObs := vwapobservables.NewLogger()
	vwapSnsObs := vwapobservables.NewPublishSns()

	traderOrderRepository := repositories.NewTradeOrderRepository(dbPool)

	vwapCalculator = vwap.NewCalculator(config, traderOrderRepository, vwapLoggerObs, vwapSnsObs)
	vwapConsumer = vwap.NewConsumer(vwapCalculator, tradesCh)
	registerTradeUC := registertradeorder.NewHandler(config.GetProductIds(), traderOrderRepository, tradesCh)

	subController := coinbase.NewSubscriptionsController(config)
	errController := coinbase.NewErrorController()
	unkController := coinbase.NewUnknownController()
	matchController := coinbase.NewMatchController(registerTradeUC)

	consumer := coinbase.NewConsumer(unkController, map[string]coinbase.Controller{
		coinbase.SubscriptionsType: subController,
		coinbase.ErrorType:         errController,
		coinbase.MatchType:         matchController,
	})

	worker = websocketserver.NewServer(config, consumer)
}

func Start(ctx context.Context) error {
	vwapConsumer.StartConsumer(ctx)
	vwapCalculator.Setup(ctx)
	return worker.Start(ctx)
}

func Close(ctx context.Context) {
	postgres.Close()
	apptracer.Close(ctx)
	vwapConsumer.Close(ctx)
	worker.Close(ctx)
}
