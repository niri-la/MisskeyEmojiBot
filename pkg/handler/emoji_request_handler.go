package handler

import (
	"MisskeyEmojiBot/pkg/entity"
	"MisskeyEmojiBot/pkg/repository"

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
	AddProcess(processor EmojiProcessHandler)
	Process(emoji *entity.Emoji, s *discordgo.Session, m *discordgo.MessageCreate) error
	ProcessRequest(emoji *entity.Emoji, s *discordgo.Session, channelID string) error
	ResetState(emoji *entity.Emoji, s *discordgo.Session)
}

type emojiRequestHandler struct {
	reverseWorkflowMap map[string]int
	processor          []EmojiProcessHandler
	emojiRepository    repository.EmojiRepository
}

func NewEmojiRequestHandler(emojiRepo repository.EmojiRepository) EmojiRequestHandler {
	handler := &emojiRequestHandler{
		emojiRepository: emojiRepo,
	}
	handler.init()
	return handler
}

func (h *emojiRequestHandler) init() {
	h.reverseWorkflowMap = make(map[string]int)
	for key, value := range workflow {
		h.reverseWorkflowMap[value] = key
	}
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

	if emoji.NowStateIndex >= len(h.processor) {
		return nil
	}

	processor := h.processor[emoji.NowStateIndex]

	if !emoji.ResponseFlag {
		r, err := processor.Request(emoji, s, m.ChannelID)
		if err != nil {
			return err
		}
		if r.IsSuccess {
			emoji.ResponseFlag = true
			h.emojiRepository.Save(emoji)
		}
	} else {
		r, err := processor.Response(emoji, s, m)
		if err != nil {
			return err
		}
		if r.IsSuccess {
			emoji.NowStateIndex++
			emoji.ResponseFlag = false
			h.emojiRepository.Save(emoji)
			return h.ProcessRequest(emoji, s, m.ChannelID)
		}
	}

	return nil
}

func (h *emojiRequestHandler) ResetState(emoji *entity.Emoji, s *discordgo.Session) {
	emoji.NowStateIndex = 0
	emoji.ResponseFlag = false
}

func (h *emojiRequestHandler) ProcessRequest(emoji *entity.Emoji, s *discordgo.Session, channelID string) error {
	if emoji.NowStateIndex >= len(h.processor) {
		return nil
	}

	processor := h.processor[emoji.NowStateIndex]

	r, err := processor.Request(emoji, s, channelID)
	if err != nil {
		return err
	}
	if r.IsSuccess {
		emoji.ResponseFlag = true
		h.emojiRepository.Save(emoji)
	}

	return nil
}
