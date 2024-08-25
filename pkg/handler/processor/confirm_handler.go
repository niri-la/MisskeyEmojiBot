package processor

import (
	"MisskeyEmojiBot/pkg/entity"
	"MisskeyEmojiBot/pkg/handler"
	"strconv"

	"github.com/bwmarrin/discordgo"
)

type confirmHandler struct {
}

func NewConfirmHandler() handler.EmojiProcessHandler {
	return &confirmHandler{}
}

func (h *confirmHandler) Request(emoji *entity.Emoji, s *discordgo.Session, cID string) (entity.Response, error) {
	response := entity.Response{
		IsSuccess: true,
	}

	s.ChannelMessageSend(cID, "# ----------\n")
	s.ChannelMessageSend(cID, "## 最終確認を行います。\n"+
		"- 名前 / Name: **"+emoji.Name+"**\n"+
		"- カテゴリ / Category: **"+emoji.Category+"**\n"+
		"- タグ / Tag: **"+emoji.Tag+"**\n"+
		"- ライセンス / License: **"+emoji.License+"**\n"+
		"- その他 / Other: **"+emoji.Other+"**\n"+
		"- NSFW: **"+strconv.FormatBool(emoji.IsSensitive)+"**\n",
	)
	_, err := s.ChannelMessageSendComplex(cID,
		&discordgo.MessageSend{
			Content: "## 以上で申請しますか?\n",
			Components: []discordgo.MessageComponent{
				discordgo.ActionsRow{
					Components: []discordgo.MessageComponent{
						&discordgo.Button{
							Label:    "はい",
							CustomID: "emoji_request",
							Style:    discordgo.PrimaryButton,
							Emoji: discordgo.ComponentEmoji{
								Name: "📨",
							},
						},
						&discordgo.Button{
							Label:    "最初からやり直す",
							CustomID: "emoji_request_retry",
							Style:    discordgo.DangerButton,
							Emoji: discordgo.ComponentEmoji{
								Name: "🗑️",
							},
						},
					},
				},
			},
		},
	)

	if err != nil {
		return entity.Response{}, err
	}

	return response, nil
}

func (h *confirmHandler) Response(emoji *entity.Emoji, s *discordgo.Session, m *discordgo.MessageCreate) (entity.Response, error) {
	response := entity.Response{
		IsSuccess: false,
	}

	return response, nil
}
