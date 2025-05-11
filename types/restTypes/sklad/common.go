package sklad

const (
	ItemUrl         = "https://api.moysklad.ru/api/remap/1.2/entity/product"
	GroupUrl        = "https://api.moysklad.ru/api/remap/1.2/entity/productfolder"
	StoreUrl        = "https://api.moysklad.ru/api/remap/1.2/entity/store"
	OrganizationUrl = "https://api.moysklad.ru/api/remap/1.2/entity/organization"
	AcceptanceUrl   = "https://api.moysklad.ru/api/remap/1.2/entity/supply"
	WoffUrl         = "https://api.moysklad.ru/api/remap/1.2/entity/loss"
	StocksUrl       = "https://api.moysklad.ru/api/remap/1.2/report/stock/bystore/current"
	ImageUrl        = "https://api.moysklad.ru/api/remap/1.2/entity/product/%s/images"
)

type MetaMiddle struct {
	Meta Meta `json:"meta"`
}

type Meta struct {
	Href         string `json:"href,omitempty"`         //URL
	MetadataHref string `json:"metadataHref,omitempty"` //URL
	Type         string `json:"type,omitempty"`
	MediaType    string `json:"mediaType,omitempty"`
	UuidHref     string `json:"uuidHref,omitempty"`     //URL
	DownloadHref string `json:"downloadHref,omitempty"` //URL
}
