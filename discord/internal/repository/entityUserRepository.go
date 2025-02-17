package repository

import (
	"context"

	"github.com/Placeblock/nostalgicraft-discord/pkg/models"
)

type EntityUserRepository interface {
	GetEntityIdByUserId(ctx context.Context, userId string) (uint, error)
	GetUserIdByEntityId(ctx context.Context, entityId uint) (string, error)
	CreateUserEntity(ctx context.Context, userEntity models.UserEntity) error
}
