package controller

import (
	"bytes"
	"clean-storemap-api/src/adapter/gateway"
	"clean-storemap-api/src/driver/db"
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

type MockUserInputPort struct {
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

func (m *MockUserDriverFactory) CreateUser(*db.User) error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockUserOutputFactoryFuncObject) OutputCreateResult() error {
	args := m.Called()
	return args.Error(0)
}

func mockUserOutputFactoryFunc(c echo.Context) port.UserOutputPort {
	return &MockUserOutputFactoryFuncObject{}
}

func (m *MockUserRepositoryFactoryFuncObject) Create(*model.User) error {
	args := m.Called()
	return args.Error(0)
}

func mockUserRepositoryFactoryFunc(userDriver gateway.UserDriver) port.UserRepository {
	return &MockUserRepositoryFactoryFuncObject{}
}

func (m *MockUserInputFactoryFuncObject) CreateUser(*model.User) error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockUserInputPort) CreateUser(*model.User) error {
	args := m.Called()
	return args.Error(0)
}

func TestCreateUser(t *testing.T) {
	/* Arrange */
	c, rec := newRouter()
	var expected error = nil
	// デフォルトでリクエストメソッドがGETのため、POSTに変更。こういうPOSTリクエストが来たことにする
	reqBody := `{"name":"noiman","email":"noiman@groovex.co.jp","age":10,"sex":0.4,"gender":-0.3}`
	req := httptest.NewRequest(http.MethodPost, "/user", bytes.NewBufferString(reqBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	c.SetRequest(req)

	// Driverだけは実体が必要
	mockUserDriverFactory := new(MockUserDriverFactory)
	mockUserDriverFactory.On("CreateUser").Return(nil)

	// InputPortのCreateUserのモックを作成
	uc := &UserController{
		userDriverFactory:     mockUserDriverFactory,
		userOutputFactory:     mockUserOutputFactoryFunc,
		userRepositoryFactory: mockUserRepositoryFactoryFunc,
	}

	// newUserInputPort.CreateUser()をするためには、CreateUser()を持つmockUserInputFactoryFuncObjectがuserInputFactoryに必要だから無名関数でreturnする必要があった
	mockUserInputFactoryFuncObject := new(MockUserInputFactoryFuncObject)
	mockUserInputFactoryFuncObject.On("CreateUser").Return(expected)
	uc.userInputFactory = func(repository port.UserRepository, output port.UserOutputPort) port.UserInputPort {
		return mockUserInputFactoryFuncObject
	}

	/* Act */
	actual := uc.CreateUser(c)

	/* Assert */
	// uc.CreateUser()がUserInputPort.CreateUser()を返すこと
	assert.Equal(t, expected, actual)
	// echoが正しく起動したか
	assert.Equal(t, http.StatusOK, rec.Code)
	// InputPortのCreateUserが1回呼ばれること
	mockUserInputFactoryFuncObject.AssertNumberOfCalls(t, "CreateUser", 1)
}
