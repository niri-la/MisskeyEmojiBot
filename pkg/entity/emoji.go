package entity

import (
	"path/filepath"
	"time"
)

var (
	validExtensions = map[string]bool{
		".png":  true,
		".jpg":  true,
		".jpeg": true,
		".gif":  true,
	}
)

type Emoji struct {
	ID          string `gorm:"primaryKey" json:"id"`
	ChannelID   string `gorm:"not null" json:"channelID"`
	RequestUser string `gorm:"not null" json:"requestUser"`
	Name        string `json:"name"`
	Category    string `json:"category"`
	Tag         string `json:"tag"`
	License     string `json:"license"`
	Other       string `json:"other"`
	FilePath    string `json:"filepath"`
	IsSensitive bool   `gorm:"default:false" json:"isSensitive"`

	IsRequested bool `gorm:"default:false;index:idx_finish_status" json:"isRequested"`
	IsAccepted  bool `gorm:"default:false;index:idx_finish_status,idx_finish_accepted" json:"isAccepted"`
	IsFinish    bool `gorm:"default:false;index:idx_finish_status,idx_finish_accepted" json:"isFinish"`

	ApproveCount    int `gorm:"default:0" json:"approveCount"`
	DisapproveCount int `gorm:"default:0" json:"disapproveCount"`

	ResponseFlag        bool      `gorm:"default:false" json:"responseState"`
	NowStateIndex       int       `gorm:"default:0" json:"nowStateIndex"`
	ModerationMessageID string    `json:"moderationMessageID"`
	UserThreadID        string    `json:"userThreadID"`
	StartAt             time.Time `gorm:"default:CURRENT_TIMESTAMP;index:idx_start_at" json:"startAt"`

	// GORM timestamps
	CreatedAt time.Time `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updatedAt"`
}

func IsValidEmojiFile(fileName string) bool {
	fileExtension := filepath.Ext(fileName)
	_, exists := validExtensions[fileExtension]
	return exists
}
