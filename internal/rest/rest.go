package http

import (
	"github.com/codelix/ems/internal/rest/middleware"
	"github.com/codelix/ems/internal/rest/routes/entities"
	"github.com/codelix/ems/internal/rest/routes/tokens"
	"github.com/codelix/ems/internal/service/entity"
	"github.com/codelix/ems/internal/service/token"
	"github.com/gin-gonic/gin"
)

type updateEntityParams struct {
	NewName string `json:"new_name" binding:"required"`
}

type createEntityParams struct {
	Name string `json:"name" binding:"required"`
}

type createTokenParams struct {
	EntityId uint `json:"entityId" binding:"required"`
}

type HttpServer struct {
	entityService entity.EntityService
	tokenService  token.TokenService
}

func NewHttpServer(entityService entity.EntityService, tokenService token.TokenService) *HttpServer {
	return &HttpServer{entityService, tokenService}
}

func (server *HttpServer) Serve() {
	r := gin.Default()
	r.SetTrustedProxies(nil)
	r.Use(middleware.ErrorHandler())

	entitiesGroup := r.Group("entities")
	entities.Handle(entitiesGroup, &server.entityService)
	tokenGroup := r.Group("token")
	tokens.Handle(tokenGroup, &server.tokenService)

	r.Run("localhost:3006")
}
