package auth

import "errors"

var (
	ErrGetEmailFailed      = errors.New("failed to get email")
	ErrFailedGenerateToken = errors.New("failed to generate token")
)
