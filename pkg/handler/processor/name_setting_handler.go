package processor

import (
	"MisskeyEmojiBot/pkg/entity"
	"MisskeyEmojiBot/pkg/handler"
	"regexp"
	"strings"

	"github.com/bwmarrin/discordgo"
)

type nameSettingHandler struct {
}

func NewNameSettingHandler() handler.EmojiProcessHandler {
	return &nameSettingHandler{}
}

func (h *nameSettingHandler) Request(emoji *entity.Emoji, s *discordgo.Session, cID string) (entity.Response, error) {

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

	return response, nil
}

func (h *nameSettingHandler) Response(emoji *entity.Emoji, s *discordgo.Session, m *discordgo.MessageCreate) (entity.Response, error) {
	response := entity.Response{
		IsSuccess: false,
	}

	input := strings.TrimPrefix(m.Content, ":")
	input = strings.TrimSuffix(input, ":")

	if len(input) <= 1 {
		_, err := s.ChannelMessageSend(m.ChannelID, ":2文字以上入力してください。")
		if err != nil {

			return response, err
		}
		return response, nil
	}
	reg := regexp.MustCompile(`[^a-z0-9_]+`)
	input = reg.ReplaceAllStringFunc(strings.ToLower(input), func(s string) string {
		return "_"
	})
	s.ChannelMessageSend(m.ChannelID, ":: 入力されたメッセージ\n [ `"+input+"` ]")
	s.ChannelMessageSend(m.ChannelID, "# ----------\n")
	emoji.Name = input
	response.IsSuccess = true
	return response, nil
}
