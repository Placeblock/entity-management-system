package rest

import (
	"github.com/codelix/ems/internal/rest/middleware"
	"github.com/codelix/ems/internal/rest/routes/entities"
	teamentities "github.com/codelix/ems/internal/rest/routes/teamEntities"
	"github.com/codelix/ems/internal/rest/routes/teams"
	"github.com/codelix/ems/internal/rest/routes/tokens"
	"github.com/codelix/ems/internal/service/entity"
	"github.com/codelix/ems/internal/service/team"
	teamentity "github.com/codelix/ems/internal/service/teamEntity"
	"github.com/codelix/ems/internal/service/token"
	"github.com/gin-gonic/gin"
)

type HttpServer struct {
	entityService      entity.EntityService
	tokenService       token.TokenService
	teamService        team.TeamService
	teamEntitiyService teamentity.TeamEntityService
}

func NewHttpServer(entityService entity.EntityService,
	tokenService token.TokenService,
	teamService team.TeamService,
	teamEntitiyService teamentity.TeamEntityService) *HttpServer {
	return &HttpServer{entityService, tokenService, teamService, teamEntitiyService}
}

func (server *HttpServer) Serve() {
	r := gin.Default()
	r.SetTrustedProxies(nil)
	r.Use(middleware.ErrorHandler())

	entitiesGroup := r.Group("entities")
	entities.Handle(entitiesGroup, &server.entityService)
	tokenGroup := r.Group("token")
	tokens.Handle(tokenGroup, &server.tokenService)
	teamsGroup := r.Group("teams")
	teams.Handle(teamsGroup, &server.teamService)
	teamEntitiesGroup := r.Group("team-entities")
	teamentities.Handle(teamEntitiesGroup, &server.teamEntitiyService)

	r.Run("localhost:3006")
}
