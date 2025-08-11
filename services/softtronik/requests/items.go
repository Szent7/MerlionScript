package requests

import (
	"MerlionScript/keeper"
	softtronikTypes "MerlionScript/services/softtronik/types"
	"MerlionScript/utils/rest"
	"encoding/json"
	"fmt"
)

func GetItemsByCatId(catId string) ([]softtronikTypes.ProductItem, error) {
	url := fmt.Sprintf(softtronikTypes.ItemUrl, keeper.GetSofttronikContractor(), catId)

	response, err := rest.CreateRequest("GET", url, nil, "")
	if err != nil {
		return nil, err
	}
	if response.StatusCode != 200 {
		fmt.Println(response.StatusCode)
		fmt.Println(string(response.Body))
		return nil, err
	}

	var items []softtronikTypes.ProductItem
	if err := json.Unmarshal(response.Body, &items); err != nil {
		return nil, fmt.Errorf("ошибка при декодировании item (getitemsbycatid): %s", err.Error())
	}

	return items, nil
}

func GetItemsAvails(catId string) (softtronikTypes.StocksItem, error) {
	url := fmt.Sprintf(softtronikTypes.StocksUrl, keeper.GetSofttronikContractor(), catId, keeper.GetSofttronikContractKey())

	response, err := rest.CreateRequest("GET", url, nil, "")
	if err != nil {
		return softtronikTypes.StocksItem{}, err
	}
	if response.StatusCode != 200 {
		fmt.Println(response.StatusCode)
		fmt.Println(string(response.Body))
		return softtronikTypes.StocksItem{}, err
	}

	var avails softtronikTypes.StocksItem
	if err := json.Unmarshal(response.Body, &avails); err != nil {
		return softtronikTypes.StocksItem{}, fmt.Errorf("ошибка при декодировании item (getitemsavails): %s", err.Error())
	}

	return avails, nil
}

func GetItemsAvailsAll(catId []softtronikTypes.CategoryItem) (softtronikTypes.StocksItem, error) {
	var availsFull softtronikTypes.StocksItem
	var avails softtronikTypes.StocksItem
	var err error = nil
	for i := range catId {
		avails, err = GetItemsAvails(catId[i].ID)
		if err != nil {
			return softtronikTypes.StocksItem{}, err
		}
		if i == 0 {
			availsFull = avails
			continue
		}
		availsFull.Body.ProductsDataWithPricesAndBalances = append(availsFull.Body.ProductsDataWithPricesAndBalances, avails.Body.ProductsDataWithPricesAndBalances...)
	}
	return availsFull, err
}
