package coinbase

import (
	"context"
)

const (
	subscribe   = "subscribe"
	unsubscribe = "unsubscribe"
)

type (
	Controller interface {
		// Handle just returns errors you want to restart the app
		Handle(ctx context.Context, message []byte) error
	}

	SubUnsubRequest struct {
		Type       string   `json:"type"`
		ProductIds []string `json:"product_ids"`
		Channels   []any    `json:"channels"`
	}

	Channel struct {
		Name       string   `json:"name"`
		ProductIds []string `json:"product_ids"`
	}
)

func NewSubscribeRequest() *SubUnsubRequest {
	return &SubUnsubRequest{
		Type: subscribe,
	}
}

func NewUnsubscribeRequest() *SubUnsubRequest {
	return &SubUnsubRequest{
		Type: unsubscribe,
	}
}

func (s *SubUnsubRequest) WithProductIds(products ...string) *SubUnsubRequest {
	s.ProductIds = products
	return s
}

func (s *SubUnsubRequest) WithChannels(channels ...any) *SubUnsubRequest {
	s.Channels = channels
	return s
}
