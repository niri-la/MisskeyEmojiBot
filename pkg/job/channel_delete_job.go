package job

import (
	"strconv"
	"time"

	"github.com/sirupsen/logrus"
)

type ChannelDeleteJob struct {
}

func (j *ChannelDeleteJob) Run() {

	cleanRequest := time.NewTicker(12 * time.Hour)
	go func() {
		for {
			select {
			case <-cleanRequest.C:
				var targetEmoji []Emoji
				for _, emoji := range emojiProcessList {
					if time.Since(emoji.StartAt) > 48*time.Hour && !emoji.IsRequested {
						targetEmoji = append(targetEmoji, emoji)
					}
				}

				for _, emoji := range targetEmoji {
					emoji.abort()
					deleteChannel(emoji)
				}

				if len(targetEmoji) != 0 {
					logrus.Warn("delete emoji request : " + strconv.Itoa(len(targetEmoji)) + " emojis")
				}
			}
		}
	}()
}
