package component

import (
	"github.com/bwmarrin/discordgo"

	"MisskeyEmojiBot/pkg/config"
	"MisskeyEmojiBot/pkg/handler"
	"MisskeyEmojiBot/pkg/repository"
)

type emojiOverwriteConfirmComponent struct {
	config          config.Config
	emojiRepository repository.EmojiRepository
	discordRepo     repository.DiscordRepository
}

func NewEmojiOverwriteConfirmComponent(cfg config.Config, emojiRepository repository.EmojiRepository, discordRepo repository.DiscordRepository) handler.Component {
	return &emojiOverwriteConfirmComponent{
		config:          cfg,
		emojiRepository: emojiRepository,
		discordRepo:     discordRepo,
	}
}

func (c *emojiOverwriteConfirmComponent) GetCommand() *discordgo.ApplicationCommand {
	return &discordgo.ApplicationCommand{
		Name: "emoji_overwrite_confirm",
	}
}

func (c *emojiOverwriteConfirmComponent) Execute(s *discordgo.Session, i *discordgo.InteractionCreate) {
	customID := i.MessageComponentData().CustomID
	
	// カスタムIDの形式: "emoji_overwrite_confirm:yes:emojiID" または "emoji_overwrite_confirm:no:emojiID"
	// この簡単な実装では、customIDから決定を抽出
	
	if len(customID) > 26 { // "emoji_overwrite_confirm:".length = 26
		decision := customID[26:29] // "yes" または "no:"
		emojiID := customID[30:]    // emoji ID部分
		
		emoji, err := c.emojiRepository.GetEmoji(emojiID)
		if err != nil {
			_ = c.discordRepo.ReturnFailedMessage(i, "絵文字データの取得に失敗しました")
			return
		}
		
		if decision == "yes" {
			// 上書きを承認
			emoji.IsOverwrite = true
			_ = c.emojiRepository.Save(emoji)
			
			_ = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseUpdateMessage,
				Data: &discordgo.InteractionResponseData{
					Content: "✅ 絵文字の上書きが承認されました。申請を続行します。",
					Components: []discordgo.MessageComponent{}, // ボタンを削除
				},
			})
			
			// 次のステップ（カテゴリ設定など）に進むメッセージを送信
			_, _ = s.ChannelMessageSend(emoji.ChannelID, "## カテゴリを設定してください。\n既存の絵文字を上書きします。")
			
		} else if decision == "no:" {
			// 上書きを拒否、名前入力からやり直し
			_ = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseUpdateMessage,
				Data: &discordgo.InteractionResponseData{
					Content: "❌ 上書きがキャンセルされました。別の絵文字名を入力してください。",
					Components: []discordgo.MessageComponent{}, // ボタンを削除
				},
			})
			
			// 名前設定に戻る
			_, _ = s.ChannelMessageSend(emoji.ChannelID, "## 絵文字名を設定してください。\n別の名前を入力してください。")
			// NowStateIndexを名前設定に戻す（通常は1）
			emoji.NowStateIndex = 1
			_ = c.emojiRepository.Save(emoji)
		}
	}
}