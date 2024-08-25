package component

import (
	"MisskeyEmojiBot/pkg/handler"
	"MisskeyEmojiBot/pkg/repository"

	"github.com/bwmarrin/discordgo"
)

type NsfwActiveComponen interface {
}

type nsfwActiveComponen struct {
	emojiRequestHandler handler.EmojiRequestHandler
	emojiRepository     repository.EmojiRepository
	discordRepo         repository.DiscordRepository
}

func NewNsfwActiveComponent(emojiRequestHandler handler.EmojiRequestHandler, emojiRepository repository.EmojiRepository, discordRepo repository.DiscordRepository) handler.Component {
	return &nsfwActiveComponen{emojiRequestHandler: emojiRequestHandler, emojiRepository: emojiRepository, discordRepo: discordRepo}
}

func (c *nsfwActiveComponen) GetCommand() *discordgo.ApplicationCommand {
	return &discordgo.ApplicationCommand{
		Name: "nsfw_yes",
	}
}

// Execute implements handler.Component.
func (c *nsfwActiveComponen) Execute(s *discordgo.Session, i *discordgo.InteractionCreate) {
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
			Content: "NSFWに設定されました\n",
		},
	})
	emoji.IsSensitive = true
	emoji.NowStateIndex++
	emoji.ResponseFlag = false
	c.emojiRequestHandler.ProcessRequest(emoji, s, i.ChannelID)
}
