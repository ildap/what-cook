package domain

import "errors"

type CrudRepository interface {
	Save(model interface{}) error
	Update(id uint, model interface{}) error
	Delete(id uint) error
	Get(id uint) (interface{}, error)
}

var ModelNotFoundError = errors.New("not found")
