package requests

import (
	elektronmirTypes "MerlionScript/services/elektronmir/types"
	"MerlionScript/utils/rest"
	"encoding/json"
	"fmt"
	"strings"
)

func GetAllCategoryCodes(token string) (elektronmirTypes.CatalogResponse, error) {
	response, err := rest.CreateRequestElektronmir("GET", elektronmirTypes.CatalogUrl, nil, token)
	if err != nil {
		return elektronmirTypes.CatalogResponse{}, err
	}

	if response.StatusCode != 200 {
		fmt.Println(response.StatusCode)
		fmt.Println(string(response.Body))
		return elektronmirTypes.CatalogResponse{}, err
	}

	var categories elektronmirTypes.CatalogResponse
	if err := json.Unmarshal(response.Body, &categories); err != nil {
		return elektronmirTypes.CatalogResponse{}, fmt.Errorf("ошибка при декодировании item (getallcategorycodes): %s", err.Error())
	}

	return categories, nil
}

//выбраны только главные категории
func GetAllCategoryCodesFormatted(token string) (elektronmirTypes.CatalogResponse, error) {
	rawCategories, err := GetAllCategoryCodes(token)
	if err != nil {
		return elektronmirTypes.CatalogResponse{}, err
	}

	var categories elektronmirTypes.CatalogResponse
	categories.Data = make([]elektronmirTypes.Item, 0, 200)
	categories.Status = rawCategories.Status
	categories.Total = rawCategories.Total

	for i := range rawCategories.Data {
		if rawCategories.Data[i].ParentID == nil &&
			(strings.Contains(rawCategories.Data[i].Name, "офис") ||
				strings.Contains(rawCategories.Data[i].Name, "видео") ||
				strings.Contains(rawCategories.Data[i].Name, "сети") ||
				strings.Contains(rawCategories.Data[i].Name, "гаджеты")) {
			categories.Data = append(categories.Data, rawCategories.Data[i])
		}
	}

	return categories, nil
}
