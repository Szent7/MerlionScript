package requests

import (
	skladTypes "MerlionScript/services/sklad/types"
	"MerlionScript/utils/rest"
	"encoding/json"
	"fmt"
	"strings"
)

func GetOrganizationMeta(OrganizationName string) (skladTypes.Meta, error) {
	//authHeader := "Basic " + base64.StdEncoding.EncodeToString([]byte(SkladCredentials))
	//response, err := rest.CreateRequest("GET", skladTypes.OrganizationUrl+"?search="+OrganizationName, nil)
	response, err := rest.CreateRequestMS("GET", skladTypes.OrganizationUrl, nil)
	if err != nil || response.StatusCode != 200 {
		return skladTypes.Meta{}, err
	}

	items := skladTypes.SearchStoreOrganization{}
	if err := json.Unmarshal(response.Body, &items); err != nil {
		return skladTypes.Meta{}, fmt.Errorf("ошибка при декодировании item (getstoremeta): %s", err.Error())
	}

	for _, item := range items.Rows {
		if item.Name == OrganizationName {
			return item.StoreMeta, nil
		}
	}
	return skladTypes.Meta{}, nil
}

func GetStoreMeta(StoreName string) (skladTypes.Meta, error) {
	//authHeader := "Basic " + base64.StdEncoding.EncodeToString([]byte(SkladCredentials))
	response, err := rest.CreateRequestMS("GET", skladTypes.StoreUrl, nil)
	if err != nil || response.StatusCode != 200 {
		return skladTypes.Meta{}, err
	}

	items := skladTypes.SearchStoreOrganization{}
	if err := json.Unmarshal(response.Body, &items); err != nil {
		return skladTypes.Meta{}, fmt.Errorf("ошибка при декодировании item (getstoremeta): %s", err.Error())
	}

	for _, item := range items.Rows {
		if item.Name == StoreName {
			return item.StoreMeta, nil
		}
	}
	return skladTypes.Meta{}, nil
}

func GetCatMeta(CatName string) (skladTypes.Meta, error) {
	//authHeader := "Basic " + base64.StdEncoding.EncodeToString([]byte(SkladCredentials))
	response, err := rest.CreateRequestMS("GET", skladTypes.GroupUrl, nil)
	if err != nil || response.StatusCode != 200 {
		return skladTypes.Meta{}, err
	}

	items := skladTypes.SearchStoreOrganization{}
	if err := json.Unmarshal(response.Body, &items); err != nil {
		return skladTypes.Meta{}, fmt.Errorf("ошибка при декодировании item (getcatmeta): %s", err.Error())
	}

	for _, item := range items.Rows {
		if item.Name == CatName {
			return item.StoreMeta, nil
		}
	}
	return skladTypes.Meta{}, nil
}

func GetItemMeta(article string) (skladTypes.Meta, error) {
	//authHeader := "Basic " + base64.StdEncoding.EncodeToString([]byte(SkladCredentials))
	response, err := rest.CreateRequestMS("GET", skladTypes.ItemUrl+"?search="+article, nil)
	if err != nil || response.StatusCode != 200 {
		return skladTypes.Meta{}, err
	}

	items := skladTypes.SearchItem{}
	if err := json.Unmarshal(response.Body, &items); err != nil {
		return skladTypes.Meta{}, fmt.Errorf("ошибка при декодировании item (getstoremeta): %s", err.Error())
	}
	for i := range items.Rows {
		if strings.Contains(items.Rows[i].Article, article) {
			return items.Rows[i].Meta, nil
		}
	}

	return skladTypes.Meta{}, nil
}
