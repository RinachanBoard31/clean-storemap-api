package port

import (
	model "clean-storemap-api/src/entity"
)

type StoreInputPort interface {
	GetStores() error
	GetNearStores() error
	GetFavoriteStores(userId string) error
	SaveFavoriteStore(store *model.Store, userId string) error
	GetTopFavoriteStores() error
}

type StoreRepository interface {
	GetAll() ([]*model.Store, error)
	GetNearStores() ([]*model.Store, error)
	ExistFavorite(store *model.Store, userId string) (bool, error)
	GetFavoriteStores(userId string) ([]*model.Store, error)
	SaveFavoriteStore(store *model.Store, userId string) error
	GetTopFavoriteStores() ([]*model.Store, error)
}

type StoreOutputPort interface {
	OutputStoresHandleError(error) error
	OutputAllStores([]*model.Store) error
	OutputSaveFavoriteStoreResult() error
}
