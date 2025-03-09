package soapTypes

import "encoding/xml"

const (
	GetCatalogUrl = "https://apitest.merlion.com/rl/mlservice3#getCatalog"
	GetItemsUrl   = "https://apitest.merlion.com/rl/mlservice3#getItems"
)

type ItemMenuReq struct {
	XMLName xml.Name `xml:"getCatalog"`
	Cat_id  string   `xml:"cat_id"`
}

type ItemCatalogReq struct {
	XMLName      xml.Name `xml:"getItems"`
	Cat_id       string   `xml:"cat_id"`
	Page         string   `xml:"page"`
	Rows_on_page string   `xml:"rows_on_page"`
}
