package coinbase

import (
	"bytes"
	"context"
)

type (
	Consumer struct {
		controllers       map[string]Controller
		unknownController Controller
	}
)

var (
	typeSeparator            = []byte(`"type":"`)
	typeEndCommaSeparator    = []byte(`",`)
	typeEndBracketsSeparator = []byte(`"}`)
)

func NewConsumer(unknownController Controller, controllers map[string]Controller) *Consumer {
	return &Consumer{
		controllers:       controllers,
		unknownController: unknownController,
	}
}

func (r *Consumer) Consume(ctx context.Context, message []byte) error {
	mtype := r.getTypeFromMessage(message)
	if c, ok := r.controllers[mtype]; ok {
		return c.Handle(ctx, message)
	}
	return r.unknownController.Handle(ctx, message)
}

func (r *Consumer) getTypeFromMessage(message []byte) string {
	const typeLen = 8

	typeStart := bytes.Index(message, typeSeparator)
	if typeStart == -1 {
		return ""
	}

	typeStart = typeStart + typeLen

	typeEnd := bytes.Index(message[typeStart:], typeEndCommaSeparator)
	if typeEnd == -1 {
		typeEnd = bytes.Index(message[typeStart:], typeEndBracketsSeparator)
		if typeEnd == -1 {
			return ""
		}
	}

	return string(message[typeStart : typeStart+typeEnd])
}
