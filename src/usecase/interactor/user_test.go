package interactor

import (
	model "clean-storemap-api/src/entity"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockUserRepository struct {
	mock.Mock
}

type MockUserOutputPort struct {
	mock.Mock
}

func (m *MockUserRepository) Create(user *model.User) (*model.User, error) {
	args := m.Called(user)
	return args.Get(0).(*model.User), args.Error(1)
}

func (m *MockUserRepository) Exist(user *model.User) error {
	args := m.Called(user)
	return args.Error(0)
}
func (m *MockUserRepository) Update(user *model.User, updateData model.ChangeForUser) error {
	args := m.Called(user, updateData)
	return args.Error(0)
}
func (m *MockUserRepository) Get(id string) (*model.User, error) {
	args := m.Called(id)
	return args.Get(0).(*model.User), args.Error(1)
}

func (m *MockUserRepository) GenerateAuthUrl() string {
	args := m.Called()
	return args.Get(0).(string)
}

func (m *MockUserRepository) FindBy(user *model.UserCredentials) (*model.User, error) {
	args := m.Called(user)
	return args.Get(0).(*model.User), args.Error(1)
}

func (m *MockUserRepository) GetUserInfoWithAuthCode(code string) (string, error) {
	args := m.Called(code)
	return args.Get(0).(string), args.Error(1)
}

func (m *MockUserRepository) GenerateAccessToken(id string) (string, error) {
	args := m.Called(id)
	return args.Get(0).(string), args.Error(1)
}

func (m *MockUserOutputPort) OutputUpdateResult() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockUserOutputPort) OutputLoginResult(token string) error {
	args := m.Called(token)
	return args.Error(0)
}

func (m *MockUserOutputPort) OutputAuthUrl(url string) error {
	args := m.Called(url)
	return args.Error(0)
}

func (m *MockUserOutputPort) OutputSignupWithAuth(id string) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockUserOutputPort) OutputAlreadySignedup() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockUserOutputPort) OutputHasEmailInRequestBody() error {
	args := m.Called()
	return args.Error(0)
}

func TestUpdateUser(t *testing.T) {
	/* Arrange */
	var expected error = nil
	id := "id_1"
	existUser := &model.User{
		Id:     id,
		Name:   "sample",
		Age:    10,
		Sex:    0.1,
		Gender: -0.1,
	}
	updateData := model.ChangeForUser{"name": "sample2", "sex": 1.0, "gender": -1.0}

	mockUserRepository := new(MockUserRepository)
	mockUserOutputPort := new(MockUserOutputPort)
	mockUserOutputPort.On("OutputHasEmailInRequestBody").Return(nil)
	mockUserRepository.On("Get", id).Return(existUser, nil)
	mockUserRepository.On("Update", existUser, updateData).Return(nil)
	mockUserOutputPort.On("OutputUpdateResult").Return(nil)

	ui := &UserInteractor{userRepository: mockUserRepository, userOutputPort: mockUserOutputPort}

	/* Act */
	actual := ui.UpdateUser(id, updateData)

	/* Assert */
	assert.Equal(t, expected, actual)
	mockUserRepository.AssertNumberOfCalls(t, "Get", 1)
	mockUserRepository.AssertNumberOfCalls(t, "Update", 1)
	mockUserOutputPort.AssertNumberOfCalls(t, "OutputUpdateResult", 1)
}

func TestLoginUser(t *testing.T) {
	/* Arrange */
	var expected error = nil
	userCredentials := &model.UserCredentials{Email: "test@example.com"}
	token := "test_token"
	user := &model.User{
		Id:     "id_1",
		Email:  userCredentials.Email,
		Age:    52,
		Sex:    -0.2,
		Gender: 1.0,
	}
	mockUserRepository := new(MockUserRepository)
	mockUserRepository.On("FindBy", userCredentials).Return(user, nil)
	mockUserRepository.On("GenerateAccessToken", user.Id).Return(token, nil)
	mockUserOutputPort := new(MockUserOutputPort)
	mockUserOutputPort.On("OutputLoginResult", token).Return(nil)
	ui := &UserInteractor{userRepository: mockUserRepository, userOutputPort: mockUserOutputPort}

	/* Act */
	actual := ui.LoginUser(userCredentials)

	/* Assert */
	assert.Equal(t, expected, actual)
	mockUserRepository.AssertNumberOfCalls(t, "FindBy", 1)
	mockUserRepository.AssertNumberOfCalls(t, "GenerateAccessToken", 1)
	mockUserOutputPort.AssertNumberOfCalls(t, "OutputLoginResult", 1)
}

func TestGetAuthUrl(t *testing.T) {
	/* Arrange */
	url := "https://www.google.com"
	var expected error = nil

	mockUserRepository := new(MockUserRepository)
	mockUserRepository.On("GenerateAuthUrl").Return(url)
	mockUserOutputPort := new(MockUserOutputPort)
	mockUserOutputPort.On("OutputAuthUrl", url).Return(nil)

	ui := &UserInteractor{
		userRepository: mockUserRepository,
		userOutputPort: mockUserOutputPort,
	}

	/* Act */
	actual := ui.GetAuthUrl()

	/* Assert */
	assert.Equal(t, expected, actual)
	mockUserRepository.AssertNumberOfCalls(t, "GenerateAuthUrl", 1)
	mockUserOutputPort.AssertNumberOfCalls(t, "OutputAuthUrl", 1)
}

func TestSignupDraft(t *testing.T) {
	/* Arrange */
	code := ""
	email := "sample@example.com"
	var expected error = nil
	err := errors.New("user is not found")

	draftUser := &model.User{
		Name:   "",
		Email:  email,
		Age:    0,
		Sex:    0.0,
		Gender: 0.0,
	}
	createdUser := &model.User{
		Id:     "id_1",
		Name:   "",
		Email:  email,
		Age:    0,
		Sex:    0.0,
		Gender: 0.0,
	}
	token := "token"

	mockUserRepository := new(MockUserRepository)
	mockUserRepository.On("GetUserInfoWithAuthCode", code).Return(email, nil)
	mockUserRepository.On("Exist", draftUser).Return(err) // 存在していない場合にエラーが返る
	mockUserRepository.On("Create", draftUser).Return(createdUser, nil)
	mockUserRepository.On("GenerateAccessToken", createdUser.Id).Return(token, nil)
	mockUserOutputPort := new(MockUserOutputPort)
	mockUserOutputPort.On("OutputAlreadySignedup").Return(nil)
	mockUserOutputPort.On("OutputSignupWithAuth", token).Return(nil)

	ui := &UserInteractor{
		userRepository: mockUserRepository,
		userOutputPort: mockUserOutputPort,
	}

	/* Act */
	actual := ui.SignupDraft(code)

	/* Assert */
	assert.Equal(t, expected, actual)
	mockUserRepository.AssertNumberOfCalls(t, "GetUserInfoWithAuthCode", 1)
	mockUserRepository.AssertNumberOfCalls(t, "Exist", 1)
	mockUserRepository.AssertNumberOfCalls(t, "Create", 1)
	mockUserRepository.AssertNumberOfCalls(t, "GenerateAccessToken", 1)
	mockUserOutputPort.AssertNumberOfCalls(t, "OutputSignupWithAuth", 1)
}
