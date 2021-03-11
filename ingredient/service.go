package ingredient

import (
	"what_cook/domain"
)

type service struct {
	ingredientRepository domain.IngredientRepository
}

func (s *service) Get(id uint) (*domain.Ingredient, error) {
	i, e := s.ingredientRepository.Get(id)
	if i == nil {
		return nil, e
	}
	return i.(*domain.Ingredient), e
}

func (s *service) Save(food *domain.Ingredient) error {
	return s.ingredientRepository.Save(food)
}

func (s *service) Update(id uint, ingredient *domain.Ingredient) error {
	return s.ingredientRepository.Update(id, ingredient)
}

func (s *service) Delete(id uint) error {
	return s.ingredientRepository.Delete(id)
}

func NewService(r domain.IngredientRepository) domain.IngredientService {
	return &service{ingredientRepository: r}
}
