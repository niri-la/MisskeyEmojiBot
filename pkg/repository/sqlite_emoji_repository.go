package repository

import (
	"MisskeyEmojiBot/pkg/config"
	"MisskeyEmojiBot/pkg/entity"
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	_ "modernc.org/sqlite"
)

type SQLiteEmojiRepository struct {
	config config.Config
	db     *sql.DB
}

func NewSQLiteEmojiRepository(config config.Config) (EmojiRepository, error) {
	db, err := sql.Open("sqlite", "emoji_requests.db")
	if err != nil {
		return nil, errors.Wrap(err, "failed to open database")
	}

	repo := &SQLiteEmojiRepository{
		config: config,
		db:     db,
	}

	if err := repo.initDB(); err != nil {
		return nil, errors.Wrap(err, "failed to initialize database")
	}

	return repo, nil
}

func (r *SQLiteEmojiRepository) initDB() error {
	query := `
	CREATE TABLE IF NOT EXISTS emoji_requests (
		id TEXT PRIMARY KEY,
		channel_id TEXT NOT NULL,
		request_user TEXT NOT NULL,
		name TEXT,
		category TEXT,
		tag TEXT,
		license TEXT,
		other TEXT,
		file_path TEXT,
		is_sensitive BOOLEAN DEFAULT 0,
		is_requested BOOLEAN DEFAULT 0,
		is_accepted BOOLEAN DEFAULT 0,
		is_finish BOOLEAN DEFAULT 0,
		approve_count INTEGER DEFAULT 0,
		disapprove_count INTEGER DEFAULT 0,
		response_flag BOOLEAN DEFAULT 0,
		now_state_index INTEGER DEFAULT 0,
		moderation_message_id TEXT,
		user_thread_id TEXT,
		start_at DATETIME NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
	)`

	_, err := r.db.Exec(query)
	return err
}


func (r *SQLiteEmojiRepository) NewEmoji(user string) *entity.Emoji {
	id := uuid.New()
	emoji := entity.Emoji{
		ID:            id.String(),
		RequestUser:   user,
		StartAt:       time.Now(),
		NowStateIndex: 0,
	}

	// Persist to database
	query := `
	INSERT INTO emoji_requests (
		id, channel_id, request_user, name, category, tag, license, other,
		file_path, is_sensitive, is_requested, is_accepted, is_finish,
		approve_count, disapprove_count, response_flag, now_state_index,
		moderation_message_id, user_thread_id, start_at
	) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	_, _ = r.db.Exec(query,
		emoji.ID, emoji.ChannelID, emoji.RequestUser, emoji.Name, emoji.Category,
		emoji.Tag, emoji.License, emoji.Other, emoji.FilePath, emoji.IsSensitive,
		emoji.IsRequested, emoji.IsAccepted, emoji.IsFinish, emoji.ApproveCount,
		emoji.DisapproveCount, emoji.ResponseFlag, emoji.NowStateIndex,
		emoji.ModerationMessageID, emoji.UserThreadID, emoji.StartAt,
	)

	return &emoji
}

func (r *SQLiteEmojiRepository) GetEmojis() []entity.Emoji {
	query := `
	SELECT id, channel_id, request_user, name, category, tag, license, other,
		   file_path, is_sensitive, is_requested, is_accepted, is_finish,
		   approve_count, disapprove_count, response_flag, now_state_index,
		   moderation_message_id, user_thread_id, start_at
	FROM emoji_requests
	WHERE is_finish = 0
	ORDER BY start_at`

	rows, err := r.db.Query(query)
	if err != nil {
		return []entity.Emoji{}
	}
	defer rows.Close()

	var emojis []entity.Emoji
	for rows.Next() {
		var emoji entity.Emoji
		err := rows.Scan(
			&emoji.ID, &emoji.ChannelID, &emoji.RequestUser, &emoji.Name,
			&emoji.Category, &emoji.Tag, &emoji.License, &emoji.Other,
			&emoji.FilePath, &emoji.IsSensitive, &emoji.IsRequested,
			&emoji.IsAccepted, &emoji.IsFinish, &emoji.ApproveCount,
			&emoji.DisapproveCount, &emoji.ResponseFlag, &emoji.NowStateIndex,
			&emoji.ModerationMessageID, &emoji.UserThreadID, &emoji.StartAt,
		)
		if err != nil {
			continue
		}
		emojis = append(emojis, emoji)
	}

	return emojis
}

func (r *SQLiteEmojiRepository) GetEmoji(id string) (*entity.Emoji, error) {
	query := `
	SELECT id, channel_id, request_user, name, category, tag, license, other,
		   file_path, is_sensitive, is_requested, is_accepted, is_finish,
		   approve_count, disapprove_count, response_flag, now_state_index,
		   moderation_message_id, user_thread_id, start_at
	FROM emoji_requests
	WHERE id = ?`

	row := r.db.QueryRow(query, id)

	var emoji entity.Emoji
	err := row.Scan(
		&emoji.ID, &emoji.ChannelID, &emoji.RequestUser, &emoji.Name,
		&emoji.Category, &emoji.Tag, &emoji.License, &emoji.Other,
		&emoji.FilePath, &emoji.IsSensitive, &emoji.IsRequested,
		&emoji.IsAccepted, &emoji.IsFinish, &emoji.ApproveCount,
		&emoji.DisapproveCount, &emoji.ResponseFlag, &emoji.NowStateIndex,
		&emoji.ModerationMessageID, &emoji.UserThreadID, &emoji.StartAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("emoji not found")
		}
		return nil, errors.Wrap(err, "failed to get emoji")
	}

	return &emoji, nil
}

func (r *SQLiteEmojiRepository) EmojiReconstruction() []entity.Emoji {
	// Get accepted emojis from database for return value
	acceptedQuery := `
	SELECT id, channel_id, request_user, name, category, tag, license, other,
		   file_path, is_sensitive, is_requested, is_accepted, is_finish,
		   approve_count, disapprove_count, response_flag, now_state_index,
		   moderation_message_id, user_thread_id, start_at
	FROM emoji_requests
	WHERE is_finish = 1 AND is_accepted = 1`

	rows, err := r.db.Query(acceptedQuery)
	if err != nil {
		return []entity.Emoji{}
	}
	defer rows.Close()

	var accepted []entity.Emoji
	for rows.Next() {
		var emoji entity.Emoji
		err := rows.Scan(
			&emoji.ID, &emoji.ChannelID, &emoji.RequestUser, &emoji.Name,
			&emoji.Category, &emoji.Tag, &emoji.License, &emoji.Other,
			&emoji.FilePath, &emoji.IsSensitive, &emoji.IsRequested,
			&emoji.IsAccepted, &emoji.IsFinish, &emoji.ApproveCount,
			&emoji.DisapproveCount, &emoji.ResponseFlag, &emoji.NowStateIndex,
			&emoji.ModerationMessageID, &emoji.UserThreadID, &emoji.StartAt,
		)
		if err != nil {
			continue
		}
		accepted = append(accepted, emoji)
	}

	// Keep finished emojis in database as history (no deletion)
	
	return accepted
}

func (r *SQLiteEmojiRepository) Approve(emoji *entity.Emoji) error {
	if emoji.IsAccepted {
		return errors.New("Emoji is already accepted")
	}

	emoji.IsAccepted = true
	emoji.IsFinish = true

	query := `UPDATE emoji_requests SET is_accepted = 1, is_finish = 1, updated_at = CURRENT_TIMESTAMP WHERE id = ?`
	_, err := r.db.Exec(query, emoji.ID)
	return err
}

func (r *SQLiteEmojiRepository) Disapprove(emoji *entity.Emoji) error {
	if emoji.IsAccepted {
		return errors.New("Emoji is already accepted")
	}

	emoji.IsAccepted = false
	emoji.IsFinish = true

	query := `UPDATE emoji_requests SET is_accepted = 0, is_finish = 1, updated_at = CURRENT_TIMESTAMP WHERE id = ?`
	_, err := r.db.Exec(query, emoji.ID)
	return err
}

func (r *SQLiteEmojiRepository) Abort(emoji *entity.Emoji) {
	r.Remove(*emoji)
	r.ResetState(emoji)
	emoji.IsFinish = true
}

func (r *SQLiteEmojiRepository) Remove(emoji entity.Emoji) {
	// Remove from database
	query := `DELETE FROM emoji_requests WHERE id = ?`
	r.db.Exec(query, emoji.ID)
}

func (r *SQLiteEmojiRepository) Save(emoji *entity.Emoji) error {
	query := `
	UPDATE emoji_requests SET
		channel_id = ?, name = ?, category = ?, tag = ?, license = ?, other = ?,
		file_path = ?, is_sensitive = ?, is_requested = ?, is_accepted = ?,
		is_finish = ?, approve_count = ?, disapprove_count = ?, response_flag = ?,
		now_state_index = ?, moderation_message_id = ?, user_thread_id = ?,
		updated_at = CURRENT_TIMESTAMP
	WHERE id = ?`

	_, err := r.db.Exec(query,
		emoji.ChannelID, emoji.Name, emoji.Category, emoji.Tag, emoji.License,
		emoji.Other, emoji.FilePath, emoji.IsSensitive, emoji.IsRequested,
		emoji.IsAccepted, emoji.IsFinish, emoji.ApproveCount, emoji.DisapproveCount,
		emoji.ResponseFlag, emoji.NowStateIndex, emoji.ModerationMessageID,
		emoji.UserThreadID, emoji.ID,
	)

	return err
}

func (r *SQLiteEmojiRepository) ResetState(emoji *entity.Emoji) error {
	emoji.IsSensitive = false
	emoji.IsAccepted = false
	emoji.IsRequested = false

	query := `
	UPDATE emoji_requests SET
		is_sensitive = 0, is_accepted = 0, is_requested = 0,
		updated_at = CURRENT_TIMESTAMP
	WHERE id = ?`

	_, err := r.db.Exec(query, emoji.ID)
	return err
}

func (r *SQLiteEmojiRepository) Close() error {
	return r.db.Close()
}