package migration

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"MisskeyEmojiBot/pkg/config"
	"MisskeyEmojiBot/pkg/entity"
	"MisskeyEmojiBot/pkg/repository"
)

type JsonMigration struct {
	config          *config.Config
	emojiRepository repository.EmojiRepository
}

func NewJsonMigration(cfg *config.Config, emojiRepo repository.EmojiRepository) *JsonMigration {
	return &JsonMigration{
		config:          cfg,
		emojiRepository: emojiRepo,
	}
}

func (m *JsonMigration) Run() error {
	if !m.config.EnableJsonMigration {
		fmt.Println("JSON migration is disabled in configuration.")
		return nil
	}

	if _, err := os.Stat(m.config.JsonMigrationPath); os.IsNotExist(err) {
		fmt.Printf("Migration path does not exist: %s\n", m.config.JsonMigrationPath)
		return nil
	}

	fmt.Printf("Starting JSON to SQLite migration from: %s\n", m.config.JsonMigrationPath)

	// JSONファイルを検索
	jsonFiles, err := m.findJsonFiles()
	if err != nil {
		return fmt.Errorf("failed to find JSON files: %w", err)
	}

	if len(jsonFiles) == 0 {
		fmt.Println("No JSON files found for migration.")
		return nil
	}

	migratedCount := 0
	skippedCount := 0

	for _, filePath := range jsonFiles {
		emoji, err := m.loadJsonFile(filePath)
		if err != nil {
			fmt.Printf("Warning: Failed to load %s: %v\n", filePath, err)
			continue
		}

		// 既存チェック
		existing, err := m.emojiRepository.GetEmoji(emoji.ID)
		if err == nil && existing != nil {
			fmt.Printf("Skipping existing emoji: %s\n", emoji.ID)
			skippedCount++
			continue
		}

		// データベースに保存
		if err := m.emojiRepository.Save(emoji); err != nil {
			fmt.Printf("Warning: Failed to save emoji %s: %v\n", emoji.ID, err)
			continue
		}

		migratedCount++
		fmt.Printf("Migrated emoji: %s (%s)\n", emoji.Name, emoji.ID)
	}

	fmt.Printf("Migration completed: %d migrated, %d skipped\n", migratedCount, skippedCount)
	return nil
}

func (m *JsonMigration) findJsonFiles() ([]string, error) {
	var jsonFiles []string

	err := filepath.Walk(m.config.JsonMigrationPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && strings.HasSuffix(strings.ToLower(path), ".json") {
			jsonFiles = append(jsonFiles, path)
		}

		return nil
	})

	return jsonFiles, err
}

func (m *JsonMigration) loadJsonFile(filePath string) (*entity.Emoji, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	var legacyEmoji entity.LegacyEmoji
	if err := json.Unmarshal(data, &legacyEmoji); err != nil {
		return nil, err
	}

	// 旧形式から新形式に変換
	return legacyEmoji.ToCurrentEmoji(), nil
}

// CleanupAfterMigration moves migrated JSON files to a backup directory
func (m *JsonMigration) CleanupAfterMigration() error {
	backupDir := filepath.Join(m.config.JsonMigrationPath, "_migrated_backup")
	if err := os.MkdirAll(backupDir, 0755); err != nil {
		return fmt.Errorf("failed to create backup directory: %w", err)
	}

	jsonFiles, err := m.findJsonFiles()
	if err != nil {
		return err
	}

	for _, filePath := range jsonFiles {
		fileName := filepath.Base(filePath)
		backupPath := filepath.Join(backupDir, fileName)
		
		if err := os.Rename(filePath, backupPath); err != nil {
			fmt.Printf("Warning: Failed to move %s to backup: %v\n", filePath, err)
		}
	}

	fmt.Printf("Moved %d JSON files to backup directory: %s\n", len(jsonFiles), backupDir)
	return nil
}