package main

import (
	//"MerlionScript/service/sklad"
	//"MerlionScript/types/restTypes"
	"MerlionScript/service/merlion"
	csvInstance "MerlionScript/utils/csv"
	"MerlionScript/utils/db"
	"fmt"
	//"fmt"
)

const (
	MerlionCredentials = `TC0051161|API:lt2iZpXb41`
	//SkladCredentials   = `admin@sandbox1244:Jf!FQBy!q5"N3]Z`
	SkladCredentials = `admin@sandbox1250:Jf!FQBy!q5"N3]Z`
)

func main() {
	//testDB()
	//testCSV()
	importCSV()
	database, _ := db.GetDBInstance()
	records, _ := database.GetCodeRecords()
	counterSuccess := 0
	for _, record := range *records {
		entity := merlion.CreateMerlionEntity(record.Merlion)
		err := entity.FillData(MerlionCredentials)
		if err != nil {
			fmt.Println("error:" + err.Error())
			fmt.Println("entity No: " + entity.No)
			continue
		}
		entity.SendDataToMoySklad(SkladCredentials)
		fmt.Println("entity No: " + entity.No)
		counterSuccess++
	}
	fmt.Printf("counterSuccess: %d\n", counterSuccess)
}

func importCSV() {
	csv, err := csvInstance.GetCSVInstance()
	if err != nil {
		panic(err)
	}
	err = csv.ImportCodes()
	if err != nil {
		panic(err)
	}
}

/*
func testDB() {
	instance, err := DBInstance.GetDBInstance()
	if err != nil {
		panic(err)
	}
	instance.Init()
	defer DBInstance.CloseDB()
	start := time.Now()
	fmt.Println("Getting catalog ids")
	catId := merlion.GetCatalogUniqueCodes(MerlionCredentials)
	fmt.Printf("Get ids: %d records\n", len(catId))
	//length := len(catId)
	bar := progressbar.Default((int64)(len(catId)))
	for _, id := range catId {
		itemsRet := merlion.GetCatItems(id, MerlionCredentials)
		//fmt.Printf("Get items: %d records\n", len(itemsRet))
		//fmt.Printf("\rProgress: %.2f%%", float64(i)/float64(length)*100)
		bar.Add(1)
		for _, item := range itemsRet {
			instance.AddCodeRecord(&typesDB.Codes{MoySklad: "", Manufacturer: item.Vendor_part, Merlion: item.No})
		}
	}
	elapsed := time.Since(start)
	fmt.Printf("Formed by %.3f seconds", elapsed.Seconds())
}
*/
/*
func testCSV() {
	inst, err := csvInstance.GetCSVInstance()
	if err != nil {
		panic(err)
	}
	inst.WriteRecord([]string{"test1", "test2", "test3"})
	inst.WriteRecord([]string{"test4", "test5", "test6"})
	inst.WriteRecord([]string{"test7", "test8", "test9"})
	records, err := inst.ReadAllRecords()
	if err != nil {
		panic(err)
	}
	for _, record := range records {
		fmt.Println("ReadAllRecords:")
		fmt.Println(record)
	}
	records, err = inst.ReadSpecificRows(1, 2)
	if err != nil {
		panic(err)
	}
	for _, record := range records {
		fmt.Println("ReadSpecificRows:")
		fmt.Println(record)
	}
	csvInstance.CloseCSV()
}

func testCreateGoodsFromMerlionToSklad() {
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
*/
