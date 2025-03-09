package types

import "encoding/xml"

type ItemMenu struct {
	XMLName     xml.Name `xml:"item"`
	ID          string   `xml:"ID"`
	ID_PARENT   string   `xml:"ID_PARENT"`
	Description string   `xml:"Description"`
}

type ItemMenuReq struct {
	XMLName xml.Name `xml:"getCatalog"`
	Cat_id  string   `xml:"cat_id"`
}

type ItemCatalog struct {
	XMLName xml.Name `xml:"item"`
	Name    string   `xml:"Name"`
	Brand   string   `xml:"Brand"`
	No      string   `xml:"No"`
	Weight  float32  `xml:"Weight"`
	VAT     int      `xml:"VAT"`
}

type ItemCatalogReq struct {
	XMLName      xml.Name `xml:"getItems"`
	Cat_id       string   `xml:"cat_id"`
	Page         string   `xml:"page"`
	Rows_on_page string   `xml:"rows_on_page"`
}
