package merlion

import (
	"MerlionScript/types/soapTypes"
	"MerlionScript/utils/soap"
	"encoding/xml"
	"fmt"
	"log"
)

func GetItemsByCatId(catId string) ([]soapTypes.ItemCatalog, error) {
	//fmt.Println(base64.StdEncoding.EncodeToString([]byte(credentials)))
	req := soapTypes.ItemCatalogReq{
		Cat_id: catId,
		//Page:   "1",
	}
	//var res = make([]types.ItemMenu, 100)
	decoder, err := soap.SoapCallHandleResponse("https://apitest.merlion.com/rl/mlservice3", soapTypes.GetItemsUrl, req)
	if err != nil {
		log.Printf("Ошибка при SOAP-запросе (GetItemsByCatId): %s\n", err)
		return nil, err
	}
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
				err := decoder.DecodeElement(&item, &start)
				//fmt.Println(item)
				//res = append(res, item)
				if err != nil {
					fmt.Println("Ошибка при декодировании item(GetItemsByCatId):", err)
					break
				}

				items = append(items, item)
			}
		}
	}
	return items, nil
}

func GetItemsAvailByItemId(itemId string) ([]soapTypes.ItemAvail, error) {
	//fmt.Println(base64.StdEncoding.EncodeToString([]byte(credentials)))
	req := soapTypes.ItemAvailReq{
		Item_id: []soapTypes.ItemId{{Item: itemId}},
		//Page:   "1",
	}
	//var res = make([]types.ItemMenu, 100)
	decoder, err := soap.SoapCallHandleResponse("https://apitest.merlion.com/rl/mlservice3", soapTypes.GetItemsAvailUrl, req)
	if err != nil {
		log.Printf("Ошибка при SOAP-запросе (GetItemsAvailByItemId): %s\n", err)
		return nil, err
	}
	var item soapTypes.ItemAvail
	var items []soapTypes.ItemAvail
	for {
		token, err := decoder.Token()
		if err != nil {
			break
		}
		switch start := token.(type) {
		case xml.StartElement:
			//fmt.Println("Xml tag:", start.Name.Local)
			if start.Name.Local == "item" {
				err := decoder.DecodeElement(&item, &start)
				//fmt.Println(item)
				//res = append(res, item)
				if err != nil {
					fmt.Println("ошибка при декодировании item(getitemsavail):", err)
					break
				}

				items = append(items, item)
			}
		}
	}
	return items, nil
}

func GetItemsByItemId(itemId string) ([]soapTypes.ItemCatalog, error) {
	//fmt.Println(base64.StdEncoding.EncodeToString([]byte(credentials)))
	req := soapTypes.ItemCatalogReq{
		Item_id: []soapTypes.ItemId{{Item: itemId}},
		//Page:   "1",
	}
	//var res = make([]types.ItemMenu, 100)
	decoder, err := soap.SoapCallHandleResponse("https://apitest.merlion.com/rl/mlservice3", soapTypes.GetItemsUrl, req)
	if err != nil {
		log.Printf("Ошибка при SOAP-запросе (GetItemsByItemId): %s\n", err)
		return nil, err
	}
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
				err := decoder.DecodeElement(&item, &start)
				//fmt.Println(item)
				//res = append(res, item)
				if err != nil {
					fmt.Println("ошибка при декодировании item(getitems):", err)
					break
				}

				items = append(items, item)
			}
		}
	}
	return items, nil
}
