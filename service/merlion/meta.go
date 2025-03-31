package merlion

import (
	"MerlionScript/types/restTypes"
	"MerlionScript/utils/rest"
	"encoding/base64"
	"encoding/json"
	"fmt"
)

func GetOrganizationMeta(SkladCredentials string, OrganizationName string) (restTypes.Meta, error) {
	authHeader := "Basic " + base64.StdEncoding.EncodeToString([]byte(SkladCredentials))
	response, err := rest.CreateRequest("GET", restTypes.OrganizationUrl+"?search="+OrganizationName, authHeader, nil)
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

func GetStoreMeta(SkladCredentials string, StoreName string) (restTypes.Meta, error) {
	authHeader := "Basic " + base64.StdEncoding.EncodeToString([]byte(SkladCredentials))
	response, err := rest.CreateRequest("GET", restTypes.StoreUrl, authHeader, nil)
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

func GetItemMeta(SkladCredentials string, article string) (restTypes.Meta, error) {
	authHeader := "Basic " + base64.StdEncoding.EncodeToString([]byte(SkladCredentials))
	response, err := rest.CreateRequest("GET", restTypes.ItemUrl+"?search="+article, authHeader, nil)
	if err != nil || response.StatusCode != 200 {
		return restTypes.Meta{}, err
	}

	items := restTypes.SearchStoreOrganization{}
	if err := json.Unmarshal(response.Body, &items); err != nil {
		return restTypes.Meta{}, fmt.Errorf("ошибка при декодировании item (getstoremeta): %s", err.Error())
	}
	if len(items.Rows) != 0 {
		return items.Rows[0].StoreMeta, nil
	}
	return restTypes.Meta{}, nil
}
