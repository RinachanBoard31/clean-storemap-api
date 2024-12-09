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

func (m *MockUserRepository) CreateUser(dbUser *db.User) (*db.User, error) {
	args := m.Called(dbUser)
	return args.Get(0).(*db.User), args.Error(1)
}

func (m *MockUserRepository) UpdateUser(*db.User, map[string]interface{}) error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockUserRepository) FindById(id string) (*db.User, error) {
	args := m.Called(id)
	return args.Get(0).(*db.User), args.Error(1)
}

func (m *MockUserRepository) FindByEmail(email string) (*db.User, error) {
	args := m.Called(email)
	return args.Get(0).(*db.User), args.Error(1)
}
func (m *MockUserRepository) GenerateUrl(actionType string) string {
	args := m.Called(actionType)
	return args.Get(0).(string)
}

func (m *MockUserRepository) GetEmail(code string, actionType string) (string, error) {
	args := m.Called(code, actionType)
	return args.Get(0).(string), args.Error(1)
}

type MockJwtRepository struct {
	mock.Mock
}

func (m *MockJwtRepository) GenerateToken(subject string) (string, error) {
	args := m.Called(subject)
	return args.Get(0).(string), args.Error(1)
}

func TestCreate(t *testing.T) {
	/* Arrange */
	user := &model.User{
		Name:   "noiman",
		Email:  "noiman@groovex.co.jp",
		Age:    35,
		Sex:    1.0,
		Gender: -0.5,
	}
	dbUser := &db.User{
		Id:     "id_1",
		Name:   user.Name,
		Email:  user.Email,
		Age:    user.Age,
		Sex:    user.Sex,
		Gender: user.Gender,
	}
	expected := user

	mockUserRepository := new(MockUserRepository)
	mockUserRepository.On("CreateUser", mock.MatchedBy(func(dbUser *db.User) bool {
		return dbUser.Name == user.Name &&
			dbUser.Email == user.Email &&
			dbUser.Age == user.Age &&
			dbUser.Sex == user.Sex &&
			dbUser.Gender == user.Gender
	})).Return(dbUser, nil)
	ug := &UserGateway{userDriver: mockUserRepository}

	/* Act */
	actual, _ := ug.Create(user)

	/* Assert */
	// 返り値が正しいこと
	assert.Equal(t, expected, actual)
	// userDriver.CreateUser()が1回呼ばれること
	mockUserRepository.AssertNumberOfCalls(t, "CreateUser", 1)
}

func TestExist(t *testing.T) {
	/* Arrange */
	user := &model.User{
		Name:   "noiman",
		Email:  "noiman@groovex.co.jp",
		Age:    35,
		Sex:    1.0,
		Gender: -0.5,
	}
	var expected error = nil
	foundUser := &db.User{}

	mockUserRepository := new(MockUserRepository)
	mockUserRepository.On("FindByEmail", user.Email).Return(foundUser, nil)
	ug := &UserGateway{userDriver: mockUserRepository}

	/* Act */
	actual := ug.Exist(user)

	/* Assert */
	// 返り値が正しいこと
	assert.Equal(t, expected, actual)
	// userDriver.CreateUser()が1回呼ばれること
	mockUserRepository.AssertNumberOfCalls(t, "FindByEmail", 1)
}
func TestUpdate(t *testing.T) {
	/* Arrange */
	var expected error = nil
	// 更新されるUser
	user := &model.User{
		Id:     "id_1",
		Name:   "sample",
		Email:  "sample@example.com",
		Age:    10,
		Sex:    0.1,
		Gender: -0.1,
	}
	// 更新後のデータ
	updatedUserData := model.ChangeForUser{"name": "sample2", "sex": 1.0, "gender": -1.0}

	mockUserRepository := new(MockUserRepository)
	mockUserRepository.On("UpdateUser").Return(nil)
	ug := &UserGateway{userDriver: mockUserRepository}

	/* Act */
	actual := ug.Update(user, updatedUserData)

	/* Assert */
	assert.Equal(t, expected, actual)
	mockUserRepository.AssertNumberOfCalls(t, "UpdateUser", 1)
}

func TestGet(t *testing.T) {
	/* Arrange */
	id := "id_1"
	dbUser := &db.User{
		Id:     id,
		Name:   "sample",
		Email:  "sample@example.com",
		Age:    20,
		Sex:    1.0,
		Gender: -0.5,
	}

	var expectedError error = nil
	expectedUser := &model.User{
		Id:     dbUser.Id,
		Name:   dbUser.Name,
		Email:  dbUser.Email,
		Age:    dbUser.Age,
		Sex:    dbUser.Sex,
		Gender: dbUser.Gender,
	}

	mockUserRepository := new(MockUserRepository)
	mockUserRepository.On("FindById", id).Return(dbUser, nil)
	ug := &UserGateway{userDriver: mockUserRepository}

	/* Act */
	actual, actualErr := ug.Get(id)

	/* Assert */
	// 返り値が正しいこと
	assert.Equal(t, expectedUser, actual)
	assert.Equal(t, expectedError, actualErr)
	// userDriver.CreateUser()が1回呼ばれること
	mockUserRepository.AssertNumberOfCalls(t, "FindById", 1)
}

func TestFindBy(t *testing.T) {
	/* Arrange */
	query := &model.UserQuery{
		Id:    "Id_1",
		Email: "noiman@groovex.co.jp",
	}
	dbUser := &db.User{
		Id:     query.Id,
		Name:   "noiman",
		Email:  query.Email,
		Age:    35,
		Sex:    1.0,
		Gender: -0.5,
	}
	user := &model.User{
		Id:     dbUser.Id,
		Name:   dbUser.Name,
		Email:  dbUser.Email,
		Age:    dbUser.Age,
		Sex:    dbUser.Sex,
		Gender: dbUser.Gender,
	}
	expected := user
	mockUserRepository := new(MockUserRepository)
	mockUserRepository.On("FindById", query.Id).Return(&db.User{}, nil)
	mockUserRepository.On("FindByEmail", query.Email).Return(dbUser, nil)
	ug := &UserGateway{userDriver: mockUserRepository}

	/* Act */
	actual, _ := ug.FindBy(query)

	/* Assert */
	// 返り値が正しいこと
	assert.Equal(t, expected, actual)
	mockUserRepository.AssertNumberOfCalls(t, "FindById", 1)
	mockUserRepository.AssertNumberOfCalls(t, "FindByEmail", 1)
}

func TestGenerateAuthUrl(t *testing.T) {
	/* Arrange */
	actionType := "signup" // もしくはlogin
	generatedUrl := "https://www.google.com"
	expected := generatedUrl
	mockUserRepository := new(MockUserRepository)
	mockUserRepository.On("GenerateUrl", actionType).Return(generatedUrl)
	ug := &UserGateway{
		googleOAuthDriver: mockUserRepository,
	}

	/* Act */
	actual := ug.GenerateAuthUrl(actionType)

	/* Assert */
	assert.Equal(t, expected, actual)
	mockUserRepository.AssertNumberOfCalls(t, "GenerateUrl", 1)
}

func TestGetUserInfoWithAuthCode(t *testing.T) {
	/* Arrange */
	code := ""
	actionType := "signup"
	email := "sample@example.com"
	expected := email
	mockUserRepository := new(MockUserRepository)
	mockUserRepository.On("GetEmail", code, actionType).Return(email, nil)
	ug := &UserGateway{
		googleOAuthDriver: mockUserRepository,
	}

	/* Act */
	actual, _ := ug.GetUserInfoWithAuthCode(code, actionType)

	/* Assert */
	assert.Equal(t, expected, actual)
	mockUserRepository.AssertNumberOfCalls(t, "GetEmail", 1)
}

func TestGenerateAccessToken(t *testing.T) {
	/* Arrange */
	id := "Id001"
	token := "token"
	expected := token
	MockJwtRepository := new(MockJwtRepository)
	MockJwtRepository.On("GenerateToken", id).Return(token, nil)
	ug := &UserGateway{
		jwtDriver: MockJwtRepository,
	}

	/* Act */
	actual, _ := ug.GenerateAccessToken(id)

	/* Assert */
	assert.Equal(t, expected, actual)
	MockJwtRepository.AssertNumberOfCalls(t, "GenerateToken", 1)
}
