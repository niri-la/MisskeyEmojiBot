package main

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

var (
	Components         = make([]*discordgo.ApplicationCommand, 0)
	ComponentsHandlers = make(map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate))
)

func addComponent(command *discordgo.ApplicationCommand, fn func(s *discordgo.Session, i *discordgo.InteractionCreate)) {
	_, exist := ComponentsHandlers[command.Name]
	if exist {
		panic(fmt.Sprintf("Error: [%s] is already existed.", command.Name))
	}
	ComponentsHandlers[command.Name] = fn
	Components = append(Components, command)
}
