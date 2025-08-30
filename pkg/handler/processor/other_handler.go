package processor

import (
	"github.com/bwmarrin/discordgo"

	"MisskeyEmojiBot/pkg/entity"
	"MisskeyEmojiBot/pkg/handler"
)

type otherHandler struct {
}

func NewOtherHandler() handler.EmojiProcessHandler {
	return &otherHandler{}
}

func (h *otherHandler) Request(emoji *entity.Emoji, s *discordgo.Session, cID string) (entity.Response, error) {

	response := entity.Response{
		IsSuccess: true,
	}

	_, err := s.ChannelMessageSend(cID, "## 備考があれば記載してください。\nこの内容はMisskey上には掲載されません。\n特にない場合は`なし`と入力してください。")
	if err != nil {
		return entity.Response{}, err
	}

	return response, nil
}

func (h *otherHandler) Response(emoji *entity.Emoji, s *discordgo.Session, m *discordgo.MessageCreate) (entity.Response, error) {
	response := entity.Response{
		IsSuccess: false,
	}

	input := m.Content
	if input == "なし" {
		input = ""
	}
	emoji.Other = input

	_, _ = s.ChannelMessageSend(m.ChannelID, ":: 入力されたメッセージ\n [ `"+input+"` ]")
	_, _ = s.ChannelMessageSend(m.ChannelID, "# ----------\n")

	response.IsSuccess = true

	return response, nil
}
