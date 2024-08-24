package command

import (
	"MisskeyEmojiBot/pkg/entity"

	"github.com/bwmarrin/discordgo"
)

func NewNirilaCommand() *entity.Command {
	return &entity.Command{
		Command: &discordgo.ApplicationCommand{
			Name:        "ni_rilana",
			Description: "Misskey Emoji Bot © 2023 KineL",
		},
		Executer: func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			Execute(s, i)
		},
	}
}

func Execute(s *discordgo.Session, i *discordgo.InteractionCreate) {
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
