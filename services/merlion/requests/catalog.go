package requests

import (
	merlionTypes "MerlionScript/services/merlion/types"
	"MerlionScript/utils/soap"
	"encoding/xml"
	"fmt"
	"log"
)

func GetCatalog(catName string) (merlionTypes.ItemMenu, bool) {
	req := merlionTypes.ItemMenuReq{
		Cat_id: "All",
	}

	decoder, err := soap.SoapCallHandleResponse(merlionTypes.MerlionMainURL, merlionTypes.GetCatalogUrl, req)
	if err != nil {
		log.Fatalf("SoapCallHandleResponse error: %s", err)
	}
	for {
		token, err := decoder.Token()
		if err != nil {
			break
		}
		var item merlionTypes.ItemMenu
		switch start := token.(type) {
		case xml.StartElement:
			if start.Name.Local == "item" {
				err := decoder.DecodeElement(&item, &start)

				if item.Description == catName {
					return item, true
				}

				if err != nil {
					fmt.Println("ошибка при декодировании item(getcatalog):", err)
					break
				}
			}
		}
	}
	return merlionTypes.ItemMenu{}, false
}

func GetAllCatalogCodes() ([]string, error) {
	req := merlionTypes.ItemMenuReq{
		Cat_id: "Order",
	}
	decoder, err := soap.SoapCallHandleResponse(merlionTypes.MerlionMainURL, merlionTypes.GetCatalogUrl, req)
	if err != nil {
		log.Printf("Ошибка при SOAP-запросе (GetAllCatalogCodes): %s\n", err)
		return nil, err
	}
	var codes []string
	for {
		token, err := decoder.Token()
		if err != nil {
			break
		}
		var item merlionTypes.ItemMenu
		switch start := token.(type) {
		case xml.StartElement:
			//fmt.Println("Xml tag:", start.Name.Local)
			if start.Name.Local == "item" {
				err := decoder.DecodeElement(&item, &start)
				if err != nil {
					log.Println("ошибка при декодировании item(getcatalog):", err)
					break
				}
				codes = append(codes, item.ID)
			}
		}
	}
	return codes, nil
}

func GetCatalogUniqueCodes() []string {
	req := merlionTypes.ItemMenuReq{
		Cat_id: "All",
	}

	decoder, err := soap.SoapCallHandleResponse(merlionTypes.MerlionMainURL, merlionTypes.GetCatalogUrl, req)
	if err != nil {
		log.Fatalf("SoapCallHandleResponse error: %s", err)
	}

	counts := make(map[string]int)
	for {
		token, err := decoder.Token()
		if err != nil {
			break
		}
		var item merlionTypes.ItemMenu
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
