package component

import (
	"MisskeyEmojiBot/pkg/handler"
	"MisskeyEmojiBot/pkg/repository"

	"github.com/bwmarrin/discordgo"
	"github.com/sirupsen/logrus"
)

type NsfwInactiveComponen interface {
}

type nsfwInactiveComponen struct {
	discordRepo repository.DiscordRepository
}

func NewNsfwInactiveComponent(discordRepo repository.DiscordRepository) handler.Component {
	return &nsfwInactiveComponen{discordRepo: discordRepo}
}

func (n *nsfwInactiveComponen) GetCommand() *discordgo.ApplicationCommand {
	return &discordgo.ApplicationCommand{
		Name: "nsfw_no",
	}
}

func (n *nsfwInactiveComponen) Execute(s *discordgo.Session, i *discordgo.InteractionCreate) {
	channel, _ := s.Channel(i.ChannelID)
	emoji, err := GetEmoji(channel.Name[6:])
	if err != nil {
		s.ChannelMessageSend(
			channel.ID,
			"設定に失敗しました。管理者に問い合わせを行ってください。 #03a\n",
		)

		logger.WithFields(logrus.Fields{
			"event": "nsfw",
			"id":    emoji.ID,
			"user":  i.Member.User,
			"name":  emoji.Name,
		}).Error(err)
		return
	}

	if emoji.IsRequested {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Flags:   discordgo.MessageFlagsEphemeral,
				Content: "既に申請は終了しています\n",
			},
		})
		return
	}

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Flags:   discordgo.MessageFlagsEphemeral,
			Content: "非NSFWに設定されました\n",
		},
	})

	emoji.IsSensitive = false
	emoji.RequestState = "Nsfw"
	emoji.ResponseState = "Nsfw"
	ProcessNextRequest(emoji, s, i.ChannelID)
}
