package command

import (
	"MisskeyEmojiBot/pkg/entity"

	"github.com/bwmarrin/discordgo"
)

type initCommand struct {
}

func NewInitCommand() *entity.Command {
	InitCommand := &initCommand{}
	return &entity.Command{
		Command: &discordgo.ApplicationCommand{
			Name:        "init",
			Description: "絵文字申請用の初期化を行います",
		},
		Executer: func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			InitCommand.Execute(s, i)
		},
	}
}

func (c *initCommand) Execute(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if !hasPermission(*i.Member.User) {
		returnFailedMessage(s, i, "No permission.")
		return
	}

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "こんにちは！絵文字申請の初期化を行います。\n" +
				"絵文字申請用のチャンネルを指定してください！",
			Flags: discordgo.MessageFlagsEphemeral,
			Components: []discordgo.MessageComponent{
				discordgo.ActionsRow{
					Components: []discordgo.MessageComponent{
						discordgo.SelectMenu{
							MenuType:     discordgo.ChannelSelectMenu,
							CustomID:     "init_channel",
							Placeholder:  "申請を行うチャンネルを選択してください",
							ChannelTypes: []discordgo.ChannelType{discordgo.ChannelTypeGuildText},
						},
					},
				},
			},
		},
	})
}
