package component

import (
	"MisskeyEmojiBot/pkg/config"
	"MisskeyEmojiBot/pkg/handler"
	"MisskeyEmojiBot/pkg/repository"
	"os"
	"strconv"

	"github.com/bwmarrin/discordgo"
)

type EmojiRequestComponen interface {
}

type emojiRequestComponen struct {
	config          config.Config
	emojiRepository repository.EmojiRepository
	discordRepo     repository.DiscordRepository
}

func NewEmojiRequestComponen(config config.Config, emojiRepository repository.EmojiRepository, discordRepo repository.DiscordRepository) handler.Component {
	return &emojiRequestComponen{config: config, emojiRepository: emojiRepository, discordRepo: discordRepo}
}
func (c *emojiRequestComponen) GetCommand() *discordgo.ApplicationCommand {
	return &discordgo.ApplicationCommand{
		Name: "emoji_request",
	}
}

func (c *emojiRequestComponen) Execute(s *discordgo.Session, i *discordgo.InteractionCreate) {
	channel, _ := s.Channel(i.ChannelID)
	emoji, err := c.emojiRepository.GetEmoji(channel.Name[6:])
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
		"## 申請をしました！\n"+
			"申請結果については追ってDMでご連絡いたします。\n"+
			"なお、申請結果について疑問がございましたら管理者へお問い合わせください！\n"+
			"この度は申請いただき大変ありがとうございました。\n",
	)

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Flags:   discordgo.MessageFlagsEphemeral,
			Content: "📨",
		},
	})

	emoji.IsRequested = true

	c.discordRepo.SendDirectMessage(*&emoji.RequestUser, "--- 申請内容 "+emoji.ID+"---\n名前: "+emoji.Name+"\nCategory: "+
		emoji.Category+"\n"+"tag:"+emoji.Tag+"\n"+"License:"+emoji.License+"\n"+"isNSFW:"+strconv.FormatBool(emoji.IsSensitive)+"\n"+
		"備考: "+emoji.Other+"\nURL: https://discordapp.com/channels/"+c.config.GuildID+"/"+emoji.ChannelID+"\n---")

	moderationChannel, err := c.discordRepo.FindChannelByName(c.config.GuildID, "emoji-moderation")
	if err != nil {
		return
	}
	send, err := s.ChannelMessageSend(moderationChannel.ID, "## 申請 "+emoji.ID+"\n- 申請者: "+i.Member.User.Username+"\n"+"- 絵文字名: "+emoji.Name)
	if err != nil {
		return
	}

	emoji.ModerationMessageID = send.ID

	thread, err := s.MessageThreadStartComplex(moderationChannel.ID, send.ID, &discordgo.ThreadStart{
		Name:                emoji.ID,
		AutoArchiveDuration: 60,
		Invitable:           false,
	})

	s.ChannelMessageSend(thread.ID, "## 申請内容\n")
	s.ChannelMessageSend(thread.ID,
		"- Name    : **"+emoji.Name+"**\n"+
			"- Category: **"+emoji.Category+"**\n"+
			"- Tag     : **"+emoji.Tag+"**\n"+
			"- License : **"+emoji.License+"**\n"+
			"- Other   : **"+emoji.Other+"**\n"+
			"- NSFW    : **"+strconv.FormatBool(emoji.IsSensitive)+"**\n"+
			"## 絵文字画像")

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
}
