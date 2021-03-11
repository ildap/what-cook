package food

import (
	"context"
	"github.com/go-kit/kit/endpoint"
	"gopkg.in/validator.v2"
	"strconv"
	"what_cook/domain"
)

type foodRequest struct {
	ID uint
}

type foodResponse struct {
	Food *domain.Food
	Err  error
}

func (f foodResponse) error() error {
	return f.Err
}

func makeFoodEndpoint(foodService domain.FoodService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, error error) {
		foodId := request.(foodRequest).ID
		food, err := foodService.Get(foodId)
		return foodResponse{food, err}, error
	}
}

type createFoodRequest struct {
	Food *domain.Food
}

type createFoodResponse struct {
	FoodId string
	Err    error `json:"err,omitempty"`
}

func (r createFoodResponse) error() error {
	return r.Err
}

func makeCreateFoodEndpoint(foodService domain.FoodService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		var error error
		req := request.(createFoodRequest)
		if error = validator.Validate(req); error != nil {
			return createFoodResponse{"", error}, err
		}
		error = foodService.Save(req.Food)
		if error != nil {
			return createFoodResponse{"", error}, err
		}
		return createFoodResponse{strconv.Itoa(int(req.Food.ID)), nil}, err
	}
}

type updateFoodRequest struct {
	Id   uint
	Food *domain.Food
}

type updateFoodResponse struct {
	Err error `json:"err,omitempty"`
}

func (r updateFoodResponse) error() error {
	return r.Err
}

func makeUpdateFoodEndpoint(foodService domain.FoodService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(updateFoodRequest)
		if validateError := validator.Validate(req); validateError != nil {
			return updateFoodResponse{validateError}, nil
		}
		updateError := foodService.Update(req.Id, req.Food)
		return updateFoodResponse{updateError}, err
	}
}

type deleteFoodRequest struct {
	Id uint
}

type deleteFoodResponse struct {
	Err error `json:"err,omitempty"`
}

func (r deleteFoodResponse) error() error {
	return r.Err
}

func makeDeleteFoodEndpoint(foodService domain.FoodService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(deleteFoodRequest)
		deleteError := foodService.Delete(req.Id)
		return deleteFoodResponse{deleteError}, err
	}
}

type foodsByIngredientsRequest struct {
	Ingredients []string
}

type foodsByIngredientsResponse struct {
	Foods []domain.FoodRecommendation
	Err   error
}

func (f *foodsByIngredientsResponse) error() error {
	return f.Err
}

func makeFoodsByIngredientEndpoint(foodService domain.FoodService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(foodsByIngredientsRequest)
		foodRecommendations, foodServiceError := foodService.FindByIngredients(req.Ingredients)
		if foodServiceError != nil {
			return foodsByIngredientsResponse{nil, foodServiceError}, err
		}
		return foodsByIngredientsResponse{foodRecommendations, nil}, err
	}
}
