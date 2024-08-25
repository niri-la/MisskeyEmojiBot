package main

import (
	"MisskeyEmojiBot/pkg/command"
	"MisskeyEmojiBot/pkg/component"
	"MisskeyEmojiBot/pkg/config"
	"MisskeyEmojiBot/pkg/handler"
	"MisskeyEmojiBot/pkg/handler/emoji"
	"MisskeyEmojiBot/pkg/handler/processor"
	"MisskeyEmojiBot/pkg/job"
	"MisskeyEmojiBot/pkg/repository"
	_ "embed"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"

	"github.com/bwmarrin/discordgo"
)

var moderationChannel *discordgo.Channel

var Session *discordgo.Session

// //go:embed message/ja-jp.yaml
// var messageJp string

func main() {
	println(":::::::::::::::::::::::")
	println(":: Misskey Emoji Bot ")
	println(":::::::::::::::::::::::")
	println(":: initializing")

	config := config.LoadConfig()

	Session, err := discordgo.New("Bot " + config.BotToken)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	// start
	Session.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		println(":: Bot starting")
		println(":::::::::::::::::::::::")
	})

	version, err := os.ReadFile("version.txt")
	if err != nil {
		log.Fatal(err)
	}

	discordRepository := repository.NewDiscordRepository(Session)
	emojiRepository := repository.NewEmojiRepository()
	misskeyRepository, err := repository.NewMisskeyRepository(config.MisskeyToken, config.MisskeyHost)

	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	emojiHandler := emoji.NewEmojiHandler(emojiRepository, discordRepository, misskeyRepository)
	emojiRequestHandler := handler.NewEmojiRequestHandler()
	emojiModerationReaction := emoji.NewEmojiModerationReactionHandler(emojiHandler, emojiRepository, discordRepository, config)
	commandHandler := handler.NewCommandHandler(config, discordRepository)
	componentHandler := handler.NewComponentHandler()

	channelDeleteJob := job.NewChannelDeleteJob(emojiRepository, discordRepository)
	emojiUpdateInfoJob := job.NewEmojiUpdateInfoJob(emojiRepository, misskeyRepository)

	// register command
	commandHandler.RegisterCommand(command.NewInitCommand(config, discordRepository))
	commandHandler.RegisterCommand(command.NewNirilaCommand(discordRepository, string(version)))

	// register component
	componentHandler.AddComponent(component.NewCreateEmojiChannelComponen(emojiRequestHandler, emojiRepository, discordRepository))
	componentHandler.AddComponent(component.NewEmojiCancelRequestComponent(emojiRepository, discordRepository))
	componentHandler.AddComponent(component.NewEmojiRequestComponen(config, emojiRepository, discordRepository))
	componentHandler.AddComponent(component.NewEmojiRequestRetryComponen(emojiRequestHandler, emojiRepository, discordRepository))
	componentHandler.AddComponent(component.NewInitComponent(config, discordRepository))
	componentHandler.AddComponent(component.NewNsfwActiveComponent(emojiRequestHandler, emojiRepository, discordRepository))
	componentHandler.AddComponent(component.NewNsfwInactiveComponent(emojiRequestHandler, emojiRepository, discordRepository))

	// register processor
	emojiRequestHandler.AddProcess(processor.NewUploadHandler())
	emojiRequestHandler.AddProcess(processor.NewNameSettingHandler())
	emojiRequestHandler.AddProcess(processor.NewCategoryHandler())
	emojiRequestHandler.AddProcess(processor.NewTagHandler())
	emojiRequestHandler.AddProcess(processor.NewLicenseHandlerHandler())
	emojiRequestHandler.AddProcess(processor.NewOtherHandler())
	emojiRequestHandler.AddProcess(processor.NewNsfwHandler())
	emojiRequestHandler.AddProcess(processor.NewConfirmHandler())

	// コンポーネントはインタラクションの一部なので、InteractionCreateHandlerを登録します。
	Session.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		switch i.Type {
		case discordgo.InteractionApplicationCommand:
			commandHandler.Handle(s, i)
		case discordgo.InteractionMessageComponent:
			componentHandler.Handle(s, i)
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

		emoji, err := emojiRepository.GetEmoji(channel.Name[6:])

		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}

		emojiRequestHandler.Process(emoji, s, m)
	})

	Session.AddHandler(emojiModerationReaction.HandleEmojiModerationReaction)

	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	err = Session.Open()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	defer Session.Close()

	channelDeleteJob.Run()
	emojiUpdateInfoJob.Run()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	<-stop
	println(":: Graceful shutdown")
	return
}
