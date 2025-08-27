package src

import (
	"MerlionScript/keeper"
	"MerlionScript/services/common"
	"MerlionScript/utils/db"
	"MerlionScript/utils/db/interfaceDB"
	"MerlionScript/utils/db/typesDB"
	"context"
	"fmt"
	"log"
	"strings"
	"time"
)

// Создаем отсутствующие позиции в ERP или создаем привязки
func CreateNewPositionsERP(ctx context.Context, dbInstance interfaceDB.DB, service common.Service, erpSystem common.ERPSystem) error {
	var ServiceName string = service.GetSystemName()
	select {
	case <-ctx.Done():
		fmt.Printf("%s (CreateNewPositionsERP): работу закончил из-за контекста\n", ServiceName)
		return nil
	default:
		fmt.Printf("%s: начал сверять новые позиции с мс\n", ServiceName)
		Items, err := service.GetItemsList(ctx, dbInstance)
		if err != nil || Items == nil {
			log.Printf("%s (CreateNewPositionsERP): ошибка при получении записей из сервиса | err = %s\n", ServiceName, err)
			return err
		}
		counter, err := dbInstance.GetLastOwnIdMS()
		if err != nil {
			log.Printf("%s (CreateNewPositionsERP): ошибка при получении counter из БД | err = %s\n", ServiceName, err)
			return err
		}
		var createdItems int = 0
		for i := range *Items {
			// if createdItems >= 0 {
			// 	break
			// }
			select {
			case <-ctx.Done():
				return nil
			default:
				//Поиск по товарам МС
				articleReplace := strings.ReplaceAll(IgnoreDHManufacturer((*Items)[i].Article), " ", "+")
				erpItems, err := erpSystem.GetItemsByArticle(articleReplace)
				if err != nil || erpItems == nil {
					log.Printf("%s (CreateNewPositionsERP): ошибка при получении записи из ERP | article = %s; err = %s\n", ServiceName, (*Items)[i].Article, err)
					continue
				}
				//Запись из БД
				erpRecord, err := dbInstance.GetCodesIDs((*Items)[i].Article)
				if err != nil || erpRecord == (typesDB.CodesIDs{}) {
					log.Printf("%s (CreateNewPositionsERP): ошибка при получении записи из БД | article = %s; err = %s\n", ServiceName, (*Items)[i].Article, err)
					continue
				}
				//Если этого товара не существует на МС
				if len(*erpItems) == 0 {
					if err = createPosition(erpRecord, &counter, dbInstance, erpSystem, ServiceName, (*Items)[i]); err != nil {
						continue
					}
					createdItems++
				} else { //Если найдены совпадения/похожие
					//Проверяем существование позиции на МС
					foundedItem, err := CompareArticle(erpItems, &(*Items)[i])
					if err != nil {
						log.Printf("%s: %s\n", ServiceName, err.Error())
						continue
					}
					//Если позиция существует
					if foundedItem != (common.ItemList{}) {
						//foundedItem.Article - код МС, а не артикул позиции
						erpRecord.MoySkladCode = foundedItem.Article
						// Извлекаем counter из артикула
						ownId, _ := db.ExtractCounterFromOwnID(erpRecord.MoySkladCode)
						if ownId > 0 {
							erpRecord.MsOwnId = ownId
						}
						if err = dbInstance.UpdateCodesIDs(erpRecord); err != nil {
							log.Printf("%s (CreateNewPositionsERP): ошибка при изменении записи в БД | serviceCode = %s; err = %s\n", ServiceName, ServiceName, err)
						}
					} else { //Если совпадения не найдены
						// if createdItems >= 10 {
						// 	continue
						// }
						if err = createPosition(erpRecord, &counter, dbInstance, erpSystem, ServiceName, (*Items)[i]); err != nil {
							continue
						}
						createdItems++
					}
				}
			}
			time.Sleep(time.Millisecond * 100)
		}
		fmt.Printf("%s: закончил сверять новые позиции с мс\n", ServiceName)
		fmt.Printf("%s: добавлено %d записей на МС\n", ServiceName, createdItems)
		return nil
	}
}

func createPosition(erpRecord typesDB.CodesIDs, counter *int64, dbInstance interfaceDB.DB, erpSystem common.ERPSystem, ServiceName string, item common.ItemList) error {
	log.Printf("%s (CreateNewPositionsERP): полных соответствий не найдено (создаю новую позицию) | article = %s ", ServiceName, item.Article)
	newId := db.GetFormatID(*counter)
	erpRecord.MoySkladCode = newId
	if err := erpSystem.CreateItem(&erpRecord, newId, item.PositionName, keeper.GetSkladCat()); err != nil {
		log.Printf("%s (CreateNewPositionsERP): ошибка при создании записи на МС | erpCode = %s; err = %s\n", ServiceName, erpRecord.MoySkladCode, err)
		return err
	}
	*counter++
	erpRecord.MsOwnId = *counter
	if err := dbInstance.UpdateCodesIDs(erpRecord); err != nil {
		log.Printf("%s (CreateNewPositionsERP): ошибка при изменении записи в БД | erpCode = %s; err = %s\n", ServiceName, erpRecord.MoySkladCode, err)
	}
	return nil
}
