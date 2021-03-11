package domain

import (
	"errors"
	"gorm.io/gorm"
)

type Food struct {
	gorm.Model
	Name              string
	Description       string
	IngredientWeights []IngredientWeight
}

func (f *Food) Equal(food interface{}) bool {
	f2, ok := food.(*Food)
	if !ok {
		return false
	}
	if f2 == f || f2.ID == f.ID {
		return true
	}
	return f2.Name == f.Name
}

type FoodRecommendation struct {
	Food              Food
	HasIngredients    []Ingredient
	AbsentIngredients []Ingredient
}

func FoodToFoodRecommendation(f Food, ingredients []Ingredient) FoodRecommendation {
	foodRecommendation := FoodRecommendation{
		Food:              f,
		HasIngredients:    make([]Ingredient, 0),
		AbsentIngredients: make([]Ingredient, 0),
	}
	// search has or absent ingredients
	for _, ingredientWeight := range f.IngredientWeights {
		absent := true
		for _, ingredient := range ingredients {
			if ingredient.ID == ingredientWeight.IngredientID {
				foodRecommendation.HasIngredients = append(foodRecommendation.HasIngredients, ingredient)
				absent = false
				break
			}
		}
		if absent {
			foodRecommendation.AbsentIngredients = append(foodRecommendation.AbsentIngredients, ingredientWeight.Ingredient)
		}
	}
	return foodRecommendation
}

var FoodNotFoundError = errors.New("food not found")

type FoodRepository interface {
	CrudRepository
	FindByIngredients(ingredients []string) ([]FoodRecommendation, error)
}

type FoodService interface {
	Save(food *Food) error
	Update(id uint, food *Food) error
	Delete(id uint) error
	Get(id uint) (*Food, error)
	FindByIngredients(ingredients []string) ([]FoodRecommendation, error)
}
