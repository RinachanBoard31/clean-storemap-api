package db

import (
	"errors"

	"gorm.io/gorm"
)

var (
	ErrRecordNotFound     = errors.New("record not found")
	ErrDuplicatedKey      = errors.New("duplicate record")
	ErrInvalidData        = errors.New("invalid data")
	ErrInvalidDB          = errors.New("invalid db")
	ErrMissingWhereClause = errors.New("missing where clause")
	ErrInvalidField       = errors.New("invalid field")
	ErrInvalidValue       = errors.New("invalid value")
	ErrUnsupported        = errors.New("unsupported")
)

func convertGormError(err error) error {
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return ErrRecordNotFound
	} else if errors.Is(err, gorm.ErrDuplicatedKey) {
		return ErrDuplicatedKey
	} else if errors.Is(err, gorm.ErrInvalidData) {
		return ErrInvalidData
	} else if errors.Is(err, gorm.ErrInvalidDB) {
		return ErrInvalidDB
	} else if errors.Is(err, gorm.ErrMissingWhereClause) {
		return ErrMissingWhereClause
	} else if errors.Is(err, gorm.ErrInvalidField) {
		return ErrInvalidField
	} else if errors.Is(err, gorm.ErrInvalidValue) {
		return ErrInvalidValue
	}
	return ErrUnsupported
}
