package controller

import (
	"bytes"
	"clean-storemap-api/src/adapter/gateway"
	api "clean-storemap-api/src/driver/api"
	db "clean-storemap-api/src/driver/db"
	model "clean-storemap-api/src/entity"
	"clean-storemap-api/src/usecase/port"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"gopkg.in/go-playground/validator.v9"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockStoreDriverFactory struct {
	mock.Mock
}

type MockGoogleMapDriverFactory struct {
	mock.Mock
}

type MockStoreOutputFactoryFuncObject struct {
	mock.Mock
}

type MockStoreRepositoryFactoryFuncObject struct {
	mock.Mock
}

type MockStoreInputFactoryFuncObject struct {
	mock.Mock
}

func (m *MockStoreDriverFactory) GetStores() ([]*db.FavoriteStore, error) {
	args := m.Called()
	return args.Get(0).([]*db.FavoriteStore), args.Error(1)
}

func (m *MockStoreDriverFactory) FindFavorite(string, string) (*db.FavoriteStore, error) {
	args := m.Called()
	return args.Get(0).(*db.FavoriteStore), args.Error(1)
}

func (m *MockStoreDriverFactory) FindFavoriteByUser(string) ([]*db.FavoriteStore, error) {
	args := m.Called()
	return args.Get(0).([]*db.FavoriteStore), args.Error(1)
}

func (m *MockStoreDriverFactory) SaveStore(*db.FavoriteStore) error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockStoreDriverFactory) GetTopStores() ([]*db.FavoriteStore, error) {
	args := m.Called()
	return args.Get(0).([]*db.FavoriteStore), args.Error(1)
}

func (m *MockGoogleMapDriverFactory) GetStores() ([]*api.Store, error) {
	args := m.Called()
	return args.Get(0).([]*api.Store), args.Error(1)
}

func (m *MockStoreOutputFactoryFuncObject) OutputStoresHandleError(err error) error {
	args := m.Called(err)
	return args.Error(0)
}

func (m *MockStoreOutputFactoryFuncObject) OutputAllStores([]*model.Store) error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockStoreOutputFactoryFuncObject) OutputSaveFavoriteStoreResult() error {
	args := m.Called()
	return args.Error(0)
}

func mockStoreOutputFactoryFunc(c echo.Context) port.StoreOutputPort {
	return &MockStoreOutputFactoryFuncObject{}
}

func (m *MockStoreRepositoryFactoryFuncObject) GetAll() ([]*model.Store, error) {
	args := m.Called()
	return args.Get(0).([]*model.Store), args.Error(1)
}

func (m *MockStoreRepositoryFactoryFuncObject) GetNearStores() ([]*model.Store, error) {
	args := m.Called()
	return args.Get(0).([]*model.Store), args.Error(1)
}

func (m *MockStoreRepositoryFactoryFuncObject) ExistFavorite(store *model.Store, userId string) (bool, error) {
	args := m.Called(store, userId)
	return args.Get(0).(bool), args.Error(1)
}

func (m *MockStoreRepositoryFactoryFuncObject) GetFavoriteStores(userId string) ([]*model.Store, error) {
	args := m.Called()
	return args.Get(0).([]*model.Store), args.Error(1)
}

func (m *MockStoreRepositoryFactoryFuncObject) SaveFavoriteStore(store *model.Store, userId string) error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockStoreRepositoryFactoryFuncObject) GetTopFavoriteStores() ([]*model.Store, error) {
	args := m.Called()
	return args.Get(0).([]*model.Store), args.Error(1)
}

func mockStoreRepositoryFactoryFunc(storeDriver gateway.StoreDriver, googleMapDriver gateway.GoogleMapDriver) port.StoreRepository {
	return &MockStoreRepositoryFactoryFuncObject{}
}

func (m *MockStoreInputFactoryFuncObject) GetStores() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockStoreInputFactoryFuncObject) GetNearStores() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockStoreInputFactoryFuncObject) GetFavoriteStores(userId string) error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockStoreInputFactoryFuncObject) SaveFavoriteStore(store *model.Store, userId string) error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockStoreInputFactoryFuncObject) GetTopFavoriteStores() error {
	args := m.Called()
	return args.Error(0)
}

// Validationのために必要なメソッド
type CustomValidator struct {
	validator *validator.Validate
}

func NewValidator() echo.Validator {
	return &CustomValidator{validator: validator.New()}
}

func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.validator.Struct(i)
}

func newRouter() (echo.Context, *httptest.ResponseRecorder) {
	e := echo.New()
	e.Validator = NewValidator()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	return e.NewContext(req, rec), rec
}

func TestGetStores(t *testing.T) {
	/* Arrange */
	c, rec := newRouter()
	expected := errors.New("")

	// Driverだけは実体が必要
	mockStoreDriverFactory := new(MockStoreDriverFactory)
	mockStoreDriverFactory.On("GetStores").Return([]*db.FavoriteStore{}, nil)
	mockStoreDriverFactory.On("SaveStore").Return(nil)

	// InputPortのGetStoresのモックを作成
	sc := &StoreController{
		storeDriverFactory:     mockStoreDriverFactory,
		storeOutputFactory:     mockStoreOutputFactoryFunc,
		storeRepositoryFactory: mockStoreRepositoryFactoryFunc,
	}

	// newStoreInputPort.GetStores()をするためには、GetStores()を持つmockStoreInputFactoryFuncObjectがstoreInputFactoryに必要だから無名関数でreturnする必要があった
	mockStoreInputFactoryFuncObject := new(MockStoreInputFactoryFuncObject)
	mockStoreInputFactoryFuncObject.On("GetStores").Return(expected)
	sc.storeInputFactory = func(repository port.StoreRepository, output port.StoreOutputPort) port.StoreInputPort {
		return mockStoreInputFactoryFuncObject
	}

	/* Act */
	actual := sc.GetStores(c)

	/* Assert */
	// sc.GetStores()がstoreInputPort.GetStores()を返すこと
	assert.Equal(t, expected, actual)
	// echoが正しく起動したか
	assert.Equal(t, http.StatusOK, rec.Code)
	// InputPortのGetStoresが1回呼ばれること
	mockStoreInputFactoryFuncObject.AssertNumberOfCalls(t, "GetStores", 1)
}

func TestGetNearStores(t *testing.T) {
	/* Arrange */
	c, rec := newRouter()
	expected := errors.New("")

	mockGoogleMapDriverFactory := new(MockGoogleMapDriverFactory)
	mockGoogleMapDriverFactory.On("GetStores").Return(expected)

	sc := &StoreController{
		googleMapDriverFactory: mockGoogleMapDriverFactory,
		storeOutputFactory:     mockStoreOutputFactoryFunc,
		storeRepositoryFactory: mockStoreRepositoryFactoryFunc,
	}

	mockStoreInputFactoryFuncObject := new(MockStoreInputFactoryFuncObject)
	mockStoreInputFactoryFuncObject.On("GetNearStores").Return(expected)
	sc.storeInputFactory = func(repository port.StoreRepository, output port.StoreOutputPort) port.StoreInputPort {
		return mockStoreInputFactoryFuncObject
	}

	/* Act */
	actual := sc.GetNearStores(c)

	/* Assert */
	assert.Equal(t, expected, actual)
	assert.Equal(t, http.StatusOK, rec.Code)
	mockStoreInputFactoryFuncObject.AssertNumberOfCalls(t, "GetNearStores", 1)
}

func TestGetFavoriteStores(t *testing.T) {
	/* Arrange */
	c, rec := newRouter()
	expected := errors.New("")
	userId := "id_1"
	req := httptest.NewRequest(http.MethodGet, "/user/favorite-store", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	c.Set("userId", userId)
	c.SetRequest(req)

	mockStoreDriverFactory := new(MockStoreDriverFactory)
	mockStoreDriverFactory.On("FindStores").Return(&db.FavoriteStore{}, nil)

	sc := &StoreController{
		storeDriverFactory:     mockStoreDriverFactory,
		storeOutputFactory:     mockStoreOutputFactoryFunc,
		storeRepositoryFactory: mockStoreRepositoryFactoryFunc,
	}

	mockStoreInputFactoryFuncObject := new(MockStoreInputFactoryFuncObject)
	mockStoreInputFactoryFuncObject.On("GetFavoriteStores").Return(expected)
	sc.storeInputFactory = func(repository port.StoreRepository, output port.StoreOutputPort) port.StoreInputPort {
		return mockStoreInputFactoryFuncObject
	}

	/* Act */
	actual := sc.GetFavoriteStores(c)

	/* Assert */
	assert.Equal(t, expected, actual)
	assert.Equal(t, http.StatusOK, rec.Code)
	mockStoreInputFactoryFuncObject.AssertNumberOfCalls(t, "GetFavoriteStores", 1)
}

func TestFavoriteSaveStore(t *testing.T) {
	/* Arrange */
	var expected error = nil
	c, rec := newRouter()
	userId := "id_1"
	reqBody := `{
		"storeId": "Id001",
		"storeName": "UEC cafe",
		"regularOpeningHours": "Sat: 06:00 - 22:00, Sun: 06:00 - 22:00",
		"priceLevel": "PRICE_LEVEL_MODERATE",
		"latitude": "35.713",
		"longitude": "139.762"
	}`

	req := httptest.NewRequest(http.MethodPost, "/user/favorite-store", bytes.NewBufferString(reqBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	c.SetRequest(req)
	c.Set("userId", userId)

	mockStoreDriverFactory := new(MockStoreDriverFactory)
	mockStoreDriverFactory.On("GetStores").Return([]*db.FavoriteStore{}, nil)
	mockStoreDriverFactory.On("SaveStore").Return(nil)

	sc := &StoreController{
		storeDriverFactory:     mockStoreDriverFactory,
		storeOutputFactory:     mockStoreOutputFactoryFunc,
		storeRepositoryFactory: mockStoreRepositoryFactoryFunc,
	}

	mockStoreInputFactoryFuncObject := new(MockStoreInputFactoryFuncObject)
	mockStoreInputFactoryFuncObject.On("SaveFavoriteStore").Return(nil)
	sc.storeInputFactory = func(repository port.StoreRepository, output port.StoreOutputPort) port.StoreInputPort {
		return mockStoreInputFactoryFuncObject
	}

	/* Act */
	actual := sc.SaveFavoriteStore(c)

	/* Assert */
	assert.Equal(t, expected, actual)
	assert.Equal(t, http.StatusOK, rec.Code)
	mockStoreInputFactoryFuncObject.AssertNumberOfCalls(t, "SaveFavoriteStore", 1)
}

func TestGetTopFavoriteStores(t *testing.T) {
	/* Arrange */
	var expected error = nil
	c, rec := newRouter()

	mockStoreDriverFactory := new(MockStoreDriverFactory)
	mockStoreDriverFactory.On("GetTopStores").Return([]*db.FavoriteStore{}, nil)

	sc := &StoreController{
		storeDriverFactory:     mockStoreDriverFactory,
		storeOutputFactory:     mockStoreOutputFactoryFunc,
		storeRepositoryFactory: mockStoreRepositoryFactoryFunc,
	}

	mockStoreInputFactoryFuncObject := new(MockStoreInputFactoryFuncObject)
	mockStoreInputFactoryFuncObject.On("GetTopFavoriteStores").Return(nil)
	sc.storeInputFactory = func(repository port.StoreRepository, output port.StoreOutputPort) port.StoreInputPort {
		return mockStoreInputFactoryFuncObject
	}

	/* Act */
	actual := sc.GetTopFavoriteStores(c)

	/* Assert */
	assert.Equal(t, expected, actual)
	assert.Equal(t, http.StatusOK, rec.Code)
	mockStoreInputFactoryFuncObject.AssertNumberOfCalls(t, "GetTopFavoriteStores", 1)
}
