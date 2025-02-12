package main

import (
	"log"

	"github.com/Placeblock/nostalgicraft-ems/internal/realtime"
	"github.com/Placeblock/nostalgicraft-ems/internal/repository/entity"
	member "github.com/Placeblock/nostalgicraft-ems/internal/repository/member"
	"github.com/Placeblock/nostalgicraft-ems/internal/repository/team"
	"github.com/Placeblock/nostalgicraft-ems/internal/repository/token"
	"github.com/Placeblock/nostalgicraft-ems/internal/rest"
	sentity "github.com/Placeblock/nostalgicraft-ems/internal/service/entity"
	smember "github.com/Placeblock/nostalgicraft-ems/internal/service/member"
	steam "github.com/Placeblock/nostalgicraft-ems/internal/service/team"
	stoken "github.com/Placeblock/nostalgicraft-ems/internal/service/token"
	"github.com/Placeblock/nostalgicraft-ems/internal/storage"
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
