package job

import "time"

type EmojiUpdateInfoJob struct {
}

func (j *EmojiUpdateInfoJob) Run() {
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
