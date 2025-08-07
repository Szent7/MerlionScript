package requests

import (
	elektronmirTypes "MerlionScript/services/elektronmir/types"
	"MerlionScript/utils/rest"
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"strings"
)

func GetItemsByCatId(catId int, token string) (elektronmirTypes.ItemResponse, error) {
	request := elektronmirTypes.ItemRequest{
		SectionID: catId,
	}
	jsonBody, err := json.MarshalIndent(request, "", "  ")
	if err != nil {
		log.Println("ошибка при преобразовании структуры в JSON:", err)
		return elektronmirTypes.ItemResponse{}, err
	}

	response, err := rest.CreateRequestElektronmir("POST", elektronmirTypes.ItemUrl, bytes.NewBuffer(jsonBody), token)
	if err != nil {
		return elektronmirTypes.ItemResponse{}, err
	}

	if response.StatusCode != 200 {
		fmt.Println(response.StatusCode)
		fmt.Println(string(response.Body))
		return elektronmirTypes.ItemResponse{}, err
	}

	var items elektronmirTypes.ItemResponse
	if err := json.Unmarshal(response.Body, &items); err != nil {
		return elektronmirTypes.ItemResponse{}, fmt.Errorf("ошибка при декодировании item (getallcategorycodes): %s", err.Error())
	}

	return items, nil
}

// func GetItemsByItemId(itemId string, token string) ([]netlabTypes.Item, error) {
// 	url := fmt.Sprintf(netlabTypes.ItemIdUrl, itemId, token)

// 	decoder, err := rest.CreateRequestXML("GET", url, nil)
// 	if err != nil {
// 		return nil, err
// 	}

// 	var item netlabTypes.Item
// 	var items []netlabTypes.Item

// 	for {
// 		token, err := decoder.Token()
// 		if err != nil {
// 			break
// 		}
// 		switch start := token.(type) {
// 		case xml.StartElement:
// 			if start.Name.Local == "goods" {
// 				err := decoder.DecodeElement(&item, &start)
// 				if err != nil {
// 					fmt.Println("Ошибка при декодировании item(GetItemsByCatId):", err)
// 					break
// 				}

// 				items = append(items, item)
// 			}
// 		}
// 	}
// 	return items, nil
// }

// фильтр по производителю
func GetItemsByCatIdFormatted(catId int, token string) (elektronmirTypes.ItemResponse, error) {
	rawItems, err := GetItemsByCatId(catId, token)
	if err != nil {
		return elektronmirTypes.ItemResponse{}, err
	}

	var items elektronmirTypes.ItemResponse
	items.Data = make([]elektronmirTypes.Data, 0, 200)
	items.Status = rawItems.Status
	items.Limit = rawItems.Limit
	items.Offset = rawItems.Offset
	items.Total = rawItems.Total

	for i := range rawItems.Data {
		lower := strings.ToLower(rawItems.Data[i].Vendor)
		if strings.Contains(lower, "dahua") || strings.Contains(lower, "tenda") {
			rawItems.Data[i].Vendor = lower
			items.Data = append(items.Data, rawItems.Data[i])
		}
	}

	return items, nil
}

// func GetItemsByItemIdFormatted(itemId string, token string) (netlabTypes.ItemNetlab, error) {
// 	rawItems, err := GetItemsByItemId(itemId, token)
// 	if err != nil {
// 		return netlabTypes.ItemNetlab{}, err
// 	}

// 	for i := range rawItems {
// 		var item netlabTypes.ItemNetlab
// 		for _, property := range rawItems[i].Properties.Property {
// 			if item.Id == itemId {
// 				switch property.Name {
// 				case "название":
// 					item.Name = property.Value
// 				case "производитель":
// 					item.ManufacturerName = property.Value
// 				case "PN":
// 					item.Manufacturer = property.Value
// 				case "количество на Лобненской":
// 					rawRemains, err := strconv.ParseFloat(property.Value, 64)
// 					if err != nil {
// 						log.Printf("Ошибка при парсинге значения (GetItemsByCatIdFormatted) %s : err = %s\n", property.Value, err.Error())
// 						break
// 					}
// 					item.Remains = int(rawRemains)
// 				case "цена по категории F":
// 					rawPrice, err := strconv.ParseFloat(property.Value, 64)
// 					if err != nil {
// 						log.Printf("Ошибка при парсинге значения (GetItemsByCatIdFormatted) %s : err = %s\n", property.Value, err.Error())
// 						break
// 					}
// 					item.Price = rawPrice
// 				}
// 				item.Id = rawItems[i].ID
// 				return item, nil
// 			}
// 		}
// 	}

// 	return netlabTypes.ItemNetlab{}, nil
// }

func GetItemsAvailByItemIdBatch(itemId []int, token string) (*elektronmirTypes.StockResponse, error) {
	request := elektronmirTypes.StockRequest{
		ProductIDs: itemId,
		Page:       1,
		Available: []string{
			"stock",
		},
	}

	jsonBody, err := json.MarshalIndent(request, "", "  ")
	if err != nil {
		log.Println("ошибка при преобразовании структуры в JSON:", err)
		return nil, err
	}

	response, err := rest.CreateRequestElektronmir("POST", elektronmirTypes.StockUrl, bytes.NewBuffer(jsonBody), token)
	if err != nil {
		return nil, err
	}

	if response.StatusCode != 200 {
		fmt.Println(response.StatusCode)
		fmt.Println(string(response.Body))
		return nil, err
	}

	var items elektronmirTypes.StockResponse
	if err := json.Unmarshal(response.Body, &items); err != nil {
		return nil, fmt.Errorf("ошибка при декодировании item (getallcategorycodes): %s", err.Error())
	}

	return &items, nil
}
