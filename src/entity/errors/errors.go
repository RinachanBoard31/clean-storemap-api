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
	// ビジネスルール
	ErrUnAuthorized = errors.New("認証に失敗しました。")
)

// 未定義のエラー
var (
	ErrUnsupported = errors.New("unsupported")
)

// userのビジネスルール
var (
	ErrUserNotFound       = errors.New("ユーザが見つかりません。")
	ErrUserAlreadyExists  = errors.New("ユーザが既に存在します。")
	ErrInvalidUserData    = errors.New("ユーザのデータが不正です。")
	ErrInvalidEmail       = errors.New("ユーザのemailが不正です。")
	ErrAgeBelowZero       = errors.New("年齢が0未満です。")
	ErrInvalidEmailAndAge = errors.New("ユーザのemailと年齢が不正です。")
	ErrUnExpected         = errors.New("予期せぬエラーが発生しました。")
)

func Is(err, target error) bool {
	return errors.Is(err, target)
}
