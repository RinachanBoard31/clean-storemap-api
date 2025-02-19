package db

import (
	"time"
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
	err := DB.Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, convertGormError(err)
	}
	return user, nil
}

func (dbu *DbUserDriver) UpdateUser(user *User, updateData map[string]interface{}) error {
	err := DB.Model(&user).Updates(updateData).Error
	if err != nil {
		return convertGormError(err)
	}
	return nil
}

func (dbu *DbUserDriver) FindById(id string) (*User, error) {
	var user *User
	err := DB.First(&user, "id = ?", id).Error
	if err != nil {
		return nil, convertGormError(err)
	}
	return user, nil
}
