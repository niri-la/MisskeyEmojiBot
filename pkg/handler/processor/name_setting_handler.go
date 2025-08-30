package processor

import (
	"regexp"
	"strings"

	"github.com/bwmarrin/discordgo"

	"MisskeyEmojiBot/pkg/entity"
	"MisskeyEmojiBot/pkg/handler"
	"MisskeyEmojiBot/pkg/repository"
)

type nameSettingHandler struct {
	misskeyRepo repository.MisskeyRepository
}

func NewNameSettingHandler() handler.EmojiProcessHandler {
	return &nameSettingHandler{}
}

func NewNameSettingHandlerWithMisskey(misskeyRepo repository.MisskeyRepository) handler.EmojiProcessHandler {
	return &nameSettingHandler{misskeyRepo: misskeyRepo}
}

func (h *nameSettingHandler) Request(emoji *entity.Emoji, s *discordgo.Session, cID string) (entity.Response, error) {

	response := entity.Response{
		IsSuccess: true,
	}

	_, err := s.ChannelMessageSend(
		cID,
		"## 絵文字の名前を入力してください。\n実際にMisskey上で絵文字を入力する際は`:emoji-name:`としますが、この`emoji-name`の部分を入力してください。\n入力可能な文字は`小文字アルファベット`, `数字`, `_`です。",
	)
	if err != nil {
		return entity.Response{}, err
	}

	return response, nil
}

func (h *nameSettingHandler) Response(emoji *entity.Emoji, s *discordgo.Session, m *discordgo.MessageCreate) (entity.Response, error) {
	response := entity.Response{
		IsSuccess: false,
	}

	input := strings.TrimPrefix(m.Content, ":")
	input = strings.TrimSuffix(input, ":")

	if len(input) <= 1 {
		_, err := s.ChannelMessageSend(m.ChannelID, ":2文字以上入力してください。")
		if err != nil {

			return response, err
		}
		return response, nil
	}
	reg := regexp.MustCompile(`[^a-z0-9_]+`)
	input = reg.ReplaceAllStringFunc(strings.ToLower(input), func(s string) string {
		return "_"
	})
	_, _ = s.ChannelMessageSend(m.ChannelID, ":: 入力されたメッセージ\n [ `"+input+"` ]")
	
	// Issue #72: 既存絵文字の重複チェック
	if h.misskeyRepo != nil {
		exists, existingEmoji, err := h.misskeyRepo.CheckEmojiExists(input)
		if err != nil {
			_, _ = s.ChannelMessageSend(m.ChannelID, "⚠️ 絵文字の重複チェックに失敗しました。処理を続行します。")
		} else if exists {
			// 重複している場合、上書き確認UIを表示
			components := []discordgo.MessageComponent{
				discordgo.ActionsRow{
					Components: []discordgo.MessageComponent{
						discordgo.Button{
							Label:    "はい、上書きします",
							Style:    discordgo.PrimaryButton,
							CustomID: "emoji_overwrite_confirm:yes:" + emoji.ID,
							Emoji: discordgo.ComponentEmoji{
								Name: "✅",
							},
						},
						discordgo.Button{
							Label:    "いいえ、別の名前にします", 
							Style:    discordgo.SecondaryButton,
							CustomID: "emoji_overwrite_confirm:no:" + emoji.ID,
							Emoji: discordgo.ComponentEmoji{
								Name: "❌",
							},
						},
					},
				},
			}
			
			overwriteMsg := "⚠️ **絵文字 `" + input + "` は既に存在します**\n\n"
			if existingEmoji != nil {
				overwriteMsg += "既存の絵文字情報:\n"
				// TODO: 既存絵文字の詳細情報を表示（カテゴリ、タグなど）
			}
			overwriteMsg += "\n上書きしますか？"
			
			_, err = s.ChannelMessageSendComplex(m.ChannelID, &discordgo.MessageSend{
				Content:    overwriteMsg,
				Components: components,
			})
			if err != nil {
				_, _ = s.ChannelMessageSend(m.ChannelID, "上書き確認の表示に失敗しました。")
				return response, nil
			}
			
			// 上書き確認待ちの状態にするため、まだ成功にしない
			emoji.Name = input // 名前は保存しておく
			response.IsSuccess = false
			return response, nil
		}
	}
	
	_, _ = s.ChannelMessageSend(m.ChannelID, "# ----------\n")
	emoji.Name = input
	response.IsSuccess = true
	return response, nil
}
