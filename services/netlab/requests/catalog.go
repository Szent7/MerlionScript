package requests

import (
	netlabTypesRest "MerlionScript/types/restTypes/netlab"
	netlabTypesSoap "MerlionScript/types/soapTypes/netlab"
	"MerlionScript/utils/rest"
	"encoding/xml"
	"fmt"
	"log"
	"strings"
)

/*
func GetAllCatalogCodes(token string) ([]string, error) {
	url := fmt.Sprintf(netlabTypes.CatalogUrl, token)

	decoder, err := rest.CreateRequestXML("GET", url, nil)
	if err != nil {
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
*/

func GetAllCategoryCodes(token string) ([]netlabTypesSoap.ItemCategory, error) {
	url := fmt.Sprintf(netlabTypesRest.CategoryUrl, netlabTypesRest.CatalogName, token)

	decoder, err := rest.CreateRequestXML("GET", url, nil)
	if err != nil {
		return nil, err
	}

	var codes []netlabTypesSoap.ItemCategory

	for {
		token, err := decoder.Token()
		if err != nil {
			break
		}
		var item netlabTypesSoap.ItemCategory
		switch start := token.(type) {
		case xml.StartElement:
			//fmt.Println("Xml tag:", start.Name.Local)
			if start.Name.Local == "category" {
				err := decoder.DecodeElement(&item, &start)
				if err != nil {
					log.Println("ошибка при декодировании item(getcatalog):", err)
					break
				}
				lower := strings.ToLower(item.Name)
				if strings.Contains(lower, "dahua") || strings.Contains(lower, "tenda") {
					codes = append(codes, item)
				}
			}
		}
	}
	return codes, nil
}
