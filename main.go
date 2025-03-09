package main

import (
	"MerlionScript/service/merlion"
	"MerlionScript/service/sklad"
	"MerlionScript/types/restTypes"
	"fmt"
)

const (
	MerlionCredentials = `TC0051161|API:lt2iZpXb41`
	SkladCredentials   = `admin@sandbox1244:Jf!FQBy!q5"N3]Z`
)

func main() {
	//getCatSklad()
	cat, exist := merlion.GetCatalog("Ноутбуки", MerlionCredentials)
	if !exist {
		fmt.Println("Такой категории не существует (MERLION)")
		return
	}
	fmt.Printf("Категория найдена (MERLION):%v\n", cat)
	metadata, err := sklad.CreateTestCatSklad(restTypes.TestProductGroup{
		Name:        cat.Description,
		Description: cat.Description,
		Code:        cat.ID,
	}, SkladCredentials)
	if err != nil {
		fmt.Printf("Ошибка при создании группы (СКЛАД):%v\n", err)
		return
	}
	fmt.Printf("Метаданные созданной группы (СКЛАД):%v\n", metadata)
	items, exist := merlion.GetCatItems(cat.ID, MerlionCredentials)
	if !exist {
		fmt.Printf("Товаров в категории %s не существует (MERLION)", cat.Description)
		return
	}
	restItems := make([]restTypes.TestProduct, 5)
	for i, rawItem := range items {
		restItems[i].Name = rawItem.Name
		restItems[i].Code = rawItem.No
		restItems[i].Vat = rawItem.VAT
		restItems[i].Weight = int(rawItem.Weight)
		restItems[i].ProductFolder.Meta = metadata.Meta

		err = sklad.CreateTestItemSklad(restItems[i], SkladCredentials)
		if err != nil {
			fmt.Printf("Ошибка при создании товара [%d]\n", i)
		}
	}
}
