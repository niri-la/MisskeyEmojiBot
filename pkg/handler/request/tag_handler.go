package request

import (
	"MisskeyEmojiBot/pkg/entity"

	"github.com/bwmarrin/discordgo"
)

type TagHandler struct {
}

func (h *TagHandler) Request(emoji *entity.Emoji, s *discordgo.Session, cID string) (entity.Response, error) {

	response := entity.Response{
		IsSuccess: true,
	}

	_, err := s.ChannelMessageSend(cID, "## 次に絵文字ファイルに設定するタグ(エイリアス)を入力してください。\n空白を間に挟むと複数設定できます。\n"+
		"これは絵文字の検索を行う際に使用されるため、漢字、ひらがな、カタカナ、ローマ字などのバリエーションがあると利用しやすくなります。\n"+
		"例: `絵文字 えもじ emoji emozi`\n必要がない場合は`tagなし`と入力してください。")

	if err != nil {
		return entity.Response{}, err
	}

	emoji.RequestState = "Tag"

	return response, nil
}

func (h *TagHandler) Response(emoji *entity.Emoji, s *discordgo.Session, m *discordgo.MessageCreate) (entity.Response, error) {

	response := entity.Response{
		IsSuccess: false,
	}

	emoji.Tag = m.Content
	if m.Content == "tagなし" {
		emoji.Tag = ""
	}
	emoji.ResponseState = "Tag"
	response.IsSuccess = true
	response.NextState = response.NextState + 1
	s.ChannelMessageSend(m.ChannelID, ":: 入力されたメッセージ\n [ `"+emoji.Tag+"` ]")
	s.ChannelMessageSend(m.ChannelID, ":---\n")

	return response, nil
}
