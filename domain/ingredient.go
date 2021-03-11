package domain

import (
	"gorm.io/gorm"
)

type Ingredient struct {
	gorm.Model
	Name     string  `validate:"nonzero"`
	Calories float64 `validate:"min=0"`
}

func (i *Ingredient) Equal(i2 interface{}) bool {
	other, ok := i2.(*Ingredient)
	if !ok {
		return false
	}
	if other == i {
		return true
	}
	if i.ID != 0 && i.ID == other.ID {
		return true
	}
	return i.Name == other.Name &&
		i.Calories == other.Calories &&
		i.UpdatedAt.Equal(other.UpdatedAt) &&
		i.CreatedAt.Equal(other.CreatedAt)
}

type IngredientWeight struct {
	gorm.Model
	FoodID       uint
	IngredientID uint
	Ingredient   Ingredient
	Weight       float64 //kg
}

type IngredientRepository interface {
	CrudRepository
}

type IngredientService interface {
	Save(ingredient *Ingredient) error
	Update(id uint, ingredient *Ingredient) error
	Delete(id uint) error
	Get(id uint) (*Ingredient, error)
}
