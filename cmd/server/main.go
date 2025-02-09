package main

import (
	"log"

	"github.com/codelix/ems/internal/realtime"
	"github.com/codelix/ems/internal/repository/entity"
	member "github.com/codelix/ems/internal/repository/member"
	"github.com/codelix/ems/internal/repository/team"
	"github.com/codelix/ems/internal/repository/token"
	"github.com/codelix/ems/internal/rest"
	sentity "github.com/codelix/ems/internal/service/entity"
	smember "github.com/codelix/ems/internal/service/member"
	steam "github.com/codelix/ems/internal/service/team"
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
	db := storage.Connect()
	entityRepo := entity.NewMysqlEntityRepository(db)
	tokenRepo := token.NewMysqlTokenRepository(db)
	teamRepo := team.NewMysqlTeamRepository(db)
	memberRepo := member.NewMysqlMemberRepository(db)

	entityService := sentity.NewEntityService(entityRepo, publisher)
	tokenService := stoken.NewTokenService(tokenRepo)
	teamService := steam.NewMysqlTeamRepository(teamRepo, publisher)
	memberService := smember.NewMemberService(memberRepo, publisher)

	httpServer := rest.NewHttpServer(*entityService, *tokenService, *teamService, *memberService)
	httpServer.Serve()
}
