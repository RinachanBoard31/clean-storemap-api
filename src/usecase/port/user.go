package port

import (
	model "clean-storemap-api/src/entity"
)

type UserInputPort interface {
	UpdateUser(string, model.ChangeForUser) error
	LoginUser(string) error
	GetAuthUrl(string) error
	SignupDraft(string) error
}

type UserRepository interface {
	Exist(*model.User) error
	Create(*model.User) (*model.User, error)
	Update(*model.User, model.ChangeForUser) error
	Get(string) (*model.User, error)
	FindBy(*model.UserQuery) (*model.User, error)
	GenerateAuthUrl(string) string
	GetUserInfoWithAuthCode(string, string) (string, error)
	GenerateAccessToken(string) (string, error)
}

type UserOutputPort interface {
	OutputUpdateResult() error
	OutputAuthUrl(string) error
	OutputLoginWithAuth(string) error
	OutputNotRegistered() error
	OutputSignupWithAuth(string) error
	OutputAlreadySignedup() error
	OutputHasEmailInRequestBody() error
}
