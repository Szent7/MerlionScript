package soapTypes

import "encoding/xml"

type ItemMenu struct {
	XMLName     xml.Name `xml:"item"`
	ID          string   `xml:"ID"`
	ID_PARENT   string   `xml:"ID_PARENT"`
	Description string   `xml:"Description"`
}

type ItemCatalog struct {
	XMLName xml.Name `xml:"item"`
	Name    string   `xml:"Name"`
	Brand   string   `xml:"Brand"`
	No      string   `xml:"No"`
	Weight  float32  `xml:"Weight"`
	VAT     int      `xml:"VAT"`
}
