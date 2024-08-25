package job

import (
	"MisskeyEmojiBot/pkg/repository"
	"time"
)

type emojiUpdateInfoJob struct {
	misskeyRepo repository.MisskeyRepository
}

func NewEmojiUpdateInfoJob(misskeyRepo repository.MisskeyRepository) Job {
	return &emojiUpdateInfoJob{misskeyRepo: misskeyRepo}
}

func (j *emojiUpdateInfoJob) Run() {
	ticker := time.NewTicker(12 * time.Hour)
	go func() {
		for {
			select {
			case <-ticker.C:
				emoji := emojiReconstruction()
				if len(emoji) != 0 {
					noteEmojiAdded(emoji)
				}
			}
		}
	}()
}
