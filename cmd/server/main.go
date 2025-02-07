package main

import (
	"log"

	"github.com/codelix/ems/internal/http"
	"github.com/codelix/ems/internal/realtime"
	"github.com/codelix/ems/internal/repository/entity"
	"github.com/codelix/ems/internal/repository/token"
	sentity "github.com/codelix/ems/internal/service/entity"
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
	entityService := sentity.NewEntityService(entityRepo, publisher)
	tokenService := stoken.NewTokenService(tokenRepo)
	httpServer := http.NewHttpServer(*entityService, *tokenService)
	httpServer.Serve()
}
