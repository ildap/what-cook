package ingredient

import (
	"context"
	"github.com/go-kit/kit/endpoint"
	"gopkg.in/validator.v2"
	"strconv"
	"what_cook/domain"
)

type ingredientRequest struct {
	ID uint
}

type ingredientResponse struct {
	Ingredient *domain.Ingredient
	Err        error `json:"err,omitempty"`
}

func (i ingredientResponse) error() error {
	return i.Err
}

func makeIngredientEndpoint(is domain.IngredientService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		var req, _ = request.(ingredientRequest)
		ingredient, e := is.Get(req.ID)
		return ingredientResponse{ingredient, e}, nil
	}
}

type createIngredientRequest struct {
	Ingredient domain.Ingredient
}

type createIngredientResponse struct {
	IngredientID string
	Err          error `json:"err,omitempty"`
}

func (c createIngredientResponse) error() error {
	return c.Err
}

func makeCreateIngredientEndpoint(is domain.IngredientService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		var req = request.(createIngredientRequest)
		if error := validator.Validate(req); error != nil {
			return createIngredientResponse{"", error}, nil
		}
		saveError := is.Save(&req.Ingredient)
		return createIngredientResponse{strconv.Itoa(int(req.Ingredient.ID)), saveError}, nil
	}
}

type updateIngredientRequest struct {
	ID         uint
	Ingredient domain.Ingredient
}

type updateIngredientResponse struct {
	Err error `json:"err,omitempty"`
}

func (u updateIngredientResponse) error() error {
	return u.Err
}

func makeUpdateIngredientEndpoint(is domain.IngredientService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		var req = request.(updateIngredientRequest)
		if error := validator.Validate(req); error != nil {
			return updateIngredientResponse{error}, nil
		}
		e := is.Update(req.ID, &req.Ingredient)
		return updateIngredientResponse{e}, nil
	}
}

type deleteIngredientRequest struct {
	ID uint
}

type deleteIngredientResponse struct {
	Err error `json:"err,omitempty"`
}

func (d deleteIngredientResponse) error() error {
	return d.Err
}

func makeDeleteIngredientEndpoint(is domain.IngredientService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		var req = request.(deleteIngredientRequest)
		deleteError := is.Delete(req.ID)
		return deleteIngredientResponse{deleteError}, nil
	}
}
