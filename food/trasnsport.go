package food

import (
	"context"
	"encoding/json"
	"errors"
	kitlog "github.com/go-kit/kit/log"
	"github.com/go-kit/kit/transport"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"net/http"
	"what_cook/domain"
	"what_cook/helper"
)

var badRequest = errors.New("bad request")

func MakeHandler(foodService domain.FoodService, logger kitlog.Logger) http.Handler {
	opts := []kithttp.ServerOption{
		kithttp.ServerErrorHandler(transport.NewLogErrorHandler(logger)),
		kithttp.ServerErrorEncoder(encodeError),
	}
	foodHandler := kithttp.NewServer(
		makeFoodEndpoint(foodService),
		decodeFoodRequest,
		encodeResponse,
		opts...,
	)
	createFoodHandler := kithttp.NewServer(
		makeCreateFoodEndpoint(foodService),
		decodeCreateFoodRequest,
		encodeResponse,
		opts...,
	)
	updateFoodHandler := kithttp.NewServer(
		makeUpdateFoodEndpoint(foodService),
		decodeUpdateFoodRequest,
		encodeResponse,
		opts...,
	)
	deleteFoodHandler := kithttp.NewServer(
		makeDeleteFoodEndpoint(foodService),
		decodeDeleteFoodRequest,
		encodeResponse,
		opts...,
	)
	foodsByIngredientsHandler := kithttp.NewServer(
		makeFoodsByIngredientEndpoint(foodService),
		decodeFoodsByIngredientsRequest,
		encodeResponse,
		opts...,
	)

	router := mux.NewRouter()
	router.Handle("/food/{id}", foodHandler).Methods("GET")
	router.Handle("/food/", createFoodHandler).Methods("POST")
	router.Handle("/food/{id}", updateFoodHandler).Methods("PUT")
	router.Handle("/food/{id}", deleteFoodHandler).Methods("DELETE")
	router.Handle("/food/byIngredients/", foodsByIngredientsHandler).Methods("GET")
	return router
}

func decodeFoodRequest(_ context.Context, r *http.Request) (interface{}, error) {
	if id, err := helper.GetRequestParam(r, "id"); err == nil {
		return foodRequest{id}, nil
	}
	return nil, badRequest
}

func decodeCreateFoodRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request createFoodRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}
	return request, nil
}

func decodeUpdateFoodRequest(_ context.Context, r *http.Request) (interface{}, error) {
	if id, err := helper.GetRequestParam(r, "id"); err == nil {
		var body struct {
			Food *domain.Food `json:"food"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err == nil {
			return updateFoodRequest{id, body.Food}, nil
		}
	}
	return nil, badRequest
}

func decodeDeleteFoodRequest(_ context.Context, r *http.Request) (interface{}, error) {
	if id, err := helper.GetRequestParam(r, "id"); err == nil {
		return deleteFoodRequest{id}, nil
	}
	return nil, badRequest
}

func decodeFoodsByIngredientsRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request foodsByIngredientsRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}
	return request, nil
}

type errorer interface {
	error() error
}

func encodeResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	e, ok := response.(errorer)
	if ok && e.error() != nil {
		encodeError(ctx, e.error(), w)
		return nil
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(w).Encode(response)
}

func encodeError(_ context.Context, err error, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	switch err {
	case badRequest:
		w.WriteHeader(http.StatusBadRequest)
	case domain.ModelNotFoundError:
		w.WriteHeader(http.StatusNotFound)
	default:
		w.WriteHeader(http.StatusInternalServerError) // TODO: debug true|false, logging
	}
	json.NewEncoder(w).Encode(map[string]interface{}{
		"error": err.Error(),
	})
}
