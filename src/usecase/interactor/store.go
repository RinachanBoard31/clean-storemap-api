package interactor

import (
	model "clean-storemap-api/src/entity"
	"clean-storemap-api/src/entity/errors"
	port "clean-storemap-api/src/usecase/port"
)

type StoreInteractor struct {
	storeRepository port.StoreRepository
	storeOutputPort port.StoreOutputPort
}

func NewStoreInputPort(storeRepository port.StoreRepository, storeOutputPort port.StoreOutputPort) port.StoreInputPort {
	return &StoreInteractor{
		storeRepository: storeRepository,
		storeOutputPort: storeOutputPort,
	}
}

func (si *StoreInteractor) GetStores() error {
	stores, err := si.storeRepository.GetAll()
	if err != nil {
		return si.storeOutputPort.OutputStoresHandleError(err)
	}
	return si.storeOutputPort.OutputAllStores(stores)
}

func (si *StoreInteractor) GetNearStores() error {
	places, err := si.storeRepository.GetNearStores()
	if err != nil {
		return si.storeOutputPort.OutputStoresHandleError(err)
	}
	return si.storeOutputPort.OutputAllStores(places)
}

func (si *StoreInteractor) GetFavoriteStores(userId string) error {
	stores, err := si.storeRepository.GetFavoriteStores(userId)
	if err != nil {
		return si.storeOutputPort.OutputStoresHandleError(err)
	}
	return si.storeOutputPort.OutputAllStores(stores)
}

func (si *StoreInteractor) SaveFavoriteStore(store *model.Store, userId string) error {
	exist, err := si.storeRepository.ExistFavorite(store, userId)
	if err != nil {
		return si.storeOutputPort.OutputStoresHandleError(err)
	}
	if exist {
		return si.storeOutputPort.OutputStoresHandleError(errors.ErrFavoriteStoreAlreadyExists)
	}
	if err := si.storeRepository.SaveFavoriteStore(store, userId); err != nil {
		return si.storeOutputPort.OutputStoresHandleError(err)
	}
	return si.storeOutputPort.OutputSaveFavoriteStoreResult()
}

func (si *StoreInteractor) GetTopFavoriteStores() error {
	stores, err := si.storeRepository.GetTopFavoriteStores()
	if err != nil {
		return si.storeOutputPort.OutputStoresHandleError(err)
	}
	return si.storeOutputPort.OutputAllStores(stores)
}
