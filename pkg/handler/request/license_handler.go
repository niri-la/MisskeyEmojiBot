package request

import (
	"MisskeyEmojiBot/pkg/entity"

	"github.com/bwmarrin/discordgo"
)

type LicenseHandler struct {
}

func (h *LicenseHandler) Request(emoji *entity.Emoji, s *discordgo.Session, cID string) (entity.Response, error) {

	response := entity.Response{
		IsSuccess: true,
	}

	_, err := s.ChannelMessageSend(cID, "## ライセンス情報を入力してください。\nこれは絵文字ファイルやその素材に関する権利/所有者を明らかにするために重要なものです。\n入力する内容がない場合は`なし`と入力してください。")
	if err != nil {
		return entity.Response{}, err
	}

	emoji.RequestState = "License"

	return response, nil
}

func (h *LicenseHandler) Response(emoji *entity.Emoji, s *discordgo.Session, m *discordgo.MessageCreate) (entity.Response, error) {
	response := entity.Response{
		IsSuccess: false,
	}

	emoji.ResponseState = "License"
	input := m.Content
	if input == "なし" {
		input = ""
	}
	emoji.License = input

	s.ChannelMessageSend(m.ChannelID, ":: 入力されたメッセージ\n [ `"+input+"` ]")
	s.ChannelMessageSend(m.ChannelID, ":---")

	response.IsSuccess = true
	response.NextState = response.NextState + 1

	return response, nil
}
