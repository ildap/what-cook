package gorm

import (
	"gorm.io/gorm"
	"testing"
	"what_cook/domain"
	"what_cook/helper"
)

var (
	db                   *gorm.DB
	testIngredient       domain.Ingredient
	testFood             domain.Food
	ingredientRepository domain.IngredientRepository
	foodRepository       domain.FoodRepository
)

func TestMain(m *testing.M) {
	// setup
	db = SqliteDbSession(DSN_SQLITE_TEST)
	ClearData(db)
	ingredientRepository = NewIngredientRepository(db)
	testIngredient = CreateRandomIngredient(db)
	foodRepository = NewFoodRepository(db)
	testFood = CreateRandomFood(db)
	// run tests
	m.Run()
}

func TestIngredientRepository_Get(t *testing.T) {
	ingredient, err := ingredientRepository.Get(testIngredient.ID)
	if err != nil {
		t.Error(err)
	}

	i, ok := ingredient.(*domain.Ingredient)
	if !ok {
		t.Error("type error")
	}

	if !testIngredient.Equal(i) {
		t.Error("test ingredient is not equal returned object")
	}
}

func TestIngredientRepository_Save(t *testing.T) {
	ingredient := RandomIngredient()
	// check save
	err := ingredientRepository.Save(&ingredient)
	if err != nil {
		t.Error(err)
	}
	// check availability
	ingredient2, err := ingredientRepository.Get(ingredient.ID)
	if err != nil {
		t.Error(err)
	}
	if !ingredient.Equal(ingredient2) {
		t.Error("ingredient is not equal returned object")
	}
}

func TestIngredientRepository_Update(t *testing.T) {
	testIngredient.Name = "test_ingredient" + helper.RandomName()
	err := ingredientRepository.Update(testIngredient.ID, testIngredient)
	if err != nil {
		t.Error(err)
	}
	// check update
	ingredient, err := ingredientRepository.Get(testIngredient.ID)
	if ingredient.(*domain.Ingredient).Name != testIngredient.Name {
		t.Error("name is not updated")
	}
}

func TestIngredientRepository_Delete(t *testing.T) {
	err := ingredientRepository.Delete(testIngredient.ID)
	if err != nil {
		t.Error(err)
	}
	// check delete
	_, err = ingredientRepository.Get(testIngredient.ID)
	if err != domain.ModelNotFoundError {
		t.Error("test ingredient is not deleted")
	}
}

func TestFoodRepository_Get(t *testing.T) {
	food, err := foodRepository.Get(testFood.ID)
	if err != nil {
		t.Error(err)
	}

	i, ok := food.(*domain.Food)
	if !ok {
		t.Error("type error")
	}

	if !testFood.Equal(i) {
		t.Error("test food is not equal returned object")
	}
}

func TestFoodRepository_Save(t *testing.T) {
	food := RandomIngredient()
	// check save
	err := foodRepository.Save(&food)
	if err != nil {
		t.Error(err)
	}
	// check availability
	food2, err := foodRepository.Get(food.ID)
	if err != nil {
		t.Error(err)
	}
	if !food.Equal(food2) {
		t.Error("food is not equal returned object")
	}
}

func TestFoodRepository_Update(t *testing.T) {
	testFood.Name = "test_food" + helper.RandomName()
	err := foodRepository.Update(testFood.ID, testFood)
	if err != nil {
		t.Error(err)
	}
	// check update
	food, err := foodRepository.Get(testFood.ID)
	if food.(*domain.Food).Name != testFood.Name {
		t.Error("name is not updated")
	}
}

func TestFoodRepository_Delete(t *testing.T) {
	err := foodRepository.Delete(testFood.ID)
	if err != nil {
		t.Error(err)
	}
	// check delete
	_, err = foodRepository.Get(testFood.ID)
	if err != domain.ModelNotFoundError {
		t.Error("test food is not deleted")
	}
}

func TestFoodRepository_FindByIngredients(t *testing.T) {
	// create test data
	ingredients := [4]domain.Ingredient{
		CreateRandomIngredient(db),
		CreateRandomIngredient(db),
		CreateRandomIngredient(db),
		CreateRandomIngredient(db),
	}
	foods := make([]*domain.Food, len(ingredients))
	for i := 0; i < len(ingredients); i++ {
		foods[i] = &domain.Food{
			Name:              helper.RandomName(),
			IngredientWeights: make([]domain.IngredientWeight, 0),
		}
		// for check order
		for j := i; j < len(ingredients); j++ {
			foods[i].IngredientWeights = append(foods[i].IngredientWeights,
				domain.IngredientWeight{IngredientID: ingredients[j].ID},
			)
		}
		foodRepository.Save(foods[i])
	}
	// test
	ingredientNames := make([]string, len(ingredients))
	for i, ingredient := range ingredients {
		ingredientNames[i] = ingredient.Name
	}
	r, err := foodRepository.FindByIngredients(ingredientNames)
	if err != nil {
		t.Error(err)
	}
	// check order
	for i, foodRecommendation := range r {
		if !foodRecommendation.Food.Equal(foods[i]) {
			t.Error("wrong order")
		}
	}
}
