module MisskeyEmojiBot

go 1.23.0

toolchain go1.23.10

require (
	github.com/bwmarrin/discordgo v0.27.1
	github.com/glebarez/sqlite v1.11.0
	github.com/google/uuid v1.3.0
	github.com/joho/godotenv v1.5.1
	github.com/sirupsen/logrus v1.9.3
	github.com/yitsushi/go-misskey v1.1.6
	golang.org/x/text v0.28.0
	gorm.io/gorm v1.30.2
)

require (
	github.com/dustin/go-humanize v1.0.1 // indirect
	github.com/glebarez/go-sqlite v1.21.2 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	github.com/mattn/go-isatty v0.0.17 // indirect
	github.com/remyoudompheng/bigfft v0.0.0-20230129092748-24d4a6f8daec // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
	modernc.org/libc v1.22.5 // indirect
	modernc.org/mathutil v1.5.0 // indirect
	modernc.org/memory v1.5.0 // indirect
	modernc.org/sqlite v1.23.1 // indirect
)

require (
	github.com/gorilla/websocket v1.4.2 // indirect
	github.com/kr/text v0.2.0 // indirect
	github.com/pkg/errors v0.9.1
	golang.org/x/crypto v0.0.0-20210421170649-83a5a9bb288b // indirect
	golang.org/x/net v0.7.0 // indirect
	golang.org/x/sys v0.7.0 // indirect
	gopkg.in/yaml.v2 v2.4.0
)

replace github.com/yitsushi/go-misskey => github.com/niwaniwa/go-misskey v0.0.0-20230710181204-1210df04cd80
