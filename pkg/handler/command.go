package handler

import (
	"MisskeyEmojiBot/pkg/config"
	"MisskeyEmojiBot/pkg/repository"
	"fmt"
	"os"

	"github.com/bwmarrin/discordgo"
)

var ()

type CommandInterface interface {
	GetCommand() *discordgo.ApplicationCommand
	Execute(s *discordgo.Session, i *discordgo.InteractionCreate)
}

type CommandHandler struct {
	config          config.Config
	discordRepo     repository.DiscordRepository
	commands        []*discordgo.ApplicationCommand
	commandHandlers map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate)
}

func NewCommandHandler(config config.Config, discordRepo repository.DiscordRepository) *CommandHandler {
	return &CommandHandler{
		config:          config,
		discordRepo:     discordRepo,
		commands:        []*discordgo.ApplicationCommand{},
		commandHandlers: map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){},
	}
}

func (c *CommandHandler) RegisterCommand(ci CommandInterface) {
	_, exist := c.commandHandlers[ci.GetCommand().Name]
	if exist {
		panic(fmt.Sprintf("Error: [%s] is already existed.", ci.GetCommand().Name))
	}

	cmd, err := c.discordRepo.GetSession().ApplicationCommandCreate(c.config.AppID, c.config.GuildID, ci.GetCommand())
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	c.commandHandlers[ci.GetCommand().Name] = ci.Execute
	c.commands = append(c.commands, cmd)
}

func (c *CommandHandler) GetCommands() []*discordgo.ApplicationCommand {
	return c.commands
}

func (c *CommandHandler) Handle(s *discordgo.Session, i *discordgo.InteractionCreate) {
	command, exist := c.commandHandlers[i.ApplicationCommandData().Name]
	if !exist {
		return
	}
	command(s, i)
}
