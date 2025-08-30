package command

import (
	"strconv"

	"github.com/bwmarrin/discordgo"

	"MisskeyEmojiBot/pkg/config"
	"MisskeyEmojiBot/pkg/handler"
	"MisskeyEmojiBot/pkg/repository"
)

type EmojiDetailChangeCommand interface {
}

type emojiDetailChangeCommand struct {
	config          config.Config
	emojiRepository repository.EmojiRepository
	discordRepo     repository.DiscordRepository
}

func NewEmojiDetailChangeCommand(config config.Config, emojiRepository repository.EmojiRepository, discordRepo repository.DiscordRepository) handler.CommandInterface {
	return &emojiDetailChangeCommand{config: config, emojiRepository: emojiRepository, discordRepo: discordRepo}
}

func (c *emojiDetailChangeCommand) GetCommand() *discordgo.ApplicationCommand {
	return &discordgo.ApplicationCommand{
		Name:        "change_emoji_detail",
		Description: "絵文字申請のプロパティを変更します",
		Type:        discordgo.ChatApplicationCommand,
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "property",
				Description: "Property to change (name, category, tag, license, other)",
				Required:    true,
			},
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "value",
				Description: "New value for the property",
				Required:    true,
			},
		},
	}
}

func (c *emojiDetailChangeCommand) Execute(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if !c.discordRepo.HasRole(c.config.GuildID, *i.Member.User, c.config.ModeratorID) {
		_ = c.discordRepo.ReturnFailedMessage(i, "No permission.")
		return
	}

	channel, _ := s.Channel(i.ChannelID)
	emoji, err := c.emojiRepository.GetEmoji(channel.Name)
	if err != nil {
		_ = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Flags:   discordgo.MessageFlagsEphemeral,
				Content: "設定に失敗しました。管理者に問い合わせを行ってください。\n",
			},
		})
	}

	options := i.ApplicationCommandData().Options
	if len(options) < 2 {
		_ = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Flags:   discordgo.MessageFlagsEphemeral,
				Content: "引数が足りません / Not enough arguments.",
			},
		})
		return
	}

	property := options[0].StringValue()
	value := options[1].StringValue()

	switch property {
	case "name":
		emoji.Name = value
	case "category":
		emoji.Category = value
	case "tag":
		emoji.Tag = value
	case "license":
		emoji.License = value
	case "other":
		emoji.Other = value
	default:
		_ = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Flags:   discordgo.MessageFlagsEphemeral,
				Content: "プロパティが見つかりません / Property not found.",
			},
		})
		return
	}

	// Save changes to database
	err = c.emojiRepository.Save(emoji)
	if err != nil {
		_ = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Flags:   discordgo.MessageFlagsEphemeral,
				Content: "保存に失敗しました / Failed to save changes.",
			},
		})
		return
	}

	_ = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Flags:   discordgo.MessageFlagsEphemeral,
			Content: "プロパティを変更しました / Changed the property.",
		},
	})

	_, _ = s.ChannelMessageSend(channel.ID, "# 変更後の絵文字\n")
	_, _ = s.ChannelMessageSend(channel.ID,
		"- Name    : **"+emoji.Name+"**\n"+
			"- Category: **"+emoji.Category+"**\n"+
			"- Tag     : **"+emoji.Tag+"**\n"+
			"- License : **"+emoji.License+"**\n"+
			"- Other   : **"+emoji.Other+"**\n"+
			"- NSFW    : **"+strconv.FormatBool(emoji.IsSensitive)+"**\n")

}
