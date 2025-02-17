package repository

import (
	"context"
	"errors"

	cerr "github.com/Placeblock/nostalgicraft-discord/pkg/errors"
	"github.com/Placeblock/nostalgicraft-discord/pkg/models"
	"gorm.io/gorm"
)

type MysqlEntityUserRepository struct {
	db *gorm.DB
}

func NewMysqlEntityUserRepository(db *gorm.DB) *MysqlEntityUserRepository {
	return &MysqlEntityUserRepository{db}
}

func (repo *MysqlEntityUserRepository) GetEntityIdByUserId(ctx context.Context, userId string) (uint, error) {
	userEntity := models.UserEntity{UserID: userId}
	if err := repo.db.WithContext(ctx).First(&userEntity).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return 0, &cerr.ErrInvalidUserID{UserID: userId}
		}
		return 0, err
	}
	return userEntity.EntityID, nil
}

func (repo *MysqlEntityUserRepository) GetUserIdByEntityId(ctx context.Context, entityId uint) (string, error) {
	userEntity := models.UserEntity{EntityID: entityId}
	if err := repo.db.WithContext(ctx).First(&userEntity).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", cerr.ErrNotFound{}
		}
		return "", err
	}
	return userEntity.UserID, nil
}

func (repo *MysqlEntityUserRepository) CreateUserEntity(ctx context.Context, userEntity models.UserEntity) error {
	if err := repo.db.WithContext(ctx).First(&userEntity).Error; err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return cerr.ErrUserEntityAlreadyExists{UserEntity: userEntity}
		}
		return err
	}
	return nil
}
