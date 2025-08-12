module MisskeyEmojiBot

go 1.23.0

toolchain go1.24.6

require (
	github.com/bwmarrin/discordgo v0.27.1
	github.com/google/uuid v1.6.0
	github.com/joho/godotenv v1.5.1
	github.com/sirupsen/logrus v1.9.3
	github.com/yitsushi/go-misskey v1.1.6
	golang.org/x/text v0.11.0
	gopkg.in/yaml.v3 v3.0.1
)

require (
	github.com/dustin/go-humanize v1.0.1 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/ncruces/go-strftime v0.1.9 // indirect
	github.com/remyoudompheng/bigfft v0.0.0-20230129092748-24d4a6f8daec // indirect
	golang.org/x/exp v0.0.0-20250620022241-b7579e27df2b // indirect
	modernc.org/libc v1.66.3 // indirect
	modernc.org/mathutil v1.7.1 // indirect
	modernc.org/memory v1.11.0 // indirect
	modernc.org/sqlite v1.38.2 // indirect
)

require (
	github.com/gorilla/websocket v1.4.2 // indirect
	github.com/kr/text v0.2.0 // indirect
	github.com/pkg/errors v0.9.1
	golang.org/x/crypto v0.0.0-20210421170649-83a5a9bb288b // indirect
	golang.org/x/net v0.7.0 // indirect
	golang.org/x/sys v0.34.0 // indirect
	gopkg.in/yaml.v2 v2.4.0
)

replace github.com/yitsushi/go-misskey => github.com/niwaniwa/go-misskey v0.0.0-20230710181204-1210df04cd80
