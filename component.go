package main

import (
	debug "github.com/sirupsen/logrus"
	"os"
	"strconv"

	"github.com/bwmarrin/discordgo"
)

var (
	Components         = make([]*discordgo.ApplicationCommand, 0)
	ComponentsHandlers = make(map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate))
)

func init() {
	initComponent()
	nsfwComponent()
	newEmojiComponent()
	newEmojiChannelComponent()
}

func addComponent(command *discordgo.ApplicationCommand, fn func(s *discordgo.Session, i *discordgo.InteractionCreate)) {
	_, exist := ComponentsHandlers[command.Name]
	if exist {
		logger.WithFields(debug.Fields{
			"event": "command",
			"name":  command.Name,
		}).Panic("command already existed.")
	}
	ComponentsHandlers[command.Name] = fn
	Components = append(Components, command)
}

func initComponent() {
	// init_channel
	addComponent(
		&discordgo.ApplicationCommand{
			Name: "init_channel",
		},
		func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			channelID := i.MessageComponentData().Values[0]

			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Flags: discordgo.MessageFlagsEphemeral,
					Content: "選択チャンネル <#" + i.MessageComponentData().Values[0] + ">\n" +
						"初期設定を行いました。",
				},
			})

			s.ChannelMessageSendComplex(channelID,
				&discordgo.MessageSend{
					Content: "こんにちは！絵文字申請チャンネルへようこそ！\n",
					Components: []discordgo.MessageComponent{
						discordgo.ActionsRow{
							Components: []discordgo.MessageComponent{
								&discordgo.Button{
									Label:    "絵文字の申請をする / Requset emoji",
									CustomID: "new_emoji_channel",
									Style:    discordgo.PrimaryButton,
									Emoji: discordgo.ComponentEmoji{
										Name: "🏗️",
									},
								},
							},
						},
					},
				},
			)

			overwrites := []*discordgo.PermissionOverwrite{
				{
					ID:   ModeratorID,
					Type: discordgo.PermissionOverwriteTypeRole,
					Allow: discordgo.PermissionViewChannel |
						discordgo.PermissionSendMessages,
				},
				{
					ID:   BotID,
					Type: discordgo.PermissionOverwriteTypeRole,
					Allow: discordgo.PermissionViewChannel |
						discordgo.PermissionSendMessages,
				},
				{
					ID:   i.GuildID,
					Type: discordgo.PermissionOverwriteTypeRole,
					Deny: discordgo.PermissionViewChannel,
				},
			}

			parent, err := s.Channel(i.ChannelID)

			if err != nil {
				returnFailedMessage(s, i, "Could not retrieve channel")
				return
			}

			channel, err := s.GuildChannelCreateComplex(GuildID, discordgo.GuildChannelCreateData{
				Type:                 discordgo.ChannelTypeGuildText,
				Name:                 ModerationChannelName,
				ParentID:             parent.ParentID,
				PermissionOverwrites: overwrites,
			})

			s.ChannelMessageSend(
				channel.ID,
				": モデレーション用チャンネルです。\nここに各種申請のスレッドが生成されます。",
			)

			logger.Debug(":: Create a moderation channel")

			return

		},
	)
}

func nsfwComponent() {
	// nsfw_yes
	addComponent(
		&discordgo.ApplicationCommand{
			Name: "nsfw_yes",
		},
		func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			channel, _ := s.Channel(i.ChannelID)
			emoji, err := GetEmoji(channel.Name[6:])
			if err != nil {
				s.ChannelMessageSend(
					channel.ID,
					"設定に失敗しました。管理者に問い合わせを行ってください。 #03a\n",
				)

				logger.WithFields(debug.Fields{
					"event": "nsfw",
					"id":    emoji.ID,
					"user":  i.Member.User,
					"name":  emoji.Name,
				}).Error(err)
				return
			}

			if emoji.IsRequested {
				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Flags:   discordgo.MessageFlagsEphemeral,
						Content: "既に申請は終了しています\n",
					},
				})
				return
			}

			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Flags:   discordgo.MessageFlagsEphemeral,
					Content: "NSFWに設定されました\n",
				},
			})
			emoji.IsSensitive = true
			emoji.State = 5
			emojiLastConfirmation(emoji, s, i.ChannelID)
		},
	)

	// nsfw_no
	addComponent(
		&discordgo.ApplicationCommand{
			Name: "nsfw_no",
		},
		func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			channel, _ := s.Channel(i.ChannelID)
			emoji, err := GetEmoji(channel.Name[6:])
			if err != nil {
				s.ChannelMessageSend(
					channel.ID,
					"設定に失敗しました。管理者に問い合わせを行ってください。 #03a\n",
				)

				logger.WithFields(debug.Fields{
					"event": "nsfw",
					"id":    emoji.ID,
					"user":  i.Member.User,
					"name":  emoji.Name,
				}).Error(err)
				return

			}

			if emoji.IsRequested {
				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Flags:   discordgo.MessageFlagsEphemeral,
						Content: "既に申請は終了しています\n",
					},
				})
				return
			}

			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Flags:   discordgo.MessageFlagsEphemeral,
					Content: "非NSFWに設定されました\n",
				},
			})

			emoji.IsSensitive = false
			emoji.State = 5
			emojiLastConfirmation(emoji, s, i.ChannelID)

		},
	)
}

func newEmojiComponent() {
	// emoji_request
	addComponent(
		&discordgo.ApplicationCommand{
			Name: "emoji_request",
		},
		func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			channel, _ := s.Channel(i.ChannelID)
			emoji, err := GetEmoji(channel.Name[6:])
			if err != nil {
				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Flags:   discordgo.MessageFlagsEphemeral,
						Content: "設定に失敗しました。管理者に問い合わせを行ってください。\n",
					},
				})
			}

			if emoji.IsRequested {
				s.ChannelMessageSend(
					channel.ID,
					"既に申請していますよ！\n",
				)
				return
			}

			s.ChannelMessageSend(
				channel.ID,
				"申請をしました！\n"+
					"なお、申請結果についてはこちらではお伝えできかねますのでご了承ください。\n"+
					"詳細な申請内容については管理者へお問い合わせください！\n"+
					"この度は申請いただき大変ありがとうございました。\n",
			)

			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Flags:   discordgo.MessageFlagsEphemeral,
					Content: "😎",
				},
			})

			emoji.IsRequested = true

			sendDirectMessage(*emoji, "--- 申請内容 "+emoji.ID+"---\n名前: "+emoji.Name+"\nCategory: "+
				emoji.Category+"\n"+"tag"+emoji.Tag+"\n"+"isNSFW:"+strconv.FormatBool(emoji.IsSensitive)+"\n---")

			send, err := s.ChannelMessageSend(moderationChannel.ID, ":作成者: "+i.Member.User.Username+"\n"+
				":: ID "+emoji.ID)
			if err != nil {
				return
			}

			thread, err := s.MessageThreadStartComplex(moderationChannel.ID, send.ID, &discordgo.ThreadStart{
				Name:                emoji.ID,
				AutoArchiveDuration: 60,
				Invitable:           false,
				RateLimitPerUser:    10,
			})

			s.ChannelMessageSend(thread.ID, ":---\n"+
				"Requested by "+i.Member.User.Username+"\n"+
				":---\n")
			s.ChannelMessageSend(thread.ID,
				"Name: "+emoji.Name+"\n"+
					"Category: "+emoji.Category+"\n"+
					"Tag: "+emoji.Tag+"\n"+
					"isNSFW: "+strconv.FormatBool(emoji.IsSensitive)+"\n")

			file, err := os.Open(emoji.FilePath)
			if err != nil {
				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Flags:   discordgo.MessageFlagsEphemeral,
						Content: "設定に失敗しました。管理者に問い合わせを行ってください。#01b\n",
					},
				})
				return
			}
			defer file.Close()

			lastMessage, err := s.ChannelFileSend(thread.ID, emoji.FilePath, file)
			if err != nil {
				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Flags:   discordgo.MessageFlagsEphemeral,
						Content: "設定に失敗しました。管理者に問い合わせを行ってください。#01d\n",
					},
				})
				return
			}

			s.MessageReactionAdd(thread.ID, lastMessage.ID, "🆗")
			s.MessageReactionAdd(thread.ID, lastMessage.ID, "🆖")

		},
	)

	// emoji_request_retry
	addComponent(
		&discordgo.ApplicationCommand{
			Name: "emoji_request_retry",
		},
		func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			channel, _ := s.Channel(i.ChannelID)
			emoji, err := GetEmoji(channel.Name[6:])
			if err != nil {
				s.ChannelMessageSend(
					channel.ID,
					"設定に失敗しました。管理者に問い合わせを行ってください。 #04a\n",
				)
			}

			if emoji.IsRequested {
				s.ChannelMessageSend(
					channel.ID,
					"既に絵文字は申請されています。新たな申請を作成してください。\n",
				)
				return
			}

			emoji.IsSensitive = false
			emoji.State = 0

			deleteEmoji(emoji.FilePath)

			s.ChannelMessageSend(
				channel.ID,
				"リクエストを初期化します。"+
					":---\n"+
					"1. 絵文字の名前について教えてください 例: 絵文字では`:emoji-name:`となりますが、この時の`emoji-name`を入力してください \n",
			)

		},
	)
}

func newEmojiChannelComponent() {
	// new_emoji_channel
	addComponent(
		&discordgo.ApplicationCommand{
			Name: "new_emoji_channel",
		},
		func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			parent, err := s.Channel(i.ChannelID)

			if err != nil {
				returnFailedMessage(s, i, "Could not retrieve channel")
				return
			}

			overwrites := []*discordgo.PermissionOverwrite{
				{
					ID:   i.Member.User.ID,
					Type: discordgo.PermissionOverwriteTypeMember,
					Allow: discordgo.PermissionViewChannel |
						discordgo.PermissionSendMessages,
				},
				{
					ID:   ModeratorID,
					Type: discordgo.PermissionOverwriteTypeRole,
					Allow: discordgo.PermissionViewChannel |
						discordgo.PermissionSendMessages,
				},
				{
					ID:   BotID,
					Type: discordgo.PermissionOverwriteTypeRole,
					Allow: discordgo.PermissionViewChannel |
						discordgo.PermissionSendMessages,
				},
				{
					ID:   i.GuildID,
					Type: discordgo.PermissionOverwriteTypeRole,
					Deny: discordgo.PermissionViewChannel,
				},
			}

			emoji := newEmojiRequest(i.Member.User.ID)

			channel, err := s.GuildChannelCreateComplex(GuildID, discordgo.GuildChannelCreateData{
				Type:                 discordgo.ChannelTypeGuildText,
				Name:                 "Emoji-" + emoji.ID,
				ParentID:             parent.ParentID,
				PermissionOverwrites: overwrites,
			})

			if err != nil {
				returnFailedMessage(s, i, "Could not create emoji channel")
				emoji.abort()
				return
			}

			s.ChannelMessageSend(
				channel.ID,
				": 絵文字申請チャンネルへようこそ！\n"+
					":---\n"+
					" ここでは絵文字に関する各種登録を行います。表示されるメッセージに従って入力を行ってください！\n"+
					" 申請は絵文字Botが担当させていただきます。Botが一度非アクティブになると設定は初期化されますのでご注意ください！\n"+
					":---\n",
			)

			s.ChannelMessageSend(
				channel.ID,
				"1. 絵文字の名前について教えてください。 例: 絵文字では`:emoji-name:`となりますが、この時の`emoji-name`を入力してください。入力可能な文字は`小文字アルファベット`, `数字`, `_`です。 \n",
			)

			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Flags:   discordgo.MessageFlagsEphemeral,
					Content: "申請チャンネルを作成しました < #" + channel.Name + " >\n",
				},
			})
		},
	)

}
