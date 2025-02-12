package commands

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

func NewVerifyCommand() Command {
	return Command{
		&verifyCommand,
		handleVerifyCommand,
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

func handleVerifyCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
	pin := i.ApplicationCommandData().Options[0].StringValue()
	fmt.Println(pin)
}
