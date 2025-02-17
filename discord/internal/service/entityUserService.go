package service

import (
	"context"

	entityuser "github.com/Placeblock/nostalgicraft-discord/internal/repository/entityUser"
	"github.com/Placeblock/nostalgicraft-discord/pkg/models"
)

type EntityUserService struct {
	repo *entityuser.EntityUserRepository
}

func NewEntityUserService(repository *entityuser.EntityUserRepository) *EntityUserService {
	return &EntityUserService{repo: repository}
}

func (service *EntityUserService) GetEntityIdByUserId(ctx context.Context, userId string) (uint, error) {
	return (*service.repo).GetEntityIdByUserId(ctx, userId)
}
func (service *EntityUserService) GetUserIdByEntityId(ctx context.Context, entityId uint) (string, error) {
	return (*service.repo).GetUserIdByEntityId(ctx, entityId)
}

func (service *EntityUserService) CreateUserEntity(ctx context.Context, userEntity models.UserEntity) error {
	return (*service.repo).CreateUserEntity(ctx, userEntity)
}
