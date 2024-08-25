package repository

import (
	"MisskeyEmojiBot/pkg/entity"
	"io"
	"os"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/yitsushi/go-misskey"
	"github.com/yitsushi/go-misskey/core"
	"github.com/yitsushi/go-misskey/models"
	"github.com/yitsushi/go-misskey/services/admin/emoji"
	"github.com/yitsushi/go-misskey/services/drive/files"
	"github.com/yitsushi/go-misskey/services/drive/folders"
	"github.com/yitsushi/go-misskey/services/notes"
)

type MisskeyRepository interface {
	Note(message notes.CreateRequest) error
	UploadEmoji(emoji *entity.Emoji) error
	GetFolder(folderName string) (models.Folder, error)
	NewString(message string) core.String
}

type misskeyRepository struct {
	client *misskey.Client
}

func NewMisskeyRepository(misskeyToken string, misskeyHost string) (MisskeyRepository, error) {
	client, err := misskey.NewClientWithOptions(
		misskey.WithAPIToken(misskeyToken),
		misskey.WithBaseURL("https", misskeyHost, ""),
		misskey.WithLogLevel(logrus.ErrorLevel),
	)

	if err != nil {
		return nil, err
	}

	return &misskeyRepository{client: client}, nil
}

func (r *misskeyRepository) UploadEmoji(userEmoji *entity.Emoji) error {
	file, err := os.Open(userEmoji.FilePath)

	if err != nil {
		return err
	}

	defer file.Close()

	fileBytes, err := io.ReadAll(file)

	if err != nil {
		return err
	}

	folder, err := r.GetFolder("Emoji")

	if err != nil {
		return err
	}

	drive, err := r.client.Drive().File().Create(files.CreateRequest{
		FolderID:    folder.ID,
		Name:        userEmoji.Name,
		IsSensitive: userEmoji.IsSensitive,
		Force:       false,
		Content:     fileBytes,
	})

	if err != nil {
		return err
	}

	add, err := r.client.Admin().Emoji().Add(emoji.AddRequest{
		Name:   userEmoji.Name,
		FileID: drive.ID,
	})

	if err != nil {
		return err
	}

	err = r.client.Admin().Emoji().Update(emoji.UpdateRequest{
		ID:                                      add,
		Name:                                    userEmoji.Name,
		Category:                                userEmoji.Category,
		Aliases:                                 strings.Split(userEmoji.Tag, " "),
		License:                                 userEmoji.License,
		IsSensitive:                             userEmoji.IsSensitive,
		LocalOnly:                               false,
		RoleIdsThatCanBeUsedThisEmojiAsReaction: []string{},
	})

	if err != nil {
		return err
	}

	userEmoji.IsAccepted = true

	return nil
}

func (r *misskeyRepository) NewString(message string) core.String {
	return core.NewString(message)
}

func (r *misskeyRepository) Note(message notes.CreateRequest) error {

	_, err := r.client.Notes().Create(message)
	if err != nil {
		return err
	}

	return nil
}

func (r *misskeyRepository) GetFolder(folderName string) (models.Folder, error) {
	find, err := r.client.Drive().Folder().Find(folders.FindRequest{
		Name: folderName,
	})
	if err == nil && len(find) != 0 {
		return find[0], nil
	}

	create, err := r.client.Drive().Folder().Create(folders.CreateRequest{
		Name: folderName,
	})

	if err != nil {
		return models.Folder{}, err
	}
	return create, nil
}
