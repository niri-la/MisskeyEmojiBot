package main

import (
	"log"
	"os"

	"MisskeyEmojiBot/pkg/bot"
	"MisskeyEmojiBot/pkg/config"
	"MisskeyEmojiBot/pkg/container"
	"MisskeyEmojiBot/pkg/errors"
)

func main() {
	println(":::::::::::::::::::::::")
	println(":: Misskey Emoji Bot ")
	println(":::::::::::::::::::::::")
	println(":: initializing")

	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	config, err := config.LoadConfig()
	if err != nil {
		return err
	}

	if err := ensureSaveDirectory(config.SavePath); err != nil {
		return err
	}

	container, err := container.NewContainer(config)
	if err != nil {
		return err
	}

	bot := bot.New(container)
	return bot.Run()
}

func ensureSaveDirectory(savePath string) error {
	_, err := os.Stat(savePath)
	if os.IsNotExist(err) {
		if err := os.MkdirAll(savePath, os.ModePerm); err != nil {
			return errors.FileOperation("failed to create save directory", err)
		}
	}
	return nil
}
