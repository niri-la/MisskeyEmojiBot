package job

import (
	"strconv"
	"time"

	"github.com/sirupsen/logrus"

	"MisskeyEmojiBot/pkg/entity"
	"MisskeyEmojiBot/pkg/repository"
)

type channelDeleteJob struct {
	emojiRepository repository.EmojiRepository
	discordRepo     repository.DiscordRepository
}

func NewChannelDeleteJob(emojiRepository repository.EmojiRepository, discordRepo repository.DiscordRepository) Job {
	return &channelDeleteJob{emojiRepository: emojiRepository, discordRepo: discordRepo}
}

func (j *channelDeleteJob) Run() {

	cleanRequest := time.NewTicker(12 * time.Hour)
	go func() {
		for {
			select {
			case <-cleanRequest.C:
				var targetEmoji []entity.Emoji
				for _, emoji := range j.emojiRepository.GetEmojis() {
					if time.Since(emoji.StartAt) > 48*time.Hour && !emoji.IsRequested {
						targetEmoji = append(targetEmoji, emoji)
					}
				}

				for _, emoji := range targetEmoji {
					j.emojiRepository.Abort(&emoji)
					j.discordRepo.DeleteChannel(emoji.ChannelID)
				}

				if len(targetEmoji) != 0 {
					logrus.Warn("delete emoji request : " + strconv.Itoa(len(targetEmoji)) + " emojis")
				}
			}
		}
	}()
}
