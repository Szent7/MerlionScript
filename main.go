package main

import (
	"MerlionScript/rest"
	"MerlionScript/soap"
	"MerlionScript/types"
	"bytes"
	"encoding/base64"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"log"
)

func main() {
	//getCatSklad()
	cat, exist := getCatalog("Ноутбуки")
	if !exist {
		fmt.Println("Такой категории не существует (MERLION)")
		return
	}
	fmt.Printf("Категория найдена (MERLION):%v\n", cat)
	metadata, err := createTestCatSklad(cat.Description, cat.Description, cat.ID)
	if err != nil {
		fmt.Printf("Ошибка при создании группы (СКЛАД):%v\n", err)
		return
	}
	fmt.Printf("Метаданные созданной группы (СКЛАД):%v\n", metadata)
	items, exist := getCatItems(cat.ID)
	if !exist {
		fmt.Printf("Товаров в категории %s не существует (MERLION)", cat.Description)
		return
	}
	restItems := make([]types.TestProduct, 5)
	for i, rawItem := range items {
		restItems[i].Name = rawItem.Name
		restItems[i].Code = rawItem.No
		restItems[i].Vat = rawItem.VAT
		restItems[i].Weight = int(rawItem.Weight)
		restItems[i].ProductFolder.Meta = metadata.Meta

		err = createTestItemSklad(restItems[i])
		if err != nil {
			fmt.Printf("Ошибка при создании товара [%d]\n", i)
		}
	}
}

func createTestItemSklad(product types.TestProduct) error {
	//jsonBody, err := json.Marshal(reqBody)
	jsonBody, err := json.MarshalIndent(product, "", "  ")
	fmt.Println("тело запроса в JSON:", string(jsonBody))
	if err != nil {
		fmt.Println("ошибка при преобразовании структуры в JSON:", err)
		return err
	}
	authHeader := "Basic " + base64.StdEncoding.EncodeToString([]byte(`admin@sandbox1244:Jf!FQBy!q5"N3]Z`))
	url := "https://api.moysklad.ru/api/remap/1.2/entity/product"
	body, err := rest.CreateRequest("POST", url, authHeader, bytes.NewBuffer(jsonBody))
	if err != nil {
		return err
	}
	var meta types.TestProductMeta
	fmt.Println("rawBody:", string(body))
	err = json.Unmarshal(body, &meta)
	if err != nil {
		return err
	}
	return nil
}

func createTestCatSklad(catName string, catDescription string, catCode string) (types.TestProductMeta, error) {
	reqBody := &types.TestProductGroup{
		Name:        catName,
		Description: catDescription,
		Code:        catCode,
	}
	//jsonBody, err := json.Marshal(reqBody)
	jsonBody, err := json.MarshalIndent(reqBody, "", "  ")
	fmt.Println("тело запроса в JSON:", string(jsonBody))
	if err != nil {
		fmt.Println("ошибка при преобразовании структуры в JSON:", err)
		return types.TestProductMeta{}, err
	}
	authHeader := "Basic " + base64.StdEncoding.EncodeToString([]byte(`admin@sandbox1244:Jf!FQBy!q5"N3]Z`))
	url := "https://api.moysklad.ru/api/remap/1.2/entity/productfolder"
	body, err := rest.CreateRequest("POST", url, authHeader, bytes.NewBuffer(jsonBody))
	if err != nil {
		return types.TestProductMeta{}, err
	}
	var meta types.TestProductMeta
	fmt.Println("rawBody:", string(body))
	err = json.Unmarshal(body, &meta)
	if err != nil {
		return types.TestProductMeta{}, err
	}
	return meta, nil
}

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

func getCatItems(catId string) ([]types.ItemCatalog, bool) {
	fmt.Println(base64.StdEncoding.EncodeToString([]byte("TC0051161|API:lt2iZpXb41")))
	req := types.ItemCatalogReq{
		Cat_id:       catId,
		Page:         "1",
		Rows_on_page: "50",
	}
	//var res = make([]types.ItemMenu, 100)
	decoder, err := soap.SoapCallHandleResponse("https://apitest.merlion.com/rl/mlservice3", "https://apitest.merlion.com/rl/mlservice3#getItems", req)
	if err != nil {
		log.Fatalf("SoapCallHandleResponse error: %s", err)
	}
	var i int = 0
	var item types.ItemCatalog
	var items []types.ItemCatalog
	for {
		token, err := decoder.Token()
		if err != nil {
			break
		}
		switch start := token.(type) {
		case xml.StartElement:
			//fmt.Println("Xml tag:", start.Name.Local)
			if start.Name.Local == "item" {
				if i >= 5 {
					break
				}
				err := decoder.DecodeElement(&item, &start)
				//fmt.Println(item)
				items = append(items, item)
				//res = append(res, item)
				if err != nil {
					fmt.Println("Ошибка при декодировании item:", err)
					break
				}
				i++
			}
		}
	}
	/*	for _, item := range res {
		if item.ID != "" && item.ID_PARENT != "" && item.Description != "" {
			fmt.Printf("%+v\n", item)
		}
		if item.Description == catName {
			return item, true
		}
	}*/
	if i == 0 {
		return nil, false
	} else {
		return items, true
	}

}

func getCatalog(catName string) (types.ItemMenu, bool) {
	fmt.Println(base64.StdEncoding.EncodeToString([]byte("TC0051161|API:lt2iZpXb41")))
	req := types.ItemMenuReq{
		Cat_id: "All",
	}
	//var res = make([]types.ItemMenu, 100)
	decoder, err := soap.SoapCallHandleResponse("https://apitest.merlion.com/rl/mlservice3", "https://apitest.merlion.com/rl/mlservice3#getCatalog", req)
	if err != nil {
		log.Fatalf("SoapCallHandleResponse error: %s", err)
	}
	for {
		token, err := decoder.Token()
		if err != nil {
			break
		}
		var item types.ItemMenu
		switch start := token.(type) {
		case xml.StartElement:
			//fmt.Println("Xml tag:", start.Name.Local)
			if start.Name.Local == "item" {
				err := decoder.DecodeElement(&item, &start)
				//fmt.Println(item)
				if item.Description == catName {
					return item, true
				}
				//res = append(res, item)
				if err != nil {
					fmt.Println("Ошибка при декодировании item:", err)
					break
				}
			}
		}
	}
	/*	for _, item := range res {
		if item.ID != "" && item.ID_PARENT != "" && item.Description != "" {
			fmt.Printf("%+v\n", item)
		}
		if item.Description == catName {
			return item, true
		}
	}*/
	return types.ItemMenu{}, false
}
