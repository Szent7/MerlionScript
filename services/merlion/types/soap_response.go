package types

import "encoding/xml"

type ItemMenu struct {
	XMLName     xml.Name `xml:"item"`
	ID          string   `xml:"ID"`
	ID_PARENT   string   `xml:"ID_PARENT"`
	Description string   `xml:"Description"`
}

type ItemCatalog struct {
	XMLName     xml.Name `xml:"item"`
	Name        string   `xml:"Name"`
	Brand       string   `xml:"Brand"`
	Vendor_part string   `xml:"Vendor_part"`
	No          string   `xml:"No"`
	Weight      float32  `xml:"Weight"`
	VAT         int      `xml:"VAT"`
}

type ItemImage struct {
	XMLName  xml.Name `xml:"item"`
	No       string   `xml:"No"`
	ViewType string   `xml:"ViewType"`
	SizeType string   `xml:"SizeType"`
	FileName string   `xml:"FileName"`
	Created  string   `xml:"Created"`
	Size     string   `xml:"Size"`
	Width    string   `xml:"Width"`
	Height   string   `xml:"Height"`
}

type ItemAvail struct {
	XMLName               xml.Name `xml:"item"`
	No                    string   `xml:"No"`
	PriceClient           float32  `xml:"PriceClient"`
	PriceClient_RG        float32  `xml:"PriceClient_RG"`
	PriceClient_MSK       float32  `xml:"PriceClient_MSK"`
	AvailableClient       int      `xml:"AvailableClient"`
	AvailableClient_RG    int      `xml:"AvailableClient_RG"`
	AvailableClient_MSK   int      `xml:"AvailableClient_MSK"`
	AvailableExpected     int      `xml:"AvailableExpected"`
	AvailableExpectedNext int      `xml:"AvailableExpectedNext"`
	DateExpectedNext      string   `xml:"DateExpectedNext"`
	RRP                   float32  `xml:"RRP"`
	RRP_Date              string   `xml:"RRP_Date"`
	PriceClientRUB        float32  `xml:"PriceClientRUB"`
	PriceClientRUB_RG     float32  `xml:"PriceClientRUB_RG"`
	PriceClientRUB_MSK    float32  `xml:"PriceClientRUB_MSK"`
	Online_Reserve        int      `xml:"Online_Reserve"`
	ReserveCost           float32  `xml:"ReserveCost"`
}

type ItemAvailPrice struct {
	PriceClientRUB_MSK  float32
	AvailableClient_MSK int
}
