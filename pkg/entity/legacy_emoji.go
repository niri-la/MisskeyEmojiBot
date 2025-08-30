package entity

import (
	"time"
)

// LegacyEmoji は旧JSONファイル形式の絵文字データ構造
type LegacyEmoji struct {
	ID                  string    `json:"id"`
	ChannelID           string    `json:"channelID"`
	RequestUser         string    `json:"requestUser"`
	Name                string    `json:"name"`
	Category            string    `json:"category"`
	Tag                 string    `json:"tag"`
	License             string    `json:"license"`
	Other               string    `json:"other"`
	Filepath            string    `json:"filepath"`
	IsSensitive         bool      `json:"isSensitive"`
	IsRequested         bool      `json:"isRequested"`
	IsAccepted          bool      `json:"isAccepted"`
	IsFinish            bool      `json:"isFinish"`
	ApproveCount        int       `json:"approveCount"`
	DisapproveCount     int       `json:"disapproveCount"`
	ResponseState       bool      `json:"responseState"`
	NowStateIndex       int       `json:"nowStateIndex"`
	ModerationMessageID string    `json:"moderationMessageID"`
	UserThreadID        string    `json:"userThreadID"`
	StartAt             time.Time `json:"startAt"`
}

// ToCurrentEmoji converts LegacyEmoji to current Emoji format
func (le *LegacyEmoji) ToCurrentEmoji() *Emoji {
	return &Emoji{
		ID:                  le.ID,
		ChannelID:           le.ChannelID,
		RequestUser:         le.RequestUser,
		Name:                le.Name,
		Category:            le.Category,
		Tag:                 le.Tag,
		License:             le.License,
		Other:               le.Other,
		FilePath:            le.Filepath, // フィールド名の違いに注意
		IsSensitive:         le.IsSensitive,
		IsRequested:         le.IsRequested,
		IsAccepted:          le.IsAccepted,
		IsFinish:            le.IsFinish,
		IsOverwrite:         false, // 新フィールドはデフォルトfalse
		ApproveCount:        le.ApproveCount,
		DisapproveCount:     le.DisapproveCount,
		ResponseFlag:        le.ResponseState, // フィールド名の違いに注意
		NowStateIndex:       le.NowStateIndex,
		ModerationMessageID: le.ModerationMessageID,
		UserThreadID:        le.UserThreadID,
		StartAt:             le.StartAt,
		// 新しいフィールドはゼロ値で初期化
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}