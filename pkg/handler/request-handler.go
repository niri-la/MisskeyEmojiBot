package handler

import (
	"MisskeyEmojiBot/pkg/entity"

	"github.com/bwmarrin/discordgo"
	"github.com/sirupsen/logrus"
)

type Processor interface {
	Request(*entity.Emoji, *discordgo.Session, string) (entity.Response, error)
	Response(*entity.Emoji, *discordgo.Session, *discordgo.MessageCreate) (entity.Response, error)
}

type ResponseProcessor interface {
	Execute(*entity.Emoji, *discordgo.Session, *discordgo.MessageCreate) (entity.Response, error)
}

// type RequestProcessor func(*entity.Emoji, *discordgo.Session, string) Response
// type ResponceProcessor func(*entity.Emoji, *discordgo.Session, *discordgo.MessageCreate) Response

var workflow = map[int]string{
	0: "Default",
	2: "SetName",
	1: "Upload",
	3: "Category",
	4: "Tag",
	5: "License",
	6: "Other",
	7: "Nsfw",
	8: "Check",
}

var request = make(map[string]RequestProcessor)
var response = make(map[string]ResponseProcessor)
var reverseWorkflowMap = make(map[string]int)

func init() {
	for key, value := range workflow {
		reverseWorkflowMap[value] = key
	}
}

func init() {

}

func ProcessNextRequest(emoji *Emoji, s *discordgo.Session, id string) bool {
	requestIndex := reverseWorkflowMap[emoji.RequestState]
	logger.WithFields(logrus.Fields{
		"emoji id":       emoji.ID,
		"request index":  requestIndex,
		"response index": requestIndex,
	}).Trace("Emoji Processing (request)...")
	r1 := request[workflow[requestIndex+1]](emoji, s, id)
	return r1.IsSuccess
}

func Process(emoji *Emoji, s *discordgo.Session, m *discordgo.MessageCreate) bool {
	// 0. まずrequestを確認する(初期はRequest及びResponseは0である)
	// 1. 両者が等しい時はRequestを1進める
	// 2. RequestよりResponseが小さい場合はResponse待ちなのでResponseに値を渡す
	// 3. Responseが完了したらResponseを1すすめる。
	// 4. 1に戻る
	// 最終的に次の値がない場合は終了する。
	requestIndex := reverseWorkflowMap[emoji.RequestState]
	responseIndex := reverseWorkflowMap[emoji.ResponseState]

	logger.WithFields(logrus.Fields{
		"emoji id":       emoji.ID,
		"user":           m.Author.Username,
		"request index":  requestIndex,
		"response index": requestIndex,
	}).Trace("Emoji Processing...")

	if requestIndex == responseIndex {
		r1 := request[workflow[requestIndex+1]](emoji, s, m.ChannelID)
		return r1.IsSuccess
	}

	if requestIndex > responseIndex {
		r2 := response[workflow[responseIndex+1]](emoji, s, m)
		if r2.IsSuccess {
			Process(emoji, s, m)
		}
		return r2.IsSuccess
	}
	return false
}

func first(emoji *Emoji, s *discordgo.Session, id string) {
	request[workflow[1]](emoji, s, id)
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

	if emoji.DisapproveCount-1 >= roleCount || (isDebug && emoji.DisapproveCount-1 >= 1) {
		emoji.disapprove()
		s.ChannelMessageSend(m.ChannelID, "## 申請は却下されました")
		closeThread(m.ChannelID, emoji.ModerationMessageID)
		return
	}

	if emoji.ApproveCount-1 >= roleCount || (isDebug && emoji.ApproveCount-1 >= 1) {
		emoji.approve()
		s.ChannelMessageSend(m.ChannelID, "## 絵文字はアップロードされました")
		closeThread(m.ChannelID, emoji.ModerationMessageID)
		return
	}

}
