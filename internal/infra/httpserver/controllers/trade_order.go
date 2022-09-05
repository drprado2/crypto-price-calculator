package controllers

import (
	"context"
	"crypto-price-calculator/internal/core/repositories"
	"crypto-price-calculator/internal/core/usecases/registertradeorder"
	"crypto-price-calculator/internal/infra/httpserver"
	"crypto-price-calculator/internal/observability/applog"
	"encoding/json"
	"errors"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
)

type (
	TradeOrderController struct {
		registerTradeOrder registertradeorder.HandlerInterface
	}
)

func NewHelloWorldController(registerTradeOrder registertradeorder.HandlerInterface) *TradeOrderController {
	return &TradeOrderController{
		registerTradeOrder: registerTradeOrder,
	}
}

func (c *TradeOrderController) RegisterRouteHandlers(router *mux.Router) {
	router.
		Path("/internal/v1/trade-order").
		HandlerFunc(c.RegisterTradeOrderAction).
		Name("Hello world").
		Methods(httpserver.Post)
}

func (c *TradeOrderController) RegisterTradeOrderAction(writter http.ResponseWriter, req *http.Request) {
	writter.Header().Set("Content-Type", "application/json")

	reqBody, _ := ioutil.ReadAll(req.Body)
	model := new(registertradeorder.Input)
	if err := json.Unmarshal(reqBody, model); err != nil {
		response := map[string]string{
			"error": err.Error(),
		}
		writter.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(writter).Encode(response)
		return
	}

	res, err := c.registerTradeOrder.Exec(req.Context(), model)

	if err != nil {
		HandleError(req.Context(), err, writter, req)
		return
	}

	writter.WriteHeader(http.StatusCreated)
	json.NewEncoder(writter).Encode(res)
}

func HandleError(ctx context.Context, err error, writter http.ResponseWriter, _ *http.Request) {
	if err, ok := err.(*registertradeorder.InvalidInputErr); ok {
		writter.WriteHeader(http.StatusBadRequest)
		response := map[string]string{
			"error": err.Error(),
		}
		if err := json.NewEncoder(writter).Encode(response); err != nil {
			panic(err)
		}
		return
	}

	if errors.Is(err, repositories.RegisterAlreadyExists) {
		writter.WriteHeader(http.StatusConflict)
		response := map[string]string{
			"error": err.Error(),
		}
		if err := json.NewEncoder(writter).Encode(response); err != nil {
			panic(err)
		}
		return
	}

	writter.WriteHeader(http.StatusInternalServerError)
	response := map[string]string{
		"error": "unexpected error happened",
	}
	if err := json.NewEncoder(writter).Encode(response); err != nil {
		applog.Logger(ctx).WithError(err).Error("error writing http response")
	}
	return
}
