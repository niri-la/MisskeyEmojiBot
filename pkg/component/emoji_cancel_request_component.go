package component

import (
	"MisskeyEmojiBot/pkg/handler"
	"MisskeyEmojiBot/pkg/repository"

	"github.com/bwmarrin/discordgo"
)

type EmojiCancelRequestComponent interface {
}

type emojiCancelRequestComponent struct {
	discordRepo repository.DiscordRepository
}

func NewEmojiCancelRequestComponent(discordRepo repository.DiscordRepository) handler.Component {
	return &emojiCancelRequestComponent{discordRepo: discordRepo}
}

func (e *emojiCancelRequestComponent) GetCommand() *discordgo.ApplicationCommand {
	return &discordgo.ApplicationCommand{
		Name: "cancel_request",
	}
}

func (e *emojiCancelRequestComponent) Execute(s *discordgo.Session, i *discordgo.InteractionCreate) {
	channel, _ := s.Channel(i.ChannelID)
	emoji, err := GetEmoji(channel.Name[6:])
	if err != nil {
		s.ChannelMessageSend(
			channel.ID,
			"設定に失敗しました。管理者に問い合わせを行ってください。 #03a\n",
		)
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
			Content: "リクエストをキャンセルしました。\n",
		},
	})
	e.discordRepo.SendDirectMessage(*emoji, "申請された絵文字はキャンセルされました。: ")
	emoji.abort()
	e.discordRepo.DeleteChannel(*emoji)
}
