package gateway

import (
	"clean-storemap-api/src/driver/db"
	model "clean-storemap-api/src/entity"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) CreateUser(*db.User) error {
	args := m.Called()
	return args.Error(0)
}

func TestCreate(t *testing.T) {
	/* Arrange */
	var expected error = nil
	mockUserRepository := new(MockUserRepository)
	mockUserRepository.On("CreateUser").Return(nil)
	ug := &UserGateway{userDriver: mockUserRepository}
	user := &model.User{
		Name:   "noiman",
		Email:  "noiman@groovex.co.jp",
		Age:    35,
		Sex:    1.0,
		Gender: -0.5,
	}

	/* Act */
	actual := ug.Create(user)

	/* Assert */
	// 返り値が正しいこと
	assert.Equal(t, expected, actual)
	// userDriver.CreateUser()が1回呼ばれること
	mockUserRepository.AssertNumberOfCalls(t, "CreateUser", 1)
}
