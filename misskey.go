package main

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/yitsushi/go-misskey"
	"github.com/yitsushi/go-misskey/core"
	"github.com/yitsushi/go-misskey/models"
	"github.com/yitsushi/go-misskey/services/admin/emoji"
	"github.com/yitsushi/go-misskey/services/drive/files"
	"github.com/yitsushi/go-misskey/services/drive/folders"
	"github.com/yitsushi/go-misskey/services/notes"
	"io"
	"log"
	"os"
	"strings"
)

func uploadToMisskey(e Emoji) bool {
	client, err := misskey.NewClientWithOptions(
		misskey.WithAPIToken(misskeyToken),
		misskey.WithBaseURL("https", misskeyHost, ""),
		misskey.WithLogLevel(logrus.ErrorLevel),
	)

	if err != nil {
		fmt.Println("[ERROR] Could not connect to misskey")
	}

	file, err := os.Open(e.FilePath)
	if err != nil {
		log.Printf("[ERROR] file not found %s", e.FilePath)
		return false
	}
	defer file.Close()

	fileBytes, err := io.ReadAll(file)
	if err != nil {
		log.Printf("[ERROR] file read error %s", e.FilePath)
		return false
	}

	folder, err := getFolder("Emoji", client)
	if err != nil {
		log.Printf("[ERROR] folder error %s", err)
		return false
	}

	drive, err := client.Drive().File().Create(files.CreateRequest{
		FolderID:    folder.ID,
		Name:        e.Name,
		IsSensitive: e.IsSensitive,
		Force:       false,
		Content:     fileBytes,
	})

	if err != nil {
		log.Printf("[Misskey] [Drive/File/Create] %s", err)
		return false
	}

	log.Printf(
		"[Misskey] [Drive/File/Create] %s file uploaded. (%s)",
		core.StringValue(drive.Name),
		drive.ID,
	)

	add, err := client.Admin().Emoji().Add(emoji.AddRequest{
		Name:   e.Name,
		FileID: drive.ID,
	})

	if err != nil {
		log.Printf("[Admin/Emoji/Add] %s", err)
		return false
	}

	log.Printf("[Admin/Emoji/Add] %s", add)

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
		log.Printf("[Admin/Emoji/Update] %s", err)
		return false
	}

	e.IsAccepted = true

	log.Printf("[Admin/Emoji/Update] Update completed: %s", e.Name)

	return true
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

func note(message notes.CreateRequest) {
	client, err := misskey.NewClientWithOptions(
		misskey.WithAPIToken(misskeyToken),
		misskey.WithBaseURL("https", misskeyHost, ""),
		misskey.WithLogLevel(logrus.ErrorLevel),
	)

	if err != nil {
		fmt.Println("[ERROR] Could not connect to misskey")
	}

	response, err := client.Notes().Create(message)

	if err != nil {
		log.Printf("[Notes] Error happened: %s", err)
		return
	}

	log.Println("Created note " + response.CreatedNote.ID)

}
