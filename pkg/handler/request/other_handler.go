package request

import (
	"MisskeyEmojiBot/pkg/entity"

	"github.com/bwmarrin/discordgo"
)

type OtherHandler struct {
}

func (h *OtherHandler) Request(emoji *entity.Emoji, s *discordgo.Session, cID string) (entity.Response, error) {

	response := entity.Response{
		IsSuccess: true,
	}

	_, err := s.ChannelMessageSend(cID, "## 備考があれば記載してください。\nこの内容はMisskey上には掲載されません。\n特にない場合は`なし`と入力してください。")
	if err != nil {
		return entity.Response{}, err
	}

	emoji.RequestState = "Other"

	return response, nil
}

func (h *OtherHandler) Response(emoji *entity.Emoji, s *discordgo.Session, m *discordgo.MessageCreate) (entity.Response, error) {
	response := entity.Response{
		IsSuccess: false,
	}

	emoji.ResponseState = "Other"
	input := m.Content
	if input == "なし" {
		input = ""
	}
	emoji.Other = input

	s.ChannelMessageSend(m.ChannelID, ":: 入力されたメッセージ\n [ `"+input+"` ]")
	s.ChannelMessageSend(m.ChannelID, ":---")

	response.IsSuccess = true
	response.NextState = response.NextState + 1

	return response, nil
}