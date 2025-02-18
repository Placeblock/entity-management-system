package commands

import (
	"fmt"
	"log"

	"github.com/bwmarrin/discordgo"
)

type Command struct {
	AppCommand *discordgo.ApplicationCommand
	Handler    func(s *discordgo.Session, i *discordgo.InteractionCreate)
}

type CommandRegistry struct {
	session  *discordgo.Session
	guild    string
	commands map[string]Command
}

func NewCommandRegistry(session *discordgo.Session, guild string) *CommandRegistry {
	return &CommandRegistry{session, guild, map[string]Command{}}
}

func (reg *CommandRegistry) RegisterDefaultHandler() {
	reg.session.AddHandler(reg.handle)
}

func (reg *CommandRegistry) Register(command Command) {
	cmd, err := reg.session.ApplicationCommandCreate(reg.session.State.User.ID, reg.guild, command.AppCommand)
	if err != nil {
		log.Panicln("Could not create command", err)
	}
	command.AppCommand = cmd
	reg.commands[command.AppCommand.Name] = command
	fmt.Printf("Registered command %s\n", command.AppCommand.Name)
}

func (reg *CommandRegistry) handle(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if i.Type == discordgo.InteractionApplicationCommand {
		if cmd, ok := reg.commands[i.ApplicationCommandData().Name]; ok {
			cmd.Handler(s, i)
		}
	}
}

func (reg *CommandRegistry) DeleteCommands() {
	for _, v := range reg.commands {
		err := reg.session.ApplicationCommandDelete(reg.session.State.User.ID, reg.guild, v.AppCommand.ID)
		if err != nil {
			log.Printf("Could not delete command %s\n", err)
		}
	}
}
