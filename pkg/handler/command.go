package handler

import (
	"MisskeyEmojiBot/pkg/entity"
	"fmt"

	"github.com/bwmarrin/discordgo"
)

var (
	Commands        = make([]*discordgo.ApplicationCommand, 0)
	CommandHandlers = make(map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate))
)

type CommandHandler struct {
}

func (c *CommandHandler) RegisterCommand(command *entity.Command) {
	_, exist := CommandHandlers[command.Command.Name]
	if exist {
		panic(fmt.Sprintf("Error: [%s] is already existed.", command.Command.Name))
	}
	CommandHandlers[command.Command.Name] = command.Executer
	Commands = append(Commands, command.Command)
}
