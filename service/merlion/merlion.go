package merlion

import (
	"MerlionScript/types/soapTypes"
	"MerlionScript/utils/soap"
	"encoding/base64"
	"encoding/xml"
	"fmt"
	"log"
)

func GetItemsByCatId(catId string, credentials string) []soapTypes.ItemCatalog {
	//fmt.Println(base64.StdEncoding.EncodeToString([]byte(credentials)))
	req := soapTypes.ItemCatalogReq{
		Cat_id: catId,
		//Page:   "1",
	}
	//var res = make([]types.ItemMenu, 100)
	decoder, err := soap.SoapCallHandleResponse("https://apitest.merlion.com/rl/mlservice3", soapTypes.GetItemsUrl, req)
	if err != nil {
		log.Fatalf("SoapCallHandleResponse error: %s", err)
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
	return items
}

func GetItemsImagesByItemId(itemId string, credentials string) []soapTypes.ItemImage {
	//fmt.Println(base64.StdEncoding.EncodeToString([]byte(credentials)))
	req := soapTypes.ItemImageReq{
		Item_id: []soapTypes.ItemId{{Item: itemId}},
		//Page:   "1",
	}
	//var res = make([]types.ItemMenu, 100)
	decoder, err := soap.SoapCallHandleResponse("https://apitest.merlion.com/rl/mlservice3", soapTypes.GetItemsImagesUrl, req)
	if err != nil {
		log.Fatalf("SoapCallHandleResponse error: %s", err)
	}
	var item soapTypes.ItemImage
	var items []soapTypes.ItemImage
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
					fmt.Println("ошибка при декодировании item(getitemsimages):", err)
					break
				}

				items = append(items, item)
			}
		}
	}
	return items
}

func GetItemsByItemId(itemId string, credentials string) []soapTypes.ItemCatalog {
	//fmt.Println(base64.StdEncoding.EncodeToString([]byte(credentials)))
	req := soapTypes.ItemCatalogReq{
		Item_id: []soapTypes.ItemId{{Item: itemId}},
		//Page:   "1",
	}
	//var res = make([]types.ItemMenu, 100)
	decoder, err := soap.SoapCallHandleResponse("https://apitest.merlion.com/rl/mlservice3", soapTypes.GetItemsUrl, req)
	if err != nil {
		log.Fatalf("SoapCallHandleResponse error: %s", err)
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
	return items
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
					fmt.Println("ошибка при декодировании item(getcatalog):", err)
					break
				}
			}
		}
	}
	return soapTypes.ItemMenu{}, false
}

func GetItemsAvailByItemId(itemId string, credentials string) []soapTypes.ItemAvail {
	//fmt.Println(base64.StdEncoding.EncodeToString([]byte(credentials)))
	req := soapTypes.ItemAvailReq{
		Item_id: []soapTypes.ItemId{{Item: itemId}},
		//Page:   "1",
	}
	//var res = make([]types.ItemMenu, 100)
	decoder, err := soap.SoapCallHandleResponse("https://apitest.merlion.com/rl/mlservice3", soapTypes.GetItemsAvailUrl, req)
	if err != nil {
		log.Fatalf("SoapCallHandleResponse error: %s", err)
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
	return items
}

func GetCatalogUniqueCodes(credentials string) []string {
	fmt.Println(base64.StdEncoding.EncodeToString([]byte(credentials)))
	req := soapTypes.ItemMenuReq{
		Cat_id: "All",
	}

	decoder, err := soap.SoapCallHandleResponse("https://apitest.merlion.com/rl/mlservice3", soapTypes.GetCatalogUrl, req)
	if err != nil {
		log.Fatalf("SoapCallHandleResponse error: %s", err)
	}

	counts := make(map[string]int)
	for {
		token, err := decoder.Token()
		if err != nil {
			break
		}
		var item soapTypes.ItemMenu
		switch start := token.(type) {
		case xml.StartElement:
			if start.Name.Local == "item" {
				err := decoder.DecodeElement(&item, &start)
				if err != nil {
					fmt.Println("ошибка при декодировании item(getcataloguniquecodes):", err)
					break
				}
				counts[item.ID]++
				counts[item.ID_PARENT]++
			}
		}
	}
	var uniqueId []string
	for str, count := range counts {
		if count == 1 {
			uniqueId = append(uniqueId, str)
		}
	}
	return uniqueId
}
