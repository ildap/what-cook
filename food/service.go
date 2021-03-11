package food

import (
	"what_cook/domain"
)

type service struct {
	repository domain.FoodRepository
}

func (s service) FindByIngredients(ingredients []string) ([]domain.FoodRecommendation, error) {
	return s.repository.FindByIngredients(ingredients)
}

func (s service) Save(food *domain.Food) error {
	return s.repository.Save(food)
}

func (s service) Update(id uint, food *domain.Food) error {
	return s.repository.Update(id, food)
}

func (s service) Delete(id uint) error {
	return s.repository.Delete(id)
}

func (s service) Get(id uint) (*domain.Food, error) {
	food, err := s.repository.Get(id)
	if food == nil {
		return nil, err
	}
	return food.(*domain.Food), err
}

func NewFoodService(repository domain.FoodRepository) domain.FoodService {
	return &service{
		repository: repository,
	}
}
