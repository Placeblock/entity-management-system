package token

import (
	"context"

	"github.com/codelix/ems/pkg/models"
)

type TokenRepository interface {
	CreateToken(ctx context.Context, token models.Token) error
	GetToken(ctx context.Context, pin string) (*models.Token, error)
}
