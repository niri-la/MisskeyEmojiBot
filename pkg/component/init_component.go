package component

import (
	"MisskeyEmojiBot/pkg/handler"
	"MisskeyEmojiBot/pkg/repository"

	"github.com/bwmarrin/discordgo"
)

type InitComponen interface {
}

type initComponen struct {
	discordRepo repository.DiscordRepository
}

func NewInitComponent(discordRepo repository.DiscordRepository) handler.Component {
	return &initComponen{discordRepo: discordRepo}
}

func (i *initComponen) GetCommand() *discordgo.ApplicationCommand {
	return &discordgo.ApplicationCommand{
		Name: "init_channel",
	}
}

func (*initComponen) Execute(s *discordgo.Session, i *discordgo.InteractionCreate) {
	channelID := i.MessageComponentData().Values[0]

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Flags: discordgo.MessageFlagsEphemeral,
			Content: "選択チャンネル <#" + i.MessageComponentData().Values[0] + ">\n" +
				"初期設定を行いました。",
		},
	})

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
			ID:   ModeratorID,
			Type: discordgo.PermissionOverwriteTypeRole,
			Allow: discordgo.PermissionViewChannel |
				discordgo.PermissionSendMessages,
		},
		{
			ID:   BotID,
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
		returnFailedMessage(s, i, "Could not retrieve channel")
		return
	}

	channel, err := s.GuildChannelCreateComplex(GuildID, discordgo.GuildChannelCreateData{
		Type:                 discordgo.ChannelTypeGuildText,
		Name:                 ModerationChannelName,
		ParentID:             parent.ParentID,
		PermissionOverwrites: overwrites,
	})

	s.ChannelMessageSend(
		channel.ID,
		": モデレーション用チャンネルです。\nここに各種申請のスレッドが生成されます。",
	)

	logger.Debug(":: Create a moderation channel")

	return
}
