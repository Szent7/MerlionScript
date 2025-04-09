package merlion

import (
	"MerlionScript/keeper"
	"MerlionScript/types/soapTypes"
	"MerlionScript/utils/soap"
	"encoding/xml"
	"fmt"
	"log"
)

func GetCatalog(catName string) (soapTypes.ItemMenu, bool) {
	//fmt.Println(base64.StdEncoding.EncodeToString([]byte(credentials)))
	req := soapTypes.ItemMenuReq{
		Cat_id: "All",
	}
	//var res = make([]soapTypes.ItemMenu, 100)
	decoder, err := soap.SoapCallHandleResponse(keeper.MerlionMainURL, soapTypes.GetCatalogUrl, req)
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

func GetAllCatalogCodes() ([]string, error) {
	req := soapTypes.ItemMenuReq{
		Cat_id: "Order",
	}
	decoder, err := soap.SoapCallHandleResponse(keeper.MerlionMainURL, soapTypes.GetCatalogUrl, req)
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
		var item soapTypes.ItemMenu
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
	//fmt.Println(base64.StdEncoding.EncodeToString([]byte(credentials)))
	req := soapTypes.ItemMenuReq{
		Cat_id: "All",
	}

	decoder, err := soap.SoapCallHandleResponse(keeper.MerlionMainURL, soapTypes.GetCatalogUrl, req)
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
