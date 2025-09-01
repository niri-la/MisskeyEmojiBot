package job

import (
	"strings"
	"time"

	"github.com/yitsushi/go-misskey/models"
	"github.com/yitsushi/go-misskey/services/notes"

	"MisskeyEmojiBot/pkg/repository"
)

type emojiUpdateInfoJob struct {
	emojiRepository repository.EmojiRepository
	misskeyRepo     repository.MisskeyRepository
}

func NewEmojiUpdateInfoJob(emojiRepository repository.EmojiRepository, misskeyRepo repository.MisskeyRepository) Job {
	return &emojiUpdateInfoJob{emojiRepository: emojiRepository, misskeyRepo: misskeyRepo}
}

func (j *emojiUpdateInfoJob) Run() {
	ticker := time.NewTicker(12 * time.Hour)
	go func() {
		for range ticker.C {
			emojiArray := j.emojiRepository.EmojiReconstruction()
			if len(emojiArray) != 0 {
				var builder strings.Builder
				for _, emoji := range emojiArray {
					if emoji.IsAccepted {
						builder.WriteString(":" + emoji.Name + ":")
					}
				}

				if builder.Len() > 0 {
					message := j.misskeyRepo.NewString("#にりらみすきー部 \n絵文字が追加されました\n" +
						builder.String())

					_ = j.misskeyRepo.Note(notes.CreateRequest{
						Visibility: models.VisibilityPublic,
						Text:       message,
						LocalOnly:  true,
					})

					// Mark emojis as notified
					for _, emoji := range emojiArray {
						if emoji.IsAccepted {
							emoji.IsNotified = true
							_ = j.emojiRepository.Save(&emoji)
						}
					}
				}
			}
		}
	}()
}
