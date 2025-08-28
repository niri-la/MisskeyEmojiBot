package errors

import (
	"fmt"
)

type ErrorCode string

const (
	ErrCodeConfig          ErrorCode = "CONFIG_ERROR"
	ErrCodeDiscord         ErrorCode = "DISCORD_ERROR"
	ErrCodeMisskey         ErrorCode = "MISSKEY_ERROR"
	ErrCodeFileOperation   ErrorCode = "FILE_ERROR"
	ErrCodeEmojiNotFound   ErrorCode = "EMOJI_NOT_FOUND"
	ErrCodeEmojiAlready    ErrorCode = "EMOJI_ALREADY_PROCESSED"
	ErrCodeValidation      ErrorCode = "VALIDATION_ERROR"
	ErrCodeUpload          ErrorCode = "UPLOAD_ERROR"
)

type AppError struct {
	Code    ErrorCode
	Message string
	Err     error
}

func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("[%s] %s: %v", e.Code, e.Message, e.Err)
	}
	return fmt.Sprintf("[%s] %s", e.Code, e.Message)
}

func (e *AppError) Unwrap() error {
	return e.Err
}

func New(code ErrorCode, message string) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
	}
}

func Wrap(err error, code ErrorCode, message string) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
		Err:     err,
	}
}

func Config(message string, err error) *AppError {
	return Wrap(err, ErrCodeConfig, message)
}

func Discord(message string, err error) *AppError {
	return Wrap(err, ErrCodeDiscord, message)
}

func Misskey(message string, err error) *AppError {
	return Wrap(err, ErrCodeMisskey, message)
}

func FileOperation(message string, err error) *AppError {
	return Wrap(err, ErrCodeFileOperation, message)
}

func EmojiNotFound(id string) *AppError {
	return New(ErrCodeEmojiNotFound, fmt.Sprintf("emoji not found: %s", id))
}

func EmojiAlready(message string) *AppError {
	return New(ErrCodeEmojiAlready, message)
}

func Validation(message string) *AppError {
	return New(ErrCodeValidation, message)
}

func Upload(message string, err error) *AppError {
	return Wrap(err, ErrCodeUpload, message)
}