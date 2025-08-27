package requests

import (
	netlabTypes "MerlionScript/services/netlab/types"
	"MerlionScript/utils/rest"
	"encoding/xml"
	"fmt"
	"log"
	"strconv"
)

func GetItemsByCatId(catId string, token string) ([]netlabTypes.Item, error) {
	url := fmt.Sprintf(netlabTypes.ItemUrl, netlabTypes.CatalogName, catId, token)

	decoder, err := rest.CreateRequestXML("GET", url, nil)
	if err != nil {
		return nil, err
	}

	var item netlabTypes.Item
	var items []netlabTypes.Item

	for {
		token, err := decoder.Token()
		if err != nil {
			break
		}
		switch start := token.(type) {
		case xml.StartElement:
			if start.Name.Local == "goods" {
				err := decoder.DecodeElement(&item, &start)
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

func GetItemsByItemId(itemId string, token string) ([]netlabTypes.ItemById, error) {
	url := fmt.Sprintf(netlabTypes.ItemIdUrl, itemId, token)

	decoder, err := rest.CreateRequestXML("GET", url, nil)
	if err != nil {
		return nil, err
	}

	var item netlabTypes.ItemById
	var items []netlabTypes.ItemById

	for {
		token, err := decoder.Token()
		if err != nil {
			break
		}
		switch start := token.(type) {
		case xml.StartElement:
			if start.Name.Local == "data" {
				err := decoder.DecodeElement(&item, &start)
				if err != nil {
					fmt.Println("Ошибка при декодировании item(GetItemsByItemId):", err)
					break
				}

				items = append(items, item)
			}
		}
	}
	return items, nil
}

func GetItemsByCatIdFormatted(catId string, token string) ([]netlabTypes.ItemNetlab, error) {
	rawItems, err := GetItemsByCatId(catId, token)
	if err != nil {
		return nil, err
	}

	var items = make([]netlabTypes.ItemNetlab, 0, len(rawItems))
	for i := range rawItems {
		var item netlabTypes.ItemNetlab
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
				item.Remains = int(rawRemains)
			case "цена по категории F":
				rawPrice, err := strconv.ParseFloat(property.Value, 64)
				if err != nil {
					log.Printf("Ошибка при парсинге значения (GetItemsByCatIdFormatted) %s : err = %s\n", property.Value, err.Error())
					break
				}
				item.Price = rawPrice
			}
		}
		item.Id = rawItems[i].ID
		items = append(items, item)
	}

	return items, nil
}

func GetItemsByItemIdFormatted(itemId string, token string) (netlabTypes.ItemNetlab, error) {
	rawItems, err := GetItemsByItemId(itemId, token)
	if err != nil {
		return netlabTypes.ItemNetlab{}, err
	}

	for i := range rawItems {
		var item netlabTypes.ItemNetlab
		if rawItems[i].ID == itemId {
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
					item.Remains = int(rawRemains)
				case "цена по категории F":
					rawPrice, err := strconv.ParseFloat(property.Value, 64)
					if err != nil {
						log.Printf("Ошибка при парсинге значения (GetItemsByCatIdFormatted) %s : err = %s\n", property.Value, err.Error())
						break
					}
					item.Price = rawPrice
				}
			}
			item.Id = rawItems[i].ID
			return item, nil
		}
	}

	return netlabTypes.ItemNetlab{}, nil
}
