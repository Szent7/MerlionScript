package requests

import (
	netlabTypesRest "MerlionScript/types/restTypes/netlab"
	netlabTypesSoap "MerlionScript/types/soapTypes/netlab"
	"MerlionScript/utils/rest"
	"encoding/xml"
	"fmt"
	"log"
	"strconv"
)

func GetItemsByCatId(catId string, token string) ([]netlabTypesSoap.Item, error) {
	url := fmt.Sprintf(netlabTypesRest.ItemUrl, netlabTypesRest.CatalogName, catId, token)

	decoder, err := rest.CreateRequestXML("GET", url, nil)
	if err != nil {
		return nil, err
	}

	var item netlabTypesSoap.Item
	var items []netlabTypesSoap.Item

	for {
		token, err := decoder.Token()
		if err != nil {
			break
		}
		switch start := token.(type) {
		case xml.StartElement:
			//fmt.Println("Xml tag:", start.Name.Local)
			if start.Name.Local == "goods" {
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

func GetItemsByCatIdFormatted(catId string, token string) ([]netlabTypesSoap.ItemNetlab, error) {
	rawItems, err := GetItemsByCatId(catId, token)
	if err != nil {
		return nil, err
	}

	var items = make([]netlabTypesSoap.ItemNetlab, 0, len(rawItems))
	for i := range rawItems {
		var item netlabTypesSoap.ItemNetlab
		var remains int = 0
		for _, property := range rawItems[i].Properties.Property {
			switch property.Name {
			case "название":
				item.Name = property.Value
			case "производитель":
				item.ManufacturerName = property.Value
			case "PN":
				item.Manufacturer = property.Value
			case "количество на Лобненской":
				rawRemains, err := strconv.ParseFloat(property.Value, 64)
				if err != nil {
					log.Printf("Ошибка при парсинге значения (GetItemsByCatIdFormatted) %s : err = %s\n", property.Value, err.Error())
					break
				}
				//remains += int(rawRemains)
				item.Remains = int(rawRemains)
			/*case "количество на Калужской":
				rawRemains, err := strconv.ParseFloat(property.Value, 64)
				if err != nil {
					log.Printf("Ошибка при парсинге значения (GetItemsByCatIdFormatted) %s : err = %s\n", property.Value, err.Error())
					break
				}
				remains += int(rawRemains)
			case "количество на Курской":
				rawRemains, err := strconv.ParseFloat(property.Value, 64)
				if err != nil {
					log.Printf("Ошибка при парсинге значения (GetItemsByCatIdFormatted) %s : err = %s\n", property.Value, err.Error())
					break
				}
				remains += int(rawRemains)*/
			case "цена по категории F":
				rawPrice, err := strconv.ParseFloat(property.Value, 64)
				if err != nil {
					log.Printf("Ошибка при парсинге значения (GetItemsByCatIdFormatted) %s : err = %s\n", property.Value, err.Error())
					break
				}
				item.Price = rawPrice
			}
		}
		item.Remains += remains
		item.Id = rawItems[i].ID
		items = append(items, item)
	}

	return items, nil
}
