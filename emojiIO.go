package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

var (
	validExtensions = map[string]bool{
		".png":  true,
		".jpg":  true,
		".jpeg": true,
		".gif":  true,
	}
)

func emojiDownload(url string, filePath string) error {
	response, err := http.Get(url)
	if err != nil {
		return err
	}

	defer response.Body.Close()

	dirPath := filepath.Dir(filePath)
	// ディレクトリが存在するかチェック
	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		// ディレクトリが存在しない場合、作成する
		os.MkdirAll(dirPath, os.ModePerm)
	}

	file, err := os.Create(filePath)
	if err != nil {
		return err
	}

	defer file.Close()

	_, err = io.Copy(file, response.Body)
	if err != nil {
		return err
	}

	return nil
}

func deleteEmoji(filePath string) {
	err := os.Remove(filePath)
	if err != nil {
		fmt.Printf("[ERROR] file not found %s\n", filePath)
	}
}

func isValidEmojiFile(fileName string) bool {
	fileExtension := filepath.Ext(fileName)
	_, exists := validExtensions[fileExtension]
	return exists
}
