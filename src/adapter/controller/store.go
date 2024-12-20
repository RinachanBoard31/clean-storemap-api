package controller

import (
	"clean-storemap-api/src/adapter/gateway"
	model "clean-storemap-api/src/entity"
	"clean-storemap-api/src/usecase/port"
	"net/http"

	"github.com/labstack/echo/v4"
	"gopkg.in/go-playground/validator.v9"
)

type StoreRequestBody struct {
	StoreId             string `json:"storeId" validate:"required"`
	StoreName           string `json:"storeName" validate:"required"`
	RegularOpeningHours string `json:"regularOpeningHours"`
	PriceLevel          string `json:"priceLevel"`
	Latitude            string `json:"latitude" validate:"required"`
	Longitude           string `json:"longitude" validate:"required"`
}

type StoreI interface {
	GetStores(c echo.Context) error
	GetNearStores(c echo.Context) error
	GetFavoriteStores(c echo.Context) error
	SaveFavoriteStore(c echo.Context) error
	GetTopFavoriteStores(c echo.Context) error
}

type StoreOutputFactory func(echo.Context) port.StoreOutputPort
type StoreInputFactory func(port.StoreRepository, port.StoreOutputPort) port.StoreInputPort
type StoreRepositoryFactory func(gateway.StoreDriver, gateway.GoogleMapDriver) port.StoreRepository
type StoreDriverFactory gateway.StoreDriver
type GoogleMapDriverFactory gateway.GoogleMapDriver

type StoreController struct {
	storeDriverFactory     StoreDriverFactory
	googleMapDriverFactory GoogleMapDriverFactory
	storeOutputFactory     StoreOutputFactory
	storeInputFactory      StoreInputFactory
	storeRepositoryFactory StoreRepositoryFactory
}

func NewStoreController(
	storeDriverFactory StoreDriverFactory,
	googleMapDriverFactory GoogleMapDriverFactory,
	storeOutputFactory StoreOutputFactory,
	storeInputFactory StoreInputFactory,
	storeRepositoryFactory StoreRepositoryFactory,
) StoreI {
	return &StoreController{
		storeDriverFactory:     storeDriverFactory,
		googleMapDriverFactory: googleMapDriverFactory,
		storeOutputFactory:     storeOutputFactory,
		storeInputFactory:      storeInputFactory,
		storeRepositoryFactory: storeRepositoryFactory,
	}
}

func (sc *StoreController) GetStores(c echo.Context) error {
	return sc.newStoreInputPort(c).GetStores()
}

func (sc *StoreController) GetNearStores(c echo.Context) error {
	return sc.newStoreInputPort(c).GetNearStores()
}

func (sc *StoreController) GetFavoriteStores(c echo.Context) error {
	userId := c.Get("userId").(string)
	if userId == "" {
		return c.JSON(http.StatusBadRequest, "user_id is required")
	}
	return sc.newStoreInputPort(c).GetFavoriteStores(userId)
}

func (sc *StoreController) SaveFavoriteStore(c echo.Context) error {
	var s StoreRequestBody
	userId := c.Get("userId").(string)
	if userId == "" {
		return c.JSON(http.StatusBadRequest, "user_id is required")
	}

	if err := c.Bind(&s); err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	if err := c.Validate(&s); err != nil {
		return c.JSON(http.StatusInternalServerError, err.(validator.ValidationErrors).Error())
	}
	store, err := model.NewStore(s.StoreId, s.StoreName, s.RegularOpeningHours, s.PriceLevel, s.Latitude, s.Longitude)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return sc.newStoreInputPort(c).SaveFavoriteStore(store, userId)
}

func (sc *StoreController) GetTopFavoriteStores(c echo.Context) error {
	return sc.newStoreInputPort(c).GetTopFavoriteStores()
}

/* ここでpresenterにecho.Contextを渡している！起爆！！！（遅延） */
/* これによって、presenterのinterface(outputport)にecho.Contextを書かなくて良くなる */
func (sc *StoreController) newStoreInputPort(c echo.Context) port.StoreInputPort {
	storeOutputPort := sc.storeOutputFactory(c)
	storeDriver := sc.storeDriverFactory
	googleMapDriver := sc.googleMapDriverFactory
	storeRepository := sc.storeRepositoryFactory(storeDriver, googleMapDriver)
	return sc.storeInputFactory(storeRepository, storeOutputPort)
}
