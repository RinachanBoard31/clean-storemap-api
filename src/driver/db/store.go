package db

import (
	"time"
)

type DbStoreDriver struct{}

func NewStoreDriver() *DbStoreDriver {
	return &DbStoreDriver{}
}

/* interfaceと型は同義。仮にgatewayがDBの型を知ったとしても、どんなDBから来たかわかるわけではないのでおk */
type FavoriteStore struct {
	Id                  string `gorm:"primaryKey"`
	UserId              string `gorm:"not null;size:64"`
	User                User   `gorm:"foreignKey:UserId;references:Id"`
	StoreId             string `gorm:"not null"`
	StoreName           string `gorm:"not null"`
	RegularOpeningHours string
	PriceLevel          string
	Latitude            string `gorm:"not null"`
	Longitude           string `gorm:"not null"`
	CreatedAt           time.Time
	UpdatedAt           time.Time
}

func (dbs *DbStoreDriver) GetStores() ([]*FavoriteStore, error) {
	var stores []*FavoriteStore
	err := DB.First(&stores).Error
	if err != nil {
		return nil, convertGormError(err)
	}
	return stores, nil
}

func (dbs *DbStoreDriver) FindFavorite(storeId string, userId string) (*FavoriteStore, error) {
	var stores []FavoriteStore
	err := DB.Where("store_id = ? AND user_id = ?", storeId, userId).Find(&stores).Error
	if err != nil {
		return nil, convertGormError(err)
	}
	if len(stores) == 0 {
		return nil, nil
	}

	return &stores[0], nil
}

func (dbs *DbStoreDriver) FindFavoriteByUser(userId string) ([]*FavoriteStore, error) {
	var stores []*FavoriteStore
	err := DB.Where("user_id = ?", userId).First(&stores).Error
	if err != nil {
		return nil, convertGormError(err)
	}
	return stores, nil
}

func (dbs *DbStoreDriver) SaveStore(dbStore *FavoriteStore) error {
	err := DB.Create(&dbStore).Error
	if err != nil {
		return convertGormError(err)
	}
	return nil
}

func (dbs *DbStoreDriver) GetTopStores() ([]*FavoriteStore, error) {
	oneWeekAgo := time.Now().AddDate(0, 0, -7)

	// store_idごとにカウント、多い順に最大10件を取得
	var topStoreIds []string
	err := DB.Model(&FavoriteStore{}).
		Select("store_id").
		Where("created_at >= ?", oneWeekAgo).
		Group("store_id").
		Order("COUNT(*) desc").
		Limit(10).
		Pluck("store_id", &topStoreIds).Error
	if err != nil {
		return nil, convertGormError(err)
	}

	// 取得したstore_idに対応するfavorite_storeを順番に取得
	var stores []*FavoriteStore
	for _, storeId := range topStoreIds {
		var store FavoriteStore
		err = DB.Where("store_id = ?", storeId).First(&store).Error
		if err != nil {
			return nil, convertGormError(err)
		}
		stores = append(stores, &store)
	}

	return stores, nil
}
