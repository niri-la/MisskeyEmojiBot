package command

import (
	"MisskeyEmojiBot/pkg/handler"
	"MisskeyEmojiBot/pkg/repository"

	"github.com/bwmarrin/discordgo"
)

type NirilaCommand interface {
}

type nirilaCommand struct {
	discordRepo repository.DiscordRepository
}

func NewNirilaCommand(discordRepo repository.DiscordRepository) handler.CommandInterface {
	return &nirilaCommand{discordRepo: discordRepo}
}

func (c *nirilaCommand) GetCommand() *discordgo.ApplicationCommand {
	return &discordgo.ApplicationCommand{
		Name:        "ni_rilana",
		Description: "Misskey Emoji Bot © 2023 KineL",
	}
}

func (c *nirilaCommand) Execute(s *discordgo.Session, i *discordgo.InteractionCreate) {
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Flags: discordgo.MessageFlagsEphemeral,
			Content: "::: Misskey Emoji Bot \n" +
				": Created by ni_rila (KineL)\n" +
				": © 2023 KineL\n" +
				":::::::::: \n",
		},
	})
}
