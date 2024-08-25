package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	GuildID               string
	BotToken              string
	AppID                 string
	ModeratorID           string
	BotID                 string
	ModerationChannelName string
	MisskeyToken          string
	MisskeyHost           string
	IsDebug               bool
}

func LoadConfig() Config {
	err := godotenv.Load("settings.env")

	if err != nil {
		panic(err)
	}

	isDebug, err := strconv.ParseBool(os.Getenv("is_debug"))
	if err != nil {
		isDebug = false
	}

	config := Config{
		GuildID:               os.Getenv("guild_id"),
		BotToken:              os.Getenv("bot_token"),
		AppID:                 os.Getenv("application_id"),
		ModeratorID:           os.Getenv("moderator_role_id"),
		BotID:                 os.Getenv("bot_role_id"),
		ModerationChannelName: os.Getenv("moderation_channel_name"),
		MisskeyToken:          os.Getenv("misskey_token"),
		MisskeyHost:           os.Getenv("misskey_host"),
		IsDebug:               isDebug,
	}

	// 全ての設定を読み込んだら、設定を返す
	return config
}
