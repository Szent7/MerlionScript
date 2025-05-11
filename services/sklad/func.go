package sklad

import (
	skladReq "MerlionScript/services/sklad/requests"
	skladTypes "MerlionScript/types/restTypes/sklad"
	"MerlionScript/utils/db"
	"MerlionScript/utils/db/typesDB"
	"fmt"
	"log"
	"time"
)

func GetTimeNow() string {
	return time.Now().Format("2006-01-02 15:04:05")
}

func CreateNewItemMS(item *typesDB.Codes, ItemName string, Catalog skladTypes.Meta, counter *int64, dbInstance *db.DBInstance, tableName string) bool {
	newId := getFormatID(*counter)
	/*itemMerlion, err := merlion.GetItemsByItemId(item.Merlion)
	if err != nil || len(itemMerlion) == 0 || itemMerlion[0].No == "" {
		log.Printf("Ошибка при получении записи из Мерлиона (createNewPositionsMS) merlionCode = %s: %s\n", item.Merlion, err)
		//continue
		return false
	}*/
	//Если мета каталога пустая
	var newItem skladTypes.CreateItem
	if Catalog.Href == "" {
		newItem = skladTypes.CreateItem{
			Name:    ItemName,
			Article: newId,
		}
	} else {
		newItem = skladTypes.CreateItem{
			Name:    ItemName,
			Article: newId,
			ProductFolder: skladTypes.MetaMiddle{
				/*Meta: skladTypes.Meta{
					Href:         "https://api.moysklad.ru/api/remap/1.2/entity/productfolder/db27556f-153a-11f0-0a80-1751000ec1e4",
					MetadataHref: "https://api.moysklad.ru/api/remap/1.2/entity/productfolder/metadata",
					Type:         "productfolder",
					MediaType:    "application/json",
					UuidHref:     "https://online.moysklad.ru/app/#good/edit?id=db27556f-153a-11f0-0a80-1751000ec1e4",
				},*/
				Meta: Catalog,
			},
		}
	}
	response, err := skladReq.CreateItem(newItem)
	if err != nil || response.StatusCode != 200 {
		log.Printf("Ошибка при создании записи на МС (createNewPositionsMS) merlionCode = %s: %s\n", item.Service, err)
		//continue
		return false
	}
	item.MoySklad = newId
	*counter++
	item.MsOwnId = *counter
	if err = dbInstance.EditCodeRecord(item, tableName); err != nil {
		log.Printf("Ошибка при изменении записи в БД (createNewPositionsMS) merlionCode = %s: %s\n", item.Service, err)
	}
	return true
}

// Увеличиваем счетчик для собственных артикулов
func getFormatID(counter int64) string {
	newNumPart := counter
	newID := fmt.Sprintf("I%05d", newNumPart)
	return newID
}
