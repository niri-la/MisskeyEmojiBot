package processor

import (
	"github.com/bwmarrin/discordgo"

	"MisskeyEmojiBot/pkg/entity"
	"MisskeyEmojiBot/pkg/handler"
)

type licenseHandler struct {
}

func NewLicenseHandlerHandler() handler.EmojiProcessHandler {
	return &licenseHandler{}
}

func (h *licenseHandler) Request(emoji *entity.Emoji, s *discordgo.Session, cID string) (entity.Response, error) {

	response := entity.Response{
		IsSuccess: true,
	}

	_, err := s.ChannelMessageSend(cID, "## ライセンス情報を入力してください。\nこれは絵文字ファイルやその素材に関する権利/所有者を明らかにするために重要なものです。\n入力する内容がない場合は`なし`と入力してください。")
	if err != nil {
		return entity.Response{}, err
	}

	return response, nil
}

func (h *licenseHandler) Response(emoji *entity.Emoji, s *discordgo.Session, m *discordgo.MessageCreate) (entity.Response, error) {
	response := entity.Response{
		IsSuccess: false,
	}

	input := m.Content
	if input == "なし" {
		input = ""
	}
	emoji.License = input

	s.ChannelMessageSend(m.ChannelID, ":: 入力されたメッセージ\n [ `"+input+"` ]")
	s.ChannelMessageSend(m.ChannelID, "# ----------\n")

	response.IsSuccess = true

	return response, nil
}
