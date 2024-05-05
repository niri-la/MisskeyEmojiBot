package main

import (
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/sirupsen/logrus"
)

type RequestProcessor func(*Emoji, *discordgo.Session, string) Response
type ResponceProcessor func(*Emoji, *discordgo.Session, *discordgo.MessageCreate) Response

type Response struct {
	NextState int
	IsSuccess bool
}

var workflow = map[int]string{
	0: "Default",
	2: "SetName",
	1: "Upload",
	3: "Category",
	4: "Tag",
	5: "License",
	6: "Other",
	7: "Nsfw",
	8: "Check",
}

var request = make(map[string]RequestProcessor)
var response = make(map[string]ResponceProcessor)
var reverseWorkflowMap = make(map[string]int)

func init() {
	for key, value := range workflow {
		reverseWorkflowMap[value] = key
	}
}

func init() {

	// Request
	request["SetName"] = func(emoji *Emoji, s *discordgo.Session, cID string) Response {

		response := Response{
			IsSuccess: true,
		}

		_, err := s.ChannelMessageSend(
			cID,
			"## 絵文字の名前を入力してください。\n実際にMisskey上で絵文字を入力する際は`:emoji-name:`としますが、この`emoji-name`の部分を入力してください。\n入力可能な文字は`小文字アルファベット`, `数字`, `_`です。",
		)
		if err != nil {
			logger.WithFields(logrus.Fields{
				"user":  emoji.RequestUser,
				"emoji": emoji.ID,
			}).Error(err)
			return Response{IsSuccess: false}
		}

		emoji.RequestState = "SetName"

		return response
	}
	request["Upload"] = func(emoji *Emoji, s *discordgo.Session, cID string) Response {

		_, err := s.ChannelMessageSend(cID, "## 絵文字ファイルをDiscord上に添付してください。\n対応ファイルは`.png`,`.jpg`,`.jpeg`,`.gif`です。")
		if err != nil {
			logger.WithFields(logrus.Fields{
				"user":  emoji.RequestUser,
				"emoji": emoji.ID,
			}).Error(err)
			return Response{IsSuccess: false}
		}

		emoji.RequestState = "Upload"

		return Response{IsSuccess: true}
	}
	request["Category"] = func(emoji *Emoji, s *discordgo.Session, cID string) Response {

		response := Response{
			IsSuccess: true,
		}

		_, err := s.ChannelMessageSend(cID, "## 絵文字のカテゴリを入力してください。\n特にない場合は「なし」と入力してください。\nカテゴリ名については絵文字やリアクションを入力する際のメニューを参考にしてください。\n例: `Moji`")
		if err != nil {
			logger.WithFields(logrus.Fields{
				"user":  emoji.RequestUser,
				"emoji": emoji.ID,
			}).Error(err)
			return Response{IsSuccess: false}
		}

		emoji.RequestState = "Category"

		return response
	}
	request["Tag"] = func(emoji *Emoji, s *discordgo.Session, cID string) Response {

		response := Response{
			IsSuccess: true,
		}

		_, err := s.ChannelMessageSend(cID, "## 次に絵文字ファイルに設定するタグ(エイリアス)を入力してください。\n空白を間に挟むと複数設定できます。\n"+
			"これは絵文字の検索を行う際に使用されるため、漢字、ひらがな、カタカナ、ローマ字などのバリエーションがあると利用しやすくなります。\n"+
			"例: `絵文字 えもじ emoji emozi`\n必要がない場合は`tagなし`と入力してください。")
		if err != nil {
			logger.WithFields(logrus.Fields{
				"user":  emoji.RequestUser,
				"emoji": emoji.ID,
			}).Error(err)
			return Response{IsSuccess: false}
		}

		emoji.RequestState = "Tag"

		return response
	}
	request["License"] = func(emoji *Emoji, s *discordgo.Session, cID string) Response {

		response := Response{
			IsSuccess: true,
		}

		_, err := s.ChannelMessageSend(cID, "## ライセンス情報を入力してください。\nこれは絵文字ファイルやその素材に関する権利/所有者を明らかにするために重要なものです。\n入力する内容がない場合は`なし`と入力してください。")
		if err != nil {
			logger.WithFields(logrus.Fields{
				"user":  emoji.RequestUser,
				"emoji": emoji.ID,
			}).Error(err)
			return Response{IsSuccess: false}
		}

		emoji.RequestState = "License"

		return response
	}
	request["Other"] = func(emoji *Emoji, s *discordgo.Session, cID string) Response {

		response := Response{
			IsSuccess: true,
		}

		_, err := s.ChannelMessageSend(cID, "## 備考があれば記載してください。\nこの内容はMisskey上には掲載されません。\n特にない場合は`なし`と入力してください。")
		if err != nil {
			logger.WithFields(logrus.Fields{
				"user":  emoji.RequestUser,
				"emoji": emoji.ID,
			}).Error(err)
			return Response{IsSuccess: false}
		}

		emoji.RequestState = "Other"

		return response
	}
	request["Nsfw"] = func(emoji *Emoji, s *discordgo.Session, cID string) Response {
		response := Response{
			IsSuccess: true,
		}
		s.ChannelMessageSendComplex(cID,
			&discordgo.MessageSend{
				Content: "## 絵文字はセンシティブですか？\n",
				Components: []discordgo.MessageComponent{
					discordgo.ActionsRow{
						Components: []discordgo.MessageComponent{
							&discordgo.Button{
								Label:    "はい",
								CustomID: "nsfw_yes",
								Style:    discordgo.DangerButton,
								Emoji: discordgo.ComponentEmoji{
									Name: "🚢",
								},
							},
							&discordgo.Button{
								Label:    "いいえ",
								CustomID: "nsfw_no",
								Style:    discordgo.PrimaryButton,
								Emoji: discordgo.ComponentEmoji{
									Name: "🚀",
								},
							},
						},
					},
				},
			},
		)
		emoji.RequestState = "Nsfw"
		return response
	}
	request["Check"] = func(emoji *Emoji, s *discordgo.Session, cID string) Response {
		response := Response{
			IsSuccess: true,
		}

		s.ChannelMessageSend(cID, ":---\n")
		s.ChannelMessageSend(cID, "## 最終確認を行います。\n"+
			"- 名前 / Name: **"+emoji.Name+"**\n"+
			"- カテゴリ / Category: **"+emoji.Category+"**\n"+
			"- タグ / Tag: **"+emoji.Tag+"**\n"+
			"- ライセンス / License: **"+emoji.License+"**\n"+
			"- その他 / Other: **"+emoji.Other+"**\n"+
			"- NSFW: **"+strconv.FormatBool(emoji.IsSensitive)+"**\n",
		)
		_, err := s.ChannelMessageSendComplex(cID,
			&discordgo.MessageSend{
				Content: "## 以上で申請しますか?\n",
				Components: []discordgo.MessageComponent{
					discordgo.ActionsRow{
						Components: []discordgo.MessageComponent{
							&discordgo.Button{
								Label:    "はい",
								CustomID: "emoji_request",
								Style:    discordgo.PrimaryButton,
								Emoji: discordgo.ComponentEmoji{
									Name: "📨",
								},
							},
							&discordgo.Button{
								Label:    "最初からやり直す",
								CustomID: "emoji_request_retry",
								Style:    discordgo.DangerButton,
								Emoji: discordgo.ComponentEmoji{
									Name: "🗑️",
								},
							},
						},
					},
				},
			},
		)

		logger.Error(err)

		logger.WithFields(logrus.Fields{
			"user":  emoji.RequestUser,
			"emoji": emoji.ID,
		}).Trace("Emoji Request Check")
		emoji.RequestState = "Check"
		return response
	}

	// Responce
	response["SetName"] = func(emoji *Emoji, s *discordgo.Session, m *discordgo.MessageCreate) Response {

		response := Response{
			IsSuccess: false,
		}

		if len(m.Content) <= 1 {
			send, err := s.ChannelMessageSend(m.ChannelID, ":2文字以上入力してください。")
			if err != nil {
				logger.WithFields(logrus.Fields{
					"event": "request-handler",
					"id":    emoji.ID,
					"user":  m.Member.User.Username,
				}).Error(err)
				return response
			}
			logger.WithFields(logrus.Fields{
				"event":      "request-handler",
				"id":         emoji.ID,
				"user":       m.Author.Username,
				"channel id": send.ChannelID,
			}).Warn("Array length shortage error.")
			return response
		}
		reg := regexp.MustCompile(`[^a-z0-9_]+`)
		input := reg.ReplaceAllStringFunc(strings.ToLower(m.Content), func(s string) string {
			return "_"
		})
		s.ChannelMessageSend(m.ChannelID, ":: 入力されたメッセージ\n [ `"+input+"` ]")
		s.ChannelMessageSend(m.ChannelID, ":---")
		emoji.Name = input
		emoji.ResponseState = "SetName"
		response.IsSuccess = true
		response.NextState = response.NextState + 1
		return response
	}
	response["Upload"] = func(emoji *Emoji, s *discordgo.Session, m *discordgo.MessageCreate) Response {

		response := Response{
			IsSuccess: false,
		}

		if len(m.Attachments) > 0 {
			attachment := m.Attachments[0]
			ext := filepath.Ext(attachment.Filename)
			if !isValidEmojiFile(attachment.Filename) {
				s.ChannelMessageSend(m.ChannelID, "画像ファイルを添付してください。"+
					"対応ファイルは`.png`,`.jpg`,`.jpeg`,`.gif`です。")
				return response
			}
			emoji.FilePath = "./Emoji/" + emoji.ID + ext
			err := emojiDownload(attachment.URL, emoji.FilePath)
			if err != nil {
				s.ChannelMessageSend(m.ChannelID, ": Error! \n"+
					"申請中にエラーが発生しました。URLを確認して再アップロードを行うか、管理者へ問い合わせを行ってください。#01a")
				logger.WithFields(logrus.Fields{
					"event": "emoji-download",
					"id":    emoji.ID,
					"user":  m.Member.User,
					"name":  emoji.Name,
				}).Warn(err)
				return response
			}

			logger.WithFields(logrus.Fields{
				"event": "download",
				"id":    emoji.ID,
				"user":  m.Member.User,
				"name":  emoji.Name,
			}).Trace("Emoji Downloaded")

			file, err := os.Open(emoji.FilePath)
			if err != nil {
				s.ChannelMessageSend(m.ChannelID, ": Error! \n"+
					"申請中にエラーが発生しました。管理者へ問い合わせを行ってください。#01b")
				logger.WithFields(logrus.Fields{
					"event": "file open",
					"id":    emoji.ID,
					"user":  m.Member.User,
					"name":  emoji.Name,
				}).Warn(err)
				return response
			}
			defer file.Close()

			_, err = s.ChannelFileSend(m.ChannelID, emoji.FilePath, file)
			if err != nil {
				s.ChannelMessageSend(m.ChannelID, ": Error! \n"+
					"申請中にエラーが発生しました。管理者へ問い合わせを行ってください。#01d")
				logger.WithFields(logrus.Fields{
					"event": "file send",
					"id":    emoji.ID,
					"user":  m.Member.User,
					"name":  emoji.Name,
				}).Error(err)
				return response
			}

			emoji.ResponseState = "Upload"
			response.IsSuccess = true
			response.NextState = response.NextState + 1

			s.ChannelMessageSend(m.ChannelID, ":---\n")

			return response
		} else {
			s.ChannelMessageSend(m.ChannelID, ": ファイルの添付を行ってください。対応ファイルは`.png`,`.jpg`,`.jpeg`,`.gif`です。")
		}
		return response
	}
	response["Category"] = func(emoji *Emoji, s *discordgo.Session, m *discordgo.MessageCreate) Response {

		response := Response{
			IsSuccess: false,
		}

		emoji.Category = m.Content
		if m.Content == "なし" || m.Content == "その他" {
			emoji.Category = ""
		}
		emoji.ResponseState = "Category"
		response.IsSuccess = true
		response.NextState = response.NextState + 1
		s.ChannelMessageSend(m.ChannelID, ":: 入力されたメッセージ\n [ `"+emoji.Category+"` ]")
		s.ChannelMessageSend(m.ChannelID, ":---\n")

		logger.WithFields(logrus.Fields{
			"event":    "emoji-category",
			"id":       emoji.ID,
			"user":     m.Member.User,
			"name":     emoji.Name,
			"category": emoji.Category,
		}).Trace("Set emoji category.")

		return response
	}
	response["Tag"] = func(emoji *Emoji, s *discordgo.Session, m *discordgo.MessageCreate) Response {

		response := Response{
			IsSuccess: false,
		}

		input := strings.Replace(m.Content, "　", " ", -1)
		if input == "tagなし" {
			input = ""
		}
		s.ChannelMessageSend(m.ChannelID, ":: 入力されたメッセージ\n [ `"+input+"` ]")
		s.ChannelMessageSend(m.ChannelID, ":---")

		emoji.ResponseState = "Tag"
		emoji.Tag = input

		response.IsSuccess = true
		response.NextState = response.NextState + 1

		logger.WithFields(logrus.Fields{
			"event": "emoji-tag",
			"id":    emoji.ID,
			"user":  m.Member.User,
			"name":  emoji.Name,
			"tag":   emoji.Tag,
		}).Trace("Set emoji tag.")

		return response
	}
	response["License"] = func(emoji *Emoji, s *discordgo.Session, m *discordgo.MessageCreate) Response {

		response := Response{
			IsSuccess: false,
		}

		emoji.ResponseState = "License"
		input := m.Content
		if input == "なし" {
			input = ""
		}
		emoji.License = input

		s.ChannelMessageSend(m.ChannelID, ":: 入力されたメッセージ\n [ `"+input+"` ]")
		s.ChannelMessageSend(m.ChannelID, ":---")

		response.IsSuccess = true
		response.NextState = response.NextState + 1

		logger.WithFields(logrus.Fields{
			"event": "emoji-license",
			"id":    emoji.ID,
			"user":  m.Member.User,
			"name":  emoji.Name,
			"tag":   emoji.Tag,
		}).Trace("Set emoji license.")

		return response
	}
	response["Other"] = func(emoji *Emoji, s *discordgo.Session, m *discordgo.MessageCreate) Response {

		response := Response{
			IsSuccess: false,
		}

		emoji.ResponseState = "Other"
		input := m.Content
		if input == "なし" {
			input = ""
		}
		emoji.Other = input

		s.ChannelMessageSend(m.ChannelID, ":: 入力されたメッセージ\n [ `"+input+"` ]")
		s.ChannelMessageSend(m.ChannelID, ":---")

		response.IsSuccess = true
		response.NextState = response.NextState + 1

		logger.WithFields(logrus.Fields{
			"event": "emoji-other",
			"id":    emoji.ID,
			"user":  m.Member.User,
			"name":  emoji.Name,
			"tag":   emoji.Tag,
		}).Trace("Set emoji Bikou.")

		return response
	}
	response["Nsfw"] = func(emoji *Emoji, s *discordgo.Session, m *discordgo.MessageCreate) Response {
		// dummy
		return Response{IsSuccess: false}
	}
	response["Check"] = func(emoji *Emoji, s *discordgo.Session, m *discordgo.MessageCreate) Response {
		// dummy
		return Response{IsSuccess: false}
	}
}

func ProcessNextRequest(emoji *Emoji, s *discordgo.Session, id string) bool {
	requestIndex := reverseWorkflowMap[emoji.RequestState]
	logger.WithFields(logrus.Fields{
		"emoji id":       emoji.ID,
		"request index":  requestIndex,
		"response index": requestIndex,
	}).Trace("Emoji Processing (request)...")
	r1 := request[workflow[requestIndex+1]](emoji, s, id)
	return r1.IsSuccess
}

func Process(emoji *Emoji, s *discordgo.Session, m *discordgo.MessageCreate) bool {
	// 0. まずrequestを確認する(初期はRequest及びResponseは0である)
	// 1. 両者が等しい時はRequestを1進める
	// 2. RequestよりResponseが小さい場合はResponse待ちなのでResponseに値を渡す
	// 3. Responseが完了したらResponseを1すすめる。
	// 4. 1に戻る
	// 最終的に次の値がない場合は終了する。
	requestIndex := reverseWorkflowMap[emoji.RequestState]
	responseIndex := reverseWorkflowMap[emoji.ResponseState]

	logger.WithFields(logrus.Fields{
		"emoji id":       emoji.ID,
		"user":           m.Author.Username,
		"request index":  requestIndex,
		"response index": requestIndex,
	}).Trace("Emoji Processing...")

	if requestIndex == responseIndex {
		r1 := request[workflow[requestIndex+1]](emoji, s, m.ChannelID)
		return r1.IsSuccess
	}

	if requestIndex > responseIndex {
		r2 := response[workflow[responseIndex+1]](emoji, s, m)
		if r2.IsSuccess {
			Process(emoji, s, m)
		}
		return r2.IsSuccess
	}
	return false
}

func first(emoji *Emoji, s *discordgo.Session, id string) {
	request[workflow[1]](emoji, s, id)
}

func emojiModerationReaction(s *discordgo.Session, m *discordgo.MessageReactionAdd) {
	if m.UserID == s.State.User.ID {
		return
	}

	channel, _ := s.Channel(m.ChannelID)
	var emoji *Emoji
	found := false

	for _, e := range emojiProcessList {
		if channel.Name == e.ID {
			emoji = &e
			found = true
			break
		}
	}

	if !found {
		return
	}

	emoji, err := GetEmoji(emoji.ID)

	if err != nil {
		return
	}

	if emoji.IsFinish {
		logger.WithFields(logrus.Fields{
			"event": "emoji",
			"id":    emoji.ID,
			"user":  m.Member.User.Username,
			"name":  emoji.Name,
		}).Error("already finished emoji request.")
		return
	}

	roleCount, err := countMembersWithRole(s, GuildID, ModeratorID)
	if err != nil {
		logger.WithFields(logrus.Fields{
			"event":         "emoji",
			"id":            emoji.ID,
			"user":          m.Member.User.Username,
			"name":          emoji.Name,
			"moderation id": ModeratorID,
		}).Error("Invalid moderation ID")
		return
	}

	msg, err := s.ChannelMessage(channel.ID, m.MessageID)
	if err != nil {
		logger.WithFields(logrus.Fields{
			"event": "emoji",
			"id":    emoji.ID,
			"user":  m.Member.User.Username,
			"name":  emoji.Name,
		}).Error(err)
		return
	}

	var apCount = 0
	var dsCount = 0

	for _, reaction := range msg.Reactions {
		if reaction.Emoji.Name == "🆗" {
			apCount = reaction.Count
		} else if reaction.Emoji.Name == "🆖" {
			dsCount = reaction.Count
		}

	}

	emoji.ApproveCount = apCount
	emoji.DisapproveCount = dsCount

	if emoji.DisapproveCount-1 >= roleCount || (isDebug && emoji.DisapproveCount-1 >= 1) {
		emoji.disapprove()
		s.ChannelMessageSend(m.ChannelID, "## 申請は却下されました")
		closeThread(m.ChannelID, emoji.ModerationMessageID)
		return
	}

	if emoji.ApproveCount-1 >= roleCount || (isDebug && emoji.ApproveCount-1 >= 1) {
		emoji.approve()
		s.ChannelMessageSend(m.ChannelID, "## 絵文字はアップロードされました")
		closeThread(m.ChannelID, emoji.ModerationMessageID)
		return
	}

}
