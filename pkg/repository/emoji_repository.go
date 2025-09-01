package repository

import (
	"encoding/json"
	"os"
	"time"

	"github.com/google/uuid"

	"MisskeyEmojiBot/pkg/config"
	"MisskeyEmojiBot/pkg/entity"
	"MisskeyEmojiBot/pkg/errors"
)

type EmojiRepository interface {
	NewEmoji(user string) *entity.Emoji
	GetEmojis() []entity.Emoji
	GetEmojisForList(limit int) []entity.Emoji
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

func (r *emojiRepository) GetEmojis() []entity.Emoji {
	return r.emojiProcessList
}

func (r *emojiRepository) GetEmojisForList(limit int) []entity.Emoji {
	emojis := r.emojiProcessList
	if limit > 0 && limit < len(emojis) {
		return emojis[:limit]
	}
	return emojis
}

func (r *emojiRepository) GetEmoji(id string) (*entity.Emoji, error) {
	for i := range r.emojiProcessList {
		if r.emojiProcessList[i].ID == id {
			return &r.emojiProcessList[i], nil
		}
	}
	return nil, errors.EmojiNotFound(id)
}

func (r *emojiRepository) EmojiReconstruction() []entity.Emoji {
	var accepted []entity.Emoji
	var reconstruction []entity.Emoji
	for _, emoji := range r.emojiProcessList {
		if emoji.IsFinish {
			if emoji.IsAccepted {
				accepted = append(accepted, emoji)
			}
		} else {
			reconstruction = append(reconstruction, emoji)
		}
	}
	r.emojiProcessList = reconstruction
	return accepted
}

func (r *emojiRepository) Approve(emoji *entity.Emoji) error {
	if emoji.IsAccepted {
		return errors.EmojiAlready("emoji is already accepted")
	}
	emoji.IsAccepted = true
	emoji.IsFinish = true
	return nil
}

func (r *emojiRepository) Disapprove(emoji *entity.Emoji) error {
	if emoji.IsAccepted {
		return errors.EmojiAlready("emoji is already accepted")
	}

	emoji.IsAccepted = false
	emoji.IsFinish = true
	return nil
}

func (r *emojiRepository) Abort(emoji *entity.Emoji) {
	r.Remove(*emoji)
	_ = r.ResetState(emoji)
	emoji.IsFinish = true
}

func (r *emojiRepository) Remove(emoji entity.Emoji) {
	var newSlice []entity.Emoji
	for _, v := range r.emojiProcessList {
		if v.ID != emoji.ID {
			newSlice = append(newSlice, v)
		}
	}
	r.emojiProcessList = newSlice
}

func (r *emojiRepository) addEmoji(emoji entity.Emoji) {
	r.emojiProcessList = append(r.emojiProcessList, emoji)
}

func (r *emojiRepository) Save(emoji *entity.Emoji) error {
	jsonData, err := json.MarshalIndent(emoji, "", "  ")
	if err != nil {
		return errors.FileOperation("failed to marshal emoji data to JSON", err)
	}

	filePath := r.config.SavePath + emoji.ID + ".json"
	if err := os.WriteFile(filePath, jsonData, 0644); err != nil {
		return errors.FileOperation("failed to save emoji data", err)
	}
	return nil
}

func (r *emojiRepository) ResetState(emoji *entity.Emoji) error {
	emoji.IsSensitive = false
	emoji.IsAccepted = false
	emoji.IsRequested = false
	return nil
}
