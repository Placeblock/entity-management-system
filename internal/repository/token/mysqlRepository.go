package token

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/codelix/ems/internal/storage"
	"github.com/codelix/ems/pkg/models"
)

type MysqlTokenRepository struct {
}

func (repo *MysqlTokenRepository) CreateToken(ctx context.Context, token models.Token) error {
	_, err := repo.GetToken(ctx, token.Pin)
	if err == nil {
		return fmt.Errorf("createToken: Could not create token")
	}
	_, err = storage.DB.Exec("INSERT INTO tokens (entityId, pin, createdAt) VALUES (?, ?, ?) ON DUPLICATE KEY UPDATE pin = ?, createdAt = ?;",
		token.EntityID, token.Pin, token.CreatedAt, token.Pin, token.CreatedAt)
	if err != nil {
		return fmt.Errorf("createToken %d: %v", token.EntityID, err)
	}
	return nil
}

func (repo *MysqlTokenRepository) GetToken(ctx context.Context, pin string) (*models.Token, error) {
	token := models.Token{Pin: pin}
	if err := storage.DB.QueryRow("SELECT entityId, createdAt FROM tokens WHERE pin = ?;", pin).Scan(&token.EntityID, &token.CreatedAt); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("getToken %s: %v", pin, err)
	}
	if token.CreatedAt.Before(time.Now().Add(-time.Duration(2) * time.Minute)) {
		storage.DB.Exec("DELETE FROM tokens WHERE entityId = ?", token.EntityID)
		return nil, nil
	}
	return &token, nil
}
