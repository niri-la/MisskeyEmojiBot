module MisskeyEmojiBot

go 1.20

require (
	github.com/bwmarrin/discordgo v0.27.1
	github.com/google/uuid v1.3.0
	github.com/joho/godotenv v1.5.1
	github.com/sirupsen/logrus v1.9.3
	github.com/yitsushi/go-misskey v1.1.6
	golang.org/x/text v0.11.0
	gopkg.in/yaml.v3 v3.0.1
)

require (
	github.com/gorilla/websocket v1.4.2 // indirect
	github.com/kr/text v0.2.0 // indirect
	github.com/pkg/errors v0.9.1
	golang.org/x/crypto v0.0.0-20210421170649-83a5a9bb288b // indirect
	golang.org/x/net v0.7.0 // indirect
	golang.org/x/sys v0.5.0 // indirect
)

replace github.com/yitsushi/go-misskey => github.com/niwaniwa/go-misskey v0.0.0-20230710181204-1210df04cd80
