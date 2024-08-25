package component

import (
	"MisskeyEmojiBot/pkg/handler"
	"MisskeyEmojiBot/pkg/repository"

	"github.com/bwmarrin/discordgo"
	"github.com/sirupsen/logrus"
)

type NsfwActiveComponen interface {
}

type nsfwActiveComponen struct {
	discordRepo repository.DiscordRepository
}

func NewNsfwActiveComponent(discordRepo repository.DiscordRepository) handler.Component {
	return &nsfwActiveComponen{discordRepo: discordRepo}
}

func (n *nsfwActiveComponen) GetCommand() *discordgo.ApplicationCommand {
	return &discordgo.ApplicationCommand{
		Name: "nsfw_yes",
	}
}

// Execute implements handler.Component.
func (n *nsfwActiveComponen) Execute(s *discordgo.Session, i *discordgo.InteractionCreate) {
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
			Content: "NSFWに設定されました\n",
		},
	})
	emoji.IsSensitive = true
	emoji.RequestState = "Nsfw"
	emoji.ResponseState = "Nsfw"
	ProcessNextRequest(emoji, s, i.ChannelID)
}
