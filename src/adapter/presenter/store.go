package presenter

import (
	model "clean-storemap-api/src/entity"
	"clean-storemap-api/src/entity/errors"
	"clean-storemap-api/src/usecase/port"
	"net/http"

	"github.com/labstack/echo/v4"
)

type StorePresenter struct {
	c echo.Context
}

func NewStoreOutputPort(c echo.Context) port.StoreOutputPort {
	return &StorePresenter{c: c}
}

type StoreOutputJson struct {
	Stores []storeForPresenter `json:"stores"`
}

type locationForPresenter struct {
	Latitude  string `json:"latitude"`
	Longitude string `json:"longitude"`
}

type storeForPresenter struct {
	Id                  string               `json:"id"`
	Name                string               `json:"name"`
	RegularOpeningHours string               `json:"regularOpeningHours"`
	PriceLevel          string               `json:"priceLevel"`
	Location            locationForPresenter `json:"location"`
}

func (sp *StorePresenter) OutputStoresHandleError(err error) error {
	if errors.Is(err, errors.ErrStoreNotFound) {
		return sp.c.JSON(http.StatusNotFound, map[string]interface{}{"error": "お店が登録されていません。"})
	} else if errors.Is(err, errors.ErrStoreAlreadyExists) {
		return sp.c.JSON(http.StatusConflict, map[string]interface{}{"error": "そのお店はすでに登録されています。"})
	} else if errors.Is(err, errors.ErrFavoriteStoreAlreadyExists) {
		return sp.c.JSON(http.StatusConflict, map[string]interface{}{"error": "そのお店はすでにお気に入りに登録されています。"})
	} else if errors.Is(err, errors.ErrStoreInvalidData) {
		return sp.c.JSON(http.StatusBadRequest, map[string]interface{}{"error": "不正なデータが含まれています。"})
	} else if errors.Is(err, errors.ErrFailedGetNearByStores) {
		return sp.c.JSON(http.StatusInternalServerError, map[string]interface{}{"error": "お店の取得に失敗しました。"})
	}
	return sp.c.JSON(http.StatusInternalServerError, map[string]interface{}{"error": "予期しないエラーが発生しました"})
}

func (sp *StorePresenter) OutputAllStores(stores []*model.Store) error {
	json_stores := make([]storeForPresenter, 0)
	for _, v := range stores {
		json_stores = append(json_stores, storeForPresenter{
			Id:                  v.Id,
			Name:                v.Name,
			RegularOpeningHours: v.RegularOpeningHours,
			PriceLevel:          v.PriceLevel,
			Location: locationForPresenter{
				Latitude:  v.Location.Lat,
				Longitude: v.Location.Lng,
			},
		})
	}
	output_json := &StoreOutputJson{Stores: json_stores}
	return sp.c.JSON(http.StatusOK, output_json)
}

func (sp *StorePresenter) OutputSaveFavoriteStoreResult() error {
	return sp.c.JSON(http.StatusOK, map[string]interface{}{})
}
