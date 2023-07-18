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
	RequestUser     string `json:"requestUser"`
	ApproveCount    int    `json:"approveCount"`
	DisapproveCount int    `json:"disapproveCount"`
	IsRequested     bool   `json:"isRequested"`
	IsAccepted      bool   `json:"isAccepted"`
	IsFinish        bool   `json:"isFinish"`
}

func newEmojiRequest(user string) Emoji {
	id, _ := uuid.NewUUID()
	emoji := Emoji{
		ID: id.String(),
	}
	emoji.RequestUser = user
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
	sendDirectMessage(emoji, "申請された絵文字は登録されました。"+emoji.ID)
	deleteChannel(emoji)
}

func disapprove(emoji Emoji) {
	if emoji.IsAccepted {
		return
	}

	emoji.IsAccepted = false
	emoji.IsFinish = true
	sendDirectMessage(emoji, "申請された絵文字は却下されました。 "+emoji.ID)
	deleteChannel(emoji)
}

func deleteChannel(emoji Emoji) {
	Session.ChannelDelete(emoji.ChannelID)
}

func sendDirectMessage(emoji Emoji, message string) {
	user, _ := Session.User(emoji.RequestUser)
	direct, _ := Session.UserChannelCreate(user.ID)
	_, err := Session.ChannelMessageSend(direct.ID, message)
	if err != nil {
		fmt.Println("Error sending message: ", err)
		return
	}
}
