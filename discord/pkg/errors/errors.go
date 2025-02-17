package errors

import (
	"fmt"

	"github.com/Placeblock/nostalgicraft-discord/pkg/models"
)

type ErrNotFound struct{}

func (e ErrNotFound) Error() string {
	return "Not found"
}

type ErrEntityNotLinked struct {
	EntityID uint
}

func (e ErrEntityNotLinked) Error() string {
	return fmt.Sprintf("UserID %d is not linked to any discord user", e.EntityID)
}

type ErrUserEntityAlreadyExists struct {
	UserEntity models.UserEntity
}

func (e ErrUserEntityAlreadyExists) Error() string {
	return fmt.Sprintf("UserEntity %+v already exists", e.UserEntity)
}
