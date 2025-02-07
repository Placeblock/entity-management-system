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
	Db *gorm.DB
}

func NewMysqlTokenRepository(db *gorm.DB) *MysqlTokenRepository {
	return &MysqlTokenRepository{db}
}

func (repo *MysqlTokenRepository) CreateToken(ctx context.Context, token models.Token) error {
	if err := repo.Db.WithContext(ctx).Delete(models.Token{EntityID: token.EntityID}).Error; err != nil {
		return fmt.Errorf("createToken %d: %v", token.EntityID, err.Error())
	}
	if err := repo.Db.WithContext(ctx).Create(&token).Error; err != nil {
		return fmt.Errorf("createToken %d: %v", token.EntityID, err.Error())
	}
	return nil
}

func (repo *MysqlTokenRepository) GetToken(ctx context.Context, pin string) (*models.Token, error) {
	var token models.Token
	if err := repo.Db.WithContext(ctx).Where("pin = ?", pin).First(&token).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, fmt.Errorf("getToken %s: %v", pin, err.Error())
	}
	if token.CreatedAt.Before(time.Now().Add(-time.Duration(2) * time.Minute)) {
		repo.Db.WithContext(ctx).Delete(token)
		return nil, nil
	}
	return &token, nil
}
