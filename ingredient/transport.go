package ingredient

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

func MakeHandler(is domain.IngredientService, logger kitlog.Logger) http.Handler {
	opts := []kithttp.ServerOption{
		kithttp.ServerErrorHandler(transport.NewLogErrorHandler(logger)),
		kithttp.ServerErrorEncoder(encodeError),
	}
	ingredientHandler := kithttp.NewServer(
		makeIngredientEndpoint(is),
		decodeIngredientRequest,
		encodeResponse,
		opts...,
	)
	createIngredientHandler := kithttp.NewServer(
		makeCreateIngredientEndpoint(is),
		decodeCreateIngredientRequest,
		encodeResponse,
		opts...,
	)
	updateIngredientHandler := kithttp.NewServer(
		makeUpdateIngredientEndpoint(is),
		decodeUpdateIngredientRequest,
		encodeResponse,
		opts...,
	)
	deleteIngredientHandler := kithttp.NewServer(
		makeDeleteIngredientEndpoint(is),
		decodeDeleteIngredientRequest,
		encodeResponse,
		opts...,
	)

	router := mux.NewRouter()
	router.Handle("/ingredient/{id}", ingredientHandler).Methods("GET")
	router.Handle("/ingredient/", createIngredientHandler).Methods("POST")
	router.Handle("/ingredient/{id}", updateIngredientHandler).Methods("PUT")
	router.Handle("/ingredient/{id}", deleteIngredientHandler).Methods("DELETE")
	return router
}

func decodeIngredientRequest(_ context.Context, r *http.Request) (interface{}, error) {
	if id, err := helper.GetRequestParam(r, "id"); err == nil {
		return ingredientRequest{id}, nil
	}
	return nil, badRequest
}

func decodeCreateIngredientRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request createIngredientRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}
	return request, nil
}

func decodeUpdateIngredientRequest(_ context.Context, r *http.Request) (interface{}, error) {
	if id, err := helper.GetRequestParam(r, "id"); err == nil {
		var body struct {
			Ingredient domain.Ingredient `json:"ingredient"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err == nil {
			return updateIngredientRequest{id, body.Ingredient}, nil
		}
	}
	return nil, badRequest
}

func decodeDeleteIngredientRequest(_ context.Context, r *http.Request) (interface{}, error) {
	if id, err := helper.GetRequestParam(r, "id"); err == nil {
		return deleteIngredientRequest{id}, nil
	}
	return nil, badRequest
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
