package bot

import (
	"fmt"
	"os"
	"os/signal"
	"strings"

	"github.com/bwmarrin/discordgo"

	"MisskeyEmojiBot/pkg/container"
	"MisskeyEmojiBot/pkg/errors"
)

type Bot struct {
	container *container.Container
}

func New(container *container.Container) *Bot {
	return &Bot{
		container: container,
	}
}

func (b *Bot) Run() error {
	if err := b.setupHandlers(); err != nil {
		return err
	}

	if err := b.container.Session.Open(); err != nil {
		return errors.Discord("failed to open Discord connection", err)
	}
	defer func() { _ = b.container.Session.Close() }()

	b.startJobs()

	println(":: Bot starting")
	println(":::::::::::::::::::::::")

	// Wait for interrupt signal
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	<-stop

	println(":: Graceful shutdown")
	return nil
}

func (b *Bot) setupHandlers() error {
	// Ready handler
	b.container.Session.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		println(":: Bot ready")
	})

	// Interaction handler
	b.container.Session.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		switch i.Type {
		case discordgo.InteractionApplicationCommand:
			b.container.CommandHandler.Handle(s, i)
		case discordgo.InteractionMessageComponent:
			b.container.ComponentHandler.Handle(s, i)
		}
	})

	// Message handler for emoji channels
	b.container.Session.AddHandler(func(s *discordgo.Session, m *discordgo.MessageCreate) {
		if m.Author.ID == s.State.User.ID {
			return
		}

		channel, _ := s.Channel(m.ChannelID)
		if !strings.Contains(channel.Name, "Emoji-") {
			return
		}

		emoji, err := b.container.EmojiRepository.GetEmoji(channel.Name[6:])
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}

		_ = b.container.EmojiRequestHandler.Process(emoji, s, m)
	})

	// Emoji moderation reaction handler
	b.container.Session.AddHandler(b.container.EmojiModerationReaction.HandleEmojiModerationReaction)

	return nil
}

func (b *Bot) startJobs() {
	b.container.ChannelDeleteJob.Run()
	b.container.EmojiUpdateInfoJob.Run()
}
