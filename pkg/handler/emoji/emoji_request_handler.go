package emoji

import (
	"MisskeyEmojiBot/pkg/entity"

	"github.com/bwmarrin/discordgo"
)

type EmojiProcessHandler interface {
	Request(*entity.Emoji, *discordgo.Session, string) (entity.Response, error)
	Response(*entity.Emoji, *discordgo.Session, *discordgo.MessageCreate) (entity.Response, error)
}

// type RequestProcessor func(*entity.Emoji, *discordgo.Session, string) Response
// type ResponceProcessor func(*entity.Emoji, *discordgo.Session, *discordgo.MessageCreate) Response

var workflow = map[int]string{
	0: "Default",
	2: "SetName",
	1: "Upload",
	3: "Category",
	4: "Tag",
	5: "License",
	6: "Other",
	7: "Nsfw",
	8: "Check",
}

type EmojiRequestHandler interface {
}

type emojiRequestHandler struct {
	processor []EmojiProcessHandler
}

func NewEmojiRequestHandler() EmojiRequestHandler {
	handler := &emojiRequestHandler{}
	return handler
}

func (h *emojiRequestHandler) AddProcess(processor EmojiProcessHandler) {
	h.processor = append(h.processor, processor)
}

func (h *emojiRequestHandler) Process(emoji *entity.Emoji, s *discordgo.Session, m *discordgo.MessageCreate) error {
	// 0. まずNowStateIndexを確認し取得する。この時indexがprocess listより大きい場合は終了する
	// 1. responseFlagがfalseの場合はRequestを実行する。
	// 2. responseFlagがtrueの場合はResponseを実行する。
	// 3. 成功したらフラグをfalseにしてindexを1進める
	// 4. 1に戻る

	processor := h.processor[emoji.NowStateIndex]

	if emoji.NowStateIndex >= len(h.processor)-1 {
		return nil
	}

	if !emoji.ResponseFlag {
		r, err := processor.Request(emoji, s, m.ChannelID)
		if err != nil {
			return err
		}
		if r.IsSuccess {
			emoji.ResponseFlag = true
		}
	} else {
		r, err := processor.Response(emoji, s, m)
		if err != nil {
			return err
		}
		if r.IsSuccess {
			emoji.NowStateIndex++
			emoji.ResponseFlag = false
		}
	}

	return nil
}

func (h *emojiRequestHandler) ResetState(emoji *entity.Emoji, s *discordgo.Session, id string) {
	emoji.NowStateIndex = 0
	emoji.ResponseFlag = false
}
