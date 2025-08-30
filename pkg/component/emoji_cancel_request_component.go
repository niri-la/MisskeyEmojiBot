package component

import (
	"github.com/bwmarrin/discordgo"

	"MisskeyEmojiBot/pkg/handler"
	"MisskeyEmojiBot/pkg/repository"
)

type EmojiCancelRequestComponent interface {
}

type emojiCancelRequestComponent struct {
	emojiRepository repository.EmojiRepository
	discordRepo     repository.DiscordRepository
}

func NewEmojiCancelRequestComponent(emojiRepository repository.EmojiRepository, discordRepo repository.DiscordRepository) handler.Component {
	return &emojiCancelRequestComponent{emojiRepository: emojiRepository, discordRepo: discordRepo}
}

func (c *emojiCancelRequestComponent) GetCommand() *discordgo.ApplicationCommand {
	return &discordgo.ApplicationCommand{
		Name: "cancel_request",
	}
}

func (c *emojiCancelRequestComponent) Execute(s *discordgo.Session, i *discordgo.InteractionCreate) {
	channel, _ := s.Channel(i.ChannelID)
	emoji, err := c.emojiRepository.GetEmoji(channel.Name[6:])
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
	c.emojiRepository.Abort(emoji)
	c.discordRepo.SendDirectMessage(*&emoji.RequestUser, "申請された絵文字はキャンセルされました。: ")
	c.discordRepo.DeleteChannel(*&emoji.ChannelID)
}
