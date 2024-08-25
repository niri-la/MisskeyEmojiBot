package component

import (
	"MisskeyEmojiBot/pkg/handler"
	"MisskeyEmojiBot/pkg/repository"

	"github.com/bwmarrin/discordgo"
	"github.com/sirupsen/logrus"
)

type EmojiRequestRetryComponen interface {
}

type emojiRequestRetryComponen struct {
	discordRepo repository.DiscordRepository
}

func NewEmojiRequestRetryComponen(discordRepo repository.DiscordRepository) handler.Component {
	return &emojiRequestRetryComponen{discordRepo: discordRepo}
}
func (c *emojiRequestRetryComponen) GetCommand() *discordgo.ApplicationCommand {
	return &discordgo.ApplicationCommand{
		Name: "emoji_request_retry",
	}
}

func (c *emojiRequestRetryComponen) Execute(s *discordgo.Session, i *discordgo.InteractionCreate) {
	channel, _ := s.Channel(i.ChannelID)
	emoji, err := GetEmoji(channel.Name[6:])
	if err != nil {
		s.ChannelMessageSend(
			channel.ID,
			"設定に失敗しました。管理者に問い合わせを行ってください。 #04a\n",
		)
	}

	if emoji.IsRequested {
		s.ChannelMessageSend(
			channel.ID,
			"既に絵文字は申請されています。新たな申請を作成してください。\n",
		)
		return
	}

	emoji.IsSensitive = false
	emoji.RequestState = workflow[0]
	emoji.ResponseState = workflow[0]

	c.discordRepo.DeleteChannel(emoji.FilePath)

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Flags:   discordgo.MessageFlagsEphemeral,
			Content: "リクエストを初期化します。\n",
		},
	})

	logger.WithFields(logrus.Fields{
		"user":    i.Member.User.Username,
		"channel": channel.Name,
		"id":      emoji.ID,
		"name":    emoji.Name,
	}).Debug("Request reset.")

	first(emoji, s, channel.ID)
}
