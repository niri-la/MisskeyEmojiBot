package main

import (
	_ "embed"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"time"
)

// Bot parameters
var (
	GuildID               string
	BotToken              string
	AppID                 string
	ModeratorID           string
	BotID                 string
	ModerationChannelName string
	misskeyToken          string
	misskeyHost           string
	isDebug               bool
	Session               *discordgo.Session
	logger                *logrus.Logger
)

var moderationChannel *discordgo.Channel

//go:embed message/ja-jp.yaml
var messageJp string

func init() {

	logger = logrus.New()
	// Log as JSON instead of the default ASCII formatter.
	//debug.SetFormatter(&debug.TextFormatter{})
	logger.SetOutput(os.Stdout)
	logger.SetFormatter(&logrus.TextFormatter{
		ForceColors: true,
	})
}

func init() {
	loadEnvironments()
	var err error
	Session, err = discordgo.New("Bot " + BotToken)
	if err != nil {
		logger.WithFields(logrus.Fields{
			"event": "init",
		}).Error(err)
		return
	}
	command()
	moderationChannel, err = findChannelByName(Session, GuildID, ModerationChannelName)
	registeredCommands := make([]*discordgo.ApplicationCommand, len(Commands))
	for i, v := range Commands {
		cmd, err := Session.ApplicationCommandCreate(AppID, GuildID, v)
		if err != nil {
			Session.Close()
			logger.WithFields(logrus.Fields{
				"event":     "commad",
				"name":      v.Name,
				"guild id":  GuildID,
				"bot token": BotToken,
			}).Panic(err)
		}
		registeredCommands[i] = cmd
	}
	if isDebug {
		logger.SetLevel(logrus.TraceLevel)
	} else {
		logger.SetLevel(logrus.InfoLevel)
	}
}

func main() {
	logger.Info(":::::::::::::::::::::::")
	logger.Info(":: Misskey Emoji Bot ")
	logger.Info(":::::::::::::::::::::::")
	logger.Info(":: initializing")
	// start
	Session.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		logger.Info(":: Bot starting")
		logger.Info(":::::::::::::::::::::::")
	})

	// コンポーネントはインタラクションの一部なので、InteractionCreateHandlerを登録します。
	Session.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		switch i.Type {
		case discordgo.InteractionApplicationCommand:
			if h, ok := CommandHandlers[i.ApplicationCommandData().Name]; ok {
				h(s, i)
			}
		case discordgo.InteractionMessageComponent:
			if h, ok := ComponentsHandlers[i.MessageComponentData().CustomID]; ok {
				h(s, i)
			}
		}
	})

	Session.AddHandler(func(s *discordgo.Session, m *discordgo.MessageCreate) {
		if m.Author.ID == s.State.User.ID {
			return
		}

		channel, _ := s.Channel(m.ChannelID)

		if !strings.Contains(channel.Name, "Emoji-") {
			return
		}

		emoji, err := GetEmoji(channel.Name[6:])

		if err != nil {
			logger.WithFields(logrus.Fields{
				"user":    m.Author.Username,
				"channel": channel.Name,
				"message": m.Content,
			}).Debug(err)
			return
		}

		logger.WithFields(logrus.Fields{
			"user":    m.Author.Username,
			"channel": channel.Name,
			"message": m.Content,
		}).Trace("Send emoji channel.")

		Process(emoji, s, m)
	})

	Session.AddHandler(emojiModerationReaction)

	_, err := Session.ApplicationCommandCreate(AppID, GuildID, &discordgo.ApplicationCommand{
		Name:        "buttons",
		Description: "Test the buttons if you got courage",
	})

	if err != nil {
		logger.WithFields(logrus.Fields{
			"event": "Session",
		}).Fatal(err)
	}

	err = Session.Open()
	if err != nil {
		logger.WithFields(logrus.Fields{
			"event": "Session",
		}).Fatal(err)
	}

	defer Session.Close()

	ticker := time.NewTicker(24 * time.Hour)
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

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	<-stop
	logger.Info(":: Graceful shutdown")
}

func findChannelByName(s *discordgo.Session, guildID string, name string) (*discordgo.Channel, error) {
	channels, err := s.GuildChannels(guildID)
	if err != nil {
		return nil, err
	}

	for _, ch := range channels {
		if ch.Name == name {
			return ch, nil
		}
	}

	return nil, fmt.Errorf("channel not found")
}

func returnFailedMessage(s *discordgo.Session, i *discordgo.InteractionCreate, reason string) {
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Flags:   discordgo.MessageFlagsEphemeral,
			Content: "新たな申請のRequestに失敗しました。管理者に問い合わせを行ってください。",
		},
	})

	logger.Error(reason)
	return
}

func loadEnvironments() {
	err := godotenv.Load("settings.env")

	if err != nil {
		panic(err)
	}

	GuildID = os.Getenv("guild_id")
	BotToken = os.Getenv("bot_token")
	AppID = os.Getenv("application_id")
	ModeratorID = os.Getenv("moderator_role_id")
	BotID = os.Getenv("bot_role_id")
	ModerationChannelName = os.Getenv("moderation_channel_name")
	misskeyToken = os.Getenv("misskey_token")
	misskeyHost = os.Getenv("misskey_host")
	isDebug, _ = strconv.ParseBool(os.Getenv("debug"))

	logger.Debug(GuildID)
	logger.Debug(BotToken)
	logger.Debug(AppID)
	logger.Debug(BotID)
	logger.Debug(ModerationChannelName)
	logger.Debug(misskeyToken)
	logger.Debug(misskeyHost)
	logger.Debug(isDebug)

}
