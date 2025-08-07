package common

import (
	skladTypes "MerlionScript/services/sklad/types"
	"MerlionScript/utils/db/interfaceDB"
	"MerlionScript/utils/db/typesDB"
	"context"
)

type BaseSystem interface {
	GetSystemName() string
	GetDBTableName() string
}

type ERPSystem interface {
	BaseSystem
	//GetCatName() string
	GetItemsByArticle(article string) (*[]ItemList, error)
	GetItemID(code string) (string, error)
	GetItemAvails(code string, storeUUID string) (StockERP, error)
	GetImagesList(id string) (skladTypes.SearchImage, error)

	CreateItem(item *typesDB.CodesIDs, newId string, itemName string, catalog string) error
	CreateAcceptance(acceptanceReq *skladTypes.Acceptance) error
	CreateWoff(woffReq *skladTypes.Acceptance) error
	UploadImage(imageData skladTypes.UploadImage, id string) error

	GetCatMeta(catalog string) (skladTypes.Meta, error)
	GetOrgMeta(organization string) (skladTypes.Meta, error)
	GetStoreMeta(store string) (skladTypes.Meta, error)
	GetStoreUUID(store string) (string, error)
}

type Service interface {
	BaseSystem

	GetArticlesList(ctx context.Context) (*[]ArticleList, error)
	GetItemsList(ctx context.Context, dbInstance interfaceDB.DB) (*[]ItemList, error)
	GetStocksList(ctx context.Context, dbInstance interfaceDB.DB) (*map[string]StockList, error)
	GetImagesList(ctx context.Context, code string) (*[]ImageList, error)

	GetOrgName() string
	GetStoreName() string
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
	Stock int
	Price float32 // в рублях
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
