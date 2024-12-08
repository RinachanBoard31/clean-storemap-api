package controller

import (
	"bytes"
	"clean-storemap-api/src/adapter/gateway"
	db "clean-storemap-api/src/driver/db"
	model "clean-storemap-api/src/entity"
	"clean-storemap-api/src/usecase/port"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockUserDriverFactory struct {
	mock.Mock
}

type MockGoogleOAuthDriverFactory struct {
	mock.Mock
}

type MockJwtDriverFactory struct {
	mock.Mock
}

type MockUserOutputFactoryFuncObject struct {
	mock.Mock
}

type MockUserRepositoryFactoryFuncObject struct {
	mock.Mock
}

type MockUserInputFactoryFuncObject struct {
	mock.Mock
}

func (m *MockUserDriverFactory) CreateUser(*db.User) (*db.User, error) {
	args := m.Called()
	return args.Get(0).(*db.User), args.Error(1)
}

func (m *MockUserDriverFactory) UpdateUser(*db.User, map[string]interface{}) error {
	args := m.Called()
	return args.Error(0)
}
func (m *MockUserDriverFactory) FindById(string) (*db.User, error) {
	args := m.Called()
	return args.Get(0).(*db.User), args.Error(1)
}
func (m *MockUserDriverFactory) FindByEmail(string) (*db.User, error) {
	args := m.Called()
	return args.Get(0).(*db.User), args.Error(1)
}

func (m *MockGoogleOAuthDriverFactory) GenerateUrl() string {
	args := m.Called()
	return args.Get(0).(string)
}

func (m *MockGoogleOAuthDriverFactory) GetEmail(string) (string, error) {
	args := m.Called()
	return args.Get(0).(string), args.Error(1)
}

func (m *MockJwtDriverFactory) GenerateToken(subject string) (string, error) {
	args := m.Called(subject)
	return args.Get(0).(string), args.Error(1)
}

func (m *MockUserOutputFactoryFuncObject) OutputUpdateResult() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockUserOutputFactoryFuncObject) OutputAuthUrl(url string) error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockUserOutputFactoryFuncObject) OutputLoginResult(string) error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockUserOutputFactoryFuncObject) OutputSignupWithAuth(string) error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockUserOutputFactoryFuncObject) OutputAlreadySignedup() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockUserOutputFactoryFuncObject) OutputHasEmailInRequestBody() error {
	args := m.Called()
	return args.Error(0)
}

func mockUserOutputFactoryFunc(c echo.Context) port.UserOutputPort {
	return &MockUserOutputFactoryFuncObject{}
}

func (m *MockUserRepositoryFactoryFuncObject) Create(*model.User) (*model.User, error) {
	args := m.Called()
	return args.Get(0).(*model.User), args.Error(1)
}

func (m *MockUserRepositoryFactoryFuncObject) Exist(*model.User) error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockUserRepositoryFactoryFuncObject) Update(*model.User, model.ChangeForUser) error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockUserRepositoryFactoryFuncObject) Get(string) (*model.User, error) {
	args := m.Called()
	return args.Get(0).(*model.User), args.Error(1)
}

func (m *MockUserRepositoryFactoryFuncObject) FindBy(*model.UserCredentials) (*model.User, error) {
	args := m.Called()
	return args.Get(0).(*model.User), args.Error(1)
}

func (m *MockUserRepositoryFactoryFuncObject) GenerateAuthUrl() string {
	args := m.Called()
	return args.Get(0).(string)
}

func (m *MockUserRepositoryFactoryFuncObject) GetUserInfoWithAuthCode(string) (string, error) {
	args := m.Called()
	return args.Get(0).(string), args.Error(1)
}

func (m *MockUserRepositoryFactoryFuncObject) GenerateAccessToken(id string) (string, error) {
	args := m.Called(id)
	return args.Get(0).(string), args.Error(1)
}

func mockUserRepositoryFactoryFunc(userDriver gateway.UserDriver, googleOAuthDriver gateway.GoogleOAuthDriver, jwtDriver gateway.JwtDriver) port.UserRepository {
	return &MockUserRepositoryFactoryFuncObject{}
}

func (m *MockUserInputFactoryFuncObject) UpdateUser(string, model.ChangeForUser) error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockUserInputFactoryFuncObject) GetAuthUrl() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockUserInputFactoryFuncObject) LoginUser(*model.UserCredentials) error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockUserInputFactoryFuncObject) SignupDraft(string) error {
	args := m.Called()
	return args.Error(0)
}

func TestUpdateUser(t *testing.T) {
	/* Arrange */
	c, rec := newRouter()
	userId := "id_1"
	var expected error = nil
	reqBody := `{"name":"test","age":10,"sex":0.4, "gender":0}`
	req := httptest.NewRequest(http.MethodPut, "/user", bytes.NewBufferString(reqBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	c.Set("userId", userId)
	c.SetRequest(req)

	// Driver用
	mockUserDriverFactory := new(MockUserDriverFactory)

	uc := &UserController{
		userDriverFactory:     mockUserDriverFactory,
		userOutputFactory:     mockUserOutputFactoryFunc,
		userRepositoryFactory: mockUserRepositoryFactoryFunc,
	}

	mockUserInputFactoryFuncObject := new(MockUserInputFactoryFuncObject)
	mockUserInputFactoryFuncObject.On("UpdateUser").Return(nil)
	uc.userInputFactory = func(repository port.UserRepository, output port.UserOutputPort) port.UserInputPort {
		return mockUserInputFactoryFuncObject
	}

	/* Act */
	actual := uc.UpdateUser(c)

	/* Assert */
	assert.Equal(t, expected, actual)
	assert.Equal(t, http.StatusOK, rec.Code)
	mockUserInputFactoryFuncObject.AssertNumberOfCalls(t, "UpdateUser", 1)
}

func TestLoginUser(t *testing.T) {
	/* Arrange */
	c, rec := newRouter()
	var expected error = nil
	// デフォルトでリクエストメソッドがGETのため、POSTに変更。こういうPOSTリクエストが来たことにする
	reqBody := `{"email":"johnathan@example.com"}`
	req := httptest.NewRequest(http.MethodPost, "/user/login", bytes.NewBufferString(reqBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	c.SetRequest(req)

	// Driverだけは実体が必要
	mockUserDriverFactory := new(MockUserDriverFactory)
	mockUserDriverFactory.On("FindByEmail").Return(nil, nil)

	// InputPortのLoginUserのモックを作成
	uc := &UserController{
		userDriverFactory:     mockUserDriverFactory,
		userOutputFactory:     mockUserOutputFactoryFunc,
		userRepositoryFactory: mockUserRepositoryFactoryFunc,
	}

	// newUserInputPort.LoginUser()をするためには、LoginUser()を持つmockUserInputFactoryFuncObjectがuserInputFactoryに必要だから無名関数でreturnする必要があった
	mockUserInputFactoryFuncObject := new(MockUserInputFactoryFuncObject)
	mockUserInputFactoryFuncObject.On("LoginUser").Return(expected)
	uc.userInputFactory = func(repository port.UserRepository, output port.UserOutputPort) port.UserInputPort {
		return mockUserInputFactoryFuncObject
	}

	/* Act */
	actual := uc.LoginUser(c)

	/* Assert */
	// uc.LoginUser()がUserInputPort.LoginUser()を返すこと
	assert.Equal(t, expected, actual)
	// echoが正しく起動したか
	assert.Equal(t, http.StatusOK, rec.Code)
	// InputPortのLoginUser()が1回呼ばれること
	mockUserInputFactoryFuncObject.AssertNumberOfCalls(t, "LoginUser", 1)
}

func TestGetAuthUrl(t *testing.T) {
	/* Arrange */
	c, _ := newRouter()
	url := "https://www.google.com"
	var expected error = nil

	// Driverだけは実体が必要
	mockGoogleOAuthDriverFactory := new(MockGoogleOAuthDriverFactory)
	mockGoogleOAuthDriverFactory.On("GenerateUrl").Return(url)

	// InputPortのGetGoogleAuthUrlのモックを作成
	uc := &UserController{
		googleOAuthDriverFactory: mockGoogleOAuthDriverFactory,
		userOutputFactory:        mockUserOutputFactoryFunc,
		userRepositoryFactory:    mockUserRepositoryFactoryFunc,
	}

	// newUserInputPort.GetAuthUrl()をするためには、GetAuthUrl()を持つmockUserInputFactoryFuncObjectがuserInputFactoryに必要だから無名関数でreturnする必要があった
	mockUserInputFactoryFuncObject := new(MockUserInputFactoryFuncObject)
	mockUserInputFactoryFuncObject.On("GetAuthUrl").Return(nil)
	uc.userInputFactory = func(repository port.UserRepository, output port.UserOutputPort) port.UserInputPort {
		return mockUserInputFactoryFuncObject
	}

	/* Act */
	actual := uc.GetAuthUrl(c)

	/* Assert */
	assert.Equal(t, expected, actual)
	mockUserInputFactoryFuncObject.AssertNumberOfCalls(t, "GetAuthUrl", 1)
}

func TestSignupWithAuth(t *testing.T) {
	/* Arrange */
	c, _ := newRouter()
	var expected error = nil
	reqBody := `{"code":"123456"}`
	req := httptest.NewRequest(http.MethodPost, "/auth/signup", bytes.NewBufferString(reqBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	c.SetRequest(req)

	// OAuth用(関数が実行されるわけではないので、mockの戻り値を設定しない)
	mockGoogleOAuthDriverFactory := new(MockGoogleOAuthDriverFactory)

	// auth用
	mockJwtDriverFactory := new(MockJwtDriverFactory)

	// DB用(関数が実行されるわけではないので、mockの戻り値を設定しない)
	mockUserDriverFactory := new(MockUserDriverFactory)

	uc := &UserController{
		googleOAuthDriverFactory: mockGoogleOAuthDriverFactory,
		jwtDriverFactory:         mockJwtDriverFactory,
		userDriverFactory:        mockUserDriverFactory,
		userOutputFactory:        mockUserOutputFactoryFunc,
		userRepositoryFactory:    mockUserRepositoryFactoryFunc,
	}

	mockUserInputFactoryFuncObject := new(MockUserInputFactoryFuncObject)
	mockUserInputFactoryFuncObject.On("SignupDraft").Return(nil)
	uc.userInputFactory = func(repository port.UserRepository, output port.UserOutputPort) port.UserInputPort {
		return mockUserInputFactoryFuncObject
	}

	/* Act */
	actual := uc.SignupWithAuth(c)

	/* Assert */
	assert.Equal(t, expected, actual)
	mockUserInputFactoryFuncObject.AssertNumberOfCalls(t, "SignupDraft", 1)
}
