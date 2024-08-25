package handler

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

var (
	Commands        = make([]*discordgo.ApplicationCommand, 0)
	CommandHandlers = make(map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate))
)

type CommandInterface interface {
	GetCommand() *discordgo.ApplicationCommand
	Execute(s *discordgo.Session, i *discordgo.InteractionCreate)
}

type CommandHandler struct {
}

func (c *CommandHandler) RegisterCommand(ci CommandInterface) {
	_, exist := CommandHandlers[ci.GetCommand().Name]
	if exist {
		panic(fmt.Sprintf("Error: [%s] is already existed.", ci.GetCommand().Name))
	}
	CommandHandlers[ci.GetCommand().Name] = ci.Execute
	Commands = append(Commands, ci.GetCommand())
}
