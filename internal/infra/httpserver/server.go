package httpserver

import (
	"context"
	"crypto-price-calculator/internal/configs"
	"crypto-price-calculator/internal/infra/httpserver/middlewares"
	"crypto-price-calculator/internal/observability/applog"
	"fmt"
	"github.com/felixge/fgprof"
	"github.com/gorilla/mux"
	"net/http"
	_ "net/http/pprof"
	"time"
)

const (
	Post   = "POST"
	Put    = "PUT"
	Delete = "DELETE"
	Get    = "GET"
)

type Server struct {
	server    *http.Server
	mainCtx   context.Context
	cancelCtx context.CancelFunc
	router    *mux.Router
}

func NewServer(ctx context.Context) *Server {
	server := new(Server)

	r := mux.NewRouter().StrictSlash(true)
	ctx, cancel := context.WithCancel(ctx)
	server.mainCtx = ctx
	server.cancelCtx = cancel

	server.router = r.PathPrefix("/api").Subrouter()

	return server
}

func (s *Server) WithRoutes(handler func(router *mux.Router)) *Server {
	applog.Logger(s.mainCtx).Info("registering request handlers")

	handler(s.router)
	return s
}

func (s *Server) Start() error {
	envs := configs.Get()
	logger := applog.Logger(s.mainCtx)
	http.Handle("/", s.router)

	s.registerMiddlewares(s.router)

	if envs.ServerEnvironment == configs.DeveloperEnvironment {
		go func() {
			logger.Info("starting pgprof server")
			http.DefaultServeMux.Handle("/debug/fgprof", fgprof.Handler())
			if err := http.ListenAndServe(":6060", nil); err != nil {
				logger.WithError(err).Error("error strarting pgprof")
			}
		}()
	}

	timeout, err := time.ParseDuration(envs.ServerEndpointTimeout)
	if err != nil {
		panic(err)
	}
	muxWithMiddlewares := http.TimeoutHandler(s.router, timeout, "timeout occurred!")

	s.server = &http.Server{
		Addr:         fmt.Sprintf(":%v", envs.ServerPort),
		Handler:      muxWithMiddlewares,
		WriteTimeout: timeout,
		ReadTimeout:  timeout,
	}

	logger.Infof("http server running at port %v, with env %v", envs.ServerPort, envs.ServerEnvironment)
	if err := s.server.ListenAndServe(); err != nil {
		logger.WithError(err).Error("Fail starting http server")
		return err
	}

	return nil
}

func (s *Server) Shutdown(ctx context.Context) {
	applog.Logger(ctx).Info("shutting down http server")
	s.server.Shutdown(ctx)
}

func (s *Server) registerMiddlewares(router *mux.Router) {
	router.Use(mux.CORSMethodMiddleware(router))
	router.Use(middlewares.PanicMiddleware)
	router.Use(middlewares.CidMiddleware)
	router.Use(middlewares.SpanMiddleware)
	router.Use(middlewares.RequestLogMiddleware)
}
