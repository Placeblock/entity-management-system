package token

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/codelix/ems/pkg/models"
	"gorm.io/gorm"
)

type MysqlTokenRepository struct {
	db *gorm.DB
}

func NewMysqlTokenRepository(db *gorm.DB) *MysqlTokenRepository {
	return &MysqlTokenRepository{db}
}

func (repo *MysqlTokenRepository) CreateToken(ctx context.Context, token models.Token) error {
	if err := repo.db.WithContext(ctx).Delete(models.Token{EntityID: token.EntityID}).Error; err != nil {
		return fmt.Errorf("createToken %d: %v", token.EntityID, err.Error())
	}
	if err := repo.db.WithContext(ctx).Create(&token).Error; err != nil {
		return fmt.Errorf("createToken %d: %v", token.EntityID, err.Error())
	}
	return nil
}

func (repo *MysqlTokenRepository) GetToken(ctx context.Context, token *models.Token) (*models.Token, error) {
	if err := repo.db.WithContext(ctx).First(&token).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, fmt.Errorf("getToken %+v: %v", token, err.Error())
	}
	if token.CreatedAt.Before(time.Now().Add(-time.Duration(2) * time.Minute)) {
		repo.db.WithContext(ctx).Delete(token)
		return nil, nil
	}
	return token, nil
}
