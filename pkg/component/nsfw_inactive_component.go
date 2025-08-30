package component

import (
	"github.com/bwmarrin/discordgo"

	"MisskeyEmojiBot/pkg/handler"
	"MisskeyEmojiBot/pkg/repository"
)

type NsfwInactiveComponen interface {
}

type nsfwInactiveComponen struct {
	emojiRequestHandler handler.EmojiRequestHandler
	emojiRepository     repository.EmojiRepository
	discordRepo         repository.DiscordRepository
}

func NewNsfwInactiveComponent(emojiRequestHandler handler.EmojiRequestHandler, emojiRepository repository.EmojiRepository, discordRepo repository.DiscordRepository) handler.Component {
	return &nsfwInactiveComponen{emojiRequestHandler: emojiRequestHandler, emojiRepository: emojiRepository, discordRepo: discordRepo}
}

func (c *nsfwInactiveComponen) GetCommand() *discordgo.ApplicationCommand {
	return &discordgo.ApplicationCommand{
		Name: "nsfw_no",
	}
}

func (c *nsfwInactiveComponen) Execute(s *discordgo.Session, i *discordgo.InteractionCreate) {
	channel, _ := s.Channel(i.ChannelID)
	emoji, err := c.emojiRepository.GetEmoji(channel.Name[6:])
	if err != nil {
		_, _ = s.ChannelMessageSend(
			channel.ID,
			"設定に失敗しました。管理者に問い合わせを行ってください。 #03a\n",
		)
	}

	if emoji.IsRequested {
		_ = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Flags:   discordgo.MessageFlagsEphemeral,
				Content: "既に申請は終了しています\n",
			},
		})
		return
	}

	_ = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Flags:   discordgo.MessageFlagsEphemeral,
			Content: "非NSFWに設定されました\n",
		},
	})

	emoji.IsSensitive = false
	emoji.NowStateIndex++
	emoji.ResponseFlag = false
	_ = c.emojiRequestHandler.ProcessRequest(emoji, s, i.ChannelID)
}
