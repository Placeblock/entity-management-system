package token

import (
	"context"
	"time"

	"github.com/codelix/ems/internal/repository/token"
	"github.com/codelix/ems/pkg/models"
	"github.com/codelix/ems/tools"
)

type TokenService struct {
	tokenRepository token.TokenRepository
}

func NewTokenService(repo token.TokenRepository) *TokenService {
	return &TokenService{repo}
}

func (service *TokenService) CreateToken(ctx context.Context, entityId uint) (*models.Token, error) {
	pin := tools.GenSix()
	token := models.Token{EntityID: entityId, CreatedAt: time.Now(), Pin: pin}
	err := service.tokenRepository.CreateToken(ctx, token)
	if err != nil {
		return nil, err
	}
	return &token, nil
}

func (service *TokenService) GetToken(ctx context.Context, pin string) (*models.Token, error) {
	return service.tokenRepository.GetToken(ctx, pin)
}
