package presenter

import (
	model "clean-storemap-api/src/entity"
	"clean-storemap-api/src/entity/errors"
	"testing"

	"net/http"
	"net/http/httptest"

	"github.com/stretchr/testify/assert"

	"github.com/labstack/echo/v4"
)

func newRouter() (echo.Context, *httptest.ResponseRecorder) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	return c, rec
}

func TestOutputStoresHandleError(t *testing.T) {
	/* Arrange */
	testCases := []struct {
		input        error
		expectedCode int
		expected     string
	}{
		{input: errors.ErrStoreNotFound, expectedCode: http.StatusNotFound, expected: "{\"error\":\"お店が登録されていません。\"}\n"},
		{input: errors.ErrStoreAlreadyExists, expectedCode: http.StatusConflict, expected: "{\"error\":\"そのお店はすでに登録されています。\"}\n"},
		{input: errors.ErrFavoriteStoreAlreadyExists, expectedCode: http.StatusConflict, expected: "{\"error\":\"そのお店はすでにお気に入りに登録されています。\"}\n"},
		{input: errors.ErrStoreInvalidData, expectedCode: http.StatusBadRequest, expected: "{\"error\":\"不正なデータが含まれています。\"}\n"},
		{input: errors.ErrFailedGetNearByStores, expectedCode: http.StatusInternalServerError, expected: "{\"error\":\"お店の取得に失敗しました。\"}\n"},
		{input: errors.ErrUnExpected, expectedCode: http.StatusInternalServerError, expected: "{\"error\":\"予期しないエラーが発生しました\"}\n"},
	}
	/* Act, Assert */
	for _, tc := range testCases {
		// recが1ループごとに切れないのでrouterは毎回作り直す。
		c, rec := newRouter()
		sp := &StorePresenter{c: c}
		actual := sp.OutputStoresHandleError(tc.input)
		if assert.NoError(t, actual) {
			assert.Equal(t, tc.expectedCode, rec.Code)
			assert.Equal(t, tc.expected, rec.Body.String())
		}
	}
}

func TestOutputAllStores(t *testing.T) {
	/* Arrange */
	expected := "{\"stores\":[{\"id\":\"Id001\",\"name\":\"UEC cafe\",\"regularOpeningHours\":\"Sat: 06:00 - 22:00, Sun: 06:00 - 22:00\",\"priceLevel\":\"PRICE_LEVEL_MODERATE\",\"location\":{\"latitude\":\"35.713\",\"longitude\":\"139.762\"}}]}\n"
	stores := make([]*model.Store, 0)
	stores = append(
		stores,
		&model.Store{
			Id:                  "Id001",
			Name:                "UEC cafe",
			RegularOpeningHours: "Sat: 06:00 - 22:00, Sun: 06:00 - 22:00",
			PriceLevel:          "PRICE_LEVEL_MODERATE",
			Location:            model.Location{Lat: "35.713", Lng: "139.762"},
		},
	)
	c, rec := newRouter()
	sp := &StorePresenter{c: c}

	/* Act */
	actual := sp.OutputAllStores(stores)

	/* Assert */
	// sp.OutputAllStores()がJSONを返すこと
	if assert.NoError(t, actual) {
		assert.Equal(t, expected, rec.Body.String())
	}
}

func TestOutputSaveFavoriteStoreResult(t *testing.T) {
	/* Arrange */
	expected := "{}\n"
	c, rec := newRouter()
	sp := &StorePresenter{c: c}

	/* Act */
	actual := sp.OutputSaveFavoriteStoreResult()

	/* Assert */
	// sp.OutputSaveFavoriteStoreResult()がJSONを返すこ
	if assert.NoError(t, actual) {
		assert.Equal(t, expected, rec.Body.String())
	}
}
