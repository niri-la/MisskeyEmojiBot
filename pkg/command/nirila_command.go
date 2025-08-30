package command

import (
	"github.com/bwmarrin/discordgo"

	"MisskeyEmojiBot/pkg/handler"
	"MisskeyEmojiBot/pkg/repository"
)

type NirilaCommand interface {
}

type nirilaCommand struct {
	discordRepo repository.DiscordRepository
	version     string
}

func NewNirilaCommand(discordRepo repository.DiscordRepository, version string) handler.CommandInterface {
	return &nirilaCommand{discordRepo: discordRepo, version: version}
}

func (c *nirilaCommand) GetCommand() *discordgo.ApplicationCommand {
	return &discordgo.ApplicationCommand{
		Name:        "ni_rilana",
		Description: "Misskey Emoji Bot © 2024 KineL",
	}
}

func (c *nirilaCommand) Execute(s *discordgo.Session, i *discordgo.InteractionCreate) {
	_ = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Flags: discordgo.MessageFlagsEphemeral,
			Content: "# ::: Misskey Emoji Bot v" + c.version + "\n" +
				"::: Created by ni_rila (KineL)\n" +
				"::: © 2024 KineL\n" +
				"### ------------ \n",
		},
	})
}
