package processor

import (
	"MisskeyEmojiBot/pkg/entity"
	"MisskeyEmojiBot/pkg/handler"

	"github.com/bwmarrin/discordgo"
)

type nsfwHandler struct {
}

func NewNsfwHandler() handler.EmojiProcessHandler {
	return &nsfwHandler{}
}

func (h *nsfwHandler) Request(emoji *entity.Emoji, s *discordgo.Session, cID string) (entity.Response, error) {
	response := entity.Response{
		IsSuccess: true,
	}
	s.ChannelMessageSendComplex(cID,
		&discordgo.MessageSend{
			Content: "## 絵文字はセンシティブですか？\n",
			Components: []discordgo.MessageComponent{
				discordgo.ActionsRow{
					Components: []discordgo.MessageComponent{
						&discordgo.Button{
							Label:    "はい",
							CustomID: "nsfw_yes",
							Style:    discordgo.DangerButton,
							Emoji: discordgo.ComponentEmoji{
								Name: "🚢",
							},
						},
						&discordgo.Button{
							Label:    "いいえ",
							CustomID: "nsfw_no",
							Style:    discordgo.PrimaryButton,
							Emoji: discordgo.ComponentEmoji{
								Name: "🚀",
							},
						},
					},
				},
			},
		},
	)
	emoji.RequestState = "Nsfw"
	return response, nil
}

func (h *nsfwHandler) Response(emoji *entity.Emoji, s *discordgo.Session, m *discordgo.MessageCreate) (entity.Response, error) {
	response := entity.Response{
		IsSuccess: false,
	}
	return response, nil
}
