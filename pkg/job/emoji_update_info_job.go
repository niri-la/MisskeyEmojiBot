package job

import (
	"MisskeyEmojiBot/pkg/repository"
	"strings"
	"time"

	"github.com/yitsushi/go-misskey/models"
	"github.com/yitsushi/go-misskey/services/notes"
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
		for {
			select {
			case <-ticker.C:
				emojiArray := j.emojiRepository.EmojiReconstruction()
				if len(emojiArray) != 0 {
					var builder strings.Builder
					for _, emoji := range emojiArray {
						builder.WriteString(":" + emoji.Name + ":")
					}

					message := j.misskeyRepo.NewString("#にりらみすきー部 \n絵文字が追加されました\n" +
						builder.String())

					j.misskeyRepo.Note(notes.CreateRequest{
						Visibility: models.VisibilityPublic,
						Text:       message,
						LocalOnly:  true,
					})
				}
			}
		}
	}()
}
