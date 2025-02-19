package api

import "errors"

var (
	ErrFailedToGetLocation     = errors.New("failed to get location")
	ErrFailedToGetNearbyStores = errors.New("failed to get nearby stores")
)
