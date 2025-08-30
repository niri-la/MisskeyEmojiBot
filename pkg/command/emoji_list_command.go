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
	// フィルター引数を取得
	filter := "all"
	if len(i.ApplicationCommandData().Options) > 0 {
		filter = i.ApplicationCommandData().Options[0].StringValue()
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
