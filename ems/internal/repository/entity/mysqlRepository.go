package entity

import (
	"context"
	"errors"
	"fmt"

	"github.com/Placeblock/nostalgicraft-ems/pkg/models"
	"gorm.io/gorm"
)

type MysqlEntityRepository struct {
	db *gorm.DB
}

func NewMysqlEntityRepository(db *gorm.DB) *MysqlEntityRepository {
	return &MysqlEntityRepository{db}
}

func (repo *MysqlEntityRepository) GetEntity(ctx context.Context, entity *models.Entity) error {
	if err := repo.db.WithContext(ctx).Where(entity).First(entity).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil
		}
		return fmt.Errorf("getEntity %+v: %v", entity, err.Error())
	}
	return nil
}

func (repo *MysqlEntityRepository) CreateEntity(ctx context.Context, entity *models.Entity) error {
	if err := repo.db.WithContext(ctx).Create(entity).Error; err != nil {
		return fmt.Errorf("createEntity %s: %v", entity.Name, err.Error())
	}
	return nil
}

func (repo *MysqlEntityRepository) UpdateEntity(ctx context.Context, entity models.Entity) error {
	if err := repo.db.WithContext(ctx).Updates(&entity).Error; err != nil {
		return fmt.Errorf("updateEntity %s: %v", entity.Name, err.Error())
	}
	return nil
}

func (repo *MysqlEntityRepository) GetEntities(ctx context.Context, filter models.Entity) (*[]models.Entity, error) {
	var entities []models.Entity
	if err := repo.db.WithContext(ctx).Where(filter).Find(&entities).Error; err != nil {
		return nil, fmt.Errorf("getEntities: %v", err.Error())
	}
	return &entities, nil
}
