package gorm

import (
	"errors"
	"gorm.io/gorm"
	"what_cook/domain"
)

type CrudRepository struct {
	Db       *gorm.DB
	newModel func() interface{}
}

func (cr *CrudRepository) Save(model interface{}) error {
	res := cr.Db.Create(model)
	return res.Error
}

func (cr *CrudRepository) Get(id uint) (interface{}, error) {
	model := cr.newModel()
	res := cr.Db.First(model, id)

	if errors.Is(res.Error, gorm.ErrRecordNotFound) {
		return nil, domain.ModelNotFoundError
	}
	return model, res.Error
}

func (cr *CrudRepository) Update(id uint, model interface{}) error {
	// check model
	currentModel, e := cr.Get(id)
	if e != nil {
		return e
	}

	res := cr.Db.Model(currentModel).Updates(model)
	return res.Error
}

func (cr *CrudRepository) Delete(id uint) error {
	// check model
	model, e := cr.Get(id)
	if e != nil {
		return e
	}

	res := cr.Db.Delete(model, id)
	return res.Error
}

type IngredientRepository struct {
	CrudRepository
}

func NewIngredientRepository(db *gorm.DB) domain.IngredientRepository {
	return &IngredientRepository{
		CrudRepository{
			Db: db,
			newModel: func() interface{} {
				return &domain.Ingredient{}
			},
		},
	}
}

type FoodRepository struct {
	CrudRepository
}

func (f *FoodRepository) FindByIngredients(ingredientNames []string) ([]domain.FoodRecommendation, error) {
	// ingredients data
	var ingredients []domain.Ingredient
	err := f.Db.Find(&ingredients, "name IN ?", ingredientNames).Error
	if err != nil {
		return nil, err
	}
	ingredientIds := make([]uint, len(ingredients))
	for i, ingredient := range ingredients {
		ingredientIds[i] = ingredient.ID
	}
	// order sorted by ingredients availability
	type foodResult struct {
		Id               uint
		IngredientsCount uint
		IngredientsHas   uint
	}
	var order []foodResult
	err = f.Db.Table("food").
		Unscoped().
		Select("food.id as id",
			"count(iw.id) as ingredients_count",
			"count(has_iw.id) as ingredients_has").
		Joins("JOIN ingredient_weights iw on food.id = iw.food_id").
		Joins("LEFT OUTER JOIN ingredient_weights has_iw on "+
			"food.id = has_iw.food_id and "+
			"iw.id = has_iw.id and "+
			"has_iw.ingredient_id IN ?", ingredientIds).
		Group("name").
		Having("ingredients_has > ?", 0).
		Order("ingredients_has DESC").Order("ingredients_count").
		Scan(&order).Error
	if err != nil {
		return nil, err
	}
	// food data
	foodIds := make([]uint, len(order))
	for i, foodResult := range order {
		foodIds[i] = foodResult.Id
	}

	foods := make([]domain.Food, len(order))
	f.Db.Table("food").Preload("IngredientWeights.Ingredient").Where("id IN ?", foodIds).Find(&foods)
	foodMap := make(map[uint]domain.Food)
	for _, food := range foods {
		foodMap[food.ID] = food
	}
	// make foodRecommendations
	foodRecommendations := make([]domain.FoodRecommendation, len(order))
	for i, foodResult := range order {
		foodRecommendations[i] = domain.FoodToFoodRecommendation(foodMap[foodResult.Id], ingredients)
	}
	return foodRecommendations, nil
}

func NewFoodRepository(db *gorm.DB) domain.FoodRepository {
	return &FoodRepository{CrudRepository{
		Db: db,
		newModel: func() interface{} {
			return &domain.Food{}
		},
	}}
}
