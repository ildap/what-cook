package gorm

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"what_cook/domain"
	"what_cook/helper"
)

const DSN_SQLITE = "what_cook.db"
const DSN_SQLITE_TEST = "what_cook_test.db"

func SqliteDbSession(dsn string) *gorm.DB {
	db, dberr := gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	if dberr != nil {
		panic(dberr)
	}

	dberr = db.AutoMigrate(&domain.Food{}, &domain.Ingredient{}, &domain.IngredientWeight{})
	if dberr != nil {
		panic(dberr)
	}

	return db
}

func CreateRandomIngredient(db *gorm.DB) domain.Ingredient {
	randomIngredient := RandomIngredient()
	db.Create(&randomIngredient)
	return randomIngredient
}

func RandomIngredient() domain.Ingredient {
	return domain.Ingredient{
		Name:     "test_ingredient" + helper.RandomName(),
		Calories: 0,
	}
}

func CreateRandomFood(db *gorm.DB) domain.Food {
	randomFood := RandomFood()
	db.Create(&randomFood)
	return randomFood
}

func RandomFood() domain.Food {
	return domain.Food{
		Name: "test_food" + helper.RandomName(),
		IngredientWeights: []domain.IngredientWeight{
			{
				Ingredient: RandomIngredient(),
				Weight:     0.1,
			},
		},
	}
}

func ClearData(db *gorm.DB) {
	db.Session(&gorm.Session{AllowGlobalUpdate: true}).
		Unscoped().Delete(&domain.Ingredient{})

	db.Session(&gorm.Session{AllowGlobalUpdate: true}).
		Unscoped().Delete(&domain.IngredientWeight{})

	db.Session(&gorm.Session{AllowGlobalUpdate: true}).
		Unscoped().Delete(&domain.Food{})
}
