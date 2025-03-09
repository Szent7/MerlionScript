package sklad

import (
	"MerlionScript/types/restTypes"
	"MerlionScript/utils/rest"
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
)

func CreateTestItemSklad(product restTypes.TestProduct, credentials string) error {
	//jsonBody, err := json.Marshal(reqBody)
	jsonBody, err := json.MarshalIndent(product, "", "  ")
	fmt.Println("тело запроса в JSON:", string(jsonBody))
	if err != nil {
		fmt.Println("ошибка при преобразовании структуры в JSON:", err)
		return err
	}
	authHeader := "Basic " + base64.StdEncoding.EncodeToString([]byte(credentials))
	//url := "https://api.moysklad.ru/api/remap/1.2/entity/product"
	body, err := rest.CreateRequest("POST", restTypes.CreateItemUrl, authHeader, bytes.NewBuffer(jsonBody))
	if err != nil {
		return err
	}
	var meta restTypes.TestProductMeta
	fmt.Println("rawBody:", string(body))
	err = json.Unmarshal(body, &meta)
	if err != nil {
		return err
	}
	return nil
}

func CreateTestCatSklad(group restTypes.TestProductGroup, credentials string) (restTypes.TestProductMeta, error) {
	//jsonBody, err := json.Marshal(reqBody)
	jsonBody, err := json.MarshalIndent(group, "", "  ")
	fmt.Println("тело запроса в JSON:", string(jsonBody))
	if err != nil {
		fmt.Println("ошибка при преобразовании структуры в JSON:", err)
		return restTypes.TestProductMeta{}, err
	}
	authHeader := "Basic " + base64.StdEncoding.EncodeToString([]byte(credentials))
	//url := "https://api.moysklad.ru/api/remap/1.2/entity/productfolder"
	body, err := rest.CreateRequest("POST", restTypes.CreateGroupUrl, authHeader, bytes.NewBuffer(jsonBody))
	if err != nil {
		return restTypes.TestProductMeta{}, err
	}
	var meta restTypes.TestProductMeta
	fmt.Println("rawBody:", string(body))
	err = json.Unmarshal(body, &meta)
	if err != nil {
		return restTypes.TestProductMeta{}, err
	}
	return meta, nil
}

/*
func getCatSklad() {
	//fmt.Println(base64.StdEncoding.EncodeToString([]byte("admin@sandbox1:wolf444466")))
	authHeader := "Basic " + base64.StdEncoding.EncodeToString([]byte("admin@sandbox1:wolf444466"))
	url := "https://api.moysklad.ru/api/remap/1.2/entity/productfolder"
	body, err := rest.CreateRequest("GET", url, authHeader, nil)
	if err != nil {
		return
	}
	fmt.Println(body)
}
*/
