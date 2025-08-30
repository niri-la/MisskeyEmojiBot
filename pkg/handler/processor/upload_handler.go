package processor

import (
	"bytes"
	"os"
	"path/filepath"

	"github.com/bwmarrin/discordgo"

	"MisskeyEmojiBot/pkg/config"
	"MisskeyEmojiBot/pkg/entity"
	"MisskeyEmojiBot/pkg/handler"
	"MisskeyEmojiBot/pkg/repository"
	"MisskeyEmojiBot/pkg/utility"
)

type uploadHandler struct {
	config config.Config
	s3Repo repository.S3Repository
}

func NewUploadHandler(cfg config.Config) handler.EmojiProcessHandler {
	var s3Repo repository.S3Repository
	if cfg.UseS3 {
		s3Repo, _ = repository.NewS3Repository(&cfg)
	}
	return &uploadHandler{config: cfg, s3Repo: s3Repo}
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
			_, _ = s.ChannelMessageSend(m.ChannelID, "画像ファイルを添付してください。"+
				"対応ファイルは`.png`,`.jpg`,`.jpeg`,`.gif`です。")
			return response, nil
		}
		if h.config.UseS3 && h.s3Repo != nil {
			fileData, err := utility.EmojiDownloadToBytes(attachment.URL)
			if err != nil {
				_, _ = s.ChannelMessageSend(m.ChannelID, ": Error! \n"+
					"申請中にエラーが発生しました。URLを確認して再アップロードを行うか、管理者へ問い合わせを行ってください。#01a")
				return response, nil
			}

			key := "emojis/" + emoji.ID + ext
			contentType := repository.GetContentTypeFromExtension(attachment.Filename)
			fileURL, err := h.s3Repo.UploadFile(key, fileData, contentType)
			if err != nil {
				_, _ = s.ChannelMessageSend(m.ChannelID, ": Error! \n"+
					"申請中にエラーが発生しました。S3へのアップロードに失敗しました。#01sa")
				return response, nil
			}

			emoji.FilePath = fileURL

			_, err = s.ChannelMessageSendComplex(m.ChannelID, &discordgo.MessageSend{
				Content: "アップロード完了: " + fileURL,
				Files: []*discordgo.File{
					{
						Name:   attachment.Filename,
						Reader: bytes.NewReader(fileData),
					},
				},
			})
			if err != nil {
				_, _ = s.ChannelMessageSend(m.ChannelID, ": Error! \n"+
					"申請中にエラーが発生しました。管理者へ問い合わせを行ってください。#01sb")
				return response, nil
			}
		} else {
			emoji.FilePath = filepath.Join(h.config.SavePath, emoji.ID+ext)
			err := utility.EmojiDownload(attachment.URL, emoji.FilePath)
			if err != nil {
				_, _ = s.ChannelMessageSend(m.ChannelID, ": Error! \n"+
					"申請中にエラーが発生しました。URLを確認して再アップロードを行うか、管理者へ問い合わせを行ってください。#01a")
				return response, nil
			}

			file, err := os.Open(emoji.FilePath)
			if err != nil {
				_, _ = s.ChannelMessageSend(m.ChannelID, ": Error! \n"+
					"申請中にエラーが発生しました。管理者へ問い合わせを行ってください。#01b")
				return response, nil
			}
			defer func() { _ = file.Close() }()

			_, err = s.ChannelFileSend(m.ChannelID, emoji.FilePath, file)
			if err != nil {
				_, _ = s.ChannelMessageSend(m.ChannelID, ": Error! \n"+
					"申請中にエラーが発生しました。管理者へ問い合わせを行ってください。#01d")
				return response, nil
			}
		}

		response.IsSuccess = true

		_, _ = s.ChannelMessageSend(m.ChannelID, "# ----------\n")

		return response, nil
	}
	_, _ = s.ChannelMessageSend(m.ChannelID, ": ファイルの添付を行ってください。対応ファイルは`.png`,`.jpg`,`.jpeg`,`.gif`です。")
	return response, nil
}
