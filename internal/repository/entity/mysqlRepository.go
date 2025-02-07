package entity

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/codelix/ems/internal/storage"
	"github.com/codelix/ems/pkg/models"
)

type MysqlEntityRepository struct {
}

func (repo *MysqlEntityRepository) GetEntity(ctx context.Context, id int64) (*models.Entity, error) {
	entity := models.Entity{ID: &id}
	if err := storage.DB.QueryRow("SELECT name FROM entities WHERE id = ?;", id).Scan(&entity.Name); err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("getEntity %d: Invalid Entity ID", id)
		}
		return nil, fmt.Errorf("getEntity %d: %v", id, err)
	}
	return &entity, nil
}

func (repo *MysqlEntityRepository) CreateEntity(ctx context.Context, entity *models.Entity) error {
	result, err := storage.DB.Exec("INSERT INTO entities (name) VALUES (?);", entity.Name)
	if err != nil {
		return fmt.Errorf("createEntity %s: %v", entity.Name, err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("createEntity %s: %v", entity.Name, err)
	}
	entity.ID = &id
	return nil
}

func (repo *MysqlEntityRepository) UpdateEntity(ctx context.Context, entity models.Entity) error {
	_, err := storage.DB.Exec("UPDATE entities SET name = ? WHERE id = ?", entity.Name, entity.ID)
	if err != nil {
		return fmt.Errorf("updateEntity %s: %v", entity.Name, err)
	}
	return nil
}

func (repo *MysqlEntityRepository) GetEntities(ctx context.Context) (*[]models.Entity, error) {
	var entities []models.Entity

	rows, err := storage.DB.Query("SELECT id, name FROM entities")
	if err != nil {
		return nil, fmt.Errorf("getEntities: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var entity models.Entity
		if err := rows.Scan(&entity.ID, &entity.Name); err != nil {
			return nil, fmt.Errorf("getEntities: %v", err)
		}
		entities = append(entities, entity)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("getEntities: %v", err)
	}
	return &entities, nil
}
