package types

import "encoding/xml"

type ItemCategory struct {
	XMLName  xml.Name `xml:"category"`
	ID       string   `xml:"id"`
	Name     string   `xml:"name"`
	ParentId string   `xml:"parentId"`
	Leaf     bool     `xml:"leaf"`
}

type Item struct {
	XMLName    xml.Name   `xml:"goods"`
	ID         string     `xml:"id"`
	Properties Properties `xml:"properties"`
}

type ItemImage struct {
	XMLName    xml.Name   `xml:"item"`
	ID         string     `xml:"id"`
	Properties Properties `xml:"properties"`
}

type Properties struct {
	XMLName  xml.Name   `xml:"properties"`
	Property []Property `xml:"property"`
}

type Property struct {
	XMLName xml.Name `xml:"property"`
	Name    string   `xml:"name"`
	Value   string   `xml:"value"`
}

type ItemNetlab struct {
	Id               string
	Name             string
	Manufacturer     string
	ManufacturerName string
	Price            float64
	Remains          int
}

type Status struct {
	XMLName xml.Name `xml:"status"`
	Code    string   `xml:"code"`
	Message string   `xml:"message"`
}
