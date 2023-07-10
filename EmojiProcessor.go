package main

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

func RunEmojiProcess(emoji *Emoji, s *discordgo.Session, m *discordgo.MessageCreate) {
	switch emoji.State {
	// first Emojiの名前を設定
	case 0:
		reg := regexp.MustCompile(`[^a-zA-Z_]+`)
		result := reg.ReplaceAllStringFunc(m.Content, func(s string) string {
			return "_"
		})
		input := strings.ToLower(result)
		s.ChannelMessageSend(m.ChannelID, ": input [ "+input+"]")
		s.ChannelMessageSend(m.ChannelID, ":---")
		s.ChannelMessageSend(m.ChannelID, "2. 次に絵文字ファイルをDiscord上に添付してください。")
		emoji.Name = input
		emoji.State = emoji.State + 1
		break
	// first Emojiのファイルを入力 // 表示させるか迷う
	case 1:

		if len(m.Attachments) > 0 {
			attachment := m.Attachments[0]
			ext := filepath.Ext(attachment.Filename)

			if ext != ".png" && ext != ".jpg" && ext != ".jpeg" && ext != ".gif" {
				s.ChannelMessageSend(m.ChannelID, "画像ファイルを添付してください。"+
					"対応ファイルは`.png`,`.jpg`,`.jpeg`,`.gif`です。")
				return
			}

			response, err := http.Get(attachment.URL)
			if err != nil {
				s.ChannelMessageSend(m.ChannelID, ": Error! \n"+
					"ファイルにアクセスできませんでした。ファイルのURLを確認するか、管理者に問い合わせを行ってください。")
				return
			}

			defer response.Body.Close()

			emoji.FilePath = emoji.ID + ext

			file, err := os.Create(emoji.FilePath)
			if err != nil {
				s.ChannelMessageSend(m.ChannelID, ": Error! \n"+
					"申請中にエラーが発生しました。管理者へ問い合わせを行ってください。#01a")
				return
			}

			defer file.Close()

			_, err = io.Copy(file, response.Body)
			if err != nil {
				s.ChannelMessageSend(m.ChannelID, ": Error! \n"+
					"申請中にエラーが発生しました。管理者へ問い合わせを行ってください。#01b")
				return
			}

			file, err = os.Open(emoji.FilePath)
			if err != nil {
				s.ChannelMessageSend(m.ChannelID, ": Error! \n"+
					"申請中にエラーが発生しました。管理者へ問い合わせを行ってください。#01c")
				return
			}
			defer file.Close()

			fmt.Printf("[Emoji] File %s downloaded. (%s)\n", attachment.Filename, emoji.ID)

			_, err = s.ChannelFileSend(m.ChannelID, emoji.FilePath, file)
			if err != nil {
				s.ChannelMessageSend(m.ChannelID, ": Error! \n"+
					"申請中にエラーが発生しました。管理者へ問い合わせを行ってください。#01d")
				return
			}

			emoji.State = emoji.State + 1

			s.ChannelMessageSend(m.ChannelID, ":---\n")
			s.ChannelMessageSend(m.ChannelID, "3. 次に絵文字ファイルに設定するタグを入力してください。空白を間に挟むと複数設定できます。例: `絵文字 えもじ エモジ `")
		}
		break
	// tagの設定
	case 2:
		s.ChannelMessageSend(m.ChannelID, ": input [ "+m.Content+"]")
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
		emoji.Tag = m.Content
		emoji.State = emoji.State + 1
		break
	// NSFWかの確認
	case 3:
		fmt.Println("[ERROR] 実装がおかしい")
		break
	// (licenseの確認) 最終確認
	case 4:
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
	s.ChannelMessageSend(ChannelID, "最終確認を行います。\n")
	s.ChannelMessageSend(ChannelID, "Name: "+emoji.Name+"\n")
	s.ChannelMessageSend(ChannelID, "Tag: "+emoji.Tag+"\n")
	s.ChannelMessageSend(ChannelID, "isNSFW: "+strconv.FormatBool(emoji.NSFW)+"\n")
	s.ChannelMessageSendComplex(ChannelID,
		&discordgo.MessageSend{
			Content: "以上で申請しますか?\n",
			Components: []discordgo.MessageComponent{
				discordgo.ActionsRow{
					Components: []discordgo.MessageComponent{
						&discordgo.Button{
							Label:    "はい",
							CustomID: "Request",
							Style:    discordgo.PrimaryButton,
							Emoji: discordgo.ComponentEmoji{
								Name: "📨",
							},
						},
						&discordgo.Button{
							Label:    "やり直す",
							CustomID: "retry",
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
