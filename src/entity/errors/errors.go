package errors

import (
	"errors"
)

// 汎用エラー
var (
	ErrDBConnection = errors.New("failed to connect database")
	ErrUnsupported  = errors.New("unsupported")
)

// 認証関連エラー
var (
	ErrUnAuthorized = errors.New("認証に失敗しました。")
)

// User関連エラー
var (
	ErrUserNotFound       = errors.New("ユーザが見つかりません。")
	ErrUserAlreadyExists  = errors.New("ユーザが既に存在します。")
	ErrUserInvalidData    = errors.New("ユーザのデータが不正です。")
	ErrInvalidEmail       = errors.New("ユーザのemailが不正です。")
	ErrAgeBelowZero       = errors.New("年齢が0未満です。")
	ErrInvalidEmailAndAge = errors.New("ユーザのemailと年齢が不正です。")
	ErrUnExpected         = errors.New("予期せぬエラーが発生しました。")
)

// Store関連エラー
var (
	ErrStoreNotFound         = errors.New("店舗が見つかりません。")
	ErrStoreAlreadyExists    = errors.New("店舗が既に存在します。")
	ErrStoreInvalidData      = errors.New("店舗のデータが不正です。")
	ErrFailedGetNearByStores = errors.New("近くの店舗の取得に失敗しました。")
)

func Is(err, target error) bool {
	return errors.Is(err, target)
}
