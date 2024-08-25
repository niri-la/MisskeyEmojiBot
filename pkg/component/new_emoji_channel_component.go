package component

import (
	"MisskeyEmojiBot/pkg/handler"
	"MisskeyEmojiBot/pkg/repository"
	"fmt"

	"github.com/bwmarrin/discordgo"
)

type CreateEmojiChannelComponen interface {
}

type createEmojiChannelComponen struct {
	emojiRequestHandler handler.EmojiRequestHandler
	emojiReposiotry     repository.EmojiRepository
	discordRepo         repository.DiscordRepository
}

func NewCreateEmojiChannelComponen(emojiRequestHandler handler.EmojiRequestHandler, emojiReposiotry repository.EmojiRepository, discordRepo repository.DiscordRepository) handler.Component {
	return &createEmojiChannelComponen{emojiRequestHandler: emojiRequestHandler, emojiReposiotry: emojiReposiotry, discordRepo: discordRepo}
}

func (c *createEmojiChannelComponen) GetCommand() *discordgo.ApplicationCommand {
	return &discordgo.ApplicationCommand{
		Name: "new_emoji_channel",
	}
}

func (c *createEmojiChannelComponen) Execute(s *discordgo.Session, i *discordgo.InteractionCreate) {
	emoji := c.emojiReposiotry.NewEmoji(i.Member.User.ID)

	channel, err := s.ThreadStartComplex(i.ChannelID, &discordgo.ThreadStart{
		Name:                "Emoji-" + emoji.ID,
		AutoArchiveDuration: 1440,
		Invitable:           false,
		Type:                discordgo.ChannelTypeGuildPrivateThread,
	})

	if err != nil {
		c.discordRepo.ReturnFailedMessage(i, fmt.Sprintf("Could not create emoji channel: %v", err))
		c.emojiReposiotry.Abort(emoji)
		return
	}

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Flags:   discordgo.MessageFlagsEphemeral,
			Content: "## 申請チャンネルを作成しました\nチャンネル: https://discordapp.com/channels/" + i.GuildID + "/" + channel.ID + "\n---",
		},
	})

	user, err := s.User(emoji.RequestUser)
	if err != nil {
		c.discordRepo.ReturnFailedMessage(i, fmt.Sprintf("Could not find user: %v", err))
		return
	}

	s.ChannelMessageSend(
		channel.ID,
		"# 絵文字申請チャンネルへようこそ！\n"+user.Mention()+"\n"+
			" ここでは絵文字に関する各種登録を行います。表示されるメッセージに従って入力を行ってください！\n"+
			" 申請は絵文字Botが担当させていただきます。Botが一度非アクティブになると設定は初期化されますのでご注意ください！\n",
	)

	s.ChannelMessageSendComplex(channel.ID,
		&discordgo.MessageSend{
			Content: "## 申請のキャンセル\n申請をキャンセルする場合は以下のボタンを押してください。\n申請後はキャンセルできませんのでご注意ください。\n",
			Components: []discordgo.MessageComponent{
				discordgo.ActionsRow{
					Components: []discordgo.MessageComponent{
						&discordgo.Button{
							Label:    "申請をキャンセルする / Cancel Request",
							CustomID: "cancel_request",
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

	emoji, _ = c.emojiReposiotry.GetEmoji(emoji.ID)
	emoji.ChannelID = channel.ID
	c.emojiRequestHandler.ResetState(emoji, s)
	c.emojiRequestHandler.ProcessRequest(emoji, s, channel.ID)
}
