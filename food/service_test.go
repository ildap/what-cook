package food

import (
	"testing"
	"what_cook/domain"
	"what_cook/gorm"
	"what_cook/helper"
)

var (
	foodService domain.FoodService
	testFood    domain.Food
)

func TestMain(m *testing.M) {
	// setup
	db := gorm.SqliteDbSession(gorm.DSN_SQLITE_TEST)
	gorm.ClearData(db)
	foodRepository := gorm.NewFoodRepository(db)
	foodService = NewFoodService(foodRepository)
	testFood = gorm.CreateRandomFood(db)
	// run tests
	m.Run()
}

func TestService_Get(t *testing.T) {
	// check not found
	_, err := foodService.Get(0)
	if err != domain.ModelNotFoundError {
		t.Error("err is not equal error ", domain.ModelNotFoundError)
	}
	// check found
	food, err := foodService.Get(testFood.ID)
	if err != nil {
		t.Error(err)
	}
	if !testFood.Equal(food) {
		t.Error("test food is not equal returned value")
	}
}

func TestService_Save(t *testing.T) {
	// check save
	food := gorm.RandomFood()
	err := foodService.Save(&food)
	if err != nil {
		t.Error(err)
	}
	// check availability
	savedfood, err := foodService.Get(food.ID)
	if err != nil {
		t.Error(err)
	}
	if !food.Equal(savedfood) {
		t.Error("food is not equal saved food")
	}
}

func TestService_Update(t *testing.T) {
	testFood.Name = "test_food" + helper.RandomName()
	err := foodService.Update(testFood.ID, &testFood)
	if err != nil {
		t.Error(err)
	}
	// check update
	food, err := foodService.Get(testFood.ID)
	if food.Name != testFood.Name {
		t.Error("name is not updated")
	}
}

func TestService_Delete(t *testing.T) {
	err := foodService.Delete(testFood.ID)
	if err != nil {
		t.Error(err)
	}
	// check delete
	_, err = foodService.Get(testFood.ID)
	if err != domain.ModelNotFoundError {
		t.Error("test food is not deleted")
	}
}
