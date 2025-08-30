package utility

import (
	"io"
	"net/http"
	"os"
	"path/filepath"
)

func EmojiDownload(url string, filePath string) error {
	response, err := http.Get(url)
	if err != nil {
		return err
	}
	defer func() { _ = response.Body.Close() }()

	dirPath := filepath.Dir(filePath)
	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		if err := os.MkdirAll(dirPath, os.ModePerm); err != nil {
			return err
		}
	}

	file, err := os.Create(filePath)
	if err != nil {
		return err
	}

	defer func() { _ = file.Close() }()

	_, err = io.Copy(file, response.Body)
	if err != nil {
		return err
	}

	return nil
}
