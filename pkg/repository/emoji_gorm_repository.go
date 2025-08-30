package repository

import (
	"time"

	"github.com/google/uuid"

	"MisskeyEmojiBot/pkg/config"
	"MisskeyEmojiBot/pkg/database"
	"MisskeyEmojiBot/pkg/entity"
	"MisskeyEmojiBot/pkg/errors"
)

type emojiGormRepository struct {
	db     *database.DB
	config config.Config
}

func NewEmojiGormRepository(db *database.DB, config config.Config) EmojiRepository {
	return &emojiGormRepository{
		db:     db,
		config: config,
	}
}

func (r *emojiGormRepository) NewEmoji(user string) *entity.Emoji {
	id := uuid.New()
	emoji := &entity.Emoji{
		ID:              id.String(),
		RequestUser:     user,
		StartAt:         time.Now(),
		NowStateIndex:   0,
		ResponseFlag:    false,
		IsRequested:     false,
		IsAccepted:      false,
		IsFinish:        false,
		ApproveCount:    0,
		DisapproveCount: 0,
		IsSensitive:     false,
	}

	// Create in database
	result := r.db.Create(emoji)
	if result.Error != nil {
		return nil
	}

	return emoji
}

func (r *emojiGormRepository) GetEmojis() []entity.Emoji {
	var emojis []entity.Emoji
	result := r.db.Where("is_finish = ?", false).Find(&emojis)
	if result.Error != nil {
		return []entity.Emoji{}
	}
	return emojis
}

func (r *emojiGormRepository) GetEmoji(id string) (*entity.Emoji, error) {
	var emoji entity.Emoji
	result := r.db.First(&emoji, "id = ?", id)
	if result.Error != nil {
		return nil, errors.EmojiNotFound(id)
	}
	return &emoji, nil
}

func (r *emojiGormRepository) EmojiReconstruction() []entity.Emoji {
	// Get all accepted emojis that are finished
	var acceptedEmojis []entity.Emoji
	r.db.Where("is_finish = ? AND is_accepted = ?", true, true).Find(&acceptedEmojis)

	// Delete finished emojis from database
	r.db.Where("is_finish = ?", true).Delete(&entity.Emoji{})

	return acceptedEmojis
}

func (r *emojiGormRepository) Approve(emoji *entity.Emoji) error {
	if emoji.IsAccepted {
		return errors.EmojiAlready("emoji is already accepted")
	}

	emoji.IsAccepted = true
	emoji.IsFinish = true

	result := r.db.Save(emoji)
	if result.Error != nil {
		return errors.FileOperation("failed to approve emoji", result.Error)
	}

	return nil
}

func (r *emojiGormRepository) Disapprove(emoji *entity.Emoji) error {
	if emoji.IsAccepted {
		return errors.EmojiAlready("emoji is already accepted")
	}

	emoji.IsAccepted = false
	emoji.IsFinish = true

	result := r.db.Save(emoji)
	if result.Error != nil {
		return errors.FileOperation("failed to disapprove emoji", result.Error)
	}

	return nil
}

func (r *emojiGormRepository) Abort(emoji *entity.Emoji) {
	r.Remove(*emoji)
	r.ResetState(emoji)
	emoji.IsFinish = true
	r.db.Save(emoji)
}

func (r *emojiGormRepository) Remove(emoji entity.Emoji) {
	r.db.Delete(&emoji)
}

func (r *emojiGormRepository) Save(emoji *entity.Emoji) error {
	result := r.db.Save(emoji)
	if result.Error != nil {
		return errors.FileOperation("failed to save emoji", result.Error)
	}
	return nil
}

func (r *emojiGormRepository) ResetState(emoji *entity.Emoji) error {
	emoji.IsSensitive = false
	emoji.IsAccepted = false
	emoji.IsRequested = false

	result := r.db.Save(emoji)
	if result.Error != nil {
		return errors.FileOperation("failed to reset emoji state", result.Error)
	}

	return nil
}
