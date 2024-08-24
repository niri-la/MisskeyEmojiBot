package entity

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/yitsushi/go-misskey/core"
	"github.com/yitsushi/go-misskey/models"
	"github.com/yitsushi/go-misskey/services/notes"
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
	ID                  string    `json:"id"`
	ChannelID           string    `json:"channelID"`
	ResponseState       string    `json:"responseState"`
	RequestState        string    `json:"requestState"`
	Name                string    `json:"name"`
	Category            string    `json:"category"`
	Tag                 string    `json:"tag"`
	License             string    `json:"license"`
	Other               string    `json:"other"`
	FilePath            string    `json:"filepath"`
	IsSensitive         bool      `json:"isSensitive"`
	RequestUser         string    `json:"requestUser"`
	ApproveCount        int       `json:"approveCount"`
	DisapproveCount     int       `json:"disapproveCount"`
	IsRequested         bool      `json:"isRequested"`
	IsAccepted          bool      `json:"isAccepted"`
	IsFinish            bool      `json:"isFinish"`
	ModerationMessageID string    `json:"moderationMessageID"`
	UserThreadID        string    `json:"userThreadID"`
	StartAt             time.Time `json:"startAt"`
}

func GetEmoji(id string) (*Emoji, error) {
	for i := range emojiProcessList {
		if emojiProcessList[i].ID == id {
			return &emojiProcessList[i], nil
		}
	}
	return nil, errors.New("emoji not found")
}

func newEmojiRequest(user string) *Emoji {
	id := uuid.New()
	emoji := Emoji{
		ID: id.String(),
	}
	emoji.RequestUser = user
	emoji.StartAt = time.Now()
	emojiProcessList = append(emojiProcessList, emoji)
	return &emoji
}

func emojiReconstruction() []Emoji {
	var accepted []Emoji
	var reconstruction []Emoji
	for _, emoji := range emojiProcessList {
		if emoji.IsFinish {
			if emoji.IsAccepted {
				accepted = append(accepted, emoji)
			}
		} else {
			reconstruction = append(reconstruction, emoji)
		}
	}
	emojiProcessList = reconstruction
	return accepted
}

func noteEmojiAdded(emojis []Emoji) {
	var builder strings.Builder
	for _, emoji := range emojis {
		builder.WriteString(":" + emoji.Name + ":")
	}

	message := core.NewString("#にりらみすきー部 \n絵文字が追加されました\n" +
		builder.String())

	note(notes.CreateRequest{
		Visibility: models.VisibilityPublic,
		Text:       message,
		LocalOnly:  true,
	})
}

func (emoji *Emoji) reset() {
	emoji.RequestState = workflow[0]
	emoji.ResponseState = workflow[0]
	emoji.IsSensitive = false
	emoji.IsAccepted = false
	emoji.IsRequested = false
}

func (emoji *Emoji) abort() {
	remove(*emoji)
	emoji.reset()
	emoji.IsFinish = true
}

func remove(val Emoji) {
	var newSlice []Emoji
	for _, v := range emojiProcessList {
		if v.ID != val.ID {
			newSlice = append(newSlice, v)
		}
	}
	emojiProcessList = newSlice
}

func deleteEmoji(filePath string) error {
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
