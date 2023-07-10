package main

import (
	"errors"
	"github.com/google/uuid"
)

var (
	emojiProcessList []Emoji
)

type Emoji struct {
	ID       string `json:"id"`
	State    int    `json:"state"`
	Name     string `json:"name"`
	Tag      string `json:"tag"`
	FilePath string `json:"filepath"`
	NSFW     bool   `json:"nsfw"`
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
