package entities

import (
	"github.com/google/uuid"
	"time"
)

type (
	TradeOrder struct {
		TradeOrderId string
		TradeId      int
		MakerOrderId string
		TakerOrderId string
		Side         string
		Size         float64
		Price        float64
		ProductId    string
		Sequence     int64
		Time         time.Time
	}
)

func NewTradeOrder() *TradeOrder {
	return &TradeOrder{
		TradeOrderId: uuid.New().String(),
	}
}
