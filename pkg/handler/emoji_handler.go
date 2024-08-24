package handler

import (
	"MisskeyEmojiBot/pkg/entity"
	"MisskeyEmojiBot/pkg/repository"
	"errors"
)

type EmojiHandler struct {
	discordRepo repository.DiscordRepository
}

func NewEmojiHandler(discordRepo repository.DiscordRepository) *EmojiHandler {
	return &EmojiHandler{discordRepo: discordRepo}
}

func (h *EmojiHandler) approve(emoji *entity.Emoji) error {
	if emoji.IsAccepted {
		return errors.New("Emoji is already accepted")
	}
	uploadToMisskey(emoji)
	emoji.IsFinish = true
	h.discordRepo.SendDirectMessage(*&emoji.RequestUser, "申請された絵文字は登録されました。"+"\n"+emoji.Name)
	h.discordRepo.DeleteChannel(*&emoji.ChannelID)
	return nil
}

func (h *EmojiHandler) disapprove(emoji *entity.Emoji) {
	if emoji.IsAccepted {
		return
	}

	emoji.IsAccepted = false
	emoji.IsFinish = true
	h.discordRepo.SendDirectMessage(*&emoji.RequestUser, "申請された絵文字は却下されました。"+"\n"+emoji.Name)
	h.discordRepo.DeleteChannel(*&emoji.ChannelID)
}
