package main

import (
	"log"
	"os"

	"os/signal"

	"github.com/Placeblock/nostalgicraft-discord/internal/commands"
	"github.com/Placeblock/nostalgicraft-discord/internal/realtime"
	entityuser "github.com/Placeblock/nostalgicraft-discord/internal/repository/entityUser"
	teamrole "github.com/Placeblock/nostalgicraft-discord/internal/repository/teamRole"
	"github.com/Placeblock/nostalgicraft-discord/internal/service"
	"github.com/Placeblock/nostalgicraft-discord/pkg/config"
	"github.com/Placeblock/nostalgicraft-discord/pkg/models"
	"github.com/Placeblock/nostalgicraft-ems/pkg/storage"
	"github.com/bwmarrin/discordgo"
	"gopkg.in/yaml.v3"
)

func main() {
	f, err := os.Open("config.yml")
	if err != nil {
		log.Panicln("Failed to load config.yml", err)
	}
	defer f.Close()

	var cfg config.Config
	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(&cfg)
	if err != nil {
		log.Panicln("Failed to parse config.yml", err)
	}

	session, err := discordgo.New("Bot " + cfg.Token)
	if err != nil {
		log.Panicln("Failed to create Discord Bot", err)
	}
	db := storage.Connect()
	db.AutoMigrate(&models.TeamRole{}, &models.UserEntity{})
	entityUserRepo := entityuser.NewMysqlEntityUserRepository(db)
	teamRoleRepo := teamrole.NewMysqlTeamRoleRepository(db)
	userEntityService := service.NewEntityUserService(entityUserRepo)
	teamRoleService := service.NewTeamRoleService(teamRoleRepo)
	subscriber := realtime.NewSubscriber(&cfg, userEntityService, teamRoleService, session)
	go subscriber.Listen()
	listen(cfg, session)
}

func listen(cfg config.Config, session *discordgo.Session) {
	session.Identify.Intents = discordgo.IntentsNone
	err := session.Open()
	if err != nil {
		log.Panicln("Failed to start Discord Bot", err)
	}

	commandRegistry := commands.NewCommandRegistry(session, cfg.Guild)
	commandRegistry.RegisterDefaultHandler()
	commandRegistry.Register(commands.NewVerifyCommand())

	sigch := make(chan os.Signal, 1)
	signal.Notify(sigch, os.Interrupt)
	log.Println("Press Ctrl+C to exit")
	<-sigch

	commandRegistry.DeleteCommands()
	err = session.Close()
	if err != nil {
		log.Printf("Could not close session gracefully: %s", err)
	}
}
