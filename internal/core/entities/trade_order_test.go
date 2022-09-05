package entities

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewTradeOrder(t *testing.T) {
	trader := NewTradeOrder()

	assert.NotEmpty(t, trader.TradeOrderId)
}
