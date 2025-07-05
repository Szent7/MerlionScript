package requests

import (
	"MerlionScript/keeper"
	softtronikTypesRest "MerlionScript/types/restTypes/softtronik"
	"MerlionScript/utils/rest"
	"encoding/json"
	"fmt"
)

func GetItemsByCatId(catId string) ([]softtronikTypesRest.ProductItem, error) {
	url := fmt.Sprintf(softtronikTypesRest.ItemUrl, keeper.K.GetSofttronikContractor(), catId)

	response, err := rest.CreateRequest("GET", url, nil, "")
	if err != nil {
		return nil, err
	}
	if response.StatusCode != 200 {
		fmt.Println(response.StatusCode)
		fmt.Println(string(response.Body))
		return nil, err
	}

	var items []softtronikTypesRest.ProductItem
	if err := json.Unmarshal(response.Body, &items); err != nil {
		return nil, fmt.Errorf("ошибка при декодировании item (getitemsbycatid): %s", err.Error())
	}

	return items, nil
}

func GetItemsAvails(catId string) (softtronikTypesRest.StocksItem, error) {
	url := fmt.Sprintf(softtronikTypesRest.StocksUrl, keeper.K.GetSofttronikContractor(), catId, keeper.K.GetSofttronikContractKey())

	response, err := rest.CreateRequest("GET", url, nil, "")
	if err != nil {
		return softtronikTypesRest.StocksItem{}, err
	}
	if response.StatusCode != 200 {
		fmt.Println(response.StatusCode)
		fmt.Println(string(response.Body))
		return softtronikTypesRest.StocksItem{}, err
	}

	var avails softtronikTypesRest.StocksItem
	if err := json.Unmarshal(response.Body, &avails); err != nil {
		return softtronikTypesRest.StocksItem{}, fmt.Errorf("ошибка при декодировании item (getitemsavails): %s", err.Error())
	}

	return avails, nil
}

func GetItemsAvailsAll(catId []softtronikTypesRest.CategoryItem) (softtronikTypesRest.StocksItem, error) {
	var availsFull softtronikTypesRest.StocksItem
	var avails softtronikTypesRest.StocksItem
	var err error = nil
	for i := range catId {
		avails, err = GetItemsAvails(catId[i].ID)
		if err != nil {
			return softtronikTypesRest.StocksItem{}, err
		}
		if i == 0 {
			availsFull = avails
			continue
		}
		availsFull.Body.ProductsDataWithPricesAndBalances = append(availsFull.Body.ProductsDataWithPricesAndBalances, avails.Body.ProductsDataWithPricesAndBalances...)
	}
	return availsFull, err
}
