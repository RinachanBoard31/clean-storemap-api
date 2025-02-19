package gateway

import (
	apiErr "clean-storemap-api/src/driver/api"
	authErr "clean-storemap-api/src/driver/auth"
	dbErr "clean-storemap-api/src/driver/db"
	"clean-storemap-api/src/entity/errors"
)

func (ug *UserGateway) convertToBusinessError(err error) error {
	// DB
	if errors.Is(err, dbErr.ErrRecordNotFound) {
		return errors.ErrUserNotFound
	} else if errors.Is(err, dbErr.ErrDuplicatedKey) {
		return errors.ErrUserAlreadyExists
	} else if errors.Is(err, dbErr.ErrInvalidData) {
		return errors.ErrInvalidUserData
	} else if errors.Is(err, dbErr.ErrInvalidDB) {
		return errors.ErrDBConnection
	} else if errors.Is(err, dbErr.ErrUnsupported) {
		return errors.ErrUnsupported
	}
	// Auth
	if errors.Is(err, authErr.ErrGetEmailFailed) {
		return errors.ErrUnAuthorized
	} else if errors.Is(err, authErr.ErrFailedGenerateToken) {
		return errors.ErrUnAuthorized
	}
	return errors.ErrUnExpected
}

func (sg *StoreGateway) convertToBusinessError(err error) error {
	// DB
	if errors.Is(err, dbErr.ErrRecordNotFound) {
		return errors.ErrStoreNotFound
	} else if errors.Is(err, dbErr.ErrDuplicatedKey) {
		return errors.ErrStoreAlreadyExists
	} else if errors.Is(err, dbErr.ErrInvalidData) {
		return errors.ErrStoreInvalidData
	} else if errors.Is(err, dbErr.ErrInvalidDB) {
		return errors.ErrDBConnection
	} else if errors.Is(err, dbErr.ErrMissingWhereClause) {
		return errors.ErrDBConnection
	} else if errors.Is(err, dbErr.ErrInvalidField) {
		return errors.ErrDBConnection
	} else if errors.Is(err, dbErr.ErrInvalidValue) {
		return errors.ErrDBConnection
	} else if errors.Is(err, dbErr.ErrUnsupported) {
		return errors.ErrUnsupported
	}
	// API
	if errors.Is(err, apiErr.ErrFailedToGetLocation) {
		return errors.ErrFailedGetNearByStores
	} else if errors.Is(err, apiErr.ErrFailedToGetNearbyStores) {
		return errors.ErrFailedGetNearByStores
	}
	return errors.ErrUnExpected
}
