package rest

import (
	"github.com/codelix/ems/internal/rest/middleware"
	"github.com/codelix/ems/internal/rest/routes/entities"
	"github.com/codelix/ems/internal/rest/routes/teams"
	"github.com/codelix/ems/internal/rest/routes/tokens"
	"github.com/codelix/ems/internal/service/entity"
	member "github.com/codelix/ems/internal/service/member"
	"github.com/codelix/ems/internal/service/team"
	"github.com/codelix/ems/internal/service/token"
	"github.com/gin-gonic/gin"
)

type HttpServer struct {
	entityService      entity.EntityService
	tokenService       token.TokenService
	teamService        team.TeamService
	teamEntitiyService member.MemberService
}

func NewHttpServer(entityService entity.EntityService,
	tokenService token.TokenService,
	teamService team.TeamService,
	memberService member.MemberService) *HttpServer {
	return &HttpServer{entityService, tokenService, teamService, memberService}
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
	teams.Handle(teamsGroup, &server.teamService, &server.teamEntitiyService)

	r.Run("localhost:3006")
}
