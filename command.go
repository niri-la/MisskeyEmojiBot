package main

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

var (
	Commands        = make([]*discordgo.ApplicationCommand, 0)
	CommandHandlers = make(map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate))
)

func addCommand(command *discordgo.ApplicationCommand, fn func(s *discordgo.Session, i *discordgo.InteractionCreate)) {
	_, exist := CommandHandlers[command.Name]
	if exist {
		panic(fmt.Sprintf("Error: [%s] is already existed.", command.Name))
	}
	CommandHandlers[command.Name] = fn
	Commands = append(Commands, command)
}
