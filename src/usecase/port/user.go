package port

import (
	model "clean-storemap-api/src/entity"
)

type UserInputPort interface {
	CreateUser(*model.User) error
	LoginUser(*model.UserCredentials) error
	GetAuthUrl() error
}

type UserRepository interface {
	Create(*model.User) error
	FindBy(*model.UserCredentials) error
	GenerateAuthUrl() string
}

type UserOutputPort interface {
	OutputCreateResult() error
	OutputLoginResult() error
	OutputAuthUrl(url string) error
}
