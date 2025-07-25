package common

import (
	skladTypes "MerlionScript/services/sklad/types"
)

type BaseSystem interface {
	GetSystemName() string
	GetDBTableName() string
}

type ERPSystem interface {
	BaseSystem
	GetCatName() string
	GetItemsByArticle(find string) (*[]ItemList, error)
	GetItemID(code string) (string, error)
	GetItemAvails(code string) (StockERP, error)
	GetImagesList(id string) (skladTypes.SearchImage, error)

	CreateAcceptance(*skladTypes.Acceptance) error
	CreateWoff(*skladTypes.Acceptance) error
	UploadImage(imageData skladTypes.UploadImage, id string) error

	GetCatMeta() (skladTypes.Meta, error)
	GetOrgMeta() (skladTypes.Meta, error)
	GetStoreMeta() (skladTypes.Meta, error)
	GetStoreUUID() (string, error)
}

type Service interface {
	BaseSystem

	Init()
	GetArticlesList() (*[]ArticleList, error)
	GetItemsList() (*[]ItemList, error)
	GetStocksList() (*map[string]StockList, error)
	GetImagesList(code string) (*[]ImageList, error)
}

type ArticleList struct {
	Article     string
	Brand       string
	ServiceCode string
}

type ItemList struct {
	Article      string
	PositionName string
}

type StockList struct {
	Article string
	Stock   int
	Price   float32 // в рублях
}

type ImageList struct {
	DownloadUrl string
	Filename    string
}

type StockERP struct {
	ItemMeta          skladTypes.Meta
	Stock             int
	IsSerialTrackable bool
}
