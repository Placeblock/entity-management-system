package main

import (
	"log"

	"github.com/codelix/ems/internal/realtime"
	"github.com/codelix/ems/internal/repository/entity"
	"github.com/codelix/ems/internal/repository/team"
	teamentity "github.com/codelix/ems/internal/repository/teamEntity"
	"github.com/codelix/ems/internal/repository/token"
	"github.com/codelix/ems/internal/rest"
	sentity "github.com/codelix/ems/internal/service/entity"
	steam "github.com/codelix/ems/internal/service/team"
	steamentity "github.com/codelix/ems/internal/service/teamEntity"
	stoken "github.com/codelix/ems/internal/service/token"
	"github.com/codelix/ems/internal/storage"
	zmq "github.com/pebbe/zmq4"
)

func main() {
	zctx, err := zmq.NewContext()
	if err != nil {
		log.Fatal("Could not create ZMQ Context")
	}
	publisher := realtime.NewPublisher(zctx)
	go publisher.Listen()
	storage.Connect()
	db := storage.Connect()
	entityRepo := entity.NewMysqlEntityRepository(db)
	tokenRepo := token.NewMysqlTokenRepository(db)
	teamRepo := team.NewMysqlTeamRepository(db)
	teamEntityRepo := teamentity.NewMysqlTeamEntityRepository(db)

	entityService := sentity.NewEntityService(entityRepo, publisher)
	tokenService := stoken.NewTokenService(tokenRepo)
	teamService := steam.NewMysqlTeamRepository(teamRepo, publisher)
	teamEntityService := steamentity.NewTeamEntityService(teamEntityRepo, publisher)

	httpServer := rest.NewHttpServer(*entityService, *tokenService, *teamService, *teamEntityService)
	httpServer.Serve()
}
