package entity

import "github.com/bwmarrin/discordgo"

type Command struct {
	Command  *discordgo.ApplicationCommand
	Executer func(s *discordgo.Session, i *discordgo.InteractionCreate)
}

type CommandInterface interface {
	Execute(s *discordgo.Session, i *discordgo.InteractionCreate)
}
