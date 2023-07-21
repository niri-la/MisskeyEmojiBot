package main

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
)

var (
	Commands        = make([]*discordgo.ApplicationCommand, 0)
	CommandHandlers = make(map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate))
)

func addCommand(command *discordgo.ApplicationCommand, fn func(s *discordgo.Session, i *discordgo.InteractionCreate)) {
	_, exist := CommandHandlers[command.Name]
	if exist {
		panic(fmt.Sprintf("Error: [%s] is already existed.", command.Name))
	}
	CommandHandlers[command.Name] = fn
	Commands = append(Commands, command)
}

func command() {
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
					Content: "::: Misskey Emoji Bot \n" +
						": Created by ni_rila (KineL)\n" +
						": © 2023 KineL\n" +
						":::::::::: \n",
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

			if !hasPermission(*i.Member.User) {
				returnFailedMessage(s, i, "No permission.")
				return
			}

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
}
