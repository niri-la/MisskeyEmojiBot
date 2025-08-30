package command

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"

	"MisskeyEmojiBot/pkg/handler"
	"MisskeyEmojiBot/pkg/repository"
)

type EmojiListCommand interface {
}

type emojiListCommand struct {
	emojiRepository repository.EmojiRepository
}

func NewEmojiListCommand(emojiRepository repository.EmojiRepository) handler.CommandInterface {
	return &emojiListCommand{emojiRepository: emojiRepository}
}

func (c *emojiListCommand) GetCommand() *discordgo.ApplicationCommand {
	return &discordgo.ApplicationCommand{
		Name:        "emoji_list",
		Description: "絵文字申請一覧を表示します",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "id",
				Description: "特定の絵文字IDの詳細を表示",
				Required:    false,
			},
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "filter",
				Description: "表示する申請の状態",
				Required:    false,
				Choices: []*discordgo.ApplicationCommandOptionChoice{
					{
						Name:  "全て",
						Value: "all",
					},
					{
						Name:  "申請前",
						Value: "before_request",
					},
					{
						Name:  "審査中",
						Value: "pending",
					},
					{
						Name:  "承認済み",
						Value: "approved",
					},
					{
						Name:  "却下済み",
						Value: "rejected",
					},
				},
			},
		},
	}
}

func (c *emojiListCommand) Execute(s *discordgo.Session, i *discordgo.InteractionCreate) {
	var emojiID string
	filter := "all"

	// オプションを解析
	for _, option := range i.ApplicationCommandData().Options {
		switch option.Name {
		case "id":
			emojiID = option.StringValue()
		case "filter":
			filter = option.StringValue()
		}
	}

	// IDが指定されている場合は詳細表示
	if emojiID != "" {
		c.showEmojiDetail(s, i, emojiID)
		return
	}

	emojis := c.emojiRepository.GetEmojis()

	var requestedEmojis []string
	var pendingEmojis []string
	var approvedEmojis []string
	var rejectedEmojis []string

	for _, emoji := range emojis {
		statusText := fmt.Sprintf("**%s** - %s (申請日: %s)",
			emoji.Name,
			emoji.ID,
			emoji.StartAt.Format("2006/01/02 15:04"))

		if emoji.IsFinish {
			if emoji.IsAccepted {
				approvedEmojis = append(approvedEmojis, statusText+" ✅")
			} else {
				rejectedEmojis = append(rejectedEmojis, statusText+" ❌")
			}
		} else if emoji.IsRequested {
			pendingEmojis = append(pendingEmojis, statusText+fmt.Sprintf(" (🆗%d/🆖%d)", emoji.ApproveCount, emoji.DisapproveCount))
		} else {
			requestedEmojis = append(requestedEmojis, statusText+" ⏳")
		}
	}

	var content strings.Builder
	content.WriteString("# 絵文字申請状況\n\n")

	// フィルターに応じて表示内容を変更
	switch filter {
	case "before_request":
		if len(requestedEmojis) > 0 {
			content.WriteString("## 🔄 申請前 (" + strconv.Itoa(len(requestedEmojis)) + "件)\n")
			for _, emoji := range requestedEmojis {
				content.WriteString("- " + emoji + "\n")
			}
		} else {
			content.WriteString("申請前の絵文字はありません。")
		}
	case "pending":
		if len(pendingEmojis) > 0 {
			content.WriteString("## ⏳ 審査中 (" + strconv.Itoa(len(pendingEmojis)) + "件)\n")
			for _, emoji := range pendingEmojis {
				content.WriteString("- " + emoji + "\n")
			}
		} else {
			content.WriteString("審査中の絵文字はありません。")
		}
	case "approved":
		if len(approvedEmojis) > 0 {
			content.WriteString("## ✅ 承認済み (" + strconv.Itoa(len(approvedEmojis)) + "件)\n")
			for _, emoji := range approvedEmojis {
				content.WriteString("- " + emoji + "\n")
			}
		} else {
			content.WriteString("承認済みの絵文字はありません。")
		}
	case "rejected":
		if len(rejectedEmojis) > 0 {
			content.WriteString("## ❌ 却下済み (" + strconv.Itoa(len(rejectedEmojis)) + "件)\n")
			for _, emoji := range rejectedEmojis {
				content.WriteString("- " + emoji + "\n")
			}
		} else {
			content.WriteString("却下済みの絵文字はありません。")
		}
	default: // "all"
		if len(requestedEmojis) > 0 {
			content.WriteString("## 🔄 申請前 (" + strconv.Itoa(len(requestedEmojis)) + "件)\n")
			for _, emoji := range requestedEmojis {
				content.WriteString("- " + emoji + "\n")
			}
			content.WriteString("\n")
		}

		if len(pendingEmojis) > 0 {
			content.WriteString("## ⏳ 審査中 (" + strconv.Itoa(len(pendingEmojis)) + "件)\n")
			for _, emoji := range pendingEmojis {
				content.WriteString("- " + emoji + "\n")
			}
			content.WriteString("\n")
		}

		if len(approvedEmojis) > 0 {
			content.WriteString("## ✅ 承認済み (" + strconv.Itoa(len(approvedEmojis)) + "件)\n")
			for _, emoji := range approvedEmojis {
				content.WriteString("- " + emoji + "\n")
			}
			content.WriteString("\n")
		}

		if len(rejectedEmojis) > 0 {
			content.WriteString("## ❌ 却下済み (" + strconv.Itoa(len(rejectedEmojis)) + "件)\n")
			for _, emoji := range rejectedEmojis {
				content.WriteString("- " + emoji + "\n")
			}
			content.WriteString("\n")
		}

		if len(requestedEmojis)+len(pendingEmojis)+len(approvedEmojis)+len(rejectedEmojis) == 0 {
			content.WriteString("絵文字申請はありません。")
		}
	}

	// Discordのメッセージ長制限(2000文字)を考慮
	contentStr := content.String()
	if len(contentStr) > 1900 {
		contentStr = contentStr[:1900] + "\n\n... (文字数制限により省略)"
	}

	_ = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Flags:   discordgo.MessageFlagsEphemeral,
			Content: contentStr,
		},
	})
}

func (c *emojiListCommand) showEmojiDetail(s *discordgo.Session, i *discordgo.InteractionCreate, emojiID string) {
	emoji, err := c.emojiRepository.GetEmoji(emojiID)
	if err != nil || emoji == nil {
		_ = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Flags:   discordgo.MessageFlagsEphemeral,
				Content: fmt.Sprintf("❌ ID `%s` の絵文字が見つかりませんでした。", emojiID),
			},
		})
		return
	}

	var content strings.Builder
	content.WriteString(fmt.Sprintf("# ★絵文字詳細情報\n"))
	content.WriteString(fmt.Sprintf("**ID:** `%s`\n", emoji.ID))
	content.WriteString(fmt.Sprintf("**名前:** `%s`\n", emoji.Name))
	content.WriteString(fmt.Sprintf("**申請者:** <@%s>\n", emoji.RequestUser))
	content.WriteString(fmt.Sprintf("**申請日:** %s\n", emoji.StartAt.Format("2006年01月02日 15:04:05")))

	if emoji.Category != "" {
		content.WriteString(fmt.Sprintf("**カテゴリ:** %s\n", emoji.Category))
	}
	if emoji.Tag != "" {
		content.WriteString(fmt.Sprintf("**タグ:** %s\n", emoji.Tag))
	}
	if emoji.License != "" {
		content.WriteString(fmt.Sprintf("**ライセンス:** %s\n", emoji.License))
	}
	if emoji.Other != "" {
		content.WriteString(fmt.Sprintf("**その他:** %s\n", emoji.Other))
	}

	content.WriteString(fmt.Sprintf("**NSFW:** %s\n", map[bool]string{true: "はい", false: "いいえ"}[emoji.IsSensitive]))
	content.WriteString(fmt.Sprintf("**上書き:** %s\n", map[bool]string{true: "はい", false: "いいえ"}[emoji.IsOverwrite]))

	// ステータス情報
	content.WriteString("\n## 📊 ステータス情報\n")
	if emoji.IsFinish {
		if emoji.IsAccepted {
			content.WriteString("**状態:** ✅ 承認済み\n")
		} else {
			content.WriteString("**状態:** ❌ 却下済み\n")
		}
	} else if emoji.IsRequested {
		content.WriteString("**状態:** ⏳ 審査中\n")
		content.WriteString(fmt.Sprintf("**承認数:** %d\n", emoji.ApproveCount))
		content.WriteString(fmt.Sprintf("**却下数:** %d\n", emoji.DisapproveCount))
	} else {
		content.WriteString("**状態:** 🔄 申請前\n")
	}

	// テクニカル情報
	content.WriteString("\n## 🔧 テクニカル情報\n")
	if emoji.ChannelID != "" {
		content.WriteString(fmt.Sprintf("**チャンネル:** <#%s>\n", emoji.ChannelID))
	}
	if emoji.UserThreadID != "" {
		content.WriteString(fmt.Sprintf("**ユーザースレッド:** <#%s>\n", emoji.UserThreadID))
	}
	if emoji.ModerationMessageID != "" {
		content.WriteString(fmt.Sprintf("**モデレーションメッセージID:** `%s`\n", emoji.ModerationMessageID))
	}
	if emoji.FilePath != "" {
		content.WriteString(fmt.Sprintf("**ファイルパス:** `%s`\n", emoji.FilePath))
	}
	content.WriteString(fmt.Sprintf("**現在のステップ:** %d\n", emoji.NowStateIndex))
	content.WriteString(fmt.Sprintf("**レスポンスフラグ:** %t\n", emoji.ResponseFlag))

	contentStr := content.String()
	if len(contentStr) > 1900 {
		contentStr = contentStr[:1900] + "\n\n... (文字数制限により省略)"
	}

	_ = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Flags:   discordgo.MessageFlagsEphemeral,
			Content: contentStr,
		},
	})
}
