package processor

import (
	"MisskeyEmojiBot/pkg/entity"
	"MisskeyEmojiBot/pkg/handler"

	"github.com/bwmarrin/discordgo"
)

type categoryHandler struct {
}

func NewCategoryHandler() handler.EmojiProcessHandler {
	return &categoryHandler{}
}

func (h *categoryHandler) Request(emoji *entity.Emoji, s *discordgo.Session, cID string) (entity.Response, error) {

	response := entity.Response{
		IsSuccess: true,
	}

	_, err := s.ChannelMessageSend(cID, "## 絵文字のカテゴリを入力してください。\n特にない場合は「なし」と入力してください。\nカテゴリ名については絵文字やリアクションを入力する際のメニューを参考にしてください。\n例: `Moji`")
	if err != nil {
		return entity.Response{}, err
	}

	emoji.RequestState = "Category"

	return response, nil
}

func (h *categoryHandler) Response(emoji *entity.Emoji, s *discordgo.Session, m *discordgo.MessageCreate) (entity.Response, error) {

	response := entity.Response{
		IsSuccess: false,
	}

	emoji.Category = m.Content
	if m.Content == "なし" || m.Content == "その他" {
		emoji.Category = ""
	}
	emoji.ResponseState = "Category"
	response.IsSuccess = true
	response.NextState = response.NextState + 1
	s.ChannelMessageSend(m.ChannelID, ":: 入力されたメッセージ\n [ `"+emoji.Category+"` ]")
	s.ChannelMessageSend(m.ChannelID, ":---\n")

	return response, nil
}
