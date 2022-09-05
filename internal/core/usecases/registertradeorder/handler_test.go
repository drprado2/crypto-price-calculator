package registertradeorder

import (
	"context"
	"crypto-price-calculator/internal/core/entities"
	"crypto-price-calculator/internal/core/repositories"
	"crypto-price-calculator/internal/core/vwap"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestHandler_Exec(t *testing.T) {
	ctx := context.Background()
	productIds := []string{
		"ETH-USD",
		"ETH-BTC",
	}
	createCount := 0
	repository := &repositories.TraderOrderRepositoryMock{
		MockCreate: func(ctx context.Context, trade *entities.TradeOrder) error {
			createCount++
			return nil
		},
	}
	tradesCh := make(chan *vwap.TradeEvent, 1)

	handler := NewHandler(productIds, repository, tradesCh)

	t.Run("Invalid input", func(t *testing.T) {
		input := &Input{}
		res, err := handler.Exec(ctx, input)
		assert.IsType(t, &InvalidInputErr{}, err)
		assert.Error(t, err)
		assert.Nil(t, res)
		assert.Equal(t, 0, createCount)
		assert.Equal(t, 0, len(tradesCh))
	})

	t.Run("Invalid size", func(t *testing.T) {
		input := generateValidInput()
		input.Size = "0.2t63"
		res, err := handler.Exec(ctx, input)
		assert.IsType(t, &InvalidInputErr{}, err)
		assert.Error(t, err)
		assert.Nil(t, res)
		assert.Equal(t, 0, createCount)
		assert.Equal(t, 0, len(tradesCh))
	})

	t.Run("Invalid price", func(t *testing.T) {
		input := generateValidInput()
		input.Price = "0.2t63"
		res, err := handler.Exec(ctx, input)
		assert.IsType(t, &InvalidInputErr{}, err)
		assert.Error(t, err)
		assert.Nil(t, res)
		assert.Equal(t, 0, createCount)
		assert.Equal(t, 0, len(tradesCh))
	})

	t.Run("Valid case", func(t *testing.T) {
		input := generateValidInput()
		res, err := handler.Exec(ctx, input)
		assert.NoError(t, err)
		assert.Equal(t, 1, createCount)
		assert.Equal(t, 1, len(tradesCh))
		assert.NotNil(t, res)
		event := <-tradesCh
		assert.Equal(t, "ETH-BTC", event.Product)
	})
}

func generateValidInput() *Input {
	return &Input{
		TradeId:      12,
		MakerOrderId: "id",
		TakerOrderId: "id",
		Side:         "sell",
		Size:         "0.26",
		Price:        "15226.056",
		ProductId:    "ETH-BTC",
		Sequence:     12,
		Time:         time.Time{},
	}
}
