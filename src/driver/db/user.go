package db

import (
	"clean-storemap-api/src/entity/errors"
	"time"

	"gorm.io/gorm"
)

type DbUserDriver struct{}

func NewUserDriver() *DbUserDriver {
	return &DbUserDriver{}
}

type User struct {
	Id        string  `gorm:"primaryKey"`
	Name      string  `gorm:"not null"`
	Email     string  `gorm:"unique"`
	Age       int     `gorm:"not null"`
	Sex       float32 `gorm:"not null"`
	Gender    float32 `gorm:"not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

func convertGormError(err error) error {
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return errors.ErrNotFound
	} else if errors.Is(err, gorm.ErrDuplicatedKey) {
		return errors.ErrDuplicate
	} else if errors.Is(err, gorm.ErrInvalidData) {
		return errors.ErrInvalidData
	}
	return err
}

func (dbu *DbUserDriver) CreateUser(user *User) (*User, error) {
	result := DB.Create(&user)
	err := result.Error
	if err != nil {
		return nil, convertGormError(err)
	}
	return user, err
}

func (dbu *DbUserDriver) FindByEmail(email string) (*User, error) {
	var user *User
	result := DB.Where("email = ?", email).First(&user)
	if result.Error != nil {
		return nil, convertGormError(result.Error)
	}
	return user, nil
}

func (dbu *DbUserDriver) UpdateUser(user *User, updateData map[string]interface{}) error {
	result := DB.Model(&user).Updates(updateData)
	if err := result.Error; err != nil {
		return convertGormError(err)
	}
	return nil
}

func (dbu *DbUserDriver) FindById(id string) (*User, error) {
	var user *User
	result := DB.First(&user, "id = ?", id)
	if result.Error != nil {
		return nil, convertGormError(result.Error)
	}
	return user, nil
}
