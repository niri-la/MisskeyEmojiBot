package container

import (
	"os"

	"github.com/bwmarrin/discordgo"

	"MisskeyEmojiBot/pkg/command"
	"MisskeyEmojiBot/pkg/component"
	"MisskeyEmojiBot/pkg/config"
	"MisskeyEmojiBot/pkg/database"
	"MisskeyEmojiBot/pkg/errors"
	"MisskeyEmojiBot/pkg/handler"
	"MisskeyEmojiBot/pkg/handler/emoji"
	"MisskeyEmojiBot/pkg/handler/processor"
	"MisskeyEmojiBot/pkg/job"
	"MisskeyEmojiBot/pkg/repository"
)

type Container struct {
	Config *config.Config

	// Repositories
	DiscordRepository repository.DiscordRepository
	EmojiRepository   repository.EmojiRepository
	MisskeyRepository repository.MisskeyRepository

	// Handlers
	EmojiHandler            emoji.EmojiHandler
	EmojiRequestHandler     handler.EmojiRequestHandler
	EmojiModerationReaction emoji.EmojiModerationReactionHandler
	CommandHandler          *handler.CommandHandler
	ComponentHandler        handler.ComponentHandler

	// Jobs
	ChannelDeleteJob   job.Job
	EmojiUpdateInfoJob job.Job

	// Discord Session
	Session *discordgo.Session

	// Version
	Version string
}

func NewContainer(cfg *config.Config) (*Container, error) {
	// Read version
	version, err := os.ReadFile("version.txt")
	if err != nil {
		return nil, errors.FileOperation("failed to read version file", err)
	}

	// Initialize Discord session
	session, err := discordgo.New("Bot " + cfg.BotToken)
	if err != nil {
		return nil, errors.Discord("failed to initialize Discord bot", err)
	}

	// Initialize database
	db, err := database.NewDatabase(cfg.DatabasePath)
	if err != nil {
		return nil, errors.FileOperation("failed to initialize database", err)
	}

	// Run migrations
	if err := db.Migrate(); err != nil {
		return nil, err
	}

	// Initialize repositories
	discordRepo := repository.NewDiscordRepository(session)
	emojiRepo := repository.NewEmojiGormRepository(db, *cfg)
	misskeyRepo, err := repository.NewMisskeyRepository(cfg.MisskeyToken, cfg.MisskeyHost)
	if err != nil {
		return nil, errors.Misskey("failed to initialize Misskey API", err)
	}

	// Initialize handlers
	emojiHandler := emoji.NewEmojiHandler(emojiRepo, discordRepo, misskeyRepo)
	emojiRequestHandler := handler.NewEmojiRequestHandler(emojiRepo)
	emojiModerationReaction := emoji.NewEmojiModerationReactionHandler(emojiHandler, emojiRepo, discordRepo, *cfg)
	commandHandler := handler.NewCommandHandler(*cfg, discordRepo)
	componentHandler := handler.NewComponentHandler()

	// Initialize jobs
	channelDeleteJob := job.NewChannelDeleteJob(emojiRepo, discordRepo)
	emojiUpdateInfoJob := job.NewEmojiUpdateInfoJob(emojiRepo, misskeyRepo)

	container := &Container{
		Config:                  cfg,
		DiscordRepository:       discordRepo,
		EmojiRepository:         emojiRepo,
		MisskeyRepository:       misskeyRepo,
		EmojiHandler:            emojiHandler,
		EmojiRequestHandler:     emojiRequestHandler,
		EmojiModerationReaction: emojiModerationReaction,
		CommandHandler:          commandHandler,
		ComponentHandler:        componentHandler,
		ChannelDeleteJob:        channelDeleteJob,
		EmojiUpdateInfoJob:      emojiUpdateInfoJob,
		Session:                 session,
		Version:                 string(version),
	}

	// Register commands and components
	container.registerCommands()
	container.registerComponents()
	container.registerProcessors()

	return container, nil
}

func (c *Container) registerCommands() {
	c.CommandHandler.RegisterCommand(command.NewInitCommand(*c.Config, c.DiscordRepository))
	c.CommandHandler.RegisterCommand(command.NewNirilaCommand(c.DiscordRepository, c.Version))
	c.CommandHandler.RegisterCommand(command.NewEmojiDetailChangeCommand(*c.Config, c.EmojiRepository, c.DiscordRepository))
	c.CommandHandler.RegisterCommand(command.NewEmojiListCommand(c.EmojiRepository))
}

func (c *Container) registerComponents() {
	_ = c.ComponentHandler.AddComponent(component.NewCreateEmojiChannelComponen(c.EmojiRequestHandler, c.EmojiRepository, c.DiscordRepository))
	_ = c.ComponentHandler.AddComponent(component.NewEmojiCancelRequestComponent(c.EmojiRepository, c.DiscordRepository))
	_ = c.ComponentHandler.AddComponent(component.NewEmojiRequestComponen(*c.Config, c.EmojiRepository, c.DiscordRepository))
	_ = c.ComponentHandler.AddComponent(component.NewEmojiRequestRetryComponen(c.EmojiRequestHandler, c.EmojiRepository, c.DiscordRepository))
	_ = c.ComponentHandler.AddComponent(component.NewInitComponent(*c.Config, c.DiscordRepository))
	_ = c.ComponentHandler.AddComponent(component.NewNsfwActiveComponent(c.EmojiRequestHandler, c.EmojiRepository, c.DiscordRepository))
	_ = c.ComponentHandler.AddComponent(component.NewNsfwInactiveComponent(c.EmojiRequestHandler, c.EmojiRepository, c.DiscordRepository))
}

func (c *Container) registerProcessors() {
	c.EmojiRequestHandler.AddProcess(processor.NewUploadHandler(*c.Config))
	c.EmojiRequestHandler.AddProcess(processor.NewNameSettingHandler())
	c.EmojiRequestHandler.AddProcess(processor.NewCategoryHandler())
	c.EmojiRequestHandler.AddProcess(processor.NewTagHandler())
	c.EmojiRequestHandler.AddProcess(processor.NewLicenseHandlerHandler())
	c.EmojiRequestHandler.AddProcess(processor.NewOtherHandler())
	c.EmojiRequestHandler.AddProcess(processor.NewNsfwHandler())
	c.EmojiRequestHandler.AddProcess(processor.NewConfirmHandler())
}
