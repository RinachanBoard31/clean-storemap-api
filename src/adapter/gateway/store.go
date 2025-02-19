package gateway

import (
	api "clean-storemap-api/src/driver/api"
	db "clean-storemap-api/src/driver/db"
	model "clean-storemap-api/src/entity"
	"clean-storemap-api/src/usecase/port"
	"fmt"
	"strings"

	"github.com/google/uuid"
)

// Dbの要素を構造体として渡す必要がある。
type StoreGateway struct {
	storeDriver     StoreDriver
	googleMapDriver GoogleMapDriver
}

type StoreDriver interface {
	GetStores() ([]*db.FavoriteStore, error)
	FindFavorite(storeId string, userId string) (*db.FavoriteStore, error)
	FindFavoriteByUser(userId string) ([]*db.FavoriteStore, error)
	SaveStore(*db.FavoriteStore) error
	GetTopStores() ([]*db.FavoriteStore, error)
}

type GoogleMapDriver interface {
	GetStores() ([]*api.Store, error)
}

func NewStoreRepository(storeDriver StoreDriver, googleMapDriver GoogleMapDriver) port.StoreRepository {
	return &StoreGateway{
		storeDriver:     storeDriver,
		googleMapDriver: googleMapDriver,
	}
}

func (sg *StoreGateway) GetAll() ([]*model.Store, error) {
	dbStores, err := sg.storeDriver.GetStores()
	if err != nil {
		return nil, sg.convertToBusinessError(err)
	}
	stores := make([]*model.Store, 0)
	for _, v := range dbStores {
		stores = append(stores, &model.Store{
			Id:                  v.Id,
			Name:                v.StoreName,
			RegularOpeningHours: v.RegularOpeningHours,
			PriceLevel:          v.PriceLevel,
			Location: model.Location{
				Lat: v.Latitude,
				Lng: v.Longitude,
			},
		})
	}
	return stores, nil
}

func (sg *StoreGateway) GetNearStores() ([]*model.Store, error) {
	apiStores, err := sg.googleMapDriver.GetStores()
	if err != nil {
		return nil, sg.convertToBusinessError(err)
	}
	stores := make([]*model.Store, 0)
	for _, v := range apiStores {
		stores = append(stores, &model.Store{
			Id:                  v.Id,
			Name:                v.Name,
			RegularOpeningHours: strings.Join(v.RegularOpeningHours, ", "),
			PriceLevel:          v.PriceLevel,
			Location: model.Location{
				Lat: fmt.Sprintf("%f", v.Location.Lat),
				Lng: fmt.Sprintf("%f", v.Location.Lng),
			},
		})
	}
	return stores, nil
}

func (sg *StoreGateway) ExistFavorite(store *model.Store, userId string) (bool, error) {
	dbStore, err := sg.storeDriver.FindFavorite(store.Id, userId)
	if err != nil {
		return false, sg.convertToBusinessError(err)
	}
	if dbStore == nil {
		return false, nil
	}
	return true, nil
}

func (sg *StoreGateway) GetFavoriteStores(userId string) ([]*model.Store, error) {
	dbStores, err := sg.storeDriver.FindFavoriteByUser(userId)
	if err != nil {
		return nil, sg.convertToBusinessError(err)
	}
	stores := make([]*model.Store, 0)
	for _, v := range dbStores {
		stores = append(stores, &model.Store{
			Id:                  v.StoreId,
			Name:                v.StoreName,
			RegularOpeningHours: v.RegularOpeningHours,
			PriceLevel:          v.PriceLevel,
			Location: model.Location{
				Lat: v.Latitude,
				Lng: v.Longitude,
			},
		})
	}
	return stores, nil
}

func (sg *StoreGateway) SaveFavoriteStore(store *model.Store, userId string) error {
	dbStore := &db.FavoriteStore{
		Id:                  uuid.New().String(),
		UserId:              userId,
		StoreId:             store.Id,
		StoreName:           store.Name,
		RegularOpeningHours: store.RegularOpeningHours,
		PriceLevel:          store.PriceLevel,
		Latitude:            store.Location.Lat,
		Longitude:           store.Location.Lng,
	}

	err := sg.storeDriver.SaveStore(dbStore)
	if err != nil {
		return sg.convertToBusinessError(err)
	}

	return nil
}

func (sg *StoreGateway) GetTopFavoriteStores() ([]*model.Store, error) {
	dbStores, err := sg.storeDriver.GetTopStores()
	if err != nil {
		return nil, sg.convertToBusinessError(err)
	}
	stores := make([]*model.Store, 0)
	for _, v := range dbStores {
		stores = append(stores, &model.Store{
			Id:                  v.StoreId,
			Name:                v.StoreName,
			RegularOpeningHours: v.RegularOpeningHours,
			PriceLevel:          v.PriceLevel,
			Location: model.Location{
				Lat: v.Latitude,
				Lng: v.Longitude,
			},
		})
	}

	return stores, nil
}
