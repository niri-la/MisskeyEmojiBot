package processor

import (
	"MisskeyEmojiBot/pkg/entity"
	"MisskeyEmojiBot/pkg/handler"
	"MisskeyEmojiBot/pkg/utility"
	"os"
	"path/filepath"

	"github.com/bwmarrin/discordgo"
)

type uploadHandler struct {
}

func NewUploadHandler() handler.EmojiProcessHandler {
	return &uploadHandler{}
}

func (h *uploadHandler) Request(emoji *entity.Emoji, s *discordgo.Session, cID string) (entity.Response, error) {
	_, err := s.ChannelMessageSend(cID, "## 絵文字ファイルをDiscord上に添付してください。\n対応ファイルは`.png`,`.jpg`,`.jpeg`,`.gif`です。")
	if err != nil {
		return entity.Response{IsSuccess: false}, err
	}

	return entity.Response{IsSuccess: true}, nil
}

func (h *uploadHandler) Response(emoji *entity.Emoji, s *discordgo.Session, m *discordgo.MessageCreate) (entity.Response, error) {
	response := entity.Response{
		IsSuccess: false,
	}

	if len(m.Attachments) > 0 {
		attachment := m.Attachments[0]
		ext := filepath.Ext(attachment.Filename)
		if !entity.IsValidEmojiFile(attachment.Filename) {
			s.ChannelMessageSend(m.ChannelID, "画像ファイルを添付してください。"+
				"対応ファイルは`.png`,`.jpg`,`.jpeg`,`.gif`です。")
			return response, nil
		}
		emoji.FilePath = "./Emoji/" + emoji.ID + ext
		err := utility.EmojiDownload(attachment.URL, emoji.FilePath)
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, ": Error! \n"+
				"申請中にエラーが発生しました。URLを確認して再アップロードを行うか、管理者へ問い合わせを行ってください。#01a")
			return response, err
		}

		file, err := os.Open(emoji.FilePath)
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, ": Error! \n"+
				"申請中にエラーが発生しました。管理者へ問い合わせを行ってください。#01b")
			return response, nil
		}
		defer file.Close()

		_, err = s.ChannelFileSend(m.ChannelID, emoji.FilePath, file)
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, ": Error! \n"+
				"申請中にエラーが発生しました。管理者へ問い合わせを行ってください。#01d")
			return response, nil
		}

		response.IsSuccess = true

		s.ChannelMessageSend(m.ChannelID, "# ----------\n")

		return response, nil
	} else {
		s.ChannelMessageSend(m.ChannelID, ": ファイルの添付を行ってください。対応ファイルは`.png`,`.jpg`,`.jpeg`,`.gif`です。")
	}
	return response, nil
}
