package main

import (
	"log"
	"os"

	"os/signal"

	"github.com/Placeblock/nostalgicraft-discord/internal/commands"
	"github.com/Placeblock/nostalgicraft-discord/internal/realtime"
	"github.com/Placeblock/nostalgicraft-discord/pkg/config"
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

	go realtime.Listen()
	startDiscordBot(cfg)
}

func startDiscordBot(cfg config.Config) {
	session, err := discordgo.New("Bot " + cfg.Token)
	if err != nil {
		log.Panicln("Failed to create Discord Bot", err)
	}
	session.Identify.Intents = discordgo.IntentsNone
	err = session.Open()
	if err != nil {
		log.Panicln("Failed to start Discord Bot", err)
	}

	commandRegistry := commands.NewCommandRegistry(session, cfg.TestGuild)
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
