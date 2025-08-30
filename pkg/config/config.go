package config

import (
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"

	"MisskeyEmojiBot/pkg/errors"
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
	SavePath              string
	DatabasePath          string
	IsDebug               bool

	// Migration Configuration
	EnableJsonMigration bool
	JsonMigrationPath   string
}

func LoadConfig() (*Config, error) {
	err := godotenv.Load("settings.env")
	if err != nil {
		return nil, errors.Config("failed to load settings.env", err)
	}

	isDebug, err := strconv.ParseBool(os.Getenv("debug"))
	if err != nil {
		isDebug = false
	}

	enableJsonMigration, err := strconv.ParseBool(os.Getenv("enable_json_migration"))
	if err != nil {
		enableJsonMigration = false
	}

	config := &Config{
		GuildID:               strings.TrimSpace(os.Getenv("guild_id")),
		BotToken:              strings.TrimSpace(os.Getenv("bot_token")),
		AppID:                 strings.TrimSpace(os.Getenv("application_id")),
		ModeratorID:           strings.TrimSpace(os.Getenv("moderator_role_id")),
		BotID:                 strings.TrimSpace(os.Getenv("bot_role_id")),
		ModerationChannelName: strings.TrimSpace(os.Getenv("moderation_channel_name")),
		MisskeyToken:          strings.TrimSpace(os.Getenv("misskey_token")),
		MisskeyHost:           strings.TrimSpace(os.Getenv("misskey_host")),
		SavePath:              strings.TrimSpace(os.Getenv("save_path")),
		DatabasePath:          strings.TrimSpace(os.Getenv("database_path")),
		IsDebug:               isDebug,

		// Migration Configuration
		EnableJsonMigration: enableJsonMigration,
		JsonMigrationPath:   strings.TrimSpace(os.Getenv("json_migration_path")),
	}

	// Set default values
	if config.DatabasePath == "" {
		config.DatabasePath = "./emoji_bot.db"
	}

	if config.JsonMigrationPath == "" {
		config.JsonMigrationPath = "./saved/"
	}

	if err := config.Validate(); err != nil {
		return nil, err
	}

	return config, nil
}

func (c *Config) Validate() error {
	requiredFields := map[string]string{
		"guild_id":                c.GuildID,
		"bot_token":               c.BotToken,
		"application_id":          c.AppID,
		"moderator_role_id":       c.ModeratorID,
		"bot_role_id":             c.BotID,
		"moderation_channel_name": c.ModerationChannelName,
		"misskey_token":           c.MisskeyToken,
		"misskey_host":            c.MisskeyHost,
		"save_path":               c.SavePath,
	}

	var missingFields []string
	for fieldName, value := range requiredFields {
		if value == "" {
			missingFields = append(missingFields, fieldName)
		}
	}

	if len(missingFields) > 0 {
		return errors.Validation("missing required environment variables: " + strings.Join(missingFields, ", "))
	}

	if c.SavePath != "" && !strings.HasSuffix(c.SavePath, "/") {
		c.SavePath += "/"
	}

	return nil
}
