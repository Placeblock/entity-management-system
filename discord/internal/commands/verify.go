package commands

import (
	"context"
	"fmt"
	"time"

	"github.com/Placeblock/nostalgicraft-discord/internal"
	"github.com/Placeblock/nostalgicraft-discord/internal/service"
	"github.com/Placeblock/nostalgicraft-discord/pkg/models"
	"github.com/Placeblock/nostalgicraft-discord/pkg/rest"
	emsmodels "github.com/Placeblock/nostalgicraft-ems/pkg/models"
	"github.com/bwmarrin/discordgo"
	"github.com/carlmjohnson/requests"
)

func NewVerifyCommand(service *service.EntityUserService) Command {
	return Command{
		&verifyCommand,
		func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			handleVerifyCommand(s, i, service)
		},
	}
}

var verifyCommand = discordgo.ApplicationCommand{
	Name:        "verify",
	Description: "Verify your Nostalgicraft-ID and link your Discord-Account",
	Options: []*discordgo.ApplicationCommandOption{
		{
			Name:        "pin",
			Description: "The pin provided to you to link a new service",
			Type:        discordgo.ApplicationCommandOptionString,
			Required:    true,
			MaxLength:   6,
			// MinLength is Int Pointer??
		},
	},
}

func handleVerifyCommand(s *discordgo.Session, i *discordgo.InteractionCreate, service *service.EntityUserService) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	pin := i.ApplicationCommandData().Options[0].StringValue()
	fmt.Println(pin)

	var response rest.APIResponse[emsmodels.Token]
	err := requests.
		URL("http://"+internal.Config.Ems.RestHost+":"+internal.Config.Ems.RestPort).
		Pathf("/tokens").
		Param("pin", pin).
		ToJSON(&response).
		Fetch(ctx)
	if err != nil {
		fmt.Println("Could not receive Token when receiving Pin", err)
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Title:   "Error",
				Content: "```ansi\n\u001b[1;31mError \u001b[0mCould not verify Pin. \u001b[1;31mPlease report this to the admin!```",
			},
		})
		return
	}
	if response.Data.EntityID == 0 {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Title:   "Fail!",
				Content: "```ansi\n\u001b[1;31mFail! \u001b[0mYou entered the \u001b[1;31mwrong cridentials\u001b[0m!```",
			},
		})
		return
	}
	err = service.CreateUserEntity(ctx, models.UserEntity{EntityID: response.Data.EntityID, UserID: i.Member.User.ID})
	if err != nil {
		fmt.Println("Could not create UserEntity when receiving Pin", err)
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Title:   "Error",
				Content: "```ansi\n\u001b[1;31mError \u001b[0mCould not link account. \u001b[1;31mPlease report this to the admin!```",
			},
		})
		return
	}
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Title:   "Success!",
			Content: "```ansi\n\u001b[1;32mSuccessfully \u001b[0mlinked account!```",
		},
	})
}
