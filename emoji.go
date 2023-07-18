package main

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
)

var (
	emojiProcessList []Emoji
)

type Emoji struct {
	ID              string `json:"id"`
	ChannelID       string `json:"channelID"`
	State           int    `json:"state"`
	Name            string `json:"name"`
	Category        string `json:"category"`
	Tag             string `json:"tag"`
	FilePath        string `json:"filepath"`
	IsSensitive     bool   `json:"isSensitive"`
	ApproveCount    int    `json:"approveCount"`
	DisapproveCount int    `json:"disapproveCount"`
	IsRequested     bool   `json:"isRequested"`
	IsAccepted      bool   `json:"isAccepted"`
	IsFinish        bool   `json:"isFinish"`
}

func newEmojiRequest() Emoji {
	id, _ := uuid.NewUUID()
	emoji := Emoji{
		ID: id.String(),
	}

	emojiProcessList = append(emojiProcessList, emoji)
	return emoji
}

func GetEmoji(id string) (*Emoji, error) {
	for i := range emojiProcessList {
		if emojiProcessList[i].ID == id {
			return &emojiProcessList[i], nil
		}
	}
	return nil, errors.New("emoji not found")
}

func approve(emoji Emoji) {
	if emoji.IsAccepted {
		fmt.Println("[ERROR] 既に絵文字はアップロードされています。")
	}
	uploadToMisskey(emoji)
	emoji.IsFinish = true
	deleteChannel(emoji)
}

func disapprove(emoji Emoji) {
	if emoji.IsAccepted {
		return
	}

	emoji.IsAccepted = false
	emoji.IsFinish = true

	deleteChannel(emoji)
}

func deleteChannel(emoji Emoji) {
	Session.ChannelDelete(emoji.ChannelID)
}
