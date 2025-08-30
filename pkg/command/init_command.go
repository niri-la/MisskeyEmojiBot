package command

import (
	"github.com/bwmarrin/discordgo"

	"MisskeyEmojiBot/pkg/config"
	"MisskeyEmojiBot/pkg/handler"
	"MisskeyEmojiBot/pkg/repository"
)

type InitCommand interface {
}

type initCommand struct {
	config      config.Config
	discordRepo repository.DiscordRepository
}

func NewInitCommand(config config.Config, discordRepo repository.DiscordRepository) handler.CommandInterface {
	return &initCommand{config: config, discordRepo: discordRepo}
}

func (c *initCommand) GetCommand() *discordgo.ApplicationCommand {
	return &discordgo.ApplicationCommand{
		Name:        "init",
		Description: "絵文字申請の初期化を行います。",
		Type:        discordgo.ChatApplicationCommand,
	}
}

func (c *initCommand) Execute(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if !c.discordRepo.HasRole(c.config.GuildID, *i.Member.User, c.config.ModeratorID) {
		_ = c.discordRepo.ReturnFailedMessage(i, "No permission.")
		return
	}

	_ = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
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
