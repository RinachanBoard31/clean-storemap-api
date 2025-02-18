package errors

import (
	"errors"
)

var (
	ErrNotFound    = errors.New("record not found")
	ErrDuplicate   = errors.New("duplicate record")
	ErrInvalidData = errors.New("invalid data")
)

// 認証関連エラー
var (
	ErrGetEmailFailed      = errors.New("failed to get email")
	ErrFailedGenerateToken = errors.New("failed to generate token")
)
var (
	ErrUnAuthorized = errors.New("failed to authorize")
)

// 未定義のエラー
var (
	ErrUnsupported = errors.New("unsupported")
)

// userのビジネスルール
var (
	ErrUserNotFound      = errors.New("user not found")
	ErrUserAlreadyExists = errors.New("user already exists")
	ErrInvalidUserData   = errors.New("invalid user data")
)

func Is(err, target error) bool {
	return errors.Is(err, target)
}
