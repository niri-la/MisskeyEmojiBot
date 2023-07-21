package main

import (
	"errors"
	"github.com/google/uuid"
	debug "github.com/sirupsen/logrus"
	"github.com/yitsushi/go-misskey/core"
	"github.com/yitsushi/go-misskey/models"
	"github.com/yitsushi/go-misskey/services/notes"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
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
	ID                  string `json:"id"`
	ChannelID           string `json:"channelID"`
	ResponseState       string `json:"responseState"`
	RequestState        string `json:"requestState"`
	Name                string `json:"name"`
	Category            string `json:"category"`
	Tag                 string `json:"tag"`
	License             string `json:"license"`
	FilePath            string `json:"filepath"`
	IsSensitive         bool   `json:"isSensitive"`
	RequestUser         string `json:"requestUser"`
	ApproveCount        int    `json:"approveCount"`
	DisapproveCount     int    `json:"disapproveCount"`
	IsRequested         bool   `json:"isRequested"`
	IsAccepted          bool   `json:"isAccepted"`
	IsFinish            bool   `json:"isFinish"`
	ModerationMessageID string `json:"moderationMessageID"`
	UserThreadID        string `json:"userThreadID"`
}

func newEmojiRequest(user string) *Emoji {
	id, _ := uuid.NewUUID()
	emoji := Emoji{
		ID: id.String(),
	}
	emoji.RequestUser = user
	emojiProcessList = append(emojiProcessList, emoji)
	return &emoji
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
		u, _ := Session.User(emoji.RequestUser)
		logger.WithFields(debug.Fields{
			"event": "nsfw",
			"id":    emoji.ID,
			"user":  u.Username,
			"name":  emoji.Name,
		}).Warn("already uploaded")
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
		u, _ := Session.User(emoji.RequestUser)
		logger.WithFields(debug.Fields{
			"event": "emoji",
			"id":    emoji.ID,
			"user":  u.Username,
			"name":  emoji.Name,
		}).Error(err)
		return
	}
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
	return reconstruction
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

func (e Emoji) reset() {
	e.RequestState = workflow[0]
	e.ResponseState = workflow[0]
	e.IsSensitive = false
	e.IsAccepted = false
	e.IsRequested = false
}

func (e Emoji) abort() {
	remove(e)
	e.reset()
	e.IsFinish = true
}

func remove(val Emoji) {
	var newSlice []Emoji
	for _, v := range emojiProcessList {
		if v != val {
			newSlice = append(newSlice, v)
		}
	}
	emojiProcessList = newSlice
}

func emojiDownload(url string, filePath string) error {
	response, err := http.Get(url)
	if err != nil {
		return err
	}

	defer response.Body.Close()

	dirPath := filepath.Dir(filePath)
	// ディレクトリが存在するかチェック
	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		// ディレクトリが存在しない場合、作成する
		os.MkdirAll(dirPath, os.ModePerm)
	}

	file, err := os.Create(filePath)
	if err != nil {
		return err
	}

	defer file.Close()

	_, err = io.Copy(file, response.Body)
	if err != nil {
		return err
	}

	return nil
}

func deleteEmoji(filePath string) {
	err := os.Remove(filePath)
	if err != nil {
		logger.WithFields(debug.Fields{
			"event": "emoji",
			"path":  filePath,
		}).Error(err)
	}
}

func isValidEmojiFile(fileName string) bool {
	fileExtension := filepath.Ext(fileName)
	_, exists := validExtensions[fileExtension]
	return exists
}
