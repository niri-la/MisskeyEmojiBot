package main

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
	"log"
	"os"
	"os/signal"
	"strconv"
)

// Bot parameters
var (
	GuildID               string
	BotToken              string
	AppID                 string
	ModeratorID           string
	BotID                 string
	ModerationChannelName string
	misskeyToken          string
	misskeyHost           string
	Session               *discordgo.Session
)

var moderationChannel *discordgo.Channel

func init() {
	loadEnvironments()
}

func init() {
	var err error
	Session, err = discordgo.New("Bot " + BotToken)
	if err != nil {
		log.Fatalf("Invalid bot parameters: %v", err)
	}

	moderationChannel, err = findChannelByName(Session, GuildID, ModerationChannelName)
}

func main() {
	log.Println("initializing...")
	// start
	Session.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		log.Println("Bot starting")
	})
	log.Println("Command register...")
	register()

	// コンポーネントはインタラクションの一部なので、InteractionCreateHandlerを登録します。
	Session.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		switch i.Type {
		case discordgo.InteractionApplicationCommand:
			if h, ok := CommandHandlers[i.ApplicationCommandData().Name]; ok {
				h(s, i)
			}
		case discordgo.InteractionMessageComponent:
			if h, ok := ComponentsHandlers[i.MessageComponentData().CustomID]; ok {
				h(s, i)
			}
		}
	})

	Session.AddHandler(func(s *discordgo.Session, m *discordgo.MessageCreate) {
		if m.Author.ID == s.State.User.ID {
			return
		}

		channel, _ := s.Channel(m.ChannelID)

		emoji, err := GetEmoji(channel.Name[6:])

		if err != nil {
			//s.ChannelMessageSend(m.ChannelID, "新たな申請のRequestに失敗しました。管理者に問い合わせを行ってください。")
			//fmt.Println("[ERROR] Reason : emoji not found")
			return
		}

		runEmojiProcess(emoji, s, m)

	})

	Session.AddHandler(emojiModerationReaction)

	_, err := Session.ApplicationCommandCreate(AppID, GuildID, &discordgo.ApplicationCommand{
		Name:        "buttons",
		Description: "Test the buttons if you got courage",
	})

	if err != nil {
		log.Fatalf("Cannot create slash command: %v", err)
	}

	err = Session.Open()
	if err != nil {
		log.Fatalf("Cannot open the session: %v", err)
	}
	defer Session.Close()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	<-stop
	log.Println("Graceful shutdown")
}

func findChannelByName(s *discordgo.Session, guildID string, name string) (*discordgo.Channel, error) {
	channels, err := s.GuildChannels(guildID)
	if err != nil {
		return nil, err
	}

	for _, ch := range channels {
		if ch.Name == name {
			return ch, nil
		}
	}

	return nil, fmt.Errorf("channel not found")
}

func register() {

	// ni_rilana
	addCommand(
		&discordgo.ApplicationCommand{
			Name:        "ni_rilana",
			Description: "Misskey Emoji Bot © 2023 KineL",
		},
		func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Flags: discordgo.MessageFlagsEphemeral,
					Content: "What is this ? \n" +
						": Misskey Emoji Bot\n" +
						": Created by ni_rila (KineL)\n" +
						": © 2023 KineL\n",
				},
			})
		},
	)

	// init
	addCommand(
		&discordgo.ApplicationCommand{
			Name:        "init",
			Description: "絵文字申請用の初期化を行います",
		},
		func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "こんにちは！絵文字申請の初期化を行います。\n" +
						"絵文字申請用のチャンネルを指定してください！",
					Flags: discordgo.MessageFlagsEphemeral,
					Components: []discordgo.MessageComponent{
						discordgo.ActionsRow{
							Components: []discordgo.MessageComponent{
								discordgo.SelectMenu{
									MenuType:     discordgo.ChannelSelectMenu,
									CustomID:     "init_channel",
									Placeholder:  "申請を行うチャンネルを選択してください",
									ChannelTypes: []discordgo.ChannelType{discordgo.ChannelTypeGuildText},
								},
							},
						},
					},
				},
			})
		},
	)

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

			returnFailedMessage(s, i, "Could not create emoji channel")
			return

		},
	)

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

			//s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			//	Type: discordgo.InteractionResponseChannelMessageWithSource,
			//	Data: &discordgo.InteractionResponseData{
			//		Flags:   discordgo.MessageFlagsEphemeral,
			//		Content: "😎",
			//	},
			//})

			emoji.IsRequested = true

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

	// new_emoji_channel
	addComponent(
		&discordgo.ApplicationCommand{
			Name: "new_emoji_channel",
		},
		func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			parent, err := s.Channel(i.ChannelID)

			emoji := newEmojiRequest()

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

			channel, err := s.GuildChannelCreateComplex(GuildID, discordgo.GuildChannelCreateData{
				Type:                 discordgo.ChannelTypeGuildText,
				Name:                 "Emoji-" + emoji.ID,
				ParentID:             parent.ParentID,
				PermissionOverwrites: overwrites,
			})

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

			if err != nil {
				returnFailedMessage(s, i, "Could not create emoji channel")
				return
			}

			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Flags:   discordgo.MessageFlagsEphemeral,
					Content: "申請チャンネルを作成しました < #" + channel.Name + " >\n",
				},
			})

		},
	)
	registeredCommands := make([]*discordgo.ApplicationCommand, len(Commands))
	for i, v := range Commands {
		cmd, err := Session.ApplicationCommandCreate(AppID, GuildID, v)
		if err != nil {
			Session.Close()
			panic(fmt.Sprintf("Cannot create '%v' command: %v", v.Name, err))
		}
		registeredCommands[i] = cmd
	}
}

func returnFailedMessage(s *discordgo.Session, i *discordgo.InteractionCreate, reason string) {
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Flags:   discordgo.MessageFlagsEphemeral,
			Content: "新たな申請のRequestに失敗しました。管理者に問い合わせを行ってください。",
		},
	})

	fmt.Println("[ERROR] Reason : " + reason)
	return
}

func loadEnvironments() {
	err := godotenv.Load("settings.env")

	if err != nil {
		panic(err)
	}

	GuildID = os.Getenv("guild_id")
	BotToken = os.Getenv("bot_token")
	AppID = os.Getenv("application_id")
	ModeratorID = os.Getenv("moderator_role_id")
	BotID = os.Getenv("bot_role_id")
	ModerationChannelName = os.Getenv("moderation_channel_name")
	misskeyToken = os.Getenv("misskey_token")
	misskeyHost = os.Getenv("misskey_host")

	fmt.Println(GuildID)
	fmt.Println(BotToken)
	fmt.Println(AppID)
	fmt.Println(BotID)
	fmt.Println(ModerationChannelName)
	fmt.Println(misskeyToken)
	fmt.Println(misskeyHost)

}
