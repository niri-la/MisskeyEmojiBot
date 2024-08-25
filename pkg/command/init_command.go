package command

import (
	"MisskeyEmojiBot/pkg/handler"
	"MisskeyEmojiBot/pkg/repository"

	"github.com/bwmarrin/discordgo"
)

type InitCommand interface {
}

type initCommand struct {
	discordRepo repository.DiscordRepository
}

func NewInitCommand(discordRepo repository.DiscordRepository) handler.CommandInterface {
	return &initCommand{discordRepo: discordRepo}
}

func (c *initCommand) GetCommand() *discordgo.ApplicationCommand {
	return &discordgo.ApplicationCommand{
		Name:        "init",
		Description: "絵文字申請の初期化を行います。",
		Type:        discordgo.ChatApplicationCommand,
	}
}

func (c *initCommand) Execute(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if !c.discordRepo.HasRole("GuildID", *i.Member.User, "admin") {
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
