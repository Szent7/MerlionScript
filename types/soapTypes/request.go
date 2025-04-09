package soapTypes

import "encoding/xml"

const (
	GetCatalogUrl     = "https://api.merlion.com/rl/mlservice3#getCatalog"
	GetItemsUrl       = "https://api.merlion.com/rl/mlservice3#getItems"
	GetItemsImagesUrl = "https://api.merlion.com/rl/mlservice3#getItemsImages"
	GetItemsAvailUrl  = "https://api.merlion.com/rl/mlservice3#getItemsAvail"

	DownloadImageUrl = "http://img.merlion.ru/items"
)

type ItemMenuReq struct {
	XMLName xml.Name `xml:"getCatalog"`
	Cat_id  string   `xml:"cat_id"`
}

type ItemCatalogReq struct {
	XMLName      xml.Name `xml:"getItems"`
	Cat_id       string   `xml:"cat_id"`
	Item_id      []ItemId `xml:"item_id"`
	Page         string   `xml:"page"`
	Rows_on_page string   `xml:"rows_on_page"`
}

type ItemId struct {
	Item string `xml:"item"`
}

type ItemImageReq struct {
	XMLName          xml.Name `xml:"getItemsImages"`
	Cat_id           string   `xml:"cat_id"`
	Item_id          []ItemId `xml:"item_id"`
	Page             string   `xml:"page"`
	Rows_on_page     string   `xml:"rows_on_page"`
	Last_time_change string   `xml:"last_time_change"`
}

type ItemAvailReq struct {
	XMLName         xml.Name `xml:"getItemsAvail"`
	Cat_id          string   `xml:"cat_id"`
	Item_id         []ItemId `xml:"item_id"`
	Shipment_method string   `xml:"shipment_method"`
	Shipment_date   string   `xml:"shipment_date"`
	Only_avail      string   `xml:"only_avail"`
}
