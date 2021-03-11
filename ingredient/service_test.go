package ingredient

import (
	"testing"
	"what_cook/domain"
	"what_cook/gorm"
	"what_cook/helper"
)

var (
	ingredientService domain.IngredientService
	testIngredient    domain.Ingredient
)

func TestMain(m *testing.M) {
	// setup
	db := gorm.SqliteDbSession(gorm.DSN_SQLITE_TEST)
	gorm.ClearData(db)
	ingredientRepository := gorm.NewIngredientRepository(db)
	testIngredient = gorm.CreateRandomIngredient(db)
	ingredientService = NewService(ingredientRepository)
	// run tests
	m.Run()
}

func TestService_Get(t *testing.T) {
	// check not found
	_, err := ingredientService.Get(0)
	if err != domain.ModelNotFoundError {
		t.Error("err is not equal error ", domain.ModelNotFoundError)
	}
	// check found
	ingredient, err := ingredientService.Get(testIngredient.ID)
	if err != nil {
		t.Error(err)
	}
	if !testIngredient.Equal(ingredient) {
		t.Error("test ingredient is not equal returned value")
	}
}

func TestService_Save(t *testing.T) {
	// check save
	ingredient := gorm.RandomIngredient()
	err := ingredientService.Save(&ingredient)
	if err != nil {
		t.Error(err)
	}
	// check availability
	savedIngredient, err := ingredientService.Get(ingredient.ID)
	if err != nil {
		t.Error(err)
	}
	if !ingredient.Equal(savedIngredient) {
		t.Error("ingredient is not equal saved ingredient")
	}
}

func TestService_Update(t *testing.T) {
	testIngredient.Name = "test_ingredient" + helper.RandomName()
	err := ingredientService.Update(testIngredient.ID, &testIngredient)
	if err != nil {
		t.Error(err)
	}
	// check update
	ingredient, err := ingredientService.Get(testIngredient.ID)
	if ingredient.Name != testIngredient.Name {
		t.Error("name is not updated")
	}
}

func TestService_Delete(t *testing.T) {
	err := ingredientService.Delete(testIngredient.ID)
	if err != nil {
		t.Error(err)
	}
	// check delete
	_, err = ingredientService.Get(testIngredient.ID)
	if err != domain.ModelNotFoundError {
		t.Error("test ingredient is not deleted")
	}
}
