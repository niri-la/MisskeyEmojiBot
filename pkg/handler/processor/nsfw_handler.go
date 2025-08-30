package processor

import (
	"github.com/bwmarrin/discordgo"

	"MisskeyEmojiBot/pkg/entity"
	"MisskeyEmojiBot/pkg/handler"
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
	return response, nil
}

// dummyなので何もしない(IsSuccess: falseに設定しないとcomponentとの連携で不整合が発生する)
func (h *nsfwHandler) Response(emoji *entity.Emoji, s *discordgo.Session, m *discordgo.MessageCreate) (entity.Response, error) {
	response := entity.Response{
		IsSuccess: false,
	}
	return response, nil
}
