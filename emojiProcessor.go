package main

import (
	"github.com/bwmarrin/discordgo"
	"github.com/sirupsen/logrus"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

func runEmojiProcess(emoji *Emoji, s *discordgo.Session, m *discordgo.MessageCreate) {
	switch emoji.State {
	// first Emojiの名前を設定
	case 0:
		if len(m.Content) <= 1 {
			s.ChannelMessageSend(m.ChannelID, ":2文字以上入力してください。")
			return
		}
		emoji.ChannelID = m.ChannelID
		reg := regexp.MustCompile(`[^a-z0-9_]+`)
		result := reg.ReplaceAllStringFunc(m.Content, func(s string) string {
			return "_"
		})
		input := strings.ToLower(result)
		s.ChannelMessageSend(m.ChannelID, ":: 入力されたメッセージ\n [ `"+input+"` ]")
		s.ChannelMessageSend(m.ChannelID, ":---")
		s.ChannelMessageSend(m.ChannelID, "2. 次に絵文字ファイルをDiscord上に添付してください。対応ファイルは`.png`,`.jpg`,`.jpeg`,`.gif`です。")
		emoji.Name = input
		emoji.State = emoji.State + 1
		break
	// first Emojiのファイルを入力 // 表示させるか迷う
	case 1:

		if len(m.Attachments) > 0 {
			attachment := m.Attachments[0]
			ext := filepath.Ext(attachment.Filename)
			if !isValidEmojiFile(attachment.Filename) {
				s.ChannelMessageSend(m.ChannelID, "画像ファイルを添付してください。"+
					"対応ファイルは`.png`,`.jpg`,`.jpeg`,`.gif`です。")
				return
			}
			emoji.FilePath = "./Emoji/" + emoji.ID + ext
			err := emojiDownload(attachment.URL, emoji.FilePath)
			if err != nil {
				s.ChannelMessageSend(m.ChannelID, ": Error! \n"+
					"申請中にエラーが発生しました。URLを確認して再アップロードを行うか、管理者へ問い合わせを行ってください。#01a")
				logger.WithFields(logrus.Fields{
					"event": "download",
					"id":    emoji.ID,
					"user":  m.Member.User,
					"name":  emoji.Name,
				}).Warn(err)
				return
			}

			logger.WithFields(logrus.Fields{
				"event": "download",
				"id":    emoji.ID,
				"user":  m.Member.User,
				"name":  emoji.Name,
			}).Debug("Emoji Downloaded")

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
				return
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
				}).Fatal(err)
				return
			}

			emoji.State = emoji.State + 1

			s.ChannelMessageSend(m.ChannelID, ":---\n")
			s.ChannelMessageSend(m.ChannelID, "3. 絵文字のカテゴリを入力してください。特にない場合は「なし」と入力してください。カテゴリ名については絵文字やリアクションを入力する際のメニューを参考にしてください。 例: `Moji`")
		} else {
			s.ChannelMessageSend(m.ChannelID, ": ファイルの添付を行ってください。対応ファイルは`.png`,`.jpg`,`.jpeg`,`.gif`です。")
		}
		break
	// Categoryの設定
	case 2:
		emoji.Category = m.Content
		if m.Content == "なし" {
			emoji.Category = ""
		}
		emoji.State = emoji.State + 1
		s.ChannelMessageSend(m.ChannelID, ":: 入力されたメッセージ\n [ `"+m.Content+"` ]")
		s.ChannelMessageSend(m.ChannelID, ":---\n")
		s.ChannelMessageSend(m.ChannelID, "4. 次に絵文字ファイルに設定するタグ(エイリアス)を入力してください。空白を間に挟むと複数設定できます。これは絵文字の検索をする際に使用されます。 例: `絵文字 えもじ エモジ `")
		break
	case 3:
		input := strings.Replace(m.Content, "　", " ", -1)
		s.ChannelMessageSend(m.ChannelID, ":: 入力されたメッセージ\n [ `"+input+"` ]")
		s.ChannelMessageSend(m.ChannelID, ":---")
		s.ChannelMessageSendComplex(m.ChannelID,
			&discordgo.MessageSend{
				Content: "4. 絵文字はセンシティブですか？\n",
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
		emoji.Tag = input
		emoji.State = emoji.State + 1
		break
	// NSFWかの確認
	case 4:
		logger.Error("実装がおかしい")
		break
	// (licenseの確認) 最終確認
	case 5:
		break
		//// 最終確認
		//case 5:
		//	break
		// タスク終了。アップロード処理へ渡す
		//case 6:
		//	break

	}

}

func emojiLastConfirmation(emoji *Emoji, s *discordgo.Session, ChannelID string) {
	s.ChannelMessageSend(ChannelID, ":---\n")
	s.ChannelMessageSend(ChannelID, "最終確認を行います。\n"+
		"Name: "+emoji.Name+"\n"+
		"Category: "+emoji.Category+"\n"+
		"Tag: "+emoji.Tag+"\n"+
		"isNSFW: "+strconv.FormatBool(emoji.IsSensitive)+"\n")
	s.ChannelMessageSendComplex(ChannelID,
		&discordgo.MessageSend{
			Content: "以上で申請しますか?\n",
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
								Name: "",
							},
						},
					},
				},
			},
		},
	)
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

	if emoji.DisapproveCount-1 >= roleCount {
		disapprove(*emoji)
		s.ChannelMessageSend(m.ChannelID, "申請は却下されました")
		closeThread(m.ChannelID)
		return
	}

	if emoji.ApproveCount-1 >= roleCount {
		approve(*emoji)
		s.ChannelMessageSend(m.ChannelID, "絵文字はアップロードされました")
		closeThread(m.ChannelID)
		return
	}

}

func closeThread(id string) {
	channel, _ := Session.Channel(id)
	if !channel.IsThread() {
		return
	}
	archived := true
	locked := true
	_, err := Session.ChannelEditComplex(channel.ID, &discordgo.ChannelEdit{
		Archived: &archived,
		Locked:   &locked,
	})
	if err != nil {
		panic(err)
	}
}
