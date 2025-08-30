package processor

import (
	"strings"

	"github.com/bwmarrin/discordgo"

	"MisskeyEmojiBot/pkg/entity"
	"MisskeyEmojiBot/pkg/handler"
)

type tagHandler struct {
}

func NewTagHandler() handler.EmojiProcessHandler {
	return &tagHandler{}
}

func (h *tagHandler) Request(emoji *entity.Emoji, s *discordgo.Session, cID string) (entity.Response, error) {

	response := entity.Response{
		IsSuccess: true,
	}

	_, err := s.ChannelMessageSend(cID, "## 次に絵文字ファイルに設定するタグ(エイリアス)を入力してください。\n空白を間に挟むと複数設定できます。\n"+
		"これは絵文字の検索を行う際に使用されるため、漢字、ひらがな、カタカナ、ローマ字などのバリエーションがあると利用しやすくなります。\n"+
		"例: `絵文字 えもじ emoji emozi`\n必要がない場合は`tagなし`と入力してください。")

	if err != nil {
		return entity.Response{}, err
	}

	return response, nil
}

func (h *tagHandler) Response(emoji *entity.Emoji, s *discordgo.Session, m *discordgo.MessageCreate) (entity.Response, error) {

	response := entity.Response{
		IsSuccess: false,
	}

	emoji.Tag = h.formatTag(m.Content)
	if m.Content == "tagなし" {
		emoji.Tag = ""
	}
	response.IsSuccess = true
	_, _ = s.ChannelMessageSend(m.ChannelID, ":: 入力されたメッセージ\n [ `"+emoji.Tag+"` ]")
	_, _ = s.ChannelMessageSend(m.ChannelID, "# ----------\n")

	return response, nil
}

func (h *tagHandler) formatTag(tag string) string {
	return strings.ReplaceAll(tag, "　", " ")
}
