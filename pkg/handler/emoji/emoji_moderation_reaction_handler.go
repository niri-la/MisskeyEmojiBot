package emoji

import (
	"MisskeyEmojiBot/pkg/config"
	"MisskeyEmojiBot/pkg/entity"
	"MisskeyEmojiBot/pkg/repository"
	"fmt"

	"github.com/bwmarrin/discordgo"
)

type EmojiModerationReactionHandler interface {
	HandleEmojiModerationReaction(s *discordgo.Session, m *discordgo.MessageReactionAdd)
}

type emojiModerationReactionHandler struct {
	emojiHandler      EmojiHandler
	emojiRepository   repository.EmojiRepository
	discordRepository repository.DiscordRepository
	config            config.Config
}

func NewEmojiModerationReactionHandler(emojiHandler EmojiHandler, emojiRepository repository.EmojiRepository, discordRepository repository.DiscordRepository, config config.Config) EmojiModerationReactionHandler {
	return &emojiModerationReactionHandler{emojiHandler: emojiHandler, emojiRepository: emojiRepository, discordRepository: discordRepository, config: config}
}

func (h *emojiModerationReactionHandler) HandleEmojiModerationReaction(s *discordgo.Session, m *discordgo.MessageReactionAdd) {
	if m.UserID == s.State.User.ID {
		return
	}

	channel, _ := s.Channel(m.ChannelID)
	var emoji *entity.Emoji
	found := false

	for _, e := range h.emojiRepository.GetEmojis() {
		if channel.Name == e.ID {
			emoji = &e
			found = true
			break
		}
	}

	if !found {
		return
	}

	emoji, err := h.emojiRepository.GetEmoji(emoji.ID)

	if err != nil {
		return
	}

	if emoji.IsFinish {
		return
	}

	roleCount, err := h.discordRepository.CountMembersWithSpecificRole(h.config.GuildID, h.config.ModeratorID)
	if err != nil {
		return
	}

	msg, err := s.ChannelMessage(channel.ID, m.MessageID)
	if err != nil {
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

	if emoji.DisapproveCount-1 >= roleCount || (h.config.IsDebug && emoji.DisapproveCount-1 >= 1) {
		err := h.emojiHandler.Disapprove(emoji)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
		}
		s.ChannelMessageSend(m.ChannelID, "## 申請は却下されました")
		err = h.discordRepository.CloseThread(m.ChannelID, emoji.ModerationMessageID)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
		}
		return
	}

	if emoji.ApproveCount-1 >= roleCount || (h.config.IsDebug && emoji.ApproveCount-1 >= 1) {
		err := h.emojiHandler.Approve(emoji)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
		}
		s.ChannelMessageSend(m.ChannelID, "## 絵文字はアップロードされました")
		h.discordRepository.CloseThread(m.ChannelID, emoji.ModerationMessageID)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
		}
		return
	}
}
