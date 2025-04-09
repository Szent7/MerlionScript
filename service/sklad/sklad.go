package sklad

import (
	"MerlionScript/types/restTypes"
	"MerlionScript/utils/rest"
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"strings"
)

func CreateItem(product restTypes.CreateItem) (restTypes.Response, error) {
	//jsonBody, err := json.Marshal(reqBody)
	jsonBody, err := json.MarshalIndent(product, "", "  ")
	//fmt.Println("тело запроса в JSON:", string(jsonBody))
	if err != nil {
		log.Println("ошибка при преобразовании структуры в JSON:", err)
		return restTypes.Response{}, err
	}
	//authHeader := "Basic " + base64.StdEncoding.EncodeToString([]byte(credentials))
	//url := "https://api.moysklad.ru/api/remap/1.2/entity/product"
	body, err := rest.CreateRequest("POST", restTypes.ItemUrl, bytes.NewBuffer(jsonBody))
	if err != nil {
		return restTypes.Response{}, err
	}
	/*var meta restTypes.TestProductMeta
	fmt.Println("rawBody:", string(body.Body))
	err = json.Unmarshal(body, &meta)
	if err != nil {
		return err
	}*/
	return *body, nil
}

func GetItemByManufacturer(manufacturer string) (restTypes.Response, error) {
	//jsonBody, err := json.Marshal(reqBody)
	//jsonBody, err := json.MarshalIndent(product, "", "  ")
	//fmt.Println("тело запроса в JSON:", string(jsonBody))
	/*if err != nil {
		log.Println("ошибка при преобразовании структуры в JSON:", err)
		return restTypes.Response{}, err
	}*/
	//authHeader := "Basic " + base64.StdEncoding.EncodeToString([]byte(credentials))
	//url := "https://api.moysklad.ru/api/remap/1.2/entity/product"
	body, err := rest.CreateRequest("GET", restTypes.ItemUrl+"?search="+manufacturer, nil)
	if err != nil {
		return restTypes.Response{}, err
	}
	/*var meta restTypes.TestProductMeta
	fmt.Println("rawBody:", string(body.Body))
	err = json.Unmarshal(body, &meta)
	if err != nil {
		return err
	}*/
	return *body, nil
}

func GetItemUUID(codeMS string) (string, bool, error) {
	//authHeader := "Basic " + base64.StdEncoding.EncodeToString([]byte(SkladCredentials))
	response, err := rest.CreateRequest("GET", restTypes.ItemUrl+"?search="+codeMS, nil)
	if err != nil {
		return "", false, err
	}
	if response.StatusCode != 200 {
		fmt.Println(response.StatusCode)
		fmt.Println(string(response.Body))
		return "", false, err
	}
	items := restTypes.SearchItem{}
	if err := json.Unmarshal(response.Body, &items); err != nil {
		return "", false, fmt.Errorf("ошибка при декодировании item (getItemUUID): %s", err.Error())
	}

	if len(items.Rows) != 0 {
		return items.Rows[0].Id, items.Rows[0].IsSerialTrackable, nil
	}
	/*for i := range items.Rows {
		if strings.Contains(items.Rows[i].Name, codeMS) {
			return items.Rows[i].Id, nil
		}
	}*/
	return "", false, nil
}

func GetStoreUUID(storeName string) (string, error) {
	//authHeader := "Basic " + base64.StdEncoding.EncodeToString([]byte(SkladCredentials))
	response, err := rest.CreateRequest("GET", restTypes.StoreUrl, nil)
	if err != nil {
		return "", err
	}
	if response.StatusCode != 200 {
		fmt.Println(response.StatusCode)
		fmt.Println(string(response.Body))
		return "", err
	}
	items := restTypes.SearchItem{}
	if err := json.Unmarshal(response.Body, &items); err != nil {
		return "", fmt.Errorf("ошибка при декодировании item (GetStoreUUID): %s", err.Error())
	}

	for i := range items.Rows {
		if СontainsSubstring(items.Rows[i].Name, storeName) {
			return items.Rows[i].Id, nil
		}
		/*if strings.Contains(items.Rows[i].Name, storeName) {
			return items.Rows[i].Id, nil
		}*/
	}
	return "", nil
}

func GetItemsAvail(itemUUID string, storeUUID string) (int, error) {
	//authHeader := "Basic " + base64.StdEncoding.EncodeToString([]byte(SkladCredentials))
	response, err := rest.CreateRequest("GET", restTypes.StocksUrl+"?filter=assortmentId="+itemUUID+";storeId="+storeUUID, nil)
	if err != nil {
		return -1, err
	}
	if response.StatusCode != 200 {
		return -1, err
	}
	items := []restTypes.SearchStock{}
	if err := json.Unmarshal(response.Body, &items); err != nil {
		return -1, fmt.Errorf("ошибка при декодировании item (GetItemsAvail): %s", err.Error())
	}

	if len(items) == 0 {
		return 0, nil
	}
	return items[0].Stock, nil
}

func IncreaseItemsAvail(request *restTypes.Acceptance) error {
	jsonBody, err := json.MarshalIndent(request, "", "  ")
	//fmt.Println("тело запроса в JSON:", string(jsonBody))
	if err != nil {
		fmt.Println("ошибка при преобразовании структуры в JSON:", err)
		return err
	}
	//authHeader := "Basic " + base64.StdEncoding.EncodeToString([]byte(SkladCredentials))
	response, err := rest.CreateRequest("POST", restTypes.AcceptanceUrl, bytes.NewBuffer(jsonBody))
	if err != nil {
		return err
	}
	if response.StatusCode != 200 {
		fmt.Println(response.StatusCode)
		fmt.Println(string(response.Body))
		return err
	}
	return nil
}

func DecreaseItemsAvail(request *restTypes.Acceptance) error {
	jsonBody, err := json.MarshalIndent(request, "", "  ")
	//fmt.Println("тело запроса в JSON:", string(jsonBody))
	if err != nil {
		fmt.Println("ошибка при преобразовании структуры в JSON:", err)
		return err
	}
	//authHeader := "Basic " + base64.StdEncoding.EncodeToString([]byte(SkladCredentials))
	response, err := rest.CreateRequest("POST", restTypes.WoffUrl, bytes.NewBuffer(jsonBody))
	if err != nil {
		return err
	}
	if response.StatusCode != 200 {
		fmt.Println(response.StatusCode)
		fmt.Println(string(response.Body))
		return err
	}
	return nil
}

// Сравнение артикула в названии
func СontainsSubstring(s string, substr string) bool {
	n := len(substr)
	if n == 0 || s == "" || n > len(s) {
		return false
	}

	for i := 0; i <= len(s)-n; i++ {
		if strings.HasPrefix(s[i:], substr) {
			// Проверяем, что после подстроки идет либо пробел, либо конец строки
			if i+n == len(s) || s[i+n] == ' ' {
				return true
			}
		}
	}

	return false
}
