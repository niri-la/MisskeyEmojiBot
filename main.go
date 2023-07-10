package main

import (
	"flag"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"log"
	"os"
	"os/signal"
)

// Bot parameters
var (
	GuildID      = flag.String("guild", "", "Test guild ID")
	BotToken     = flag.String("token", "", "Bot access token")
	AppID        = flag.String("app", "", "Application ID")
	ModerationID = flag.String("moderation", "", "Moderation ID")
)

var s *discordgo.Session

func init() { flag.Parse() }

func init() {
	var err error
	s, err = discordgo.New("Bot " + *BotToken)
	if err != nil {
		log.Fatalf("Invalid bot parameters: %v", err)
	}
}

func main() {
	log.Println("initializing...")
	// start
	s.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		log.Println("Bot starting")
	})
	log.Println("Command register...")
	register()

	// コンポーネントはインタラクションの一部なので、InteractionCreateHandlerを登録します。
	s.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
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

	s.AddHandler(func(s *discordgo.Session, m *discordgo.MessageCreate) {
		if m.Author.ID == s.State.User.ID {
			return
		}

		channel, _ := s.Channel(m.ChannelID)

		emoji, err := GetEmoji(channel.Name[6:])

		if err != nil {
			s.ChannelMessageSend(m.ChannelID, "新たな申請のRequestに失敗しました。管理者に問い合わせを行ってください。")
			fmt.Println("[ERROR] Reason : emoji not found")
			return
		}

		RunEmojiProcess(emoji, s, m)

	})

	_, err := s.ApplicationCommandCreate(*AppID, *GuildID, &discordgo.ApplicationCommand{
		Name:        "buttons",
		Description: "Test the buttons if you got courage",
	})

	if err != nil {
		log.Fatalf("Cannot create slash command: %v", err)
	}

	err = s.Open()
	if err != nil {
		log.Fatalf("Cannot open the session: %v", err)
	}
	defer s.Close()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	<-stop
	log.Println("Graceful shutdown")
}

func register() {

	addCommand(
		&discordgo.ApplicationCommand{
			Name:        "ni_rilana",
			Description: "Misskey Emoji Bot © 2023 KineL",
		},
		func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "What is this ? \n" +
						": Misskey Emoji Bot\n" +
						": Created by ni_rila (KineL)\n" +
						": © 2023 KineL\n",
				},
			})
		},
	)

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
		},
	)

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
			emoji.NSFW = true
			emoji.State = 4
			emojiLastConfirmation(emoji, s, i.ChannelID)
		},
	)
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

			emoji.NSFW = false
			emoji.State = 4
			emojiLastConfirmation(emoji, s, i.ChannelID)

		},
	)

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
					ID:   *ModerationID,
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

			channel, err := s.GuildChannelCreateComplex(*GuildID, discordgo.GuildChannelCreateData{
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
				"1. まず絵文字の名前について教えてください 例: 絵文字では`:emoji-name:`となりますが、この時の`emoji-name`を入力してください \n",
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
		cmd, err := s.ApplicationCommandCreate(*AppID, *GuildID, v)
		if err != nil {
			s.Close()
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
