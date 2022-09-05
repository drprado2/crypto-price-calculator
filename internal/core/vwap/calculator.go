package vwap

import (
	"context"
	"crypto-price-calculator/internal/configs"
	"crypto-price-calculator/internal/core/entities"
	"crypto-price-calculator/internal/core/repositories"
	"crypto-price-calculator/internal/observability/applog"
	"crypto-price-calculator/internal/observability/apptracer"
	"math"
	"sync"
)

type (
	Calculator struct {
		vwapMutex            *sync.Mutex
		vwapObservables      []VwapObservable
		tradeOrderRepository repositories.TradeOrderRepository
		products             []string
		vwapByProduct        map[string]*Vwap
		windowSize           int
	}

	CalculatorInterface interface {
		UpdateVwapPrice(ctx context.Context, tradeEvent *TradeEvent)
		Setup(ctx context.Context) error
	}

	VwapUpdatedEvent struct {
		ProductId                  string
		VolumeWeightedAveragePrice float64
	}

	TradeEvent struct {
		Price   float64
		Size    float64
		Product string
	}

	Vwap struct {
		Items            chan *vwapItem
		TotalVolume      float64
		TotalWeightPrice float64
	}

	vwapItem struct {
		vWPrice float64
		volume  float64
	}

	VwapObservable interface {
		HandleNewVwap(ctx context.Context, event *VwapUpdatedEvent)
	}
)

func vwapItemFromTrade(trade *entities.TradeOrder) *vwapItem {
	return &vwapItem{
		vWPrice: trade.Price * trade.Size,
		volume:  trade.Size,
	}
}

func vwapItemFromEvent(trade *TradeEvent) *vwapItem {
	return &vwapItem{
		vWPrice: trade.Price * trade.Size,
		volume:  trade.Size,
	}
}

func NewCalculator(config *configs.Configuration, tradeOrderRepository repositories.TradeOrderRepository, vwapObservables ...VwapObservable) CalculatorInterface {
	return &Calculator{
		tradeOrderRepository: tradeOrderRepository,
		vwapObservables:      vwapObservables,
		windowSize:           config.VwapWindowSize,
		products:             config.GetProductIds(),
		vwapByProduct:        make(map[string]*Vwap),
		vwapMutex:            &sync.Mutex{},
	}
}

func (c *Calculator) Setup(ctx context.Context) error {
	ctx, span := apptracer.StartOperation(ctx, "Calculator:Setup", apptracer.SpanKindInternal)
	defer span.Finish()

	logger := applog.Logger(ctx)

	for _, product := range c.products {
		trades, err := c.tradeOrderRepository.RetrieveLastNTradesByProduct(ctx, product, c.windowSize)
		if err != nil {
			logger.WithError(err).WithField("Product", product).Errorf("Error retrieving trades in order to calculate VWAP for product %s", product)
			return err
		}

		vwap := &Vwap{
			Items:            make(chan *vwapItem, c.windowSize),
			TotalVolume:      0,
			TotalWeightPrice: 0,
		}

		for i := len(trades) - 1; i >= 0; i-- {
			trade := trades[i]
			vwap.Items <- vwapItemFromTrade(trade)
			vwap.TotalWeightPrice += trade.Price * trade.Size
			vwap.TotalVolume += trade.Size
		}

		c.vwapByProduct[product] = vwap
	}

	return nil
}

func (c *Calculator) UpdateVwapPrice(ctx context.Context, event *TradeEvent) {
	ctx, span := apptracer.StartOperation(ctx, "Calculator:UpdateVwapPrice", apptracer.SpanKindInternal)
	defer span.Finish()

	vwap := c.appendTrade(event)

	vwapEvent := &VwapUpdatedEvent{
		ProductId:                  event.Product,
		VolumeWeightedAveragePrice: vwap.CalculateVwap(),
	}

	c.notifyObservables(ctx, vwapEvent)
}

func (c *Calculator) appendTrade(event *TradeEvent) *Vwap {
	item := vwapItemFromEvent(event)
	vwap := c.vwapByProduct[event.Product]
	c.vwapMutex.Lock()
	defer c.vwapMutex.Unlock()

	if len(vwap.Items) == c.windowSize {
		removedItem := <-vwap.Items
		vwap.TotalWeightPrice -= removedItem.vWPrice
		vwap.TotalVolume -= removedItem.volume
	}

	vwap.Items <- item

	vwap.TotalVolume += item.volume
	vwap.TotalWeightPrice += item.vWPrice

	return vwap
}

func (c *Calculator) notifyObservables(ctx context.Context, event *VwapUpdatedEvent) {
	syncGroup := &sync.WaitGroup{}
	for _, obs := range c.vwapObservables {
		syncGroup.Add(1)
		go func(oobs VwapObservable) {
			defer syncGroup.Done()
			oobs.HandleNewVwap(ctx, event)
		}(obs)
	}
	syncGroup.Wait()
}

func (v *Vwap) CalculateVwap() float64 {
	value := v.TotalWeightPrice / v.TotalVolume
	return math.Round(value*100) / 100
}
