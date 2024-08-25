package handler

import (
	"MisskeyEmojiBot/pkg/entity"
	"MisskeyEmojiBot/pkg/repository"
	"errors"
)

type EmojiHandler interface {
	GetEmoji(id string) (*entity.Emoji, error)
	approve(emoji *entity.Emoji) error
	disapprove(emoji *entity.Emoji)
}

type emojiHandler struct {
	discordRepo      repository.DiscordRepository
	misskeyRepo      repository.MisskeyRepository
	emojiProcessList []entity.Emoji
}

func NewEmojiHandler(discordRepo repository.DiscordRepository, misskeyRepo repository.MisskeyRepository) EmojiHandler {
	return &emojiHandler{discordRepo: discordRepo, misskeyRepo: misskeyRepo}
}

func (h *emojiHandler) GetEmoji(id string) (*entity.Emoji, error) {
	for i := range h.emojiProcessList {
		if h.emojiProcessList[i].ID == id {
			return &h.emojiProcessList[i], nil
		}
	}
	return nil, errors.New("emoji not found")
}

func (h *emojiHandler) approve(emoji *entity.Emoji) error {
	if emoji.IsAccepted {
		return errors.New("Emoji is already accepted")
	}
	h.misskeyRepo.UploadEmoji(emoji)
	emoji.IsFinish = true
	h.discordRepo.SendDirectMessage(*&emoji.RequestUser, "申請された絵文字は登録されました。"+"\n"+emoji.Name)
	h.discordRepo.DeleteChannel(*&emoji.ChannelID)
	return nil
}

func (h *emojiHandler) disapprove(emoji *entity.Emoji) {
	if emoji.IsAccepted {
		return
	}

	emoji.IsAccepted = false
	emoji.IsFinish = true
	h.discordRepo.SendDirectMessage(*&emoji.RequestUser, "申請された絵文字は却下されました。"+"\n"+emoji.Name)
	h.discordRepo.DeleteChannel(*&emoji.ChannelID)
}
