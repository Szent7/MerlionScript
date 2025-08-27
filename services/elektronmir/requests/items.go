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

func GetItemsByCatId(catId int, limit int, offset int, token string) (*elektronmirTypes.ItemResponse, error) {
	request := elektronmirTypes.ItemRequest{
		SectionID: catId,
	}
	jsonBody, err := json.MarshalIndent(request, "", "  ")
	if err != nil {
		log.Println("ошибка при преобразовании структуры в JSON:", err)
		return nil, err
	}

	querryParam := fmt.Sprintf("?limit=%d&offset=%d", limit, offset)
	url := fmt.Sprintf("%s%s", elektronmirTypes.ItemUrl, querryParam)

	response, err := rest.CreateRequestElektronmir("POST", url, bytes.NewBuffer(jsonBody), token)
	if err != nil {
		return nil, err
	}

	if response.StatusCode != 200 {
		fmt.Println(response.StatusCode)
		fmt.Println(string(response.Body))
		return nil, err
	}

	var items elektronmirTypes.ItemResponse
	if err := json.Unmarshal(response.Body, &items); err != nil {
		return nil, fmt.Errorf("ошибка при декодировании item (getallcategorycodes): %s", err.Error())
	}

	return &items, nil
}

// фильтр по производителю
func GetItemsByCatIdFormatted(catId int, token string) (*[]elektronmirTypes.Data, error) {
	limit := 1000
	offset := 0
	total := 1000
	var rawItems []elektronmirTypes.Data

	for {
		resp, err := GetItemsByCatId(catId, limit, offset, token)
		if err != nil || resp == nil {
			return nil, err
		}

		rawItems = append(rawItems, resp.Data...)
		if offset+limit >= resp.Total {
			total = resp.Total
			break
		}
		offset += limit
	}

	var items = make([]elektronmirTypes.Data, 0, total)

	for i := range rawItems {
		lower := strings.ToLower(rawItems[i].Vendor)
		if strings.Contains(lower, "dahua") || strings.Contains(lower, "tenda") {
			rawItems[i].Vendor = lower
			items = append(items, rawItems[i])
		}
	}

	return &items, nil
}

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
