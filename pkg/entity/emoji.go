package entity

import (
	"os"
	"path/filepath"
	"time"
)

var (
	emojiProcessList []Emoji
	validExtensions  = map[string]bool{
		".png":  true,
		".jpg":  true,
		".jpeg": true,
		".gif":  true,
	}
)

type Emoji struct {
	ID          string `json:"id"`
	ChannelID   string `json:"channelID"`
	RequestUser string `json:"requestUser"`
	Name        string `json:"name"`
	Category    string `json:"category"`
	Tag         string `json:"tag"`
	License     string `json:"license"`
	Other       string `json:"other"`
	FilePath    string `json:"filepath"`
	IsSensitive bool   `json:"isSensitive"`

	IsRequested bool `json:"isRequested"`
	IsAccepted  bool `json:"isAccepted"`
	IsFinish    bool `json:"isFinish"`

	ApproveCount    int `json:"approveCount"`
	DisapproveCount int `json:"disapproveCount"`

	ResponseFlag        bool      `json:"responseState"`
	NowStateIndex       int       `json:"nowStateIndex"`
	ModerationMessageID string    `json:"moderationMessageID"`
	UserThreadID        string    `json:"userThreadID"`
	StartAt             time.Time `json:"startAt"`
}

func (emoji *Emoji) Reset() {
	emoji.IsSensitive = false
	emoji.IsAccepted = false
	emoji.IsRequested = false
}

func (emoji *Emoji) deleteEmoji(filePath string) error {
	err := os.Remove(filePath)
	if err != nil {
		return err
	}
	return nil
}

func IsValidEmojiFile(fileName string) bool {
	fileExtension := filepath.Ext(fileName)
	_, exists := validExtensions[fileExtension]
	return exists
}
