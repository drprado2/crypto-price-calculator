package websocketserver

import (
	"bytes"
	"context"
	"crypto-price-calculator/internal/configs"
	"crypto-price-calculator/internal/ctxutils"
	"crypto-price-calculator/internal/infra/websocketserver/coinbase"
	"crypto-price-calculator/internal/observability/applog"
	"github.com/felixge/fgprof"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"net/http"
	"time"
)

type (
	Server struct {
		socketUrl         string
		connByProduct     map[string]*websocket.Conn
		config            *configs.Configuration
		productIds        []string
		consumeChannel    string
		consumer          *coinbase.Consumer
		healthCheckTicker *time.Ticker
		closed            bool
	}
)

const (
	writeWait      = 6 * time.Second
	maxMessageSize = 1000000 // bytes
)

var (
	newline       = []byte{'\n'}
	space         = []byte{' '}
	serverHeaders = map[string][]string{
		"Sec-WebSocket-Extensions": {
			"permessage-deflate",
		},
	}
)

func NewServer(config *configs.Configuration, consumer *coinbase.Consumer) *Server {
	return &Server{
		connByProduct:  make(map[string]*websocket.Conn),
		socketUrl:      config.CoinbaseWsEndpoint,
		config:         config,
		productIds:     config.GetProductIds(),
		consumeChannel: config.MatchesChannel,
		consumer:       consumer,
	}
}

func (s *Server) Start(ctx context.Context) error {
	logger := applog.Logger(ctx)
	defer func() {
		for _, conn := range s.connByProduct {
			conn.Close()
		}
	}()

	if s.config.ServerEnvironment == configs.DeveloperEnvironment {
		go func() {
			logger.Info("starting pgprof server")
			http.DefaultServeMux.Handle("/debug/fgprof", fgprof.Handler())
			if err := http.ListenAndServe(":6060", nil); err != nil {
				logger.WithError(err).Error("error strarting pgprof")
			}
		}()
	}

	for _, p := range s.productIds {
		conn, _, err := websocket.DefaultDialer.Dial(s.socketUrl, serverHeaders)
		if err != nil {
			logger.WithError(err).Error("Error connecting to Websocket Server")
			return err
		}
		s.connByProduct[p] = conn
	}

	if err := s.startConsumeChannels(ctx); err != nil {
		logger.WithError(err).Error("Error connecting to Websocket Server")
		return err
	}

	return s.startConsumer(ctx)
}

func (s *Server) startConsumeChannels(ctx context.Context) error {
	logger := applog.Logger(ctx)

	for product, conn := range s.connByProduct {
		logger.WithField("ProductId", product).Infof("sending subscribe request for product %s", product)

		subscribeRequest := coinbase.NewSubscribeRequest().
			WithProductIds(product).
			WithChannels(s.consumeChannel)

		conn.SetWriteDeadline(time.Now().Add(writeWait))
		if err := conn.WriteJSON(subscribeRequest); err != nil {
			logger.WithError(err).Error("Error writing subscribe request to coinbase WS")
			return err
		}
	}

	return nil
}

func (s *Server) startConsumer(ctx context.Context) error {
	logger := applog.Logger(ctx)

	errCh := make(chan error, len(s.productIds))

	for product, conn := range s.connByProduct {
		logger.WithField("ProductId", product).Infof("starting consumer for product %s", product)

		go func(iconn *websocket.Conn) {
			for {
				iconn.SetReadLimit(maxMessageSize)

				_, msg, err := iconn.ReadMessage()
				if err != nil {
					logger.WithError(err).WithField("ProductId", product).Error("Error happened reading WS connection")
					errCh <- err
					return
				}
				ctx := ctxutils.WithCid(ctx, uuid.New().String())
				msg = bytes.TrimSpace(bytes.Replace(msg, newline, space, -1))
				if err := s.consumer.Consume(ctx, msg); err != nil {
					errCh <- err
					return
				}
			}
		}(conn)
	}

	return <-errCh
}

func (s *Server) Close(ctx context.Context) {
	if s.closed {
		return
	}

	s.closed = true

	logger := applog.Logger(ctx)

	for product, conn := range s.connByProduct {
		logger.WithField("ProductId", product).Infof("closing connection of product %s", product)

		unsubscribeRequest := coinbase.NewUnsubscribeRequest().
			WithProductIds(product).
			WithChannels(s.consumeChannel)

		conn.SetWriteDeadline(time.Now().Add(writeWait))
		if err := conn.WriteJSON(unsubscribeRequest); err != nil {
			logger.WithError(err).Error("Error writing unsubscribe request to coinbase WS")
		}
		if err := conn.WriteMessage(websocket.CloseMessage, []byte{}); err != nil {
			logger.WithError(err).Error("Error writing close message to coinbase WS")
		}
	}
}
