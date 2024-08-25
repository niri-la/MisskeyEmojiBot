package emoji

import (
	"MisskeyEmojiBot/pkg/entity"
	"MisskeyEmojiBot/pkg/repository"
)

type EmojiHandler interface {
	GetEmoji(id string) (*entity.Emoji, error)
	Approve(emoji *entity.Emoji) error
	Disapprove(emoji *entity.Emoji) error
	EmojiReconstruction() []entity.Emoji
}

type emojiHandler struct {
	emojiRepository repository.EmojiRepository
	discordRepo     repository.DiscordRepository
	misskeyRepo     repository.MisskeyRepository
}

func NewEmojiHandler(emojiRepository repository.EmojiRepository, discordRepo repository.DiscordRepository, misskeyRepo repository.MisskeyRepository) EmojiHandler {
	return &emojiHandler{emojiRepository: emojiRepository, discordRepo: discordRepo, misskeyRepo: misskeyRepo}
}

func (h *emojiHandler) GetEmoji(id string) (*entity.Emoji, error) {
	return h.emojiRepository.GetEmoji(id)
}

func (h *emojiHandler) Approve(emoji *entity.Emoji) error {
	err := h.emojiRepository.Save(emoji)
	if err != nil {
		return err
	}
	err = h.emojiRepository.Approve(emoji)
	if err != nil {
		return err
	}

	err = h.misskeyRepo.UploadEmoji(emoji)
	if err != nil {
		return err
	}
	h.discordRepo.SendDirectMessage(emoji.RequestUser, "申請された絵文字は登録されました。"+"\n"+emoji.Name)
	h.discordRepo.DeleteChannel(emoji.ChannelID)
	return nil
}

func (h *emojiHandler) Disapprove(emoji *entity.Emoji) error {
	err := h.emojiRepository.Disapprove(emoji)
	if err != nil {
		return err
	}

	h.discordRepo.SendDirectMessage(emoji.RequestUser, "申請された絵文字は却下されました。"+"\n"+emoji.Name)
	h.discordRepo.DeleteChannel(emoji.ChannelID)
	return nil
}

func (h *emojiHandler) EmojiReconstruction() []entity.Emoji {
	return h.emojiRepository.EmojiReconstruction()
}

func (h *emojiHandler) Abort(emoji entity.Emoji) {
	h.emojiRepository.Abort(&emoji)
}

func (h *emojiHandler) Remove(emoji entity.Emoji) {
	h.emojiRepository.Remove(emoji)
}
