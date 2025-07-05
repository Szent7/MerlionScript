package requests

import (
	"MerlionScript/keeper"
	softtronikTypesRest "MerlionScript/types/restTypes/softtronik"
	"MerlionScript/utils/rest"
	"encoding/json"
	"fmt"
	"strings"
)

func GetAllCategoryCodes() ([]softtronikTypesRest.CategoryItem, error) {
	url := fmt.Sprintf(softtronikTypesRest.CategoryUrl, keeper.K.GetSofttronikContractor())

	response, err := rest.CreateRequest("GET", url, nil, "")
	if err != nil {
		return nil, err
	}
	if response.StatusCode != 200 {
		fmt.Println(response.StatusCode)
		fmt.Println(string(response.Body))
		return nil, err
	}

	var categories []softtronikTypesRest.CategoryItem
	if err := json.Unmarshal(response.Body, &categories); err != nil {
		return nil, fmt.Errorf("ошибка при декодировании item (getallcategorycodes): %s", err.Error())
	}

	filterManufacturer(&categories, "dahua")
	return categories, nil
}

func filterManufacturer(categories *[]softtronikTypesRest.CategoryItem, substr string) {
	var filteredCat []softtronikTypesRest.CategoryItem
	for _, cat := range *categories {
		if strings.Contains(strings.ToLower(cat.Name), substr) {
			filteredCat = append(filteredCat, cat)
		}
	}
	*categories = filteredCat
}
