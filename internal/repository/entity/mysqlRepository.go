package entity

import (
	"context"
	"fmt"

	"github.com/codelix/ems/pkg/models"
	"gorm.io/gorm"
)

type MysqlEntityRepository struct {
	db *gorm.DB
}

func NewMysqlEntityRepository(db *gorm.DB) *MysqlEntityRepository {
	return &MysqlEntityRepository{db}
}

func (repo *MysqlEntityRepository) GetEntity(ctx context.Context, id uint) (*models.Entity, error) {
	var entity models.Entity
	if err := repo.db.WithContext(ctx).First(&entity, &id).Error; err != nil {
		return nil, fmt.Errorf("getEntity %d: %s", id, err.Error())
	}
	return &entity, nil
}

func (repo *MysqlEntityRepository) CreateEntity(ctx context.Context, entity *models.Entity) error {
	if err := repo.db.WithContext(ctx).Create(entity).Error; err != nil {
		return fmt.Errorf("createEntity %s: %v", entity.Name, err.Error())
	}
	return nil
}

func (repo *MysqlEntityRepository) UpdateEntity(ctx context.Context, entity models.Entity) error {
	if err := repo.db.WithContext(ctx).Save(&entity).Error; err != nil {
		return fmt.Errorf("updateEntity %s: %v", entity.Name, err.Error())
	}
	return nil
}

func (repo *MysqlEntityRepository) GetEntities(ctx context.Context) (*[]models.Entity, error) {
	var entities []models.Entity
	if err := repo.db.WithContext(ctx).Find(&entities).Error; err != nil {
		return nil, fmt.Errorf("getEntities: %v", err.Error())
	}
	return &entities, nil
}
