package repository

import (
	"MisskeyEmojiBot/pkg/config"
	"MisskeyEmojiBot/pkg/entity"
	"MisskeyEmojiBot/pkg/errors"
	"encoding/json"
	"os"
	"time"

	"github.com/google/uuid"
)

type EmojiRepository interface {
	NewEmoji(user string) *entity.Emoji
	GetEmojis() []entity.Emoji
	GetEmoji(id string) (*entity.Emoji, error)
	EmojiReconstruction() []entity.Emoji
	Approve(emoji *entity.Emoji) error
	Disapprove(emoji *entity.Emoji) error
	Abort(emoji *entity.Emoji)
	Remove(emoji entity.Emoji)
	Save(emoji *entity.Emoji) error
	ResetState(emoji *entity.Emoji) error
}

type emojiRepository struct {
	config           config.Config
	emojiProcessList []entity.Emoji
}

func NewEmojiRepository(config config.Config) EmojiRepository {
	return &emojiRepository{config: config}
}

func (r *emojiRepository) NewEmoji(user string) *entity.Emoji {
	id := uuid.New()
	emoji := entity.Emoji{
		ID: id.String(),
	}
	emoji.RequestUser = user
	emoji.StartAt = time.Now()
	emoji.NowStateIndex = 0
	r.addEmoji(emoji)
	return &emoji
}

func (h *emojiRepository) GetEmojis() []entity.Emoji {
	return h.emojiProcessList
}

func (h *emojiRepository) GetEmoji(id string) (*entity.Emoji, error) {
	for i := range h.emojiProcessList {
		if h.emojiProcessList[i].ID == id {
			return &h.emojiProcessList[i], nil
		}
	}
	return nil, errors.EmojiNotFound(id)
}

func (h *emojiRepository) EmojiReconstruction() []entity.Emoji {
	var accepted []entity.Emoji
	var reconstruction []entity.Emoji
	for _, emoji := range h.emojiProcessList {
		if emoji.IsFinish {
			if emoji.IsAccepted {
				accepted = append(accepted, emoji)
			}
		} else {
			reconstruction = append(reconstruction, emoji)
		}
	}
	h.emojiProcessList = reconstruction
	return accepted
}

func (h *emojiRepository) Approve(emoji *entity.Emoji) error {
	if emoji.IsAccepted {
		return errors.EmojiAlready("emoji is already accepted")
	}
	emoji.IsAccepted = true
	emoji.IsFinish = true
	return nil
}

func (h *emojiRepository) Disapprove(emoji *entity.Emoji) error {
	if emoji.IsAccepted {
		return errors.EmojiAlready("emoji is already accepted")
	}

	emoji.IsAccepted = false
	emoji.IsFinish = true
	return nil
}

func (h *emojiRepository) Abort(emoji *entity.Emoji) {
	h.Remove(*emoji)
	h.ResetState(emoji)
	emoji.IsFinish = true
}

func (h *emojiRepository) Remove(emoji entity.Emoji) {
	var newSlice []entity.Emoji
	for _, v := range h.emojiProcessList {
		if v.ID != emoji.ID {
			newSlice = append(newSlice, v)
		}
	}
	h.emojiProcessList = newSlice
}

func (h *emojiRepository) addEmoji(emoji entity.Emoji) {
	h.emojiProcessList = append(h.emojiProcessList, emoji)
}

func (h *emojiRepository) Save(emoji *entity.Emoji) error {
	jsonData, err := json.MarshalIndent(emoji, "", "  ")
	if err != nil {
		return errors.FileOperation("failed to marshal emoji data to JSON", err)
	}
	
	filePath := h.config.SavePath + emoji.ID + ".json"
	if err := os.WriteFile(filePath, jsonData, 0644); err != nil {
		return errors.FileOperation("failed to save emoji data", err)
	}
	return nil
}

func (h *emojiRepository) ResetState(emoji *entity.Emoji) error {
	emoji.IsSensitive = false
	emoji.IsAccepted = false
	emoji.IsRequested = false
	return nil
}
