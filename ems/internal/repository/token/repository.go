package token

import (
	"context"

	"github.com/Placeblock/nostalgicraft-ems/pkg/models"
)

type TokenRepository interface {
	CreateToken(ctx context.Context, token models.Token) error
	GetToken(ctx context.Context, token *models.Token) (*models.Token, error)
}
