package job

import (
	"fmt"
	"io"
	"os"
	"time"

	"MisskeyEmojiBot/pkg/config"
	"MisskeyEmojiBot/pkg/repository"
)

type DatabaseBackupJob interface {
	Run() error
}

type databaseBackupJob struct {
	config *config.Config
	s3Repo repository.S3Repository
}

func NewDatabaseBackupJob(cfg *config.Config) (DatabaseBackupJob, error) {
	job := &databaseBackupJob{config: cfg}

	if cfg.UseS3 {
		s3Repo, err := repository.NewS3Repository(cfg)
		if err != nil {
			return nil, fmt.Errorf("failed to create S3 repository: %w", err)
		}
		job.s3Repo = s3Repo
	}

	return job, nil
}

func (j *databaseBackupJob) Run() error {

	if !j.config.EnableDatabaseBackup {
		fmt.Println("Database backup job is disabled in configuration.")
		return nil
	}

	if j.config.DatabasePath == "" {
		return fmt.Errorf("database path is not configured")
	}

	cleanRequest := time.NewTicker(12 * time.Hour)
	go func() {
		for range cleanRequest.C {
			if !j.config.UseS3 || j.s3Repo == nil {
				// S3が有効でない場合はローカルバックアップのみ
				j.createLocalBackup()
			}

			// データベースファイルを読み込む
			dbData, err := j.readDatabaseFile()
			if err != nil {
				fmt.Errorf("failed to read database file: %w", err)
			}

			// S3にバックアップをアップロード
			timestamp := time.Now().Format("2006-01-02_15-04-05")
			backupKey := fmt.Sprintf("backups/database_%s.db", timestamp)

			_, err = j.s3Repo.UploadFile(backupKey, dbData, "application/octet-stream")
			if err != nil {
				fmt.Errorf("failed to upload backup to S3: %w", err)
			}

			fmt.Printf("Database backup uploaded to S3: %s\n", backupKey)

			// ローカルバックアップも作成（オプション）
			if err := j.createLocalBackup(); err != nil {
				fmt.Printf("Warning: Local backup failed: %v\n", err)
			}
		}
	}()

	return nil
}

func (j *databaseBackupJob) readDatabaseFile() ([]byte, error) {
	file, err := os.Open(j.config.DatabasePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (j *databaseBackupJob) createLocalBackup() error {
	timestamp := time.Now().Format("2006-01-02_15-04-05")
	backupPath := fmt.Sprintf("%sbackup_%s.db", j.config.SavePath, timestamp)

	// 元のファイルを開く
	src, err := os.Open(j.config.DatabasePath)
	if err != nil {
		return err
	}
	defer src.Close()

	// バックアップディレクトリを作成
	if err := os.MkdirAll(j.config.SavePath, os.ModePerm); err != nil {
		return err
	}

	// バックアップファイルを作成
	dst, err := os.Create(backupPath)
	if err != nil {
		return err
	}
	defer dst.Close()

	// ファイルをコピー
	_, err = io.Copy(dst, src)
	if err != nil {
		return err
	}

	fmt.Printf("Local database backup created: %s\n", backupPath)
	return nil
}
