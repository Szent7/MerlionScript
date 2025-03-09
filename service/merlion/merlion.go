package merlion

import (
	"MerlionScript/types/soapTypes"
	"MerlionScript/utils/soap"
	"encoding/base64"
	"encoding/xml"
	"fmt"
	"log"
)

func GetCatItems(catId string, credentials string) ([]soapTypes.ItemCatalog, bool) {
	fmt.Println(base64.StdEncoding.EncodeToString([]byte(credentials)))
	req := soapTypes.ItemCatalogReq{
		Cat_id:       catId,
		Page:         "1",
		Rows_on_page: "50",
	}
	//var res = make([]types.ItemMenu, 100)
	decoder, err := soap.SoapCallHandleResponse("https://apitest.merlion.com/rl/mlservice3", soapTypes.GetItemsUrl, req)
	if err != nil {
		log.Fatalf("SoapCallHandleResponse error: %s", err)
	}
	var i int = 0
	var item soapTypes.ItemCatalog
	var items []soapTypes.ItemCatalog
	for {
		token, err := decoder.Token()
		if err != nil {
			break
		}
		switch start := token.(type) {
		case xml.StartElement:
			//fmt.Println("Xml tag:", start.Name.Local)
			if start.Name.Local == "item" {
				if i >= 5 {
					break
				}
				err := decoder.DecodeElement(&item, &start)
				//fmt.Println(item)
				items = append(items, item)
				//res = append(res, item)
				if err != nil {
					fmt.Println("Ошибка при декодировании item:", err)
					break
				}
				i++
			}
		}
	}
	if i == 0 {
		return nil, false
	} else {
		return items, true
	}
}

func GetCatalog(catName string, credentials string) (soapTypes.ItemMenu, bool) {
	fmt.Println(base64.StdEncoding.EncodeToString([]byte(credentials)))
	req := soapTypes.ItemMenuReq{
		Cat_id: "All",
	}
	//var res = make([]soapTypes.ItemMenu, 100)
	decoder, err := soap.SoapCallHandleResponse("https://apitest.merlion.com/rl/mlservice3", soapTypes.GetCatalogUrl, req)
	if err != nil {
		log.Fatalf("SoapCallHandleResponse error: %s", err)
	}
	for {
		token, err := decoder.Token()
		if err != nil {
			break
		}
		var item soapTypes.ItemMenu
		switch start := token.(type) {
		case xml.StartElement:
			//fmt.Println("Xml tag:", start.Name.Local)
			if start.Name.Local == "item" {
				err := decoder.DecodeElement(&item, &start)
				//fmt.Println(item)
				if item.Description == catName {
					return item, true
				}
				//res = append(res, item)
				if err != nil {
					fmt.Println("Ошибка при декодировании item:", err)
					break
				}
			}
		}
	}
	return soapTypes.ItemMenu{}, false
}
