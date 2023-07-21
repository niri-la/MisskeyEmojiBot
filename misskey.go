package main

import (
	"github.com/sirupsen/logrus"
	"github.com/yitsushi/go-misskey"
	"github.com/yitsushi/go-misskey/models"
	"github.com/yitsushi/go-misskey/services/admin/emoji"
	"github.com/yitsushi/go-misskey/services/drive/files"
	"github.com/yitsushi/go-misskey/services/drive/folders"
	"github.com/yitsushi/go-misskey/services/notes"
	"io"
	"os"
	"strings"
)

func uploadToMisskey(e Emoji) error {
	client, err := misskey.NewClientWithOptions(
		misskey.WithAPIToken(misskeyToken),
		misskey.WithBaseURL("https", misskeyHost, ""),
		misskey.WithLogLevel(logrus.ErrorLevel),
	)

	if err != nil {
		return err
	}

	file, err := os.Open(e.FilePath)

	if err != nil {
		return err
	}

	defer file.Close()

	fileBytes, err := io.ReadAll(file)

	if err != nil {
		return err
	}

	folder, err := getFolder("Emoji", client)

	if err != nil {
		return err
	}

	drive, err := client.Drive().File().Create(files.CreateRequest{
		FolderID:    folder.ID,
		Name:        e.Name,
		IsSensitive: e.IsSensitive,
		Force:       false,
		Content:     fileBytes,
	})

	if err != nil {
		return err
	}

	add, err := client.Admin().Emoji().Add(emoji.AddRequest{
		Name:   e.Name,
		FileID: drive.ID,
	})

	if err != nil {
		return err
	}

	err = client.Admin().Emoji().Update(emoji.UpdateRequest{
		ID:                                      add,
		Name:                                    e.Name,
		Category:                                e.Category,
		Aliases:                                 strings.Split(e.Tag, " "),
		License:                                 "",
		IsSensitive:                             e.IsSensitive,
		LocalOnly:                               false,
		RoleIdsThatCanBeUsedThisEmojiAsReaction: []string{},
	})

	if err != nil {
		return err
	}

	e.IsAccepted = true

	return nil
}

func getFolder(folderName string, client *misskey.Client) (models.Folder, error) {
	find, err := client.Drive().Folder().Find(folders.FindRequest{
		Name: folderName,
	})
	if err == nil && len(find) != 0 {
		return find[0], nil
	}

	create, err := client.Drive().Folder().Create(folders.CreateRequest{
		Name: folderName,
	})

	if err != nil {
		return models.Folder{}, err
	}
	return create, nil
}

func note(message notes.CreateRequest) error {
	client, err := misskey.NewClientWithOptions(
		misskey.WithAPIToken(misskeyToken),
		misskey.WithBaseURL("https", misskeyHost, ""),
		misskey.WithLogLevel(logrus.ErrorLevel),
	)

	if err != nil {
		return err
	}

	response, err := client.Notes().Create(message)

	if err != nil {
		return err
	}

	logger.WithFields(logrus.Fields{
		"event":   "misskey",
		"id":      response.CreatedNote.ID,
		"message": message,
	}).Debug("note complete.")
	return nil
}
