package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"os/signal"

	"github.com/Placeblock/nostalgicraft-discord/internal/commands"
	"github.com/Placeblock/nostalgicraft-discord/internal/realtime"
	entityuser "github.com/Placeblock/nostalgicraft-discord/internal/repository/entityUser"
	teamrole "github.com/Placeblock/nostalgicraft-discord/internal/repository/teamRole"
	"github.com/Placeblock/nostalgicraft-discord/internal/service"
	"github.com/Placeblock/nostalgicraft-discord/pkg/config"
	"github.com/Placeblock/nostalgicraft-discord/pkg/models"
	"github.com/Placeblock/nostalgicraft-discord/pkg/rest"
	emsmodels "github.com/Placeblock/nostalgicraft-ems/pkg/models"
	"github.com/Placeblock/nostalgicraft-ems/pkg/storage"
	"github.com/bwmarrin/discordgo"
	"gopkg.in/yaml.v3"

	"github.com/carlmjohnson/requests"
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
	db.AutoMigrate(&models.TeamData{}, &models.UserEntity{})
	entityUserRepo := entityuser.NewMysqlEntityUserRepository(db)
	teamDataRepo := teamrole.NewMysqlTeamDataRepository(db)
	userEntityService := service.NewEntityUserService(entityUserRepo)
	teamDataService := service.NewTeamDataService(teamDataRepo)
	subscriber := realtime.NewSubscriber(&cfg, userEntityService, teamDataService, session)
	go subscriber.Listen()
	listen(cfg, session, userEntityService, teamDataService)
}

func listen(cfg config.Config, session *discordgo.Session, entityUserService *service.EntityUserService, teamDataService *service.TeamDataService) {
	session.Identify.Intents = discordgo.IntentsGuildMessages | discordgo.IntentsMessageContent
	err := session.Open()
	if err != nil {
		log.Panicln("Failed to start Discord Bot", err)
	}

	commandRegistry := commands.NewCommandRegistry(session, cfg.Guild)
	commandRegistry.RegisterDefaultHandler()
	commandRegistry.Register(commands.NewVerifyCommand())

	session.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if i.Type == discordgo.InteractionMessageComponent {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			customId := i.MessageComponentData().CustomID
			accepted := strings.Split(customId, "-")[0] == "accept"
			serializedInviteId := strings.Split(customId, "-")[2]
			requestURL := fmt.Sprintf("http://%s:%s/invites/%s", cfg.Ems.RestHost, cfg.Ems.RestPort, serializedInviteId)
			var method string
			if accepted {
				method = http.MethodPost
			} else {
				method = http.MethodDelete
			}
			req, err := http.NewRequestWithContext(ctx, method, requestURL, nil)
			if err != nil {
				fmt.Println("Could not create Request when receiving Message interaction", err)
				return
			}
			_, err = http.DefaultClient.Do(req)
			if err != nil {
				fmt.Println("Could not request when receiving Message interaction", err)
				return
			}
		}
	})

	session.AddHandler(func(s *discordgo.Session, i *discordgo.MessageCreate) {
		if i.Author.Bot {
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		teamData, err := teamDataService.GetTeamDataByChannelId(ctx, i.ChannelID)
		if err != nil {
			fmt.Println("Could not get TeamData when receiving Message", err)
			return
		}
		if teamData.TeamID == 0 {
			return
		}
		err = s.ChannelMessageDelete(teamData.ChannelID, i.Message.ID)
		if err != nil {
			fmt.Println("Could not delete Message when receiving Message", err)
			return
		}

		entityId, err := entityUserService.GetEntityIdByUserId(ctx, i.Author.ID)
		if err != nil {
			fmt.Println("Could not get EntityID when receiving Message", err)
			return
		}
		if entityId == 0 {
			return
		}

		var response rest.APIResponse[emsmodels.Member]
		err = requests.
			URL("http://"+cfg.Ems.RestHost+":"+cfg.Ems.RestPort).
			Pathf("/entities/%d/member", entityId).
			ToJSON(&response).Fetch(ctx)
		if err != nil {
			fmt.Println("Could not fetch Member when receiving Message", err)
			return
		}
		err = requests.
			URL("http://"+cfg.Ems.RestHost+":"+cfg.Ems.RestPort).
			Pathf("/members/%d/message", response.Data.ID).
			BodyJSON(rest.CreateMessageDto{Message: i.Message.Content}).
			Fetch(ctx)
		if err != nil {
			fmt.Println("Could not send Message when receiving Message", err)
			return
		}
	})

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
