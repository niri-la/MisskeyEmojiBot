package request

import (
	"MisskeyEmojiBot/pkg/entity"
	"regexp"
	"strings"

	"github.com/bwmarrin/discordgo"
)

type NameSettingHandler struct {
}

func (h *NameSettingHandler) Request(emoji *entity.Emoji, s *discordgo.Session, cID string) (entity.Response, error) {

	response := entity.Response{
		IsSuccess: true,
	}

	_, err := s.ChannelMessageSend(
		cID,
		"## 絵文字の名前を入力してください。\n実際にMisskey上で絵文字を入力する際は`:emoji-name:`としますが、この`emoji-name`の部分を入力してください。\n入力可能な文字は`小文字アルファベット`, `数字`, `_`です。",
	)
	if err != nil {
		return entity.Response{}, err
	}

	emoji.RequestState = "SetName"

	return response, nil
}

func (h *NameSettingHandler) Response(emoji *entity.Emoji, s *discordgo.Session, m *discordgo.MessageCreate) (entity.Response, error) {
	response := entity.Response{
		IsSuccess: false,
	}

	if len(m.Content) <= 1 {
		_, err := s.ChannelMessageSend(m.ChannelID, ":2文字以上入力してください。")
		if err != nil {

			return response, err
		}
		return response, nil
	}
	reg := regexp.MustCompile(`[^a-z0-9_]+`)
	input := reg.ReplaceAllStringFunc(strings.ToLower(m.Content), func(s string) string {
		return "_"
	})
	s.ChannelMessageSend(m.ChannelID, ":: 入力されたメッセージ\n [ `"+input+"` ]")
	s.ChannelMessageSend(m.ChannelID, ":---")
	emoji.Name = input
	emoji.ResponseState = "SetName"
	response.IsSuccess = true
	response.NextState = response.NextState + 1
	return response, nil
}
