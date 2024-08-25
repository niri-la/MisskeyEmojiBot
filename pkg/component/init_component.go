package component

import (
	"MisskeyEmojiBot/pkg/config"
	"MisskeyEmojiBot/pkg/handler"
	"MisskeyEmojiBot/pkg/repository"

	"github.com/bwmarrin/discordgo"
)

type InitComponen interface {
}

type initComponen struct {
	config      config.Config
	discordRepo repository.DiscordRepository
}

func NewInitComponent(config config.Config, discordRepo repository.DiscordRepository) handler.Component {
	return &initComponen{config: config, discordRepo: discordRepo}
}

func (c *initComponen) GetCommand() *discordgo.ApplicationCommand {
	return &discordgo.ApplicationCommand{
		Name: "init_channel",
	}
}

func (c *initComponen) Execute(s *discordgo.Session, i *discordgo.InteractionCreate) {
	channelID := i.MessageComponentData().Values[0]

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Flags: discordgo.MessageFlagsEphemeral,
			Content: "選択チャンネル <#" + i.MessageComponentData().Values[0] + ">\n" +
				"初期設定を行いました。",
		},
	})

	println("init component")

	s.ChannelMessageSendComplex(channelID,
		&discordgo.MessageSend{
			Content: "こんにちは！絵文字申請チャンネルへようこそ！\n",
			Components: []discordgo.MessageComponent{
				discordgo.ActionsRow{
					Components: []discordgo.MessageComponent{
						&discordgo.Button{
							Label:    "絵文字の申請をする / Requset emoji",
							CustomID: "new_emoji_channel",
							Style:    discordgo.PrimaryButton,
							Emoji: discordgo.ComponentEmoji{
								Name: "🏗️",
							},
						},
					},
				},
			},
		},
	)

	overwrites := []*discordgo.PermissionOverwrite{
		{
			ID:   c.config.ModeratorID,
			Type: discordgo.PermissionOverwriteTypeRole,
			Allow: discordgo.PermissionViewChannel |
				discordgo.PermissionSendMessages,
		},
		{
			ID:   c.config.BotID,
			Type: discordgo.PermissionOverwriteTypeRole,
			Allow: discordgo.PermissionViewChannel |
				discordgo.PermissionSendMessages,
		},
		{
			ID:   i.GuildID,
			Type: discordgo.PermissionOverwriteTypeRole,
			Deny: discordgo.PermissionViewChannel,
		},
	}

	parent, err := s.Channel(i.ChannelID)

	if err != nil {
		c.discordRepo.ReturnFailedMessage(i, "Could not retrieve channel")
		return
	}

	channel, err := c.discordRepo.FindChannelByName(i.GuildID, c.config.ModerationChannelName)

	if err != nil {
		channel, err = s.GuildChannelCreateComplex(i.GuildID, discordgo.GuildChannelCreateData{
			Type:                 discordgo.ChannelTypeGuildText,
			Name:                 c.config.ModerationChannelName,
			ParentID:             parent.ParentID,
			PermissionOverwrites: overwrites,
		})
		if err != nil {
			c.discordRepo.ReturnFailedMessage(i, "Could not create channel")
			return
		}

		s.ChannelMessageSend(
			channel.ID,
			": モデレーション用チャンネルです。\nここに各種申請のスレッドが生成されます。",
		)
		return
	}
}
