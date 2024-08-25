package job

import (
	"MisskeyEmojiBot/pkg/entity"
	"MisskeyEmojiBot/pkg/repository"
	"strconv"
	"time"

	"github.com/sirupsen/logrus"
)

type channelDeleteJob struct {
	discordRepo repository.DiscordRepository
}

func NewChannelDeleteJob(discordRepo repository.DiscordRepository) Job {
	return &channelDeleteJob{discordRepo: discordRepo}
}

func (j *channelDeleteJob) Run() {

	cleanRequest := time.NewTicker(12 * time.Hour)
	go func() {
		for {
			select {
			case <-cleanRequest.C:
				var targetEmoji []entity.Emoji
				for _, emoji := range emojiProcessList {
					if time.Since(emoji.StartAt) > 48*time.Hour && !emoji.IsRequested {
						targetEmoji = append(targetEmoji, emoji)
					}
				}

				for _, emoji := range targetEmoji {
					emoji.abort()
					j.discordRepo.DeleteChannel(emoji)
				}

				if len(targetEmoji) != 0 {
					logrus.Warn("delete emoji request : " + strconv.Itoa(len(targetEmoji)) + " emojis")
				}
			}
		}
	}()
}
