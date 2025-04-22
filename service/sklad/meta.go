package sklad

import (
	"MerlionScript/types/restTypes"
	"MerlionScript/utils/rest"
	"encoding/json"
	"fmt"
	"strings"
)

func GetOrganizationMeta(OrganizationName string) (restTypes.Meta, error) {
	//authHeader := "Basic " + base64.StdEncoding.EncodeToString([]byte(SkladCredentials))
	//response, err := rest.CreateRequest("GET", restTypes.OrganizationUrl+"?search="+OrganizationName, nil)
	response, err := rest.CreateRequest("GET", restTypes.OrganizationUrl, nil)
	if err != nil || response.StatusCode != 200 {
		return restTypes.Meta{}, err
	}

	items := restTypes.SearchStoreOrganization{}
	if err := json.Unmarshal(response.Body, &items); err != nil {
		return restTypes.Meta{}, fmt.Errorf("ошибка при декодировании item (getstoremeta): %s", err.Error())
	}

	for _, item := range items.Rows {
		if item.Name == OrganizationName {
			return item.StoreMeta, nil
		}
	}
	return restTypes.Meta{}, nil
}

func GetStoreMeta(StoreName string) (restTypes.Meta, error) {
	//authHeader := "Basic " + base64.StdEncoding.EncodeToString([]byte(SkladCredentials))
	response, err := rest.CreateRequest("GET", restTypes.StoreUrl, nil)
	if err != nil || response.StatusCode != 200 {
		return restTypes.Meta{}, err
	}

	items := restTypes.SearchStoreOrganization{}
	if err := json.Unmarshal(response.Body, &items); err != nil {
		return restTypes.Meta{}, fmt.Errorf("ошибка при декодировании item (getstoremeta): %s", err.Error())
	}

	for _, item := range items.Rows {
		if item.Name == StoreName {
			return item.StoreMeta, nil
		}
	}
	return restTypes.Meta{}, nil
}

func GetCatMeta(CatName string) (restTypes.Meta, error) {
	//authHeader := "Basic " + base64.StdEncoding.EncodeToString([]byte(SkladCredentials))
	response, err := rest.CreateRequest("GET", restTypes.GroupUrl, nil)
	if err != nil || response.StatusCode != 200 {
		return restTypes.Meta{}, err
	}

	items := restTypes.SearchStoreOrganization{}
	if err := json.Unmarshal(response.Body, &items); err != nil {
		return restTypes.Meta{}, fmt.Errorf("ошибка при декодировании item (getcatmeta): %s", err.Error())
	}

	for _, item := range items.Rows {
		if item.Name == CatName {
			return item.StoreMeta, nil
		}
	}
	return restTypes.Meta{}, nil
}

func GetItemMeta(article string) (restTypes.Meta, error) {
	//authHeader := "Basic " + base64.StdEncoding.EncodeToString([]byte(SkladCredentials))
	response, err := rest.CreateRequest("GET", restTypes.ItemUrl+"?search="+article, nil)
	if err != nil || response.StatusCode != 200 {
		return restTypes.Meta{}, err
	}

	items := restTypes.SearchItem{}
	if err := json.Unmarshal(response.Body, &items); err != nil {
		return restTypes.Meta{}, fmt.Errorf("ошибка при декодировании item (getstoremeta): %s", err.Error())
	}
	for i := range items.Rows {
		if strings.Contains(items.Rows[i].Article, article) {
			return items.Rows[i].Meta, nil
		}
	}

	return restTypes.Meta{}, nil
}
